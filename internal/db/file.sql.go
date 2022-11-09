// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: file.sql

package db

import (
	"context"
)

const createFilePlaceholder = `-- name: CreateFilePlaceholder :execrows
INSERT INTO file ("filename", "url", "mime_type")
  VALUES ($1, $2, $3)
ON CONFLICT (url)
  DO NOTHING
`

type CreateFilePlaceholderParams struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Type     string `json:"type"`
}

func (q *Queries) CreateFilePlaceholder(ctx context.Context, arg CreateFilePlaceholderParams) (int64, error) {
	result, err := q.db.Exec(ctx, createFilePlaceholder, arg.Filename, arg.URL, arg.Type)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const listFiles = `-- name: ListFiles :many
SELECT
  id, url, filename, mime_type, description, is_uploaded, created_at, updated_at
FROM
  file
WHERE
  is_uploaded = TRUE
ORDER BY
  created_at DESC
LIMIT $1 OFFSET $2
`

type ListFilesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListFiles(ctx context.Context, arg ListFilesParams) ([]File, error) {
	rows, err := q.db.Query(ctx, listFiles, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []File
	for rows.Next() {
		var i File
		if err := rows.Scan(
			&i.ID,
			&i.URL,
			&i.Filename,
			&i.MimeType,
			&i.Description,
			&i.IsUploaded,
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

const updateFile = `-- name: UpdateFile :one
UPDATE
  file
SET
  description = CASE WHEN $1::boolean THEN
    $2
  ELSE
    description
  END,
  is_uploaded = TRUE
WHERE
  url = $3
RETURNING
  id, url, filename, mime_type, description, is_uploaded, created_at, updated_at
`

type UpdateFileParams struct {
	SetDescription bool   `json:"set_description"`
	Description    string `json:"description"`
	URL            string `json:"url"`
}

func (q *Queries) UpdateFile(ctx context.Context, arg UpdateFileParams) (File, error) {
	row := q.db.QueryRow(ctx, updateFile, arg.SetDescription, arg.Description, arg.URL)
	var i File
	err := row.Scan(
		&i.ID,
		&i.URL,
		&i.Filename,
		&i.MimeType,
		&i.Description,
		&i.IsUploaded,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
