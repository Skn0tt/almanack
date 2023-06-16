// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: image.sql

package db

import (
	"context"
)

const getImageByMD5 = `-- name: GetImageByMD5 :one
SELECT
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at, md5, bytes, keywords, deleted_at
FROM
  image
WHERE
  md5 = $1
ORDER BY
  created_at DESC
LIMIT 1
`

func (q *Queries) GetImageByMD5(ctx context.Context, md5 []byte) (Image, error) {
	row := q.db.QueryRow(ctx, getImageByMD5, md5)
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
		&i.MD5,
		&i.Bytes,
		&i.Keywords,
		&i.DeletedAt,
	)
	return i, err
}

const getImageByPath = `-- name: GetImageByPath :one
SELECT
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at, md5, bytes, keywords, deleted_at
FROM
  "image"
WHERE
  "path" = $1
`

func (q *Queries) GetImageByPath(ctx context.Context, path string) (Image, error) {
	row := q.db.QueryRow(ctx, getImageByPath, path)
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
		&i.MD5,
		&i.Bytes,
		&i.Keywords,
		&i.DeletedAt,
	)
	return i, err
}

const getImageBySourceURL = `-- name: GetImageBySourceURL :one
SELECT
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at, md5, bytes, keywords, deleted_at
FROM
  image
WHERE
  src_url = $1
ORDER BY
  updated_at DESC
LIMIT 1
`

func (q *Queries) GetImageBySourceURL(ctx context.Context, srcUrl string) (Image, error) {
	row := q.db.QueryRow(ctx, getImageBySourceURL, srcUrl)
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
		&i.MD5,
		&i.Bytes,
		&i.Keywords,
		&i.DeletedAt,
	)
	return i, err
}

const getImageTypeForExtension = `-- name: GetImageTypeForExtension :one
SELECT
  name, mime, extensions
FROM
  image_type
WHERE
  $1::text = ANY (extensions)
`

func (q *Queries) GetImageTypeForExtension(ctx context.Context, extension string) (ImageType, error) {
	row := q.db.QueryRow(ctx, getImageTypeForExtension, extension)
	var i ImageType
	err := row.Scan(&i.Name, &i.Mime, &i.Extensions)
	return i, err
}

const listImageWhereNotUploaded = `-- name: ListImageWhereNotUploaded :many
SELECT
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at, md5, bytes, keywords, deleted_at
FROM
  image
WHERE
  is_uploaded = FALSE
  AND src_url <> ''
  AND deleted_at IS NULL
`

