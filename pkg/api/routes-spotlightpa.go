package api

import (
	"fmt"
	"net/http"

	"github.com/carlmjohnson/emailx"
	"github.com/carlmjohnson/errutil"
	"github.com/carlmjohnson/resperr"
	"github.com/go-chi/chi"
	"github.com/spotlightpa/almanack/internal/arc"
	"github.com/spotlightpa/almanack/internal/db"
	"github.com/spotlightpa/almanack/pkg/almanack"
)

func (app *appEnv) listAllArcStories(w http.ResponseWriter, r *http.Request) {
	var page int32
	app.mustIntParam(r, "page", &page)
	app.Printf("start listAllArcStories page %d", page)

	var (
		resp struct {
			Contents []almanack.ArcStory `json:"contents"`
			NextPage int32               `json:"next_page,omitempty"`
		}
		err error
	)
	resp.Contents, resp.NextPage, err = app.svc.ListAllArcStories(r.Context(), page)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, &resp)
}

func (app *appEnv) listWithArcRefresh(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting listWithArcRefresh")
	type response struct {
		Contents []almanack.ArcStory `json:"contents"`
		NextPage int32               `json:"next_page,omitempty"`
	}
	var (
		feed *arc.API
		err  error
	)
	if feed, err = app.FetchFeed(r.Context()); err != nil {
		// Keep trucking even if you can't load feed
		app.logErr(r.Context(), err)
	} else if err = app.svc.StoreFeed(r.Context(), feed); err != nil {
		app.replyErr(w, r, err)
		return
	}
	var resp response
	resp.Contents, resp.NextPage, err = app.svc.ListAllArcStories(r.Context(), 0)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, resp)
}

func (app *appEnv) postAlmanackArcStory(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting postAlmanackArcStory")

	var userData struct {
		ID         string          `json:"_id"`
		Note       string          `json:"almanack-note,omitempty"`
		Status     almanack.Status `json:"almanack-status,omitempty"`
		RefreshArc bool            `json:"almanack-refresh-arc"`
	}
	if !app.readJSON(w, r, &userData) {
		return
	}

	var (
		story        almanack.ArcStory
		refreshStory bool
	)
	if userData.RefreshArc {
		var (
			feed *arc.API
			err  error
		)
		if feed, err = app.FetchFeed(r.Context()); err != nil {
			app.replyErr(w, r, err)
			return
		}
		if err := app.svc.StoreFeed(r.Context(), feed); err != nil {
			app.replyErr(w, r, err)
			return
		}
		for i := range feed.Contents {
			if feed.Contents[i].ID == userData.ID {
				story.FeedItem = feed.Contents[i]
				refreshStory = true
			}
		}
	}
	story.ID = userData.ID
	story.Note = userData.Note
	story.Status = userData.Status

	if err := app.svc.SaveAlmanackArticle(r.Context(), &story, refreshStory); err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusAccepted, w, &userData)
}

func (app *appEnv) postMessage(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting postMessage")
	type request struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	var req request
	if !app.readJSON(w, r, &req) {
		return
	}
	if err := app.email.SendEmail(req.Subject, req.Body); err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusAccepted, w, http.StatusText(http.StatusAccepted))
}

var supportedContentTypes = map[string]string{
	"image/jpeg": "jpeg",
	"image/png":  "png",
	"image/tiff": "tiff",
}

func (app *appEnv) postSignedUpload(w http.ResponseWriter, r *http.Request) {
	app.Printf("start postSignedUpload")
	var userData struct {
		Type string `json:"type"`
	}
	if !app.readJSON(w, r, &userData) {
		return
	}

	ext, ok := supportedContentTypes[userData.Type]
	if !ok {
		app.replyErr(w, r, resperr.New(
			http.StatusBadRequest, "unsupported content type: %q", ext,
		))
		return
	}

	type response struct {
		SignedURL string `json:"signed-url"`
		FileName  string `json:"filename"`
	}
	var (
		res response
		err error
	)
	res.SignedURL, res.FileName, err = almanack.GetSignedImageUpload(
		r.Context(), app.svc.ImageStore, userData.Type)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	if n, err := app.svc.Queries.CreateImagePlaceholder(r.Context(), db.CreateImagePlaceholderParams{
		Path: res.FileName,
		Type: ext,
	}); err != nil {
		app.replyErr(w, r, err)
		return
	} else if n != 1 {
		// Log and continue
		app.logErr(r.Context(),
			fmt.Errorf("creating image %q but it already exists", res.FileName))
	}
	app.replyJSON(http.StatusOK, w, &res)
}

