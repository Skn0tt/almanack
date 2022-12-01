package api

import (
	"net/http"

	"github.com/carlmjohnson/resperr"
	"github.com/go-chi/chi/v5"
	"github.com/spotlightpa/almanack/internal/db"
	"github.com/spotlightpa/almanack/internal/netlifyid"
	"github.com/spotlightpa/almanack/internal/paginate"
	"github.com/spotlightpa/almanack/pkg/almanack"
)

func (app *appEnv) userInfo(w http.ResponseWriter, r *http.Request) {
	app.Println("start userInfo")
	userinfo := netlifyid.FromContext(r.Context())
	app.replyJSON(http.StatusOK, w, userinfo)
}

func (app *appEnv) listAvailableArcStories(w http.ResponseWriter, r *http.Request) {
	var page int32
	mustIntParam(r, "page", &page)

	app.Printf("starting listAvailableArcStories page %d", page)

	var (
		res struct {
			Contents []almanack.ArcStory `json:"contents"`
			NextPage int32               `json:"next_page,string,omitempty"`
		}
		err error
	)
	if res.Contents, res.NextPage, err = app.svc.ListAvailableArcStories(
		r.Context(), page,
	); err != nil {
		app.replyErr(w, r, err)
		return
	}

	app.replyJSON(http.StatusOK, w, res)
}

func (app *appEnv) getArcStory(w http.ResponseWriter, r *http.Request) {
	articleID := chi.URLParam(r, "id")
	app.Printf("starting getArcStory %s", articleID)

	article, err := app.svc.GetArcStory(r.Context(), articleID)
	if err != nil {
		app.replyErr(w, r, err)
		return
	}

	if article.Status != almanack.StatusAvailable &&
		article.Status != almanack.StatusPlanned {
		// Let Spotlight PA users get article regardless of its status
		if err := app.auth.HasRole(r, "Spotlight PA"); err != nil {
			app.replyNewErr(http.StatusNotFound, w, r,
				"user unauthorized to view article: %w", err)
			return
		}
	}

	app.replyJSON(http.StatusOK, w, article)
}

func (app *appEnv) getSignupURL(w http.ResponseWriter, r *http.Request) {
	app.Println("start getSignupURL")
	app.replyJSON(http.StatusOK, w, app.svc.MailchimpSignupURL)
}

func (app *appEnv) listSharedArticles(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting listSharedArticles")

	var page int32
	_ = intFromQuery(r, "page", &page)
	if page < 0 {
		app.replyErr(w, r, resperr.WithUserMessage(nil, "Invalid page"))
		return
	}

	// Spotlight PA users can add ?show=all
	if r.URL.Query().Get("show") == "all" {
		if err := app.auth.HasRole(r, "Spotlight PA"); err != nil {
			app.replyErr(w, r, err)
			return
		}

		pager := paginate.PageNumber(page)
		pager.PageSize = 100
		stories, err := paginate.List(pager, r.Context(),
			app.svc.Queries.ListSharedArticlesWhereActive,
			db.ListSharedArticlesWhereActiveParams{
				Offset: pager.Offset(),
				Limit:  pager.Limit(),
			})
		if err != nil {
			app.replyErr(w, r, err)
			return
		}

		app.replyJSON(http.StatusOK, w, struct {
			Stories  []db.SharedArticle `json:"stories"`
			NextPage int32              `json:"next_page,string,omitempty"`
		}{
			Stories:  stories,
			NextPage: pager.NextPage,
		})
	}

	pager := paginate.PageNumber(page)
	pager.PageSize = 20
	stories, err := paginate.List(pager, r.Context(),
		app.svc.Queries.ListSharedArticles,
		db.ListSharedArticlesParams{
			Offset: pager.Offset(),
			Limit:  pager.Limit(),
		})
	if err != nil {
		app.replyErr(w, r, err)
		return
	}

	app.replyJSON(http.StatusOK, w, struct {
		Stories  []db.SharedArticle `json:"stories"`
		NextPage int32              `json:"next_page,string,omitempty"`
	}{
		Stories:  stories,
		NextPage: pager.NextPage,
	})
}

func (app *appEnv) getSharedArticle(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting getSharedArticle")
	var id int64
	if !intFromQuery(r, "id", &id) {
		app.replyErr(w, r, resperr.WithUserMessage(nil,
			"Must provide article ID"))
		return
	}

	article, err := app.svc.Queries.GetSharedArticleByID(r.Context(), id)
	if err != nil {
		err = db.NoRowsAs404(err,
			"missing shared_article id = %d", id)
		app.replyErr(w, r, err)
		return
	}

	if article.Status != "S" &&
		article.Status != "P" {
		// Let Spotlight PA users get article regardless of its status
		if err := app.auth.HasRole(r, "Spotlight PA"); err != nil {
			app.replyNewErr(http.StatusNotFound, w, r,
				"user unauthorized to view article: %w", err)
			return
		}
	}

	app.replyJSON(http.StatusOK, w, article)
}
