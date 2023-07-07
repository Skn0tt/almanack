// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: arc.sql

package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const getArcByArcID = `-- name: GetArcByArcID :one
SELECT
  id, arc_id, raw_data, last_updated, created_at, updated_at
FROM
  arc
WHERE
  arc_id = $1
`

func (q *Queries) GetArcByArcID(ctx context.Context, arcID string) (Arc, error) {
	row := q.db.QueryRow(ctx, getArcByArcID, arcID)
	var i Arc
	err := row.Scan(
		&i.ID,
		&i.ArcID,
		&i.RawData,
		&i.LastUpdated,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listArcByLastUpdated = `-- name: ListArcByLastUpdated :many
SELECT
  arc.id, arc.arc_id, arc.raw_data, arc.last_updated, arc.created_at, arc.updated_at,
  shared_article.id AS shared_article_id,
  coalesce(shared_article.status, ''),
  shared_article.embargo_until
FROM
  arc
  LEFT JOIN shared_article ON (arc.arc_id = shared_article.source_id
      AND shared_article.source_type = 'arc')
  ORDER BY
    last_updated DESC
  LIMIT $1 OFFSET $2
`

type ListArcByLastUpdatedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListArcByLastUpdatedRow struct {
	ID              int64              `json:"id"`
	ArcID           string             `json:"arc_id"`
	RawData         json.RawMessage    `json:"raw_data"`
	LastUpdated     pgtype.Timestamptz `json:"last_updated"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	SharedArticleID pgtype.Int8        `json:"shared_article_id"`
	Status          string             `json:"status"`
	EmbargoUntil    pgtype.Timestamptz `json:"embargo_until"`
}

func (q *Queries) ListArcByLastUpdated(ctx context.Context, arg ListArcByLastUpdatedParams) ([]ListArcByLastUpdatedRow, error) {
	rows, err := q.db.Query(ctx, listArcByLastUpdated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListArcByLastUpdatedRow
	for rows.Next() {
		var i ListArcByLastUpdatedRow
		if err := rows.Scan(
			&i.ID,
			&i.ArcID,
			&i.RawData,
			&i.LastUpdated,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SharedArticleID,
			&i.Status,
			&i.EmbargoUntil,
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

const updateArc = `-- name: UpdateArc :exec
WITH arc_temp AS (
  SELECT
    jsonb_array_elements($1::jsonb) AS temp_data)
INSERT INTO arc (arc_id, raw_data)
SELECT
  temp_data ->> '_id',
  temp_data
FROM
  arc_temp
ON CONFLICT (arc_id)
  DO UPDATE SET
    raw_data = excluded.raw_data
`

func (q *Queries) UpdateArc(ctx context.Context, arcItems []byte) error {
	_, err := q.db.Exec(ctx, updateArc, arcItems)
	return err
}