func (app *appEnv) postImageUpdate(w http.ResponseWriter, r *http.Request) {
	app.Println("start postImageUpdate")

	var userData db.UpdateImageParams
	if !app.readJSON(w, r, &userData) {
		return
	}
	var (
		res db.Image
		err error
	)
	if res, err = app.svc.Queries.UpdateImage(r.Context(), userData); err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, &res)
}

func (app *appEnv) listDomains(w http.ResponseWriter, r *http.Request) {
	app.Println("start listDomains")
	type response struct {
		Domains []string `json:"domains"`
	}

	domains, err := app.svc.Queries.ListDomainsWithRole(r.Context(), "editor")
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, response{
		domains,
	})
}

func (app *appEnv) postDomain(w http.ResponseWriter, r *http.Request) {
	app.Println("start postDomain")
	type request struct {
		Domain string `json:"domain"`
		Remove bool   `json:"remove"`
	}
	type response struct {
		Domains []string `json:"domains"`
	}
	var req request
	if !app.readJSON(w, r, &req) {
		return
	}

	if req.Domain == "" {
		app.replyErr(w, r, resperr.New(http.StatusBadRequest,
			"can't add nothing!"))
		return
	}

	if req.Domain == "spotlightpa.org" {
		app.replyErr(w, r, resperr.New(http.StatusBadRequest,
			"can't change spotlightpa.org!"))
		return
	}

	var roles []string
	if !req.Remove {
		roles = []string{"editor"}
	}

	if _, err := app.svc.Queries.SetRolesForDomain(
		r.Context(),
		db.SetRolesForDomainParams{
			Domain: req.Domain,
			Roles:  roles,
		},
	); err != nil {
		app.replyErr(w, r, err)
		return
	}

	domains, err := app.svc.Queries.ListDomainsWithRole(r.Context(), "editor")
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, response{
		domains,
	})
}

func (app *appEnv) listAddresses(w http.ResponseWriter, r *http.Request) {
	app.Println("start listAddresses")
	var (
		resp struct {
			Addresses []string `json:"addresses"`
		}
		err error
	)
	resp.Addresses, err = app.svc.Queries.ListAddressesWithRole(r.Context(), "editor")
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, resp)
}

func (app *appEnv) postAddress(w http.ResponseWriter, r *http.Request) {
	app.Println("start postAddresses")
	type request struct {
		Address string `json:"address"`
		Remove  bool   `json:"remove"`
	}
	type response struct {
		Addresses []string `json:"addresses"`
	}
	var req request
	if !app.readJSON(w, r, &req) {
		return
	}

	if !emailx.Valid(req.Address) {
		app.replyErr(w, r, resperr.New(http.StatusBadRequest,
			"invalid email address: %q", req.Address))
		return
	}

	var roles []string
	if !req.Remove {
		roles = []string{"editor"}
	}

	if _, err := app.svc.Queries.SetRolesForAddress(
		r.Context(),
		db.SetRolesForAddressParams{
			EmailAddress: req.Address,
			Roles:        roles,
		},
	); err != nil {
		app.replyErr(w, r, err)
		return
	}

	var (
		resp response
		err  error
	)
	resp.Addresses, err = app.svc.Queries.ListAddressesWithRole(r.Context(), "editor")
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, resp)
}

func (app *appEnv) listImages(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting listImages")
	type response struct {
		Images []db.Image `json:"images"`
	}
	var (
		res response
		err error
	)

	if res.Images, err = app.svc.Queries.ListImages(r.Context(), db.ListImagesParams{
		Offset: 0,
		Limit:  100,
	}); err != nil {
		app.replyErr(w, r, err)
		return
	}

	app.replyJSON(http.StatusOK, w, res)
}

