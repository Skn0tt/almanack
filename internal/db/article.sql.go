// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: article.sql

package db

import (
	"context"

	"github.com/jackc/pgtype"
)

const getArticleByArcID = `-- name: GetArticleByArcID :one
SELECT
  id, arc_id, arc_data, spotlightpa_path, note, status, created_at, updated_at
FROM
  article
WHERE
  arc_id = $1::text
`

func (q *Queries) GetArticleByArcID(ctx context.Context, arcID string) (Article, error) {
	row := q.db.QueryRow(ctx, getArticleByArcID, arcID)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.ArcID,
		&i.ArcData,
		&i.SpotlightPAPath,
		&i.Note,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAllArcArticles = `-- name: ListAllArcArticles :many
SELECT
  id, arc_id, arc_data, spotlightpa_path, note, status, created_at, updated_at
FROM
  article
WHERE
  arc_id IS NOT NULL
ORDER BY
  arc_data ->> 'last_updated_date' DESC
LIMIT $1 OFFSET $2
`

type ListAllArcArticlesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAllArcArticles(ctx context.Context, arg ListAllArcArticlesParams) ([]Article, error) {
	rows, err := q.db.Query(ctx, listAllArcArticles, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Article
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.ID,
			&i.ArcID,
			&i.ArcData,
			&i.SpotlightPAPath,
			&i.Note,
			&i.Status,
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

const listAvailableArticles = `-- name: ListAvailableArticles :many
SELECT
  id, arc_id, arc_data, spotlightpa_path, note, status, created_at, updated_at
FROM
  article
WHERE
  status <> 'U'
ORDER BY
  CASE status
  WHEN 'P' THEN
    '0'
  WHEN 'A' THEN
    '1'
  END ASC,
  arc_data ->> 'last_updated_date' DESC
LIMIT $1 OFFSET $2
`

type ListAvailableArticlesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAvailableArticles(ctx context.Context, arg ListAvailableArticlesParams) ([]Article, error) {
	rows, err := q.db.Query(ctx, listAvailableArticles, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Article
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.ID,
			&i.ArcID,
			&i.ArcData,
			&i.SpotlightPAPath,
			&i.Note,
			&i.Status,
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

const updateAlmanackArticle = `-- name: UpdateAlmanackArticle :one
UPDATE
  article
SET
  status = $1,
  note = $2,
  arc_data = CASE WHEN $3::bool THEN
    $4::jsonb
  ELSE
    arc_data
  END
WHERE
  arc_id = $5::text
RETURNING
  id, arc_id, arc_data, spotlightpa_path, note, status, created_at, updated_at
`

type UpdateAlmanackArticleParams struct {
	Status     string       `json:"status"`
	Note       string       `json:"note"`
	SetArcData bool         `json:"set_arc_data"`
	ArcData    pgtype.JSONB `json:"arc_data"`
	ArcID      string       `json:"arc_id"`
}

func (q *Queries) UpdateAlmanackArticle(ctx context.Context, arg UpdateAlmanackArticleParams) (Article, error) {
	row := q.db.QueryRow(ctx, updateAlmanackArticle,
		arg.Status,
		arg.Note,
		arg.SetArcData,
		arg.ArcData,
		arg.ArcID,
	)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.ArcID,
		&i.ArcData,
		&i.SpotlightPAPath,
		&i.Note,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateArcArticleSpotlightPAPath = `-- name: UpdateArcArticleSpotlightPAPath :one
UPDATE
  article
SET
  spotlightpa_path = $1::text
WHERE
  arc_id = $2::text
RETURNING
  id, arc_id, arc_data, spotlightpa_path, note, status, created_at, updated_at
`

type UpdateArcArticleSpotlightPAPathParams struct {
	SpotlightPAPath string `json:"spotlightpa_path"`
	ArcID           string `json:"arc_id"`
}

func (q *Queries) UpdateArcArticleSpotlightPAPath(ctx context.Context, arg UpdateArcArticleSpotlightPAPathParams) (Article, error) {
	row := q.db.QueryRow(ctx, updateArcArticleSpotlightPAPath, arg.SpotlightPAPath, arg.ArcID)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.ArcID,
		&i.ArcData,
		&i.SpotlightPAPath,
		&i.Note,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateArcArticles = `-- name: UpdateArcArticles :exec
WITH arc_table AS (
  SELECT
    jsonb_array_elements($1::jsonb) AS article_data)
INSERT INTO article (arc_id, arc_data)
SELECT
  article_data ->> '_id',
  article_data
FROM
  arc_table
ON CONFLICT (arc_id)
  DO UPDATE SET
    arc_data = excluded.arc_data
  WHERE
    article.status = 'U'
`

func (q *Queries) UpdateArcArticles(ctx context.Context, arcItems pgtype.JSONB) error {
	_, err := q.db.Exec(ctx, updateArcArticles, arcItems)
	return err
}
