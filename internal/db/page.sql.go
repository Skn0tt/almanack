// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: page.sql

package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPage = `-- name: CreatePage :exec
INSERT INTO page ("file_path", "source_type", "source_id")
  VALUES ($1, $2, $3)
ON CONFLICT (file_path)
  DO NOTHING
`

type CreatePageParams struct {
	FilePath   string `json:"file_path"`
	SourceType string `json:"source_type"`
	SourceID   string `json:"source_id"`
}

func (q *Queries) CreatePage(ctx context.Context, arg CreatePageParams) error {
	_, err := q.db.Exec(ctx, createPage, arg.FilePath, arg.SourceType, arg.SourceID)
	return err
}

const getArchiveURLForPageID = `-- name: GetArchiveURLForPageID :one
SELECT
  coalesce(archive_url, '')
FROM
  page
  LEFT JOIN newsletter ON page.source_id = newsletter.id::text
    AND page.source_type = 'mailchimp'
WHERE
  page.id = $1
`

func (q *Queries) GetArchiveURLForPageID(ctx context.Context, id int64) (string, error) {
	row := q.db.QueryRow(ctx, getArchiveURLForPageID, id)
	var archive_url string
	err := row.Scan(&archive_url)
	return archive_url, err
}

const getPageByFilePath = `-- name: GetPageByFilePath :one
SELECT
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path, source_type, source_id, published_at
FROM
  "page"
WHERE
  file_path = $1
`

func (q *Queries) GetPageByFilePath(ctx context.Context, filePath string) (Page, error) {
	row := q.db.QueryRow(ctx, getPageByFilePath, filePath)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.FilePath,
		&i.Frontmatter,
		&i.Body,
		&i.ScheduleFor,
		&i.LastPublished,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.URLPath,
		&i.SourceType,
		&i.SourceID,
		&i.PublishedAt,
	)
	return i, err
}

const getPageByID = `-- name: GetPageByID :one
SELECT
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path, source_type, source_id, published_at
FROM
  "page"
WHERE
  id = $1
`

func (q *Queries) GetPageByID(ctx context.Context, id int64) (Page, error) {
	row := q.db.QueryRow(ctx, getPageByID, id)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.FilePath,
		&i.Frontmatter,
		&i.Body,
		&i.ScheduleFor,
		&i.LastPublished,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.URLPath,
		&i.SourceType,
		&i.SourceID,
		&i.PublishedAt,
	)
	return i, err
}

const getPageByURLPath = `-- name: GetPageByURLPath :one
SELECT
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path, source_type, source_id, published_at
FROM
  page
WHERE
  url_path ILIKE $1::text
`

func (q *Queries) GetPageByURLPath(ctx context.Context, urlPath string) (Page, error) {
	row := q.db.QueryRow(ctx, getPageByURLPath, urlPath)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.FilePath,
		&i.Frontmatter,
		&i.Body,
		&i.ScheduleFor,
		&i.LastPublished,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.URLPath,
		&i.SourceType,
		&i.SourceID,
		&i.PublishedAt,
	)
	return i, err
}

const listAllPages = `-- name: ListAllPages :many
SELECT
  id,
  file_path,
  coalesce(frontmatter ->> 'internal-id', '')::text AS internal_id,
  coalesce(frontmatter ->> 'title', '')::text AS hed,
  ARRAY (
    SELECT
      jsonb_array_elements_text(
        CASE WHEN frontmatter ->> 'authors' IS NOT NULL THEN
          frontmatter -> 'authors'
        ELSE
          '[]'::jsonb
        END))::text[] AS authors,
  published_at::timestamptz AS pub_date
FROM
  page
WHERE
  published_at IS NOT NULL
ORDER BY
  published_at DESC
`

type ListAllPagesRow struct {
	ID         int64                `json:"id"`
	FilePath   string               `json:"file_path"`
	InternalID string               `json:"internal_id"`
	Hed        string               `json:"hed"`
	Authors    pgtype.Array[string] `json:"authors"`
	PubDate    time.Time            `json:"pub_date"`
}

