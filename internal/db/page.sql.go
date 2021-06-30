// Code generated by sqlc. DO NOT EDIT.
// source: page.sql

package db

import (
	"context"
	"database/sql"
)

const ensurePage = `-- name: EnsurePage :exec
INSERT INTO page ("path")
  VALUES ($1)
ON CONFLICT (path)
  DO NOTHING
`

func (q *Queries) EnsurePage(ctx context.Context, path string) error {
	_, err := q.db.ExecContext(ctx, ensurePage, path)
	return err
}

const getPage = `-- name: GetPage :one
SELECT
  id, path, frontmatter, body, schedule_for, last_published, created_at, updated_at
FROM
  "page"
WHERE
  path = $1
`

func (q *Queries) GetPage(ctx context.Context, path string) (Page, error) {
	row := q.db.QueryRowContext(ctx, getPage, path)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.Path,
		&i.Frontmatter,
		&i.Body,
		&i.ScheduleFor,
		&i.LastPublished,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
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
  id, path, frontmatter, body, schedule_for, last_published, created_at, updated_at
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
			&i.Path,
			&i.Frontmatter,
			&i.Body,
			&i.ScheduleFor,
			&i.LastPublished,
			&i.CreatedAt,
			&i.UpdatedAt,
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
  last_published = CASE WHEN $7::boolean THEN
    CURRENT_TIMESTAMP
  ELSE
    last_published
  END
WHERE
  path = $8
RETURNING
  id, path, frontmatter, body, schedule_for, last_published, created_at, updated_at
`

type UpdatePageParams struct {
	SetFrontmatter   bool         `json:"set_frontmatter"`
	Frontmatter      Map          `json:"frontmatter"`
	SetBody          bool         `json:"set_body"`
	Body             string       `json:"body"`
	SetScheduleFor   bool         `json:"set_schedule_for"`
	ScheduleFor      sql.NullTime `json:"schedule_for"`
	SetLastPublished bool         `json:"set_last_published"`
	Path             string       `json:"path"`
}

func (q *Queries) UpdatePage(ctx context.Context, arg UpdatePageParams) (Page, error) {
	row := q.db.QueryRowContext(ctx, updatePage,
		arg.SetFrontmatter,
		arg.Frontmatter,
		arg.SetBody,
		arg.Body,
		arg.SetScheduleFor,
		arg.ScheduleFor,
		arg.SetLastPublished,
		arg.Path,
	)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.Path,
		&i.Frontmatter,
		&i.Body,
		&i.ScheduleFor,
		&i.LastPublished,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
