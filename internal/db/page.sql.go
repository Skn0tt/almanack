// Code generated by sqlc. DO NOT EDIT.
// source: page.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const ensurePage = `-- name: EnsurePage :exec
INSERT INTO page ("file_path")
  VALUES ($1)
ON CONFLICT (file_path)
  DO NOTHING
`

func (q *Queries) EnsurePage(ctx context.Context, filePath string) error {
	_, err := q.db.ExecContext(ctx, ensurePage, filePath)
	return err
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
	row := q.db.QueryRowContext(ctx, getPageByID, id)
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

const getPageByPath = `-- name: GetPageByPath :one
SELECT
  id, file_path, frontmatter, body, schedule_for, last_published, created_at, updated_at, url_path
FROM
  "page"
WHERE
  file_path = $1
`

func (q *Queries) GetPageByPath(ctx context.Context, filePath string) (Page, error) {
	row := q.db.QueryRowContext(ctx, getPageByPath, filePath)
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

const listPages = `-- name: ListPages :many
SELECT
  "id",
  "file_path",
  (frontmatter ->> 'internal-id')::text AS "internal_id",
  (frontmatter ->> 'title')::text AS "title",
  (frontmatter ->> 'description')::text AS "description",
  (frontmatter ->> 'blurb')::text AS "blurb",
  (frontmatter ->> 'image')::text AS "image",
  coalesce("url_path", ''),
  "last_published",
  "created_at",
  "updated_at",
  "schedule_for",
  to_timestamp(frontmatter ->> 'published',
    -- ISO date
    'YYYY-MM-DD"T"HH24:MI:SS"Z"')::timestamptz AS "published_at"
FROM
  page
WHERE
  "file_path" ILIKE $1
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
	ID            int64        `json:"id"`
	FilePath      string       `json:"file_path"`
	InternalID    string       `json:"internal_id"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Blurb         string       `json:"blurb"`
	Image         string       `json:"image"`
	URLPath       string       `json:"url_path"`
	LastPublished sql.NullTime `json:"last_published"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	ScheduleFor   sql.NullTime `json:"schedule_for"`
	PublishedAt   time.Time    `json:"published_at"`
}

func (q *Queries) ListPages(ctx context.Context, arg ListPagesParams) ([]ListPagesRow, error) {
	rows, err := q.db.QueryContext(ctx, listPages, arg.FilePath, arg.Limit, arg.Offset)
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
	if err := rows.Close(); err != nil {
		return nil, err
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
	rows, err := q.db.QueryContext(ctx, popScheduledPages)
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
	if err := rows.Close(); err != nil {
		return nil, err
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
	SetFrontmatter   bool         `json:"set_frontmatter"`
	Frontmatter      Map          `json:"frontmatter"`
	SetBody          bool         `json:"set_body"`
	Body             string       `json:"body"`
	SetScheduleFor   bool         `json:"set_schedule_for"`
	ScheduleFor      sql.NullTime `json:"schedule_for"`
	URLPath          string       `json:"url_path"`
	SetLastPublished bool         `json:"set_last_published"`
	FilePath         string       `json:"file_path"`
}

func (q *Queries) UpdatePage(ctx context.Context, arg UpdatePageParams) (Page, error) {
	row := q.db.QueryRowContext(ctx, updatePage,
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
