package almanack

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/carlmjohnson/errutil"
	"github.com/jackc/pgx/v4"
	"github.com/spotlightpa/almanack/internal/arc"
	"github.com/spotlightpa/almanack/internal/db"
	"github.com/spotlightpa/almanack/internal/slack"
	"github.com/spotlightpa/almanack/internal/stringx"
	"github.com/spotlightpa/almanack/internal/timex"
	"github.com/spotlightpa/almanack/pkg/common"
)

func (svc Services) PublishPage(ctx context.Context, q *db.Queries, page *db.Page) (err, warning error) {
	defer errutil.Prefix(&err, "Service.PublishPage(%d)", page.ID)

	page.SetURLPath()
	data, err := page.ToTOML()
	if err != nil {
		return
	}

	err = errutil.ExecParallel(func() error {
		internalID, _ := page.Frontmatter["internal-id"].(string)
		title := stringx.First(internalID, page.FilePath)
		msg := fmt.Sprintf("Content: publishing %q", title)
		return svc.ContentStore.UpdateFile(ctx, msg, page.FilePath, []byte(data))
	}, func() error {
		_, warning = svc.Indexer.SaveObject(page.ToIndex(), ctx)
		return nil
	})
	if err != nil {
		return
	}

	p2, err := q.UpdatePage(ctx, db.UpdatePageParams{
		FilePath:         page.FilePath,
		URLPath:          page.URLPath.String,
		SetLastPublished: true,
		SetFrontmatter:   false,
		SetBody:          false,
		SetScheduleFor:   false,
		ScheduleFor:      db.NullTime,
	})
	if err != nil {
		return
	}
	*page = p2
	return
}

func (svc Services) RefreshPageFromContentStore(ctx context.Context, page *db.Page) (err error) {
	defer errutil.Prefix(&err, "Service.RefreshPageFromContentStore(%d)", page.ID)

	if db.IsNull(page.LastPublished) {
		return
	}
	content, err := svc.ContentStore.GetFile(ctx, page.FilePath)
	if err != nil {
		return err
	}
	if err = page.FromTOML(content); err != nil {
		return err
	}
	return nil
}

func (svc Services) PopScheduledPages(ctx context.Context) (err, warning error) {
	var warnings errutil.Slice
	err = svc.Tx.Begin(ctx, pgx.TxOptions{}, func(q *db.Queries) (txerr error) {
		defer errutil.Trace(&txerr)

		pages, txerr := q.PopScheduledPages(ctx)
		if txerr != nil {
			return
		}
		var errs errutil.Slice
		for _, page := range pages {
			txerr, warning = svc.PublishPage(ctx, q, &page)
			errs.Push(txerr)
			warnings.Push(warning)
		}
		return errs.Merge()
	})
	return err, warnings.Merge()
}

func (svc Services) RefreshPageContents(ctx context.Context, id int64) (err error) {
	defer errutil.Trace(&err)

	page, err := svc.Queries.GetPageByID(ctx, id)
	if err != nil {
		return err
	}
	defer errutil.Prefix(&err, fmt.Sprintf("problem refreshing contents of %s", page.FilePath))

	oldURLPath := page.URLPath.String
	contentBefore, err := page.ToTOML()
	if err != nil {
		return err
	}
	err = svc.RefreshPageFromContentStore(ctx, &page)
	if err != nil {
		return err
	}
	contentAfter, err := page.ToTOML()
	if err != nil {
		return err
	}

	if _, err = svc.Indexer.SaveObject(page.ToIndex(), ctx); err != nil {
		return err
	}

	page.SetURLPath()
	newURLPath := page.URLPath.String
	if contentBefore == contentAfter && oldURLPath == newURLPath {
		return nil
	}

	common.Logger.Printf("%s changed", page.FilePath)

	_, err = svc.Queries.UpdatePage(ctx, db.UpdatePageParams{
		FilePath:       page.FilePath,
		SetFrontmatter: true,
		Frontmatter:    page.Frontmatter,
		SetBody:        true,
		Body:           page.Body,
		URLPath:        page.URLPath.String,
		ScheduleFor:    db.NullTime,
	})

	return err
}

