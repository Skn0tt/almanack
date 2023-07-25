// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: shared-article.sql

package db

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

const getSharedArticleByID = `-- name: GetSharedArticleByID :one
SELECT
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at, publication_date, internal_id, byline, budget, hed, description, lede_image, lede_image_credit, lede_image_description, lede_image_caption, blurb
FROM
  shared_article
WHERE
  id = $1
`

func (q *Queries) GetSharedArticleByID(ctx context.Context, id int64) (SharedArticle, error) {
	row := q.db.QueryRow(ctx, getSharedArticleByID, id)
	var i SharedArticle
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.EmbargoUntil,
		&i.Note,
		&i.SourceType,
		&i.SourceID,
		&i.RawData,
		&i.PageID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PublicationDate,
		&i.InternalID,
		&i.Byline,
		&i.Budget,
		&i.Hed,
		&i.Description,
		&i.LedeImage,
		&i.LedeImageCredit,
		&i.LedeImageDescription,
		&i.LedeImageCaption,
		&i.Blurb,
	)
	return i, err
}

const getSharedArticleBySource = `-- name: GetSharedArticleBySource :one
SELECT
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at, publication_date, internal_id, byline, budget, hed, description, lede_image, lede_image_credit, lede_image_description, lede_image_caption, blurb
FROM
  shared_article
WHERE
  source_type = $1
  AND source_id = $2
`

type GetSharedArticleBySourceParams struct {
	SourceType string `json:"source_type"`
	SourceID   string `json:"source_id"`
}

func (q *Queries) GetSharedArticleBySource(ctx context.Context, arg GetSharedArticleBySourceParams) (SharedArticle, error) {
	row := q.db.QueryRow(ctx, getSharedArticleBySource, arg.SourceType, arg.SourceID)
	var i SharedArticle
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.EmbargoUntil,
		&i.Note,
		&i.SourceType,
		&i.SourceID,
		&i.RawData,
		&i.PageID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PublicationDate,
		&i.InternalID,
		&i.Byline,
		&i.Budget,
		&i.Hed,
		&i.Description,
		&i.LedeImage,
		&i.LedeImageCredit,
		&i.LedeImageDescription,
		&i.LedeImageCaption,
		&i.Blurb,
	)
	return i, err
}

const listSharedArticles = `-- name: ListSharedArticles :many
SELECT
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at, publication_date, internal_id, byline, budget, hed, description, lede_image, lede_image_credit, lede_image_description, lede_image_caption, blurb
FROM
  shared_article
ORDER BY
  publication_date DESC
LIMIT $1 OFFSET $2
`

type ListSharedArticlesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListSharedArticles(ctx context.Context, arg ListSharedArticlesParams) ([]SharedArticle, error) {
	rows, err := q.db.Query(ctx, listSharedArticles, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SharedArticle
	for rows.Next() {
		var i SharedArticle
		if err := rows.Scan(
			&i.ID,
			&i.Status,
			&i.EmbargoUntil,
			&i.Note,
			&i.SourceType,
			&i.SourceID,
			&i.RawData,
			&i.PageID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PublicationDate,
			&i.InternalID,
			&i.Byline,
			&i.Budget,
			&i.Hed,
			&i.Description,
			&i.LedeImage,
			&i.LedeImageCredit,
			&i.LedeImageDescription,
			&i.LedeImageCaption,
			&i.Blurb,
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

const listSharedArticlesWhereActive = `-- name: ListSharedArticlesWhereActive :many
SELECT
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at, publication_date, internal_id, byline, budget, hed, description, lede_image, lede_image_credit, lede_image_description, lede_image_caption, blurb
FROM
  shared_article
WHERE
  status <> 'U'
ORDER BY
  publication_date DESC
LIMIT $1 OFFSET $2
`

type ListSharedArticlesWhereActiveParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListSharedArticlesWhereActive(ctx context.Context, arg ListSharedArticlesWhereActiveParams) ([]SharedArticle, error) {
	rows, err := q.db.Query(ctx, listSharedArticlesWhereActive, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SharedArticle
	for rows.Next() {
		var i SharedArticle
		if err := rows.Scan(
			&i.ID,
			&i.Status,
			&i.EmbargoUntil,
			&i.Note,
			&i.SourceType,
			&i.SourceID,
			&i.RawData,
			&i.PageID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PublicationDate,
			&i.InternalID,
			&i.Byline,
			&i.Budget,
			&i.Hed,
			&i.Description,
			&i.LedeImage,
			&i.LedeImageCredit,
			&i.LedeImageDescription,
			&i.LedeImageCaption,
			&i.Blurb,
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

const updateSharedArticle = `-- name: UpdateSharedArticle :one
UPDATE
  shared_article
SET
  "embargo_until" = $1,
  "status" = $2,
  "note" = $3,
  "publication_date" = $4,
  "internal_id" = $5,
  "byline" = $6,
  "budget" = $7,
  "hed" = $8,
  "description" = $9,
  "blurb" = $10,
  "lede_image" = $11,
  "lede_image_credit" = $12,
  "lede_image_description" = $13,
  "lede_image_caption" = $14
WHERE
  id = $15
RETURNING
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at, publication_date, internal_id, byline, budget, hed, description, lede_image, lede_image_credit, lede_image_description, lede_image_caption, blurb
`

type UpdateSharedArticleParams struct {
	EmbargoUntil         pgtype.Timestamptz `json:"embargo_until"`
	Status               string             `json:"status"`
	Note                 string             `json:"note"`
	PublicationDate      pgtype.Timestamptz `json:"publication_date"`
	InternalID           string             `json:"internal_id"`
	Byline               string             `json:"byline"`
	Budget               string             `json:"budget"`
	Hed                  string             `json:"hed"`
	Description          string             `json:"description"`
	Blurb                string             `json:"blurb"`
	LedeImage            string             `json:"lede_image"`
	LedeImageCredit      string             `json:"lede_image_credit"`
	LedeImageDescription string             `json:"lede_image_description"`
	LedeImageCaption     string             `json:"lede_image_caption"`
	ID                   int64              `json:"id"`
}

func (q *Queries) UpdateSharedArticle(ctx context.Context, arg UpdateSharedArticleParams) (SharedArticle, error) {
	row := q.db.QueryRow(ctx, updateSharedArticle,
		arg.EmbargoUntil,
		arg.Status,
		arg.Note,
		arg.PublicationDate,
		arg.InternalID,
		arg.Byline,
		arg.Budget,
		arg.Hed,
		arg.Description,
		arg.Blurb,
		arg.LedeImage,
		arg.LedeImageCredit,
		arg.LedeImageDescription,
		arg.LedeImageCaption,
		arg.ID,
	)
	var i SharedArticle
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.EmbargoUntil,
		&i.Note,
		&i.SourceType,
		&i.SourceID,
		&i.RawData,
		&i.PageID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PublicationDate,
		&i.InternalID,
		&i.Byline,
		&i.Budget,
		&i.Hed,
		&i.Description,
		&i.LedeImage,
		&i.LedeImageCredit,
		&i.LedeImageDescription,
		&i.LedeImageCaption,
		&i.Blurb,
	)
	return i, err
}

const updateSharedArticleFromGDocs = `-- name: UpdateSharedArticleFromGDocs :one
UPDATE
  shared_article
SET
  "raw_data" = $1,
  "internal_id" = $2,
  "byline" = $3,
  "budget" = $4,
  "hed" = $5,
  "description" = $6,
  "blurb" = $7,
  "lede_image" = $8,
  "lede_image_credit" = $9,
  "lede_image_description" = $10,
  "lede_image_caption" = $11
WHERE
  source_type = 'gdocs'
  AND source_id = $12
RETURNING
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at, publication_date, internal_id, byline, budget, hed, description, lede_image, lede_image_credit, lede_image_description, lede_image_caption, blurb
`

type UpdateSharedArticleFromGDocsParams struct {
	RawData              json.RawMessage `json:"raw_data"`
	InternalID           string          `json:"internal_id"`
	Byline               string          `json:"byline"`
	Budget               string          `json:"budget"`
	Hed                  string          `json:"hed"`
	Description          string          `json:"description"`
	Blurb                string          `json:"blurb"`
	LedeImage            string          `json:"lede_image"`
	LedeImageCredit      string          `json:"lede_image_credit"`
	LedeImageDescription string          `json:"lede_image_description"`
	LedeImageCaption     string          `json:"lede_image_caption"`
	ExternalID           string          `json:"external_id"`
}

func (q *Queries) UpdateSharedArticleFromGDocs(ctx context.Context, arg UpdateSharedArticleFromGDocsParams) (SharedArticle, error) {
	row := q.db.QueryRow(ctx, updateSharedArticleFromGDocs,
		arg.RawData,
		arg.InternalID,
		arg.Byline,
		arg.Budget,
		arg.Hed,
		arg.Description,
		arg.Blurb,
		arg.LedeImage,
		arg.LedeImageCredit,
		arg.LedeImageDescription,
		arg.LedeImageCaption,
		arg.ExternalID,
	)
	var i SharedArticle
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.EmbargoUntil,
		&i.Note,
		&i.SourceType,
		&i.SourceID,
		&i.RawData,
		&i.PageID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PublicationDate,
		&i.InternalID,
		&i.Byline,
		&i.Budget,
		&i.Hed,
		&i.Description,
		&i.LedeImage,
		&i.LedeImageCredit,
		&i.LedeImageDescription,
		&i.LedeImageCaption,
		&i.Blurb,
	)
	return i, err
}

const updateSharedArticlePage = `-- name: UpdateSharedArticlePage :one
UPDATE
  shared_article
SET
  page_id = $1
WHERE
  id = $2
RETURNING
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at, publication_date, internal_id, byline, budget, hed, description, lede_image, lede_image_credit, lede_image_description, lede_image_caption, blurb
`

type UpdateSharedArticlePageParams struct {
	PageID          pgtype.Int8 `json:"page_id"`
	SharedArticleID int64       `json:"shared_article_id"`
}

func (q *Queries) UpdateSharedArticlePage(ctx context.Context, arg UpdateSharedArticlePageParams) (SharedArticle, error) {
	row := q.db.QueryRow(ctx, updateSharedArticlePage, arg.PageID, arg.SharedArticleID)
	var i SharedArticle
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.EmbargoUntil,
		&i.Note,
		&i.SourceType,
		&i.SourceID,
		&i.RawData,
		&i.PageID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PublicationDate,
		&i.InternalID,
		&i.Byline,
		&i.Budget,
		&i.Hed,
		&i.Description,
		&i.LedeImage,
		&i.LedeImageCredit,
		&i.LedeImageDescription,
		&i.LedeImageCaption,
		&i.Blurb,
	)
	return i, err
}

const upsertSharedArticleFromArc = `-- name: UpsertSharedArticleFromArc :one
INSERT INTO shared_article (status, source_type, source_id, raw_data,
  publication_date, budget, description, hed, internal_id)
SELECT
  'U',
  'arc',
  arc.arc_id,
  arc.raw_data,
  iso_to_timestamptz ( --
    arc.raw_data -> 'planning' -> 'scheduling' ->> 'planned_publish_date'),
  arc.raw_data -> 'planning' ->> 'budget_line',
  arc.raw_data -> 'description' ->> 'basic',
  arc.raw_data -> 'headlines' ->> 'basic',
  arc.raw_data ->> 'slug'
FROM
  arc
WHERE
  arc_id = $1
ON CONFLICT (source_type,
  source_id)
  DO UPDATE SET
    raw_data = excluded.raw_data,
    "publication_date" = iso_to_timestamptz ( --
      excluded.raw_data -> 'planning' -> 'scheduling' ->> 'planned_publish_date'),
    "budget" = excluded.raw_data -> 'planning' ->> 'budget_line',
    "description" = excluded.raw_data -> 'description' ->> 'basic',
    "hed" = excluded.raw_data -> 'headlines' ->> 'basic',
    "internal_id" = excluded.raw_data ->> 'slug'
  RETURNING
    id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at, publication_date, internal_id, byline, budget, hed, description, lede_image, lede_image_credit, lede_image_description, lede_image_caption, blurb
`

func (q *Queries) UpsertSharedArticleFromArc(ctx context.Context, arcID string) (SharedArticle, error) {
	row := q.db.QueryRow(ctx, upsertSharedArticleFromArc, arcID)
	var i SharedArticle
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.EmbargoUntil,
		&i.Note,
		&i.SourceType,
		&i.SourceID,
		&i.RawData,
		&i.PageID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PublicationDate,
		&i.InternalID,
		&i.Byline,
		&i.Budget,
		&i.Hed,
		&i.Description,
		&i.LedeImage,
		&i.LedeImageCredit,
		&i.LedeImageDescription,
		&i.LedeImageCaption,
		&i.Blurb,
	)
	return i, err
}

const upsertSharedArticleFromGDocs = `-- name: UpsertSharedArticleFromGDocs :one
INSERT INTO shared_article (status, source_type, source_id, raw_data,
  internal_id, byline, budget, hed, description, blurb, lede_image,
  lede_image_credit, lede_image_description, lede_image_caption)
  VALUES ('U', 'gdocs', $1, $2::jsonb,
    $3, $4, $5, $6, $7, $8, $9,
    $10, $11, $12)
ON CONFLICT (source_type, source_id)
  DO UPDATE SET
    raw_data = excluded.raw_data
  RETURNING
    id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at, publication_date, internal_id, byline, budget, hed, description, lede_image, lede_image_credit, lede_image_description, lede_image_caption, blurb
`

type UpsertSharedArticleFromGDocsParams struct {
	ExternalID           string `json:"external_id"`
	RawData              []byte `json:"raw_data"`
	InternalID           string `json:"internal_id"`
	Byline               string `json:"byline"`
	Budget               string `json:"budget"`
	Hed                  string `json:"hed"`
	Description          string `json:"description"`
	Blurb                string `json:"blurb"`
	LedeImage            string `json:"lede_image"`
	LedeImageCredit      string `json:"lede_image_credit"`
	LedeImageDescription string `json:"lede_image_description"`
	LedeImageCaption     string `json:"lede_image_caption"`
}

func (q *Queries) UpsertSharedArticleFromGDocs(ctx context.Context, arg UpsertSharedArticleFromGDocsParams) (SharedArticle, error) {
	row := q.db.QueryRow(ctx, upsertSharedArticleFromGDocs,
		arg.ExternalID,
		arg.RawData,
		arg.InternalID,
		arg.Byline,
		arg.Budget,
		arg.Hed,
		arg.Description,
		arg.Blurb,
		arg.LedeImage,
		arg.LedeImageCredit,
		arg.LedeImageDescription,
		arg.LedeImageCaption,
	)
	var i SharedArticle
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.EmbargoUntil,
		&i.Note,
		&i.SourceType,
		&i.SourceID,
		&i.RawData,
		&i.PageID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PublicationDate,
		&i.InternalID,
		&i.Byline,
		&i.Budget,
		&i.Hed,
		&i.Description,
		&i.LedeImage,
		&i.LedeImageCredit,
		&i.LedeImageDescription,
		&i.LedeImageCaption,
		&i.Blurb,
	)
	return i, err
}
