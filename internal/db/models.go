// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	docs "google.golang.org/api/docs/v1"
)

type AddressRole struct {
	ID           int64     `json:"id"`
	EmailAddress string    `json:"email_address"`
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Arc struct {
	ID          int64              `json:"id"`
	ArcID       string             `json:"arc_id"`
	RawData     json.RawMessage    `json:"raw_data"`
	LastUpdated pgtype.Timestamptz `json:"last_updated"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type DomainRole struct {
	ID        int64     `json:"id"`
	Domain    string    `json:"domain"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type File struct {
	ID          int64              `json:"id"`
	URL         string             `json:"url"`
	Filename    string             `json:"filename"`
	MimeType    string             `json:"mime_type"`
	Description string             `json:"description"`
	IsUploaded  bool               `json:"is_uploaded"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	MD5         []byte             `json:"md5"`
	Bytes       int64              `json:"bytes"`
	DeletedAt   pgtype.Timestamptz `json:"deleted_at"`
}

type GDocsDoc struct {
	ID              int64              `json:"id"`
	ExternalID      string             `json:"external_id"`
	Document        docs.Document      `json:"document"`
	Metadata        GDocsMetadata      `json:"metadata"`
	Embeds          []Embed            `json:"embeds"`
	RichText        string             `json:"rich_text"`
	RawHtml         string             `json:"raw_html"`
	ArticleMarkdown string             `json:"article_markdown"`
	WordCount       int32              `json:"word_count"`
	Warnings        []string           `json:"warnings"`
	ProcessedAt     pgtype.Timestamptz `json:"processed_at"`
	CreatedAt       time.Time          `json:"created_at"`
}

type GDocsImage struct {
	ID          int64     `json:"id"`
	ExternalID  string    `json:"external_id"`
	DocObjectID string    `json:"doc_object_id"`
	ImageID     int64     `json:"image_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Image struct {
	ID          int64              `json:"id"`
	Path        string             `json:"path"`
	Type        string             `json:"type"`
	Description string             `json:"description"`
	Credit      string             `json:"credit"`
	SourceURL   string             `json:"src_url"`
	IsUploaded  bool               `json:"is_uploaded"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	MD5         []byte             `json:"md5"`
	Bytes       int64              `json:"bytes"`
	Keywords    string             `json:"keywords"`
	DeletedAt   pgtype.Timestamptz `json:"deleted_at"`
	IsLicensed  bool               `json:"is_licensed"`
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

type Option struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Page struct {
	ID              int64              `json:"id"`
	FilePath        string             `json:"file_path"`
	Frontmatter     Map                `json:"frontmatter"`
	Body            string             `json:"body"`
	ScheduleFor     pgtype.Timestamptz `json:"schedule_for"`
	LastPublished   pgtype.Timestamptz `json:"last_published"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	URLPath         pgtype.Text        `json:"url_path"`
	SourceType      string             `json:"source_type"`
	SourceID        string             `json:"source_id"`
	PublicationDate pgtype.Timestamptz `json:"publication_date"`
}

type SharedArticle struct {
	ID                   int64              `json:"id"`
	Status               string             `json:"status"`
	EmbargoUntil         pgtype.Timestamptz `json:"embargo_until"`
	Note                 string             `json:"note"`
	SourceType           string             `json:"source_type"`
	SourceID             string             `json:"source_id"`
	RawData              json.RawMessage    `json:"raw_data"`
	PageID               pgtype.Int8        `json:"page_id"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
	PublicationDate      pgtype.Timestamptz `json:"publication_date"`
	InternalID           string             `json:"internal_id"`
	Byline               string             `json:"byline"`
	Budget               string             `json:"budget"`
	Hed                  string             `json:"hed"`
	Description          string             `json:"description"`
	LedeImage            string             `json:"lede_image"`
	LedeImageCredit      string             `json:"lede_image_credit"`
	LedeImageDescription string             `json:"lede_image_description"`
	LedeImageCaption     string             `json:"lede_image_caption"`
}

type SharedStatus struct {
	ID          string `json:"id"`
	Description string `json:"description"`
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