func (svc Services) RefreshPageFromArcStory(ctx context.Context, page *db.Page, story *db.Arc) (warnings []string, err error) {
	defer errutil.Trace(&err)

	var feedItem arc.FeedItem
	if err = story.RawData.AssignTo(&feedItem); err != nil {
		return nil, err
	}
	body, warnings, err := ArcFeedItemToBody(ctx, svc, &feedItem)
	if err != nil {
		return nil, err
	}

	page.Body = body
	return warnings, nil
}

func (svc Services) CreatePageFromArcSource(ctx context.Context, shared *db.SharedArticle, kind string) (warnings []string, err error) {
	defer errutil.Trace(&err)

	if shared.SourceType != "arc" {
		return nil, fmt.Errorf(
			"can't create new page for %d; wrong source type %q %q",
			shared.ID, shared.SourceType, shared.SourceID)
	}

	var feedItem arc.FeedItem
	if err = shared.RawData.AssignTo(&feedItem); err != nil {
		return nil, err
	}
	body, warnings, err := ArcFeedItemToBody(ctx, svc, &feedItem)
	if err != nil {
		return nil, err
	}

	fm, err := ArcFeedItemToFrontmatter(ctx, svc, &feedItem)
	if err != nil {
		return nil, err
	}

	filepath := buildFilePath(fm, kind)

	err = svc.Tx.Begin(ctx, pgx.TxOptions{}, func(q *db.Queries) (txerr error) {
		defer errutil.Trace(&txerr)

		if txerr = q.CreatePage(ctx, db.CreatePageParams{
			FilePath:   filepath,
			SourceType: shared.SourceType,
			SourceID:   shared.SourceID,
		}); txerr != nil {
			return txerr
		}

		page, txerr := q.UpdatePage(ctx, db.UpdatePageParams{
			FilePath:         filepath,
			SetFrontmatter:   true,
			Frontmatter:      fm,
			SetBody:          true,
			Body:             body,
			SetScheduleFor:   false,
			ScheduleFor:      db.NullTime,
			SetLastPublished: false,
		})
		if txerr != nil {
			return txerr
		}

		newSharedArt, txerr := q.UpdateSharedArticlePage(ctx, db.UpdateSharedArticlePageParams{
			PageID:          sql.NullInt64{Int64: page.ID, Valid: true},
			SharedArticleID: shared.ID,
		})
		if txerr != nil {
			return txerr
		}

		*shared = newSharedArt
		return nil
	})
	if err != nil {
		return nil, err
	}
	return warnings, nil
}

func buildFilePath(fm map[string]any, kind string) string {
	date := "1999-01-01"
	if t, ok := timex.Unwrap(fm["published"]); ok {
		date = timex.ToEST(t).Format("2006-01-02")
	}
	slug, _ := fm["internal-id"].(string)
	slug = stringx.First(slug, "SPLXXX")
	filepath := fmt.Sprintf("content/%s/%s-%s.md", kind, date, slug)
	return filepath
}

func (svc Services) Notify(ctx context.Context, page *db.Page, publishingNow bool) (err error) {
	defer errutil.Trace(&err)

	const (
		green  = "#78bc20"
		yellow = "#ffcb05"
	)
	text := "New page publishing now…"
	color := green

	if !publishingNow {
		t := timex.ToEST(page.ScheduleFor.Time)
		text = t.Format("New article scheduled for Mon, Jan 2 at 3:04pm MST…")
		color = yellow
	}

	hed, _ := page.Frontmatter["title"].(string)
	summary := page.Frontmatter["description"].(string)
	url := page.FullURL()
	return svc.SlackSocial.Post(ctx, slack.Message{
		Text: text,
		Attachments: []slack.Attachment{
			{
				Color: color,
				Fallback: fmt.Sprintf("%s\n%s\n%s",
					hed, summary, url),
				Title:     hed,
				TitleLink: url,
				Text: fmt.Sprintf(
					"%s\n%s",
					summary, url),
			},
		},
	})
}
