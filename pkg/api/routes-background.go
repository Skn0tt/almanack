package api

import (
	"net/http"
	"time"

	"github.com/carlmjohnson/errutil"
	"github.com/go-chi/chi/v5"
	"github.com/spotlightpa/almanack/internal/db"
	"github.com/spotlightpa/almanack/internal/paginate"
	"github.com/spotlightpa/almanack/pkg/almanack"
	"golang.org/x/sync/errgroup"
)

func (app *appEnv) backgroundSleep(w http.ResponseWriter, r *http.Request) {
	app.Println("start backgroundSleep")
	if deadline, ok := r.Context().Deadline(); ok {
		app.Printf("deadline: %s", deadline.Format(time.RFC1123))
	} else {
		app.Printf("no deadline")
	}
	durationStr := chi.URLParam(r, "duration")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	time.Sleep(duration)
	app.replyJSON(http.StatusOK, w, struct {
		SleptFor time.Duration `json:"slept-for"`
	}{duration})
}

func (app *appEnv) backgroundCron(w http.ResponseWriter, r *http.Request) {
	app.Println("start background cron")

	if err := errutil.ExecParallel(func() error {
		var errs errutil.Slice
		// Publish any scheduled pages before pushing new site config
		poperr, warning := app.svc.PopScheduledPages(r.Context())
		if warning != nil {
			app.logErr(r.Context(), warning)
		}
		errs.Push(poperr)
		// TODO: Query all locations from DB side
		errs.Push(app.svc.PopScheduledSiteChanges(r.Context(), almanack.ElectionFeatLoc))
		errs.Push(app.svc.PopScheduledSiteChanges(r.Context(), almanack.HomepageLoc))
		errs.Push(app.svc.PopScheduledSiteChanges(r.Context(), almanack.SidebarLoc))
		errs.Push(app.svc.PopScheduledSiteChanges(r.Context(), almanack.SiteParamsLoc))
		errs.Push(app.svc.PopScheduledSiteChanges(r.Context(), almanack.StateCollegeLoc))
		return errs.Merge()
	}, func() error {
		return app.svc.UpdateMostPopular(r.Context())
	}, func() error {
		types, err := app.svc.Queries.ListNewsletterTypes(r.Context())
		if err != nil {
			return err
		}
		var errs errutil.Slice
		// Update newsletter archives first and then import anything new
		errs.Push(app.svc.UpdateNewsletterArchives(r.Context(), types))
		errs.Push(app.svc.ImportNewsletterPages(r.Context(), types))
		return errs.Merge()
	}, func() error {
		app.backgroundImages(w, r)
		return nil
	}); err != nil {
		// reply shows up in dev only
		app.replyErr(w, r, err)
		return
	}

	app.replyJSON(http.StatusAccepted, w, "OK")
}

func (app *appEnv) backgroundRefreshPages(w http.ResponseWriter, r *http.Request) {
	app.Println("start backgroundRefreshPages")

	hasMore := true
	for queryPage := int32(0); hasMore; queryPage++ {
		pager := paginate.PageNumber(queryPage)
		pager.PageSize = 10
		pageIDs, err := paginate.List(
			pager, r.Context(),
			app.svc.Queries.ListPageIDs,
			db.ListPageIDsParams{
				FilePath: "content/news/%",
				Offset:   pager.Offset(),
				Limit:    pager.Limit(),
			})
		if err != nil {
			app.replyErr(w, r, err)
			return
		}
		for _, id := range pageIDs {
			if err := app.svc.RefreshPageContents(r.Context(), id); err != nil {
				app.replyErr(w, r, err)
				return
			}
		}
		hasMore = pager.HasMore()
	}

	app.replyJSON(http.StatusAccepted, w, "OK")
}

func (app *appEnv) backgroundImages(w http.ResponseWriter, r *http.Request) {
	app.Println("start backgroundImages")

	images, err := app.svc.Queries.ListImageWhereNotUploaded(r.Context())
	if err != nil {
		app.replyErr(w, r, err)
		return
	}

	var eg errgroup.Group
	eg.SetLimit(5)
	for i := range images {
		image := images[i]
		eg.Go(func() error {
			_, _, err := almanack.UploadFromURL(
				r.Context(),
				app.svc.Client,
				app.svc.ImageStore,
				image.Path,
				image.SourceURL)
			if err != nil {
				app.logErr(r.Context(), err)
				return nil
			}
			if _, err = app.svc.Queries.UpdateImage(r.Context(),
				db.UpdateImageParams{
					Path:      image.Path,
					SourceURL: image.SourceURL,
				}); err != nil {
				app.logErr(r.Context(), err)
				return nil
			}
			return nil
		})
	}
	eg.Wait()
	app.replyJSON(http.StatusAccepted, w, "OK")
}