func (app *appEnv) getEditorsPicks(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting getEditorsPicks")
	type response struct {
		Datums []db.SiteDatum `json:"datums"`
	}
	var (
		res response
		err error
	)
	res.Datums, err = app.svc.Queries.GetSiteData(r.Context(), almanack.EditorsPicksLoc)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, res)
}

func (app *appEnv) postEditorsPicks(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting postEditorsPicks")

	var req struct {
		Datums []almanack.ScheduledSiteConfig `json:"datums"`
	}
	if !app.readJSON(w, r, &req) {
		return
	}
	if len(req.Datums) < 1 {
		app.replyErr(w, r, resperr.New(
			http.StatusBadRequest, "no schedulable items provided"))
		return
	}

	var (
		res struct {
			Datums []db.SiteDatum `json:"datums"`
		}
		err error
	)
	res.Datums, err = app.svc.UpdateSiteConfig(r.Context(), almanack.EditorsPicksLoc, req.Datums)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}

	app.replyJSON(http.StatusOK, w, res)
}

func (app *appEnv) listAllTopics(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting listAllTopics")
	t, err := app.svc.Queries.ListAllTopics(r.Context())
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, struct {
		Topics []string `json:"topics"`
	}{t})
}

func (app *appEnv) listAllSeries(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting listAllSeries")
	s, err := app.svc.Queries.ListAllSeries(r.Context())
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, struct {
		Series []string `json:"series"`
	}{s})
}

func (app *appEnv) listFiles(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting listFiles")
	type response struct {
		Files []db.File `json:"files"`
	}
	var (
		res response
		err error
	)

	if res.Files, err = app.svc.Queries.ListFiles(r.Context(), db.ListFilesParams{
		Offset: 0,
		Limit:  100,
	}); err != nil {
		app.replyErr(w, r, err)
		return
	}

	app.replyJSON(http.StatusOK, w, res)
}

func (app *appEnv) postFileCreate(w http.ResponseWriter, r *http.Request) {
	app.Printf("start postFileCreate")
	var userData struct {
		MimeType string `json:"mimeType"`
		FileName string `json:"filename"`
	}
	if !app.readJSON(w, r, &userData) {
		return
	}
	type response struct {
		SignedURL    string `json:"signed-url"`
		FileURL      string `json:"file-url"`
		Disposition  string `json:"disposition"`
		CacheControl string `json:"cache-control"`
	}
	var (
		res response
		err error
	)

	res.SignedURL, res.FileURL, res.Disposition, res.CacheControl, err = almanack.GetSignedFileUpload(
		r.Context(),
		app.svc.FileStore,
		userData.FileName,
		userData.MimeType,
	)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	if n, err := app.svc.Queries.CreateFilePlaceholder(r.Context(),
		db.CreateFilePlaceholderParams{
			Filename: userData.FileName,
			Type:     userData.MimeType,
			URL:      res.FileURL,
		}); err != nil {
		app.replyErr(w, r, err)
		return
	} else if n != 1 {
		// Log and continue
		app.logErr(r.Context(),
			fmt.Errorf("creating file %q but it already exists", res.FileURL))
	}
	app.replyJSON(http.StatusOK, w, &res)
}

func (app *appEnv) postFileUpdate(w http.ResponseWriter, r *http.Request) {
	app.Println("start postFileUpdate")

	var userData db.UpdateFileParams
	if !app.readJSON(w, r, &userData) {
		return
	}
	var (
		res db.File
		err error
	)
	if res, err = app.svc.Queries.UpdateFile(r.Context(), userData); err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, &res)
}

