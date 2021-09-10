// Code generated by sqlc. DO NOT EDIT.
// source: page.sql

package db

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
)

const ensurePage = `-- name: EnsurePage :exec
INSERT INTO page ("file_path")
  VALUES ($1)
ON CONFLICT (file_path)
  DO NOTHING
`

func (q *Queries) EnsurePage(ctx context.Context, filePath string) error {
	_, err := q.db.Exec(ctx, ensurePage, filePath)
	return err
}

const getPageByFilePath = `-- name: GetPageByFilePath :one
SELECT
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path
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
	)
	return i, err
}

const getPageByID = `-- name: GetPageByID :one
SELECT
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path
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
	)
	return i, err
}

const getPageByURLPath = `-- name: GetPageByURLPath :one
SELECT
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path
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
	)
	return i, err
}

const listAllPages = `-- name: ListAllPages :many
WITH ROWS AS (
  SELECT
    id,
    file_path,
    coalesce(frontmatter ->> 'internal-id', '') AS internal_id,
    coalesce(frontmatter ->> 'title', '') AS hed,
    ARRAY (
      SELECT
        jsonb_array_elements_text(
          CASE WHEN frontmatter ->> 'authors' IS NOT NULL THEN
            frontmatter -> 'authors'
          ELSE
            '[]'::jsonb
          END)) AS authors,
      to_timestamp(frontmatter ->> 'published'::text,
        -- ISO date
        'YYYY-MM-DD"T"HH24:MI:SS"Z"')::timestamptz AS pub_date
    FROM
      page
    ORDER BY
      pub_date DESC
)
SELECT
  id,
  file_path,
  internal_id::text AS internal_id,
  hed::text AS hed,
  authors::text[] AS authors,
  pub_date::timestamptz AS pub_date
FROM
  ROWS
WHERE
  pub_date IS NOT NULL
`

type ListAllPagesRow struct {
	ID         int64     `json:"id"`
	FilePath   string    `json:"file_path"`
	InternalID string    `json:"internal_id"`
	Hed        string    `json:"hed"`
	Authors    []string  `json:"authors"`
	PubDate    time.Time `json:"pub_date"`
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
WITH series_dates AS (
  SELECT
    jsonb_array_elements_text(frontmatter -> 'series') AS series,
    frontmatter ->> 'published' AS pub_date
  FROM
    page
  WHERE
    frontmatter ->> 'series' IS NOT NULL
  ORDER BY
    pub_date DESC,
    series DESC
),
distinct_series_dates AS (
  SELECT DISTINCT ON (series)
    series, pub_date
  FROM
    series_dates
  ORDER BY
    series DESC,
    pub_date DESC
)
SELECT
  series::text
FROM
  distinct_series_dates
ORDER BY
  pub_date DESC
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
    jsonb_array_elements_text(frontmatter -> 'topics') AS topic
  FROM
    page
  WHERE
    frontmatter ->> 'topics' IS NOT NULL
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
SELECT
  "id",
  "file_path",
  (
    CASE WHEN frontmatter ->> 'internal-id' IS NOT NULL THEN
      frontmatter ->> 'internal-id'
    ELSE
      ''
    END)::text AS "internal_id",
  (
    CASE WHEN frontmatter ->> 'title' IS NOT NULL THEN
      frontmatter ->> 'title'
    ELSE
      ''
    END)::text AS "title",
  (
    CASE WHEN frontmatter ->> 'description' IS NOT NULL THEN
      frontmatter ->> 'description'
    ELSE
      ''
    END)::text AS "description",
  (
    CASE WHEN frontmatter ->> 'blurb' IS NOT NULL THEN
      frontmatter ->> 'blurb'
    ELSE
      ''
    END)::text AS "blurb",
  (
    CASE WHEN frontmatter ->> 'image' IS NOT NULL THEN
      frontmatter ->> 'image'
    ELSE
      ''
    END)::text AS "image",
  coalesce("url_path", ''),
  "last_published",
  "created_at",
  "updated_at",
  "schedule_for",
  (
    CASE WHEN frontmatter ->> 'published' IS NOT NULL THEN
      frontmatter ->> 'published'
    ELSE
      ''
    END)::text AS "published_at"
FROM
  page
WHERE
  "file_path" LIKE $1
ORDER BY
  frontmatter ->> 'published' DESC
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
	PublishedAt   string             `json:"published_at"`
}

// Cannot use coalesce, see https://github.com/kyleconroy/sqlc/issues/780.
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

const popScheduledPages = `-- name: PopScheduledPages :many
UPDATE
  page
SET
  last_published = CURRENT_TIMESTAMP
WHERE
  last_published IS NULL
  AND schedule_for < (CURRENT_TIMESTAMP + '5 minutes'::interval)
RETURNING
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path
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
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path
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
	)
	return i, err
}
