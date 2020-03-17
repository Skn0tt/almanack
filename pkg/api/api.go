package api

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/carlmjohnson/flagext"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/peterbourgon/ff/v2"
	"github.com/piotrkubisa/apigo"

	"github.com/spotlightpa/almanack/internal/aws"
	"github.com/spotlightpa/almanack/internal/db"
	"github.com/spotlightpa/almanack/internal/github"
	"github.com/spotlightpa/almanack/internal/herokuapi"
	"github.com/spotlightpa/almanack/internal/httpcache"
	"github.com/spotlightpa/almanack/internal/httpjson"
	"github.com/spotlightpa/almanack/internal/mailchimp"
	"github.com/spotlightpa/almanack/internal/netlifyid"
	"github.com/spotlightpa/almanack/pkg/almanack"
	"github.com/spotlightpa/almanack/pkg/errutil"
)

const AppName = "almanack-api"

func CLI(args []string) error {
	var app appEnv
	if err := app.parseArgs(args); err != nil {
		fmt.Fprintf(os.Stderr, "Startup error: %v\n", err)
		return err
	}
	if err := app.exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return err
	}
	return nil
}

func (app *appEnv) parseArgs(args []string) error {
	fl := flag.NewFlagSet(AppName, flag.ContinueOnError)

	pg := db.FlagVar(fl, "postgres", "PostgreSQL database `URL`")
	fl.StringVar(&app.srcFeedURL, "src-feed", "", "source `URL` for Arc feed")
	cache := fl.Bool("cache", false, "use in-memory cache for fetched JSON")
	fl.BoolVar(&app.isLambda, "lambda", false, "use AWS Lambda rather than HTTP")
	fl.StringVar(&app.port, "port", ":3001", "listen on port (HTTP only)")
	fl.StringVar(&app.mailchimpSignupURL, "mc-signup-url", "http://example.com", "`URL` to redirect users to for MailChimp signup")
	checkHerokuPG := herokuapi.FlagVar(fl, "postgres")
	app.Logger = log.New(nil, AppName+" ", log.LstdFlags)
	fl.Var(
		flagext.Logger(app.Logger, flagext.LogSilent),
		"silent",
		`don't log debug output`,
	)
	getImageStore := aws.FlagVar(fl)
	mcAPIKey := fl.String("mc-api-key", "", "API `key` for MailChimp")
	mcListID := fl.String("mc-list-id", "", "List `ID` MailChimp campaign")
	sentryDSN := fl.String("sentry-dsn", "", "DSN `pseudo-URL` for Sentry")
	getGithub := github.FlagVar(fl)
	fl.Usage = func() {
		fmt.Fprintf(fl.Output(), "almanack-api help\n\n")
		fl.PrintDefaults()
	}
	if err := ff.Parse(fl, args, ff.WithEnvVarPrefix("ALMANACK")); err != nil {
		return err
	}

	if err := app.initSentry(*sentryDSN); err != nil {
		return err
	}

	// Get PostgreSQL URL from Heroku if possible, else get it from flag
	if usedHeroku, err := checkHerokuPG(); err != nil {
		return err
	} else if usedHeroku {
		app.Logger.Printf("got credentials from Heroku")
	} else {
		app.Logger.Printf("did not get credentials Heroku")
	}

	if *pg == nil {
		err := errors.New("must set postgres URL")
		app.Logger.Printf("starting up: %v", err)
		return err
	}

	app.email = mailchimp.NewMailService(*mcAPIKey, *mcListID, app.Logger)
	app.imageStore = getImageStore(app.Logger)
	app.auth = netlifyid.NewService(app.isLambda, app.Logger)
	app.c = http.DefaultClient
	if *cache {
		httpcache.SetRounderTripper(app.c, app.Logger)
	}
	if gh, err := getGithub(app.Logger); err != nil {
		app.Logger.Printf("could not connect to Github: %v", err)
		return err
	} else {
		app.gh = gh
	}
	app.svc = almanack.Service{
		Querier:      *pg,
		Logger:       app.Logger,
		ContentStore: app.gh,
	}

	return nil
}

type appEnv struct {
	srcFeedURL         string
	port               string
	isLambda           bool
	mailchimpSignupURL string
	c                  *http.Client
	auth               almanack.AuthService
	gh                 almanack.ContentStore
	imageStore         almanack.ImageStore
	email              almanack.EmailService
	svc                almanack.Service
	*log.Logger
}