func (app *appEnv) listNewsPages(w http.ResponseWriter, r *http.Request) {
	var page int32
	app.mustIntParam(r, "page", &page)

	app.Printf("start listNewsPages page %d", page)

	var (
		resp struct {
			Pages    []db.ListPagesRow `json:"pages"`
			NextPage int32             `json:"next_page,omitempty"`
		}
		err error
	)
	const pageSize = 100
	offset := page * pageSize

	resp.Pages, err = app.svc.Queries.ListPages(r.Context(), db.ListPagesParams{
		FilePath: "content/news/%",
		Limit:    pageSize + 1,
		Offset:   offset,
	})

	if len(resp.Pages) > pageSize {
		resp.Pages = resp.Pages[:pageSize]
		resp.NextPage = page + 1
	}
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, &resp)
}

func (app *appEnv) listNewsletterPages(w http.ResponseWriter, r *http.Request) {
	var page int32
	app.mustIntParam(r, "page", &page)
	app.Printf("start listNewsletterPages page %d", page)

	var (
		resp struct {
			Pages    []db.ListPagesRow `json:"pages"`
			NextPage int32             `json:"next_page,omitempty"`
		}
		err error
	)
	const pageSize = 100
	offset := page * pageSize

	resp.Pages, err = app.svc.Queries.ListPages(r.Context(), db.ListPagesParams{
		FilePath: "content/newsletters/%",
		Limit:    pageSize + 1,
		Offset:   offset,
	})

	if len(resp.Pages) > pageSize {
		resp.Pages = resp.Pages[:pageSize]
		resp.NextPage = page + 1
	}
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, &resp)
}

func (app *appEnv) getPage(w http.ResponseWriter, r *http.Request) {
	var id int64
	app.mustIntParam(r, "id", &id)
	app.Printf("start getPage for %d", id)
	page, err := app.svc.Queries.GetPageByID(r.Context(), id)
	if err != nil {
		err = db.NoRowsAs404(err, "could not find page ID %d", id)
		app.replyErr(w, r, err)
		return
	}

	app.replyJSON(http.StatusOK, w, page)
}

func (app *appEnv) getPageWithContent(w http.ResponseWriter, r *http.Request) {
	var id int64
	app.mustIntParam(r, "id", &id)
	app.Printf("start getPage for %d", id)
	page, err := app.svc.Queries.GetPageByID(r.Context(), id)
	if err != nil {
		err = db.NoRowsAs404(err, "could not find page ID %d", id)
		app.replyErr(w, r, err)
		return
	}
	if warning := app.svc.RefreshPageFromContentStore(r.Context(), &page); warning != nil {
		app.logErr(r.Context(), warning)
	}
	app.replyJSON(http.StatusOK, w, page)
}

func (app *appEnv) postPage(w http.ResponseWriter, r *http.Request) {
	app.Printf("start postPage")

	var userUpdate db.UpdatePageParams
	if !app.readJSON(w, r, &userUpdate) {
		return
	}

	oldPage, err := app.svc.Queries.GetPageByFilePath(r.Context(), userUpdate.FilePath)
	if err != nil {
		errutil.Prefix(&err, "postPage connection problem")
		app.replyErr(w, r, err)
		return
	}

	res, err := app.svc.Queries.UpdatePage(r.Context(), userUpdate)
	if err != nil {
		errutil.Prefix(&err, "postPage update problem")
		app.replyErr(w, r, err)
		return
	}
	shouldPublish := res.ShouldPublish()
	shouldNotify := res.ShouldNotify(&oldPage)
	if shouldNotify {
		if err = app.svc.Notify(r.Context(), &res, shouldPublish); err != nil {
			app.logErr(r.Context(), err)
		}
	}
	if shouldPublish {
		err, warning := app.svc.PublishPage(r.Context(), &res)
		if warning != nil {
			app.logErr(r.Context(), warning)
		}
		if err != nil {
			errutil.Prefix(&err, "postPage publish problem")
			app.replyErr(w, r, err)
			return
		}
	}
	app.replyJSON(http.StatusOK, w, &res)
}

