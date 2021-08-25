package api

import (
	"html/template"
	"net/http"

	"github.com/carlmjohnson/errutil"
	"github.com/carlmjohnson/resperr"
	"github.com/spotlightpa/almanack/internal/db"
	"github.com/spotlightpa/almanack/layouts"
)

func (app *appEnv) renderNotFound(w http.ResponseWriter, r *http.Request) {
	app.replyHTMLErr(w, r, resperr.NotFound(r))
}

func (app *appEnv) renderPage(w http.ResponseWriter, r *http.Request) {
	id, err := app.getIntParam(r, "id")
	if err != nil {
		errutil.Prefix(&err, "bad argument to renderPage")
		app.replyHTMLErr(w, r, err)
		return
	}
	app.Printf("start renderPage for %d", id)
	page, err := app.svc.Queries.GetPageByID(r.Context(), id)
	if err != nil {
		err = db.NoRowsAs404(err, "could not find page ID %d", id)
		app.replyHTMLErr(w, r, err)
		return
	}
	title, _ := page.Frontmatter["title"].(string)
	body, _ := page.Frontmatter["raw-content"].(string)

	app.replyHTML(w, r, layouts.MailChimp, struct {
		Title string
		Body  template.HTML
	}{
		Title: title,
		Body:  template.HTML(body),
	})
}
