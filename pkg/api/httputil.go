package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/carlmjohnson/resperr"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi"

	"github.com/spotlightpa/almanack/internal/httpjson"
	"github.com/spotlightpa/almanack/pkg/almanack"
)

func (app *appEnv) replyJSON(statusCode int, w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		app.Printf("jsonResponse problem: %v", err)
	}
}

func (app *appEnv) replyErr(w http.ResponseWriter, r *http.Request, err error) {
	app.logErr(r.Context(), err)
	code := resperr.StatusCode(err)
	msg := resperr.UserMessage(err)
	app.replyJSON(code, w, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		code,
		msg,
	})
}

func (app *appEnv) logErr(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException(err)
	} else {
		app.Printf("sentry not in context")
	}

	app.Printf("err: %v", err)
}

func (app *appEnv) versionMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Almanack-App-Version", almanack.BuildVersion)
		h.ServeHTTP(w, r)
	})
}

func (app *appEnv) authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, err := app.auth.AddToRequest(r)
		if err != nil {
			app.replyErr(w, r, err)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (app *appEnv) hasRoleMiddleware(role string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := app.auth.HasRole(r, role); err != nil {
				app.replyErr(w, r, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *appEnv) maxSizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			megabyte = 1 << 20
			maxSize  = 5 * megabyte
		)
		r2 := *r // shallow copy
		r2.Body = http.MaxBytesReader(w, r.Body, maxSize)
		next.ServeHTTP(w, &r2)
	})
}

func (app *appEnv) getPage(r *http.Request, route string) (page int, err error) {
	if pageStr := chi.URLParam(r, "page"); pageStr != "" {
		if page, err = strconv.Atoi(pageStr); err != nil {
			err = resperr.New(http.StatusBadRequest,
				"bad argument to %s: %w", route, err)
			return
		}
	}
	return
}

func (app *appEnv) FetchFeed(ctx context.Context) (*almanack.ArcAPI, error) {
	var feed almanack.ArcAPI
	// Timeout needs to leave enough time to report errors to Sentry before
	// AWS kills the Lambda…
	ctx, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()

	if err := httpjson.Get(ctx, app.svc.Client, app.srcFeedURL, &feed); err != nil {
		return nil, resperr.New(
			http.StatusBadGateway, "could not fetch Arc feed: %w", err)
	}
	return &feed, nil
}
