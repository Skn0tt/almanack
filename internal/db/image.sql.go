// Code generated by sqlc. DO NOT EDIT.
// source: image.sql

package db

import (
	"context"
)

const createImage = `-- name: CreateImage :execrows
INSERT INTO image ("path", "type", "description", "credit", "src_url", "is_uploaded")
  VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (path)
  DO UPDATE SET
    credit = CASE WHEN image.credit = '' THEN
      excluded.credit
    ELSE
      image.credit
    END, description = CASE WHEN image.description = '' THEN
      excluded.description
    ELSE
      image.description
    END, src_url = CASE WHEN image.src_url = '' THEN
      excluded.src_url
    ELSE
      image.src_url
    END
`

type CreateImageParams struct {
	Path        string `json:"path"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Credit      string `json:"credit"`
	SourceURL   string `json:"src_url"`
	IsUploaded  bool   `json:"is_uploaded"`
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, createImage,
		arg.Path,
		arg.Type,
		arg.Description,
		arg.Credit,
		arg.SourceURL,
		arg.IsUploaded,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const createImagePlaceholder = `-- name: CreateImagePlaceholder :execrows
INSERT INTO image ("path", "type")
  VALUES ($1, $2)
ON CONFLICT (path)
  DO NOTHING
`

type CreateImagePlaceholderParams struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

func (q *Queries) CreateImagePlaceholder(ctx context.Context, arg CreateImagePlaceholderParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, createImagePlaceholder, arg.Path, arg.Type)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const getImageBySourceURL = `-- name: GetImageBySourceURL :one
SELECT
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at
FROM
  image
WHERE
  src_url = $1
ORDER BY
  updated_at DESC
LIMIT 1
`

func (q *Queries) GetImageBySourceURL(ctx context.Context, srcUrl string) (Image, error) {
	row := q.db.QueryRowContext(ctx, getImageBySourceURL, srcUrl)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Path,
		&i.Type,
		&i.Description,
		&i.Credit,
		&i.SourceURL,
		&i.IsUploaded,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listImages = `-- name: ListImages :many
SELECT
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at
FROM
  image
WHERE
  is_uploaded = TRUE
ORDER BY
  updated_at DESC
LIMIT $1 OFFSET $2
`

type ListImagesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListImages(ctx context.Context, arg ListImagesParams) ([]Image, error) {
	rows, err := q.db.QueryContext(ctx, listImages, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ID,
			&i.Path,
			&i.Type,
			&i.Description,
			&i.Credit,
			&i.SourceURL,
			&i.IsUploaded,
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

const updateImage = `-- name: UpdateImage :one
UPDATE
  image
SET
  credit = CASE WHEN $1::boolean THEN
    $2
  ELSE
    credit
  END,
  description = CASE WHEN $3::boolean THEN
    $4
  ELSE
    description
  END,
  src_url = CASE WHEN src_url = '' THEN
    $5
  ELSE
    src_url
  END,
  is_uploaded = TRUE
WHERE
  path = $6
RETURNING
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at
`

type UpdateImageParams struct {
	SetCredit      bool   `json:"set_credit"`
	Credit         string `json:"credit"`
	SetDescription bool   `json:"set_description"`
	Description    string `json:"description"`
	SourceURL      string `json:"src_url"`
	Path           string `json:"path"`
}

func (q *Queries) UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error) {
	row := q.db.QueryRowContext(ctx, updateImage,
		arg.SetCredit,
		arg.Credit,
		arg.SetDescription,
		arg.Description,
		arg.SourceURL,
		arg.Path,
	)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.Path,
		&i.Type,
		&i.Description,
		&i.Credit,
		&i.SourceURL,
		&i.IsUploaded,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
