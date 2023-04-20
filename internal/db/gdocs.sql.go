// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: gdocs.sql

package db

import (
	"context"

	docs "google.golang.org/api/docs/v1"
)

const createGDocsDoc = `-- name: CreateGDocsDoc :one
INSERT INTO g_docs_doc ("g_docs_id", "document")
  VALUES ($1, $2)
RETURNING
  id, g_docs_id, document, embeds, rich_text, raw_html, article_markdown, word_count, warnings, processed_at, created_at
`

type CreateGDocsDocParams struct {
	GDocsID  string        `json:"g_docs_id"`
	Document docs.Document `json:"document"`
}

func (q *Queries) CreateGDocsDoc(ctx context.Context, arg CreateGDocsDocParams) (GDocsDoc, error) {
	row := q.db.QueryRow(ctx, createGDocsDoc, arg.GDocsID, arg.Document)
	var i GDocsDoc
	err := row.Scan(
		&i.ID,
		&i.GDocsID,
		&i.Document,
		&i.Embeds,
		&i.RichText,
		&i.RawHtml,
		&i.ArticleMarkdown,
		&i.WordCount,
		&i.Warnings,
		&i.ProcessedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getGDocsByGDocIDWhereProcessed = `-- name: GetGDocsByGDocIDWhereProcessed :one
SELECT
  id, g_docs_id, document, embeds, rich_text, raw_html, article_markdown, word_count, warnings, processed_at, created_at
FROM
  g_docs_doc
WHERE
  g_docs_id = $1
  AND processed_at IS NOT NULL
ORDER BY
  processed_at DESC
LIMIT 1
`

func (q *Queries) GetGDocsByGDocIDWhereProcessed(ctx context.Context, gDocsID string) (GDocsDoc, error) {
	row := q.db.QueryRow(ctx, getGDocsByGDocIDWhereProcessed, gDocsID)
	var i GDocsDoc
	err := row.Scan(
		&i.ID,
		&i.GDocsID,
		&i.Document,
		&i.Embeds,
		&i.RichText,
		&i.RawHtml,
		&i.ArticleMarkdown,
		&i.WordCount,
		&i.Warnings,
		&i.ProcessedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getGDocsByID = `-- name: GetGDocsByID :one
SELECT
  id, g_docs_id, document, embeds, rich_text, raw_html, article_markdown, word_count, warnings, processed_at, created_at
FROM
  g_docs_doc
WHERE
  id = $1
`

func (q *Queries) GetGDocsByID(ctx context.Context, id int64) (GDocsDoc, error) {
	row := q.db.QueryRow(ctx, getGDocsByID, id)
	var i GDocsDoc
	err := row.Scan(
		&i.ID,
		&i.GDocsID,
		&i.Document,
		&i.Embeds,
		&i.RichText,
		&i.RawHtml,
		&i.ArticleMarkdown,
		&i.WordCount,
		&i.Warnings,
		&i.ProcessedAt,
		&i.CreatedAt,
	)
	return i, err
}

const listGDocsImagesByGDocsID = `-- name: ListGDocsImagesByGDocsID :many
SELECT
  "doc_object_id",
  "path"::text,
  "type"::text
FROM
  g_docs_image
  LEFT JOIN image ON (image_id = image.id)
WHERE
  g_docs_id = $1
`

type ListGDocsImagesByGDocsIDRow struct {
	DocObjectID string `json:"doc_object_id"`
	Path        string `json:"path"`
	Type        string `json:"type"`
}

func (q *Queries) ListGDocsImagesByGDocsID(ctx context.Context, gDocsID string) ([]ListGDocsImagesByGDocsIDRow, error) {
	rows, err := q.db.Query(ctx, listGDocsImagesByGDocsID, gDocsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListGDocsImagesByGDocsIDRow
	for rows.Next() {
		var i ListGDocsImagesByGDocsIDRow
		if err := rows.Scan(&i.DocObjectID, &i.Path, &i.Type); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listGDocsWhereUnprocessed = `-- name: ListGDocsWhereUnprocessed :many
SELECT
  id, g_docs_id, document, embeds, rich_text, raw_html, article_markdown, word_count, warnings, processed_at, created_at
FROM
  g_docs_doc
WHERE
  processed_at IS NULL
`

func (q *Queries) ListGDocsWhereUnprocessed(ctx context.Context) ([]GDocsDoc, error) {
	rows, err := q.db.Query(ctx, listGDocsWhereUnprocessed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GDocsDoc
	for rows.Next() {
		var i GDocsDoc
		if err := rows.Scan(
			&i.ID,
			&i.GDocsID,
			&i.Document,
			&i.Embeds,
			&i.RichText,
			&i.RawHtml,
			&i.ArticleMarkdown,
			&i.WordCount,
			&i.Warnings,
			&i.ProcessedAt,
			&i.CreatedAt,
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

const updateGDocsDoc = `-- name: UpdateGDocsDoc :one
UPDATE
  g_docs_doc
SET
  "embeds" = $1,
  "rich_text" = $2,
  "raw_html" = $3,
  "article_markdown" = $4,
  "word_count" = $5,
  "warnings" = $6,
  "processed_at" = CURRENT_TIMESTAMP
WHERE
  id = $7
RETURNING
  id, g_docs_id, document, embeds, rich_text, raw_html, article_markdown, word_count, warnings, processed_at, created_at
`

type UpdateGDocsDocParams struct {
	Embeds          []Embed  `json:"embeds"`
	RichText        string   `json:"rich_text"`
	RawHtml         string   `json:"raw_html"`
	ArticleMarkdown string   `json:"article_markdown"`
	WordCount       int32    `json:"word_count"`
	Warnings        []string `json:"warnings"`
	ID              int64    `json:"id"`
}

func (q *Queries) UpdateGDocsDoc(ctx context.Context, arg UpdateGDocsDocParams) (GDocsDoc, error) {
	row := q.db.QueryRow(ctx, updateGDocsDoc,
		arg.Embeds,
		arg.RichText,
		arg.RawHtml,
		arg.ArticleMarkdown,
		arg.WordCount,
		arg.Warnings,
		arg.ID,
	)
	var i GDocsDoc
	err := row.Scan(
		&i.ID,
		&i.GDocsID,
		&i.Document,
		&i.Embeds,
		&i.RichText,
		&i.RawHtml,
		&i.ArticleMarkdown,
		&i.WordCount,
		&i.Warnings,
		&i.ProcessedAt,
		&i.CreatedAt,
	)
	return i, err
}

const upsertGDocsImage = `-- name: UpsertGDocsImage :exec
INSERT INTO g_docs_image (g_docs_id, doc_object_id, image_id)
  VALUES ($1, $2, $3)
ON CONFLICT (g_docs_id, doc_object_id)
  DO NOTHING
`

type UpsertGDocsImageParams struct {
	GDocsID     string `json:"g_docs_id"`
	DocObjectID string `json:"doc_object_id"`
	ImageID     int64  `json:"image_id"`
}

func (q *Queries) UpsertGDocsImage(ctx context.Context, arg UpsertGDocsImageParams) error {
	_, err := q.db.Exec(ctx, upsertGDocsImage, arg.GDocsID, arg.DocObjectID, arg.ImageID)
	return err
}