func (q *Queries) ListAllPages(ctx context.Context) ([]ListAllPagesRow, error) {
	rows, err := q.db.Query(ctx, listAllPages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAllPagesRow
	for rows.Next() {
		var i ListAllPagesRow
		if err := rows.Scan(
			&i.ID,
			&i.FilePath,
			&i.InternalID,
			&i.Hed,
			&i.Authors,
			&i.PubDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAllSeries = `-- name: ListAllSeries :many
WITH page_series AS (
  SELECT
    json_array_elements_text( --
      coalesce((frontmatter ->> 'series'), '[]')::json) AS series,
    published_at
  FROM
    page
),
series_dates AS (
  SELECT
    series, published_at
  FROM
    page_series
  ORDER BY
    published_at DESC,
    series DESC
),
distinct_series_dates AS (
  SELECT DISTINCT ON (series)
    series, published_at
  FROM
    series_dates
  ORDER BY
    series DESC,
    published_at DESC
)
SELECT
  series::text
FROM
  distinct_series_dates
ORDER BY
  published_at DESC
`

func (q *Queries) ListAllSeries(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, listAllSeries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var series string
		if err := rows.Scan(&series); err != nil {
			return nil, err
		}
		items = append(items, series)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAllTopics = `-- name: ListAllTopics :many
WITH topics AS (
  SELECT
    json_array_elements_text( --
      coalesce((frontmatter ->> 'topics'), '[]')::json) AS topic
  FROM
    page
)
SELECT DISTINCT ON (upper(topic)
)
  topic::text
FROM
  topics
ORDER BY
  upper(topic) ASC
`

func (q *Queries) ListAllTopics(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, listAllTopics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var topic string
		if err := rows.Scan(&topic); err != nil {
			return nil, err
		}
		items = append(items, topic)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPageIDs = `-- name: ListPageIDs :many
SELECT
  "id"
FROM
  page
WHERE
  "file_path" LIKE $1
ORDER BY
  id ASC
LIMIT $2 OFFSET $3
`

type ListPageIDsParams struct {
	FilePath string `json:"file_path"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

func (q *Queries) ListPageIDs(ctx context.Context, arg ListPageIDsParams) ([]int64, error) {
	rows, err := q.db.Query(ctx, listPageIDs, arg.FilePath, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPages = `-- name: ListPages :many
WITH paths AS (
  SELECT
    id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path, source_type, source_id, published_at
  FROM
    page
  WHERE
    "file_path" LIKE $1
),
ordered AS (
  SELECT
    id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path, source_type, source_id, published_at
  FROM
    paths
  ORDER BY
    published_at DESC
)
SELECT
  id,
  file_path::text,
  coalesce(frontmatter ->> 'internal-id', '')::text AS "internal_id",
  coalesce(frontmatter ->> 'title', '')::text AS "title",
  coalesce(frontmatter ->> 'description', '')::text AS "description",
  coalesce(frontmatter ->> 'blurb', '')::text AS "blurb",
  coalesce(frontmatter ->> 'image', '')::text AS "image",
  coalesce(url_path, '')::text AS "url_path",
  last_published,
  created_at,
  updated_at,
  schedule_for,
  published_at
FROM
  ordered
LIMIT $2 OFFSET $3
`

type ListPagesParams struct {
	FilePath string `json:"file_path"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type ListPagesRow struct {
	ID            int64              `json:"id"`
	FilePath      string             `json:"file_path"`
	InternalID    string             `json:"internal_id"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Blurb         string             `json:"blurb"`
	Image         string             `json:"image"`
	URLPath       string             `json:"url_path"`
	LastPublished pgtype.Timestamptz `json:"last_published"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	ScheduleFor   pgtype.Timestamptz `json:"schedule_for"`
	PublishedAt   pgtype.Timestamptz `json:"published_at"`
}

// Treating published_at as text because it sorts faster and we don't do
// date stuff on the backend, just frontend.
func (q *Queries) ListPages(ctx context.Context, arg ListPagesParams) ([]ListPagesRow, error) {
	rows, err := q.db.Query(ctx, listPages, arg.FilePath, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPagesRow
	for rows.Next() {
		var i ListPagesRow
		if err := rows.Scan(
			&i.ID,
			&i.FilePath,
			&i.InternalID,
			&i.Title,
			&i.Description,
			&i.Blurb,
			&i.Image,
			&i.URLPath,
			&i.LastPublished,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ScheduleFor,
			&i.PublishedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPagesByFTS = `-- name: ListPagesByFTS :many
WITH query AS (
  SELECT
    id,
    ts_rank(fts_doc_en, tsq) AS rank
  FROM
    page,
    websearch_to_tsquery('english', $2::text) tsq
  WHERE
    fts_doc_en @@ tsq
  ORDER BY
    rank DESC
  LIMIT $1
)
SELECT
  page.id, page.file_path, page.frontmatter, page.body, page.schedule_for, page.last_published, page.created_at, page.updated_at, page.url_path, page.source_type, page.source_id, page.published_at
FROM
  page
  JOIN query USING (id)
ORDER BY
  published_at DESC
`

type ListPagesByFTSParams struct {
	Limit int32  `json:"limit"`
	Query string `json:"query"`
}

func (q *Queries) ListPagesByFTS(ctx context.Context, arg ListPagesByFTSParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesByFTS, arg.Limit, arg.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.FilePath,
			&i.Frontmatter,
			&i.Body,
			&i.ScheduleFor,
			&i.LastPublished,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.URLPath,
			&i.SourceType,
			&i.SourceID,
			&i.PublishedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPagesByInternalID = `-- name: ListPagesByInternalID :many
WITH query AS (
  SELECT
    id,
    ts_rank(fts_doc_en, id_tsq) AS rank
  FROM
    page,
    tsquery ($2::text) id_tsq
  WHERE
    internal_id_fts @@ id_tsq
  ORDER BY
    rank DESC
  LIMIT $1
)
SELECT
  page.id, page.file_path, page.frontmatter, page.body, page.schedule_for, page.last_published, page.created_at, page.updated_at, page.url_path, page.source_type, page.source_id, page.published_at
FROM
  page
  JOIN query USING (id)
ORDER BY
  published_at DESC
`

type ListPagesByInternalIDParams struct {
	Limit int32  `json:"limit"`
	Query string `json:"query"`
}

func (q *Queries) ListPagesByInternalID(ctx context.Context, arg ListPagesByInternalIDParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesByInternalID, arg.Limit, arg.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.FilePath,
			&i.Frontmatter,
			&i.Body,
			&i.ScheduleFor,
			&i.LastPublished,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.URLPath,
			&i.SourceType,
			&i.SourceID,
			&i.PublishedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPagesByPublished = `-- name: ListPagesByPublished :many
SELECT
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path, source_type, source_id, published_at
FROM
  page
ORDER BY
  published_at DESC
LIMIT $1 OFFSET $2
`

type ListPagesByPublishedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListPagesByPublished(ctx context.Context, arg ListPagesByPublishedParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesByPublished, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.FilePath,
			&i.Frontmatter,
			&i.Body,
			&i.ScheduleFor,
			&i.LastPublished,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.URLPath,
			&i.SourceType,
			&i.SourceID,
			&i.PublishedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPagesByURLPaths = `-- name: ListPagesByURLPaths :many
WITH query_paths AS (
  SELECT
    $1::text[] AS "paths"
),
page_paths AS (
  SELECT
    id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path, source_type, source_id, published_at
  FROM
    page
  WHERE
    url_path IN (
      SELECT
        unnest("paths")::citext
      FROM
        query_paths))
SELECT
  "file_path"::text,
  coalesce(frontmatter ->> 'internal-id', '')::text AS "internal_id",
  coalesce(frontmatter ->> 'title', '')::text AS "title",
  coalesce(frontmatter ->> 'link-title', '')::text AS "link_title",
  coalesce(frontmatter ->> 'description', '')::text AS "description",
  coalesce(frontmatter ->> 'blurb', '')::text AS "blurb",
  coalesce(frontmatter ->> 'image', '')::text AS "image",
  coalesce(url_path, '')::text AS "url_path",
  published_at::timestamptz
FROM
  page_paths
  CROSS JOIN query_paths
ORDER BY
  array_position(query_paths.paths, url_path::text)
`

type ListPagesByURLPathsRow struct {
	FilePath    string    `json:"file_path"`
	InternalID  string    `json:"internal_id"`
	Title       string    `json:"title"`
	LinkTitle   string    `json:"link_title"`
	Description string    `json:"description"`
	Blurb       string    `json:"blurb"`
	Image       string    `json:"image"`
	URLPath     string    `json:"url_path"`
	PublishedAt time.Time `json:"published_at"`
}

func (q *Queries) ListPagesByURLPaths(ctx context.Context, paths pgtype.Array[string]) ([]ListPagesByURLPathsRow, error) {
	rows, err := q.db.Query(ctx, listPagesByURLPaths, paths)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPagesByURLPathsRow
	for rows.Next() {
		var i ListPagesByURLPathsRow
		if err := rows.Scan(
			&i.FilePath,
			&i.InternalID,
			&i.Title,
			&i.LinkTitle,
			&i.Description,
			&i.Blurb,
			&i.Image,
			&i.URLPath,
			&i.PublishedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const popScheduledPages = `-- name: PopScheduledPages :many
UPDATE
  page
SET
  last_published = CURRENT_TIMESTAMP
WHERE
  last_published IS NULL
  AND schedule_for < (CURRENT_TIMESTAMP + '5 minutes'::interval)
RETURNING
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path, source_type, source_id, published_at
`

func (q *Queries) PopScheduledPages(ctx context.Context) ([]Page, error) {
	rows, err := q.db.Query(ctx, popScheduledPages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.FilePath,
			&i.Frontmatter,
			&i.Body,
			&i.ScheduleFor,
			&i.LastPublished,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.URLPath,
			&i.SourceType,
			&i.SourceID,
			&i.PublishedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePage = `-- name: UpdatePage :one
UPDATE
  page
SET
  frontmatter = CASE WHEN $1::boolean THEN
    $2
  ELSE
    frontmatter
  END,
  body = CASE WHEN $3::boolean THEN
    $4
  ELSE
    body
  END,
  schedule_for = CASE WHEN $5::boolean THEN
    $6
  ELSE
    schedule_for
  END,
  url_path = CASE WHEN $7::text != '' THEN
    $7
  ELSE
    url_path
  END,
  last_published = CASE WHEN $8::boolean THEN
    CURRENT_TIMESTAMP
  ELSE
    last_published
  END
WHERE
  file_path = $9
RETURNING
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path, source_type, source_id, published_at
`

type UpdatePageParams struct {
	SetFrontmatter   bool               `json:"set_frontmatter"`
	Frontmatter      Map                `json:"frontmatter"`
	SetBody          bool               `json:"set_body"`
	Body             string             `json:"body"`
	SetScheduleFor   bool               `json:"set_schedule_for"`
	ScheduleFor      pgtype.Timestamptz `json:"schedule_for"`
	URLPath          string             `json:"url_path"`
	SetLastPublished bool               `json:"set_last_published"`
	FilePath         string             `json:"file_path"`
}

func (q *Queries) UpdatePage(ctx context.Context, arg UpdatePageParams) (Page, error) {
	row := q.db.QueryRow(ctx, updatePage,
		arg.SetFrontmatter,
		arg.Frontmatter,
		arg.SetBody,
		arg.Body,
		arg.SetScheduleFor,
		arg.ScheduleFor,
		arg.URLPath,
		arg.SetLastPublished,
		arg.FilePath,
	)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.FilePath,
		&i.Frontmatter,
		&i.Body,
		&i.ScheduleFor,
		&i.LastPublished,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.URLPath,
		&i.SourceType,
		&i.SourceID,
		&i.PublishedAt,
	)
	return i, err
}

const updatePageRawContent = `-- name: UpdatePageRawContent :one
UPDATE
  page
SET
  frontmatter = frontmatter || jsonb_build_object('raw-content', $1::text)
WHERE
  id = $2
RETURNING
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path, source_type, source_id, published_at
`

type UpdatePageRawContentParams struct {
	RawContent string `json:"raw_content"`
	ID         int64  `json:"id"`
}

func (q *Queries) UpdatePageRawContent(ctx context.Context, arg UpdatePageRawContentParams) (Page, error) {
	row := q.db.QueryRow(ctx, updatePageRawContent, arg.RawContent, arg.ID)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.FilePath,
		&i.Frontmatter,
		&i.Body,
		&i.ScheduleFor,
		&i.LastPublished,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.URLPath,
		&i.SourceType,
		&i.SourceID,
		&i.PublishedAt,
	)
	return i, err
}
