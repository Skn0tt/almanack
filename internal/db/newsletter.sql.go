// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: newsletter.sql

package db

import (
	"context"

	"github.com/jackc/pgtype"
)

const listNewsletterTypes = `-- name: ListNewsletterTypes :many
SELECT
  shortname, name, description
FROM
  newsletter_type
`

func (q *Queries) ListNewsletterTypes(ctx context.Context) ([]NewsletterType, error) {
	rows, err := q.db.Query(ctx, listNewsletterTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []NewsletterType
	for rows.Next() {
		var i NewsletterType
		if err := rows.Scan(&i.Shortname, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listNewslettersWithoutPage = `-- name: ListNewslettersWithoutPage :many
SELECT
  subject, archive_url, published_at, type, created_at, updated_at, id, description, blurb, spotlightpa_path
FROM
  newsletter
WHERE
  "spotlightpa_path" IS NULL
ORDER BY
  published_at DESC
LIMIT $1 OFFSET $2
`

type ListNewslettersWithoutPageParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListNewslettersWithoutPage(ctx context.Context, arg ListNewslettersWithoutPageParams) ([]Newsletter, error) {
	rows, err := q.db.Query(ctx, listNewslettersWithoutPage, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Newsletter
	for rows.Next() {
		var i Newsletter
		if err := rows.Scan(
			&i.Subject,
			&i.ArchiveURL,
			&i.PublishedAt,
			&i.Type,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID,
			&i.Description,
			&i.Blurb,
			&i.SpotlightPAPath,
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

const setNewsletterPage = `-- name: SetNewsletterPage :one
UPDATE
  newsletter
SET
  "spotlightpa_path" = $2
WHERE
  id = $1
RETURNING
  subject, archive_url, published_at, type, created_at, updated_at, id, description, blurb, spotlightpa_path
`

type SetNewsletterPageParams struct {
	ID              int64       `json:"id"`
	SpotlightPAPath pgtype.Text `json:"spotlightpa_path"`
}

func (q *Queries) SetNewsletterPage(ctx context.Context, arg SetNewsletterPageParams) (Newsletter, error) {
	row := q.db.QueryRow(ctx, setNewsletterPage, arg.ID, arg.SpotlightPAPath)
	var i Newsletter
	err := row.Scan(
		&i.Subject,
		&i.ArchiveURL,
		&i.PublishedAt,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ID,
		&i.Description,
		&i.Blurb,
		&i.SpotlightPAPath,
	)
	return i, err
}

const updateNewsletterArchives = `-- name: UpdateNewsletterArchives :execrows
WITH raw_json AS (
  SELECT
    jsonb_array_elements($2::jsonb) AS data
),
campaign AS (
  SELECT
    data ->> 'subject' AS subject,
    data ->> 'blurb' AS blurb,
    data ->> 'description' AS description,
    data ->> 'archive_url' AS archive_url,
    iso_to_timestamptz (data ->> 'published_at')::timestamptz AS published_at
  FROM
    raw_json)
  INSERT INTO newsletter ("subject", "blurb", "description", "archive_url",
    "published_at", "type")
  SELECT
    subject,
    blurb,
    description,
    archive_url,
    published_at,
    $1
  FROM
    campaign
  ON CONFLICT
    DO NOTHING
`

type UpdateNewsletterArchivesParams struct {
	Type string       `json:"type"`
	Data pgtype.JSONB `json:"data"`
}

func (q *Queries) UpdateNewsletterArchives(ctx context.Context, arg UpdateNewsletterArchivesParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateNewsletterArchives, arg.Type, arg.Data)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
