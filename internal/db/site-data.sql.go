// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: site-data.sql

package db

import (
	"context"
	"time"
)

const cleanSiteData = `-- name: CleanSiteData :exec
DELETE FROM site_data
WHERE key = $1::text
  AND published_at < (
    SELECT
      max(published_at) AS max
    FROM
      site_data
    WHERE
      key = $1::text
    GROUP BY
      key)
`

func (q *Queries) CleanSiteData(ctx context.Context, key string) error {
	_, err := q.db.Exec(ctx, cleanSiteData, key)
	return err
}

const deleteSiteData = `-- name: DeleteSiteData :exec
DELETE FROM site_data
WHERE "key" = $1
  AND "schedule_for" > (CURRENT_TIMESTAMP + '5 minutes'::interval)
`

// DeleteSiteData only removes future scheduled items.
// To remove past scheduled items, use CleanSiteData
func (q *Queries) DeleteSiteData(ctx context.Context, key string) error {
	_, err := q.db.Exec(ctx, deleteSiteData, key)
	return err
}

const getSiteData = `-- name: GetSiteData :many
SELECT
  id, key, data, created_at, updated_at, schedule_for, published_at
FROM
  site_data
WHERE
  key = $1::text
  AND (published_at IS NULL
    OR published_at = (
      SELECT
        max(published_at) AS max
      FROM
        site_data
      WHERE
        key = $1::text
      GROUP BY
        key))
ORDER BY
  schedule_for ASC
`

func (q *Queries) GetSiteData(ctx context.Context, key string) ([]SiteDatum, error) {
	rows, err := q.db.Query(ctx, getSiteData, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SiteDatum
	for rows.Next() {
		var i SiteDatum
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ScheduleFor,
			&i.PublishedAt,
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

const popScheduledSiteChanges = `-- name: PopScheduledSiteChanges :many
UPDATE
  site_data
SET
  published_at = CURRENT_TIMESTAMP
WHERE
  key = $1::text
  AND published_at IS NULL
  AND schedule_for < (CURRENT_TIMESTAMP + '5 minutes'::interval)
RETURNING
  id, key, data, created_at, updated_at, schedule_for, published_at
`

func (q *Queries) PopScheduledSiteChanges(ctx context.Context, key string) ([]SiteDatum, error) {
	rows, err := q.db.Query(ctx, popScheduledSiteChanges, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SiteDatum
	for rows.Next() {
		var i SiteDatum
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ScheduleFor,
			&i.PublishedAt,
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

const upsertSiteData = `-- name: UpsertSiteData :exec
INSERT INTO site_data ("key", "data", "schedule_for")
  VALUES ($1, $2, $3)
ON CONFLICT ("key", "schedule_for")
  DO UPDATE SET
    data = excluded.data
`

type UpsertSiteDataParams struct {
	Key         string    `json:"key"`
	Data        Map       `json:"data"`
	ScheduleFor time.Time `json:"schedule_for"`
}

func (q *Queries) UpsertSiteData(ctx context.Context, arg UpsertSiteDataParams) error {
	_, err := q.db.Exec(ctx, upsertSiteData, arg.Key, arg.Data, arg.ScheduleFor)
	return err
}
