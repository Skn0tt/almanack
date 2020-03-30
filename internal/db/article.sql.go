// Code generated by sqlc. DO NOT EDIT.
// source: article.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

const getArticle = `-- name: GetArticle :one
SELECT
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
FROM
  article
WHERE
  arc_id = $1
`

func (q *Queries) GetArticle(ctx context.Context, arcID sql.NullString) (Article, error) {
	row := q.db.QueryRowContext(ctx, getArticle, arcID)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.ArcID,
		&i.ArcData,
		&i.SpotlightPAPath,
		&i.SpotlightPAData,
		&i.ScheduleFor,
		&i.LastPublished,
		&i.Note,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getArticleByDBID = `-- name: GetArticleByDBID :one
SELECT
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
FROM
  article
WHERE
  id = $1
`

func (q *Queries) GetArticleByDBID(ctx context.Context, id int32) (Article, error) {
	row := q.db.QueryRowContext(ctx, getArticleByDBID, id)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.ArcID,
		&i.ArcData,
		&i.SpotlightPAPath,
		&i.SpotlightPAData,
		&i.ScheduleFor,
		&i.LastPublished,
		&i.Note,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAllArticles = `-- name: ListAllArticles :many
SELECT
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
FROM
  article
ORDER BY
  arc_data -> 'last_updated_date' DESC
`

func (q *Queries) ListAllArticles(ctx context.Context) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, listAllArticles)
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
			&i.SpotlightPAData,
			&i.ScheduleFor,
			&i.LastPublished,
			&i.Note,
			&i.Status,
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

const listAvailableArticles = `-- name: ListAvailableArticles :many
SELECT
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
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
  arc_data -> 'last_updated_date' DESC
`

func (q *Queries) ListAvailableArticles(ctx context.Context) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, listAvailableArticles)
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
			&i.SpotlightPAData,
			&i.ScheduleFor,
			&i.LastPublished,
			&i.Note,
			&i.Status,
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

const listSpotlightPAArticles = `-- name: ListSpotlightPAArticles :many
SELECT
  id,
  arc_id::text,
  (spotlightpa_data ->> 'internal-id')::text AS internal_id,
  (spotlightpa_data ->> 'hed')::text AS hed,
  ARRAY (
    SELECT
      jsonb_array_elements_text(spotlightpa_data -> 'authors'))::text[] AS authors,
  to_timestamp(spotlightpa_data ->> 'pub-date'::text,
    -- ISO date
    'YYYY-MM-DD"T"HH24:MI:SS"Z"')::timestamp AS pub_date
FROM
  article
WHERE
  spotlightpa_path IS NOT NULL
ORDER BY
  pub_date DESC
`

type ListSpotlightPAArticlesRow struct {
	ID         int32     `json:"id"`
	ArcID      string    `json:"arc_id"`
	InternalID string    `json:"internal_id"`
	Hed        string    `json:"hed"`
	Authors    []string  `json:"authors"`
	PubDate    time.Time `json:"pub_date"`
}

func (q *Queries) ListSpotlightPAArticles(ctx context.Context) ([]ListSpotlightPAArticlesRow, error) {
	rows, err := q.db.QueryContext(ctx, listSpotlightPAArticles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListSpotlightPAArticlesRow
	for rows.Next() {
		var i ListSpotlightPAArticlesRow
		if err := rows.Scan(
			&i.ID,
			&i.ArcID,
			&i.InternalID,
			&i.Hed,
			pq.Array(&i.Authors),
			&i.PubDate,
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

const listUpcoming = `-- name: ListUpcoming :many
SELECT
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
FROM
  article
ORDER BY
  arc_data -> 'last_updated_date' DESC
`

func (q *Queries) ListUpcoming(ctx context.Context) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, listUpcoming)
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
			&i.SpotlightPAData,
			&i.ScheduleFor,
			&i.LastPublished,
			&i.Note,
			&i.Status,
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

const popScheduled = `-- name: PopScheduled :many
UPDATE
  article
SET
  last_published = CURRENT_TIMESTAMP
WHERE
  last_published IS NULL
  AND schedule_for < CURRENT_TIMESTAMP
RETURNING
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
`

func (q *Queries) PopScheduled(ctx context.Context) ([]Article, error) {
	rows, err := q.db.QueryContext(ctx, popScheduled)
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
			&i.SpotlightPAData,
			&i.ScheduleFor,
			&i.LastPublished,
			&i.Note,
			&i.Status,
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
  arc_id = $5
RETURNING
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
`

type UpdateAlmanackArticleParams struct {
	Status     string          `json:"status"`
	Note       string          `json:"note"`
	SetArcData bool            `json:"set_arc_data"`
	ArcData    json.RawMessage `json:"arc_data"`
	ArcID      sql.NullString  `json:"arc_id"`
}

func (q *Queries) UpdateAlmanackArticle(ctx context.Context, arg UpdateAlmanackArticleParams) (Article, error) {
	row := q.db.QueryRowContext(ctx, updateAlmanackArticle,
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
		&i.SpotlightPAData,
		&i.ScheduleFor,
		&i.LastPublished,
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
    article.status <> 'A'
`

func (q *Queries) UpdateArcArticles(ctx context.Context, arcItems json.RawMessage) error {
	_, err := q.db.ExecContext(ctx, updateArcArticles, arcItems)
	return err
}

const updateSpotlightPAArticle = `-- name: UpdateSpotlightPAArticle :one
UPDATE
  article
SET
  spotlightpa_data = $1,
  schedule_for = $2,
  spotlightpa_path = CASE WHEN spotlightpa_path IS NULL THEN
    $3
  ELSE
    spotlightpa_path
  END,
  last_published = CASE WHEN $4::bool THEN
    CURRENT_TIMESTAMP
  ELSE
    last_published
  END
WHERE
  arc_id = $5
RETURNING
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
`

type UpdateSpotlightPAArticleParams struct {
	SpotlightPAData  json.RawMessage `json:"spotlightpa_data"`
	ScheduleFor      sql.NullTime    `json:"schedule_for"`
	SpotlightPAPath  sql.NullString  `json:"spotlightpa_path"`
	SetLastPublished bool            `json:"set_last_published"`
	ArcID            sql.NullString  `json:"arc_id"`
}

func (q *Queries) UpdateSpotlightPAArticle(ctx context.Context, arg UpdateSpotlightPAArticleParams) (Article, error) {
	row := q.db.QueryRowContext(ctx, updateSpotlightPAArticle,
		arg.SpotlightPAData,
		arg.ScheduleFor,
		arg.SpotlightPAPath,
		arg.SetLastPublished,
		arg.ArcID,
	)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.ArcID,
		&i.ArcData,
		&i.SpotlightPAPath,
		&i.SpotlightPAData,
		&i.ScheduleFor,
		&i.LastPublished,
		&i.Note,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
