// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: shared-article.sql

package db

import (
	"context"
	"database/sql"

	"github.com/jackc/pgtype"
)

const getSharedArticleByID = `-- name: GetSharedArticleByID :one
SELECT
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at
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
	)
	return i, err
}

const getSharedArticleBySource = `-- name: GetSharedArticleBySource :one
SELECT
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at
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
	)
	return i, err
}

const insertSharedArticleFromArc = `-- name: InsertSharedArticleFromArc :one
INSERT INTO shared_article (status, source_type, source_id, raw_data)
SELECT
  'U',
  'arc',
  arc.arc_id,
  arc.raw_data
FROM
  arc
WHERE
  arc_id = $1
ON CONFLICT (source_type,
  source_id)
  DO UPDATE SET
    raw_data = excluded.raw_data
  RETURNING
    id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at
`

func (q *Queries) InsertSharedArticleFromArc(ctx context.Context, arcID string) (SharedArticle, error) {
	row := q.db.QueryRow(ctx, insertSharedArticleFromArc, arcID)
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
	)
	return i, err
}

const listSharedArticles = `-- name: ListSharedArticles :many
SELECT
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at
FROM
  shared_article
ORDER BY
  updated_at DESC
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
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at
FROM
  shared_article
WHERE
  status <> 'U'
ORDER BY
  CASE status
  WHEN 'P' THEN
    '0'
  WHEN 'S' THEN
    '1'
  END ASC,
  updated_at DESC
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
  embargo_until = $1,
  status = $2,
  note = $3
WHERE
  id = $4
RETURNING
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at
`

type UpdateSharedArticleParams struct {
	EmbargoUntil pgtype.Timestamptz `json:"embargo_until"`
	Status       string             `json:"status"`
	Note         string             `json:"note"`
	ID           int64              `json:"id"`
}

func (q *Queries) UpdateSharedArticle(ctx context.Context, arg UpdateSharedArticleParams) (SharedArticle, error) {
	row := q.db.QueryRow(ctx, updateSharedArticle,
		arg.EmbargoUntil,
		arg.Status,
		arg.Note,
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
  id, status, embargo_until, note, source_type, source_id, raw_data, page_id, created_at, updated_at
`

type UpdateSharedArticlePageParams struct {
	PageID          sql.NullInt64 `json:"page_id"`
	SharedArticleID int64         `json:"shared_article_id"`
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
	)
	return i, err
}