// ListImageWhereNotUploaded has no limit
// because we want them all uploaded,
// but revisit if queue gets too long.
func (q *Queries) ListImageWhereNotUploaded(ctx context.Context) ([]Image, error) {
	rows, err := q.db.Query(ctx, listImageWhereNotUploaded)
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
			&i.MD5,
			&i.Bytes,
			&i.Keywords,
			&i.DeletedAt,
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

const listImages = `-- name: ListImages :many
SELECT
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at, md5, bytes, keywords, deleted_at
FROM
  image
WHERE
  is_uploaded = TRUE
  AND deleted_at IS NULL
ORDER BY
  updated_at DESC
LIMIT $1 OFFSET $2
`

type ListImagesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListImages(ctx context.Context, arg ListImagesParams) ([]Image, error) {
	rows, err := q.db.Query(ctx, listImages, arg.Limit, arg.Offset)
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
			&i.MD5,
			&i.Bytes,
			&i.Keywords,
			&i.DeletedAt,
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

const listImagesByFTS = `-- name: ListImagesByFTS :many
SELECT
  image.id, image.path, image.type, image.description, image.credit, image.src_url, image.is_uploaded, image.created_at, image.updated_at, image.md5, image.bytes, image.keywords, image.deleted_at
FROM
  image,
  websearch_to_tsquery('english', $3) tsq
WHERE
  fts @@ tsq
  AND is_uploaded
  AND deleted_at IS NULL
ORDER BY
  ts_rank(fts, tsq) DESC,
  updated_at DESC
LIMIT $1 OFFSET $2
`

type ListImagesByFTSParams struct {
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Query  string `json:"query"`
}

func (q *Queries) ListImagesByFTS(ctx context.Context, arg ListImagesByFTSParams) ([]Image, error) {
	rows, err := q.db.Query(ctx, listImagesByFTS, arg.Limit, arg.Offset, arg.Query)
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
			&i.MD5,
			&i.Bytes,
			&i.Keywords,
			&i.DeletedAt,
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

const listImagesWhereNoMD5 = `-- name: ListImagesWhereNoMD5 :many
SELECT
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at, md5, bytes, keywords, deleted_at
FROM
  image
WHERE
  md5 = ''
  AND is_uploaded
  AND deleted_at IS NULL
ORDER BY
  created_at ASC
LIMIT $1
`

func (q *Queries) ListImagesWhereNoMD5(ctx context.Context, limit int32) ([]Image, error) {
	rows, err := q.db.Query(ctx, listImagesWhereNoMD5, limit)
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
			&i.MD5,
			&i.Bytes,
			&i.Keywords,
			&i.DeletedAt,
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
  keywords = CASE WHEN $5::boolean THEN
    $6
  ELSE
    keywords
  END,
  is_uploaded = TRUE
WHERE
  path = $7
RETURNING
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at, md5, bytes, keywords, deleted_at
`

type UpdateImageParams struct {
	SetCredit      bool   `json:"set_credit"`
	Credit         string `json:"credit"`
	SetDescription bool   `json:"set_description"`
	Description    string `json:"description"`
	SetKeywords    bool   `json:"set_keywords"`
	Keywords       string `json:"keywords"`
	Path           string `json:"path"`
}

func (q *Queries) UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error) {
	row := q.db.QueryRow(ctx, updateImage,
		arg.SetCredit,
		arg.Credit,
		arg.SetDescription,
		arg.Description,
		arg.SetKeywords,
		arg.Keywords,
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
		&i.MD5,
		&i.Bytes,
		&i.Keywords,
		&i.DeletedAt,
	)
	return i, err
}

const updateImageMD5Size = `-- name: UpdateImageMD5Size :one
UPDATE
  image
SET
  md5 = $1,
  bytes = $2
WHERE
  id = $3
RETURNING
  id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at, md5, bytes, keywords, deleted_at
`

type UpdateImageMD5SizeParams struct {
	MD5   []byte `json:"md5"`
	Bytes int64  `json:"bytes"`
	ID    int64  `json:"id"`
}

func (q *Queries) UpdateImageMD5Size(ctx context.Context, arg UpdateImageMD5SizeParams) (Image, error) {
	row := q.db.QueryRow(ctx, updateImageMD5Size, arg.MD5, arg.Bytes, arg.ID)
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
		&i.MD5,
		&i.Bytes,
		&i.Keywords,
		&i.DeletedAt,
	)
	return i, err
}

const upsertImage = `-- name: UpsertImage :one
INSERT INTO image ("path", "type", "description", "credit", "keywords",
  "src_url", "is_uploaded")
  VALUES ($1, $2, $3, $4, $5, $6, $7)
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
    END, keywords = CASE WHEN image.keywords = '' THEN
      excluded.keywords
    ELSE
      image.keywords
    END, src_url = CASE WHEN image.src_url = '' THEN
      excluded.src_url
    ELSE
      image.src_url
    END
  RETURNING
    id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at, md5, bytes, keywords, deleted_at
`

type UpsertImageParams struct {
	Path        string `json:"path"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Credit      string `json:"credit"`
	Keywords    string `json:"keywords"`
	SourceURL   string `json:"src_url"`
	IsUploaded  bool   `json:"is_uploaded"`
}

func (q *Queries) UpsertImage(ctx context.Context, arg UpsertImageParams) (Image, error) {
	row := q.db.QueryRow(ctx, upsertImage,
		arg.Path,
		arg.Type,
		arg.Description,
		arg.Credit,
		arg.Keywords,
		arg.SourceURL,
		arg.IsUploaded,
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
		&i.MD5,
		&i.Bytes,
		&i.Keywords,
		&i.DeletedAt,
	)
	return i, err
}

const upsertImageWithMD5 = `-- name: UpsertImageWithMD5 :one
INSERT INTO image ("path", "type", "description", "credit", "keywords",
  "src_url", "md5", "bytes", "is_uploaded")
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE)
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
    END, keywords = CASE WHEN image.keywords = '' THEN
      excluded.keywords
    ELSE
      image.keywords
    END, src_url = CASE WHEN image.src_url = '' THEN
      excluded.src_url
    ELSE
      image.src_url
    END, md5 = CASE WHEN image.md5 = '' THEN
      excluded.md5
    ELSE
      image.md5
    END, bytes = CASE WHEN image.bytes = 0 THEN
      excluded.bytes
    ELSE
      image.bytes
    END
  RETURNING
    id, path, type, description, credit, src_url, is_uploaded, created_at, updated_at, md5, bytes, keywords, deleted_at
`

type UpsertImageWithMD5Params struct {
	Path        string `json:"path"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Credit      string `json:"credit"`
	Keywords    string `json:"keywords"`
	SourceURL   string `json:"src_url"`
	MD5         []byte `json:"md5"`
	Bytes       int64  `json:"bytes"`
}

func (q *Queries) UpsertImageWithMD5(ctx context.Context, arg UpsertImageWithMD5Params) (Image, error) {
	row := q.db.QueryRow(ctx, upsertImageWithMD5,
		arg.Path,
		arg.Type,
		arg.Description,
		arg.Credit,
		arg.Keywords,
		arg.SourceURL,
		arg.MD5,
		arg.Bytes,
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
		&i.MD5,
		&i.Bytes,
		&i.Keywords,
		&i.DeletedAt,
	)
	return i, err
}
