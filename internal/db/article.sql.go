// Code generated by sqlc. DO NOT EDIT.
// source: article.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"
)

const getArticleByArcID = `-- name: GetArticleByArcID :one
SELECT
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
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
	row := q.db.QueryRow(ctx, getArticleByDBID, id)
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

const listAllArcArticles = `-- name: ListAllArcArticles :many
SELECT
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
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
  arc_data ->> 'last_updated_date' DESC
`

func (q *Queries) ListUpcoming(ctx context.Context) ([]Article, error) {
	rows, err := q.db.Query(ctx, listUpcoming)
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
  AND schedule_for < (CURRENT_TIMESTAMP + '5 minutes'::interval)
RETURNING
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
`

func (q *Queries) PopScheduled(ctx context.Context) ([]Article, error) {
	rows, err := q.db.Query(ctx, popScheduled)
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

const updateArcArticleSpotlightPAPath = `-- name: UpdateArcArticleSpotlightPAPath :one
UPDATE
  article
SET
  spotlightpa_path = $1::text
WHERE
  arc_id = $2::text
RETURNING
  id, arc_id, arc_data, spotlightpa_path, spotlightpa_data, schedule_for, last_published, note, status, created_at, updated_at
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
    article.status = 'U'
`

func (q *Queries) UpdateArcArticles(ctx context.Context, arcItems json.RawMessage) error {
	_, err := q.db.Exec(ctx, updateArcArticles, arcItems)
	return err
}

const updateSpotlightPAArticle = `-- name: UpdateSpotlightPAArticle :one
UPDATE
  article AS "new"
SET
  spotlightpa_data = $1,
  schedule_for = $2,
  spotlightpa_path = CASE WHEN "new".spotlightpa_path IS NULL THEN
    $3
  ELSE
    "new".spotlightpa_path
  END
FROM
  article AS "old"
WHERE
  "new".id = "old".id
  AND "old".arc_id = $4
RETURNING
  "old".schedule_for
`

type UpdateSpotlightPAArticleParams struct {
	SpotlightPAData json.RawMessage `json:"spotlightpa_data"`
	ScheduleFor     sql.NullTime    `json:"schedule_for"`
	SpotlightPAPath sql.NullString  `json:"spotlightpa_path"`
	ArcID           sql.NullString  `json:"arc_id"`
}

func (q *Queries) UpdateSpotlightPAArticle(ctx context.Context, arg UpdateSpotlightPAArticleParams) (sql.NullTime, error) {
	row := q.db.QueryRow(ctx, updateSpotlightPAArticle,
		arg.SpotlightPAData,
		arg.ScheduleFor,
		arg.SpotlightPAPath,
		arg.ArcID,
	)
	var schedule_for sql.NullTime
	err := row.Scan(&schedule_for)
	return schedule_for, err
}

const updateSpotlightPAArticleLastPublished = `-- name: UpdateSpotlightPAArticleLastPublished :one
UPDATE
  article AS "new"
SET
  last_published = CURRENT_TIMESTAMP
FROM
  article AS "old"
WHERE
  "new".id = "old".id
  AND "old".arc_id = $1::text
RETURNING
  "old".last_published
`

func (q *Queries) UpdateSpotlightPAArticleLastPublished(ctx context.Context, arcID string) (sql.NullTime, error) {
	row := q.db.QueryRow(ctx, updateSpotlightPAArticleLastPublished, arcID)
	var last_published sql.NullTime
	err := row.Scan(&last_published)
	return last_published, err
}