func (app *appEnv) exec() error {
	app.Printf("starting %s (%s)", AppName, almanack.BuildVersion)

	listener := http.ListenAndServe
	if app.isLambda {
		app.Printf("starting on AWS Lambda")
		apigo.ListenAndServe("", app.routes())
		panic("unreachable")
	}

	app.Printf("starting on port %s", app.port)
	return listener(app.port, app.routes())
}

func (app *appEnv) initSentry(dsn string) error {
	var transport sentry.Transport
	if app.isLambda {
		transport = &sentry.HTTPSyncTransport{Timeout: 1 * time.Second}
	}
	return sentry.Init(sentry.ClientOptions{
		Dsn:       dsn,
		Release:   almanack.BuildVersion,
		Transport: transport,
	})
}

func (app *appEnv) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: app.Logger}))
	r.Use(app.versionMiddleware)
	r.Get("/api/healthcheck", app.ping)
	r.Route("/api", func(r chi.Router) {
		r.Use(app.authMiddleware)
		r.Get("/user-info", app.userInfo)
		r.With(
			app.hasRoleMiddleware("editor"),
		).Group(func(r chi.Router) {
			r.Get("/available-articles", app.listAvailable)
			r.Get("/available-articles/{id}", app.getAvailable)
			r.Get("/mailchimp-signup-url", app.getSignupURL)
		})
		r.With(
			app.hasRoleMiddleware("Spotlight PA"),
		).Group(func(r chi.Router) {
			r.Get("/upcoming-articles", app.listUpcoming)
			r.Post("/available-articles", app.postAvailable)
			r.Post("/message", app.postMessage)
			r.Get("/scheduled-articles/{id}", app.getScheduledArticle)
			r.Post("/scheduled-articles", app.postScheduledArticle)
			r.Post("/get-signed-upload", app.getSignedUpload)
		})
	})

	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	return sentryHandler.Handle(r)
}

func (app *appEnv) versionMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Almanack-App-Version", almanack.BuildVersion)
		h.ServeHTTP(w, r)
	})
}

func (app *appEnv) jsonResponse(statusCode int, w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		app.Printf("jsonResponse problem: %v", err)
	}
}

func (app *appEnv) errorResponse(ctx context.Context, w http.ResponseWriter, err error) {
	var errResp errutil.Response
	if !errors.As(err, &errResp) {
		errResp.StatusCode = http.StatusInternalServerError
		errResp.Message = "internal error"
		errResp.Log = err.Error()
	}
	app.Println(errResp.Log)
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException(errResp)
	}
	app.jsonResponse(errResp.StatusCode, w, errResp)
}

func (app *appEnv) ping(w http.ResponseWriter, r *http.Request) {
	app.Println("start ping")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", "public, max-age=60")
	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}
	w.Write(b)
}

func (app *appEnv) authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, err := app.auth.AddToRequest(r)
		if err != nil {
			app.errorResponse(r.Context(), w, err)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (app *appEnv) userInfo(w http.ResponseWriter, r *http.Request) {
	app.Println("start userInfo")
	userinfo, err := netlifyid.FromRequest(r)
	if err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}
	app.jsonResponse(http.StatusOK, w, userinfo)
}

func (app *appEnv) hasRoleMiddleware(role string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := app.auth.HasRole(r, role); err != nil {
				app.errorResponse(r.Context(), w, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *appEnv) listUpcoming(w http.ResponseWriter, r *http.Request) {
	app.Println("start listUpcoming")

	var (
		feed almanack.ArcAPI
		err  error
	)
	feed.Contents, err = app.svc.ListAllArticles(r.Context())
	if err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}
	app.jsonResponse(http.StatusOK, w, feed)
}

func (app *appEnv) postAvailable(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting postAvailable")

	var userData struct {
		ID         string          `json:"_id"`
		Note       string          `json:"almanack-note,omitempty"`
		Status     almanack.Status `json:"almanack-status,omitempty"`
		RefreshArc bool            `json:"almanack-refresh-arc"`
	}
	if err := httpjson.DecodeRequest(w, r, &userData); err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}

	var story almanack.ArcStory
	if userData.RefreshArc {
		var feed almanack.ArcAPI
		if err := httpjson.Get(r.Context(), app.c, app.srcFeedURL, &feed); err != nil {
			app.errorResponse(r.Context(), w, err)
			return
		}
		if err := app.svc.StoreFeed(r.Context(), &feed); err != nil {
			app.errorResponse(r.Context(), w, err)
			return
		}
		for i := range feed.Contents {
			if feed.Contents[i].ID == userData.ID {
				story = feed.Contents[i]
			}
		}
	}
	story.ID = userData.ID
	story.Note = userData.Note
	story.Status = userData.Status

	if err := app.svc.SaveAlmanackArticle(r.Context(), &story, userData.RefreshArc); err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}
	app.jsonResponse(http.StatusAccepted, w, &userData)
}