func (app *appEnv) postPageForArcID(w http.ResponseWriter, r *http.Request) {
	arcID := chi.URLParam(r, "arcID")
	app.Printf("start postPageForArcID for %s", arcID)

	dbArt, err := app.svc.Queries.GetArticleByArcID(r.Context(), arcID)
	if err != nil {
		err = db.NoRowsAs404(err, "could not find Arc ID %q", arcID)
		app.replyErr(w, r, err)
		return
	}
	filepath := dbArt.SpotlightPAPath.String
	firstConversion := filepath == ""
	if !firstConversion {
		page, err := app.svc.Queries.GetPageByFilePath(r.Context(), filepath)
		if err != nil {
			app.replyErr(w, r, err)
			return
		}
		app.replyJSON(http.StatusOK, w, page.ID)
		return
	}
	app.Printf("saving page from %s", arcID)
	page, err := app.svc.PageFromArcArticle(r.Context(), &dbArt)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}

	app.replyJSON(http.StatusOK, w, page.ID)
}

func (app *appEnv) postRefreshPageFromArc(w http.ResponseWriter, r *http.Request) {
	var id int64
	app.mustIntParam(r, "id", &id)
	app.Printf("start postRefreshPageFromArc for %d", id)

	page, err := app.svc.Queries.GetPageByID(r.Context(), id)
	if err != nil {
		err = db.NoRowsAs404(err, "could not find page ID %d", id)
		app.replyErr(w, r, err)
		return
	}

	arcID, _ := page.Frontmatter["arc-id"].(string)
	if arcID == "" {
		app.replyErr(w, r, resperr.New(http.StatusBadRequest,
			"no arc-id on page %d", id))
		return
	}
	var story *almanack.ArcStory
	feed, feedErr := app.FetchFeed(r.Context())
	if feedErr != nil {
		// Keep trucking even if you can't load feed
		app.logErr(r.Context(), feedErr)
	}
	if feedErr == nil {
		if err = app.svc.StoreFeed(r.Context(), feed); err != nil {
			app.replyErr(w, r, err)
			return
		}
		for _, rawStory := range feed.Contents {
			if rawStory.ID == arcID {
				app.Printf("Arc ID %q found in feed", arcID)
				story = &almanack.ArcStory{FeedItem: rawStory}
				break
			}
		}
	}
	if story == nil {
		app.Printf("Arc ID %q not found in feed; trying DB", arcID)
		dbArt, err := app.svc.Queries.GetArticleByArcID(r.Context(), arcID)
		if err != nil {
			err = db.NoRowsAs404(err, "could not find Arc ID %q", arcID)
			app.replyErr(w, r, err)
			return
		}
		story, err = almanack.ArcStoryFromDB(&dbArt)
		if err != nil {
			app.replyErr(w, r, err)
			return
		}
	}

	if err = app.svc.RefreshPageFromArcStory(r.Context(), story, &page); err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, page)
}

func (app *appEnv) listAllPages(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting listSpotlightPAArticles")
	type response struct {
		Pages []db.ListAllPagesRow `json:"pages"`
	}
	var (
		res response
		err error
	)

	if res.Pages, err = app.svc.Queries.ListAllPages(r.Context()); err != nil {
		app.replyErr(w, r, err)
		return
	}

	app.replyJSON(http.StatusOK, w, res)
}

func (app *appEnv) getSiteParams(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting getSiteParams")

	type response struct {
		Params []db.SiteDatum `json:"params"`
	}
	var (
		res response
		err error
	)
	res.Params, err = app.svc.Queries.GetSiteData(r.Context(), almanack.SiteParamsLoc)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}
	app.replyJSON(http.StatusOK, w, res)
}

func (app *appEnv) postSiteParams(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting postSiteParams")

	var req struct {
		Params []almanack.ScheduledSiteConfig `json:"params"`
	}
	if !app.readJSON(w, r, &req) {
		return
	}
	if len(req.Params) < 1 {
		app.replyErr(w, r, resperr.New(
			http.StatusBadRequest, "no schedulable items provided"))
		return
	}

	var (
		res struct {
			Params []db.SiteDatum `json:"params"`
		}
		err error
	)
	res.Params, err = app.svc.UpdateSiteConfig(r.Context(), almanack.SiteParamsLoc, req.Params)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}

	app.replyJSON(http.StatusOK, w, res)
}
