// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Article struct {
	ID              int32           `json:"id"`
	ArcID           sql.NullString  `json:"arc_id"`
	ArcData         json.RawMessage `json:"arc_data"`
	SpotlightPAPath sql.NullString  `json:"spotlightpa_path"`
	SpotlightPAData json.RawMessage `json:"spotlightpa_data"`
	ScheduleFor     sql.NullTime    `json:"schedule_for"`
	LastPublished   sql.NullTime    `json:"last_published"`
	Note            string          `json:"note"`
	Status          string          `json:"status"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type ArticleStatus struct {
	StatusID    string `json:"status_id"`
	Description string `json:"description"`
}

type DomainRole struct {
	ID        int32     `json:"id"`
	Domain    string    `json:"domain"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type File struct {
	ID          int32     `json:"id"`
	URL         string    `json:"url"`
	Filename    string    `json:"filename"`
	MimeType    string    `json:"mime_type"`
	Description string    `json:"description"`
	IsUploaded  bool      `json:"is_uploaded"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Image struct {
	ID          int32     `json:"id"`
	Path        string    `json:"path"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Credit      string    `json:"credit"`
	SourceURL   string    `json:"src_url"`
	IsUploaded  bool      `json:"is_uploaded"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ImageType struct {
	Name       string   `json:"name"`
	Mime       string   `json:"mime"`
	Extensions []string `json:"extensions"`
}

type SiteDatum struct {
	ID        int32           `json:"id"`
	Key       string          `json:"key"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