func (app *appEnv) listAvailable(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting listAvailable")
	type response struct {
		Contents []almanack.ArcStory `json:"contents"`
	}
	var (
		res response
		err error
	)
	if res.Contents, err = app.svc.GetAvailableFeed(r.Context()); err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}

	app.jsonResponse(http.StatusOK, w, res)
}

func (app *appEnv) getAvailable(w http.ResponseWriter, r *http.Request) {
	articleID := chi.URLParam(r, "id")
	app.Printf("starting getAvailable %s", articleID)

	article, err := app.svc.GetArcStory(r.Context(), articleID)
	if err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}

	if article.Status != almanack.StatusAvailable {
		// Let Spotlight PA users get article regardless of its status
		if err := app.auth.HasRole(r, "Spotlight PA"); err != nil {
			app.errorResponse(r.Context(), w, errutil.NotFound)
			return
		}
	}

	app.jsonResponse(http.StatusOK, w, article)
}

func (app *appEnv) postMessage(w http.ResponseWriter, r *http.Request) {
	app.Printf("starting postMessage")
	type request struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	var req request
	if err := httpjson.DecodeRequest(w, r, &req); err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}
	if err := app.email.SendEmail(req.Subject, req.Body); err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}
	app.jsonResponse(http.StatusAccepted, w, http.StatusText(http.StatusAccepted))
}

func (app *appEnv) getScheduledArticle(w http.ResponseWriter, r *http.Request) {
	articleID := chi.URLParam(r, "id")
	app.Printf("start getScheduledArticle %s", articleID)

	article, err := app.svc.GetScheduledArticle(r.Context(), articleID)
	if err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}

	app.jsonResponse(http.StatusOK, w, article)
}

func (app *appEnv) postScheduledArticle(w http.ResponseWriter, r *http.Request) {
	app.Println("start postScheduledArticle")

	var userData struct {
		almanack.SpotlightPAArticle
		RefreshArc bool `json:"almanack-refresh-arc"`
	}
	if err := httpjson.DecodeRequest(w, r, &userData); err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}

	if userData.RefreshArc {
		if err := app.svc.ResetSpotlightPAArticleArcData(r.Context(), &userData.SpotlightPAArticle); err != nil {
			app.errorResponse(r.Context(), w, err)
			return
		}
	}

	if strings.HasPrefix(userData.ImageURL, "http") {
		if imageurl, err := almanack.UploadFromURL(
			r.Context(), app.c, app.imageStore, userData.ImageURL,
		); err != nil {
			// Keep trucking, but don't publish
			userData.ImageURL = ""
			userData.ScheduleFor = nil
		} else {
			userData.ImageURL = imageurl
		}
	}

	if err := app.svc.SaveScheduledArticle(r.Context(), &userData.SpotlightPAArticle); err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}

	app.jsonResponse(http.StatusAccepted, w, &userData.SpotlightPAArticle)
}

func (app *appEnv) getSignedUpload(w http.ResponseWriter, r *http.Request) {
	app.Printf("start getSignedUpload")
	type response struct {
		SignedURL string `json:"signed-url"`
		FileName  string `json:"filename"`
	}
	var (
		res response
		err error
	)
	res.SignedURL, res.FileName, err = almanack.GetSignedUpload(app.imageStore)
	if err != nil {
		app.errorResponse(r.Context(), w, err)
		return
	}
	app.jsonResponse(http.StatusOK, w, &res)
}

func (app *appEnv) getSignupURL(w http.ResponseWriter, r *http.Request) {
	app.Println("start getSignupURL")
	app.jsonResponse(http.StatusOK, w, app.mailchimpSignupURL)
}
