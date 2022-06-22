// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"time"

	"github.com/jackc/pgtype"
	"github.com/spotlightpa/almanack/internal/arc"
)

type AddressRole struct {
	ID           int64     `json:"id"`
	EmailAddress string    `json:"email_address"`
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Article struct {
	ID              int32        `json:"id"`
	ArcID           pgtype.Text  `json:"arc_id"`
	ArcData         arc.FeedItem `json:"arc_data"`
	SpotlightPAPath pgtype.Text  `json:"spotlightpa_path"`
	Note            string       `json:"note"`
	Status          string       `json:"status"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type ArticleStatus struct {
	StatusID    string `json:"status_id"`
	Description string `json:"description"`
}

type DomainRole struct {
	ID        int64     `json:"id"`
	Domain    string    `json:"domain"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type File struct {
	ID          int64     `json:"id"`
	URL         string    `json:"url"`
	Filename    string    `json:"filename"`
	MimeType    string    `json:"mime_type"`
	Description string    `json:"description"`
	IsUploaded  bool      `json:"is_uploaded"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Image struct {
	ID          int64     `json:"id"`
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

type Newsletter struct {
	Subject         string      `json:"subject"`
	ArchiveURL      string      `json:"archive_url"`
	PublishedAt     time.Time   `json:"published_at"`
	Type            string      `json:"type"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	ID              int64       `json:"id"`
	Description     string      `json:"description"`
	Blurb           string      `json:"blurb"`
	SpotlightPAPath pgtype.Text `json:"spotlightpa_path"`
}

type NewsletterType struct {
	Shortname   string `json:"shortname"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Page struct {
	ID            int64              `json:"id"`
	FilePath      string             `json:"file_path"`
	Frontmatter   Map                `json:"frontmatter"`
	Body          string             `json:"body"`
	ScheduleFor   pgtype.Timestamptz `json:"schedule_for"`
	LastPublished pgtype.Timestamptz `json:"last_published"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	URLPath       pgtype.Text        `json:"url_path"`
}

type SiteDatum struct {
	ID          int64              `json:"id"`
	Key         string             `json:"key"`
	Data        Map                `json:"data"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	ScheduleFor time.Time          `json:"schedule_for"`
	PublishedAt pgtype.Timestamptz `json:"published_at"`
}
