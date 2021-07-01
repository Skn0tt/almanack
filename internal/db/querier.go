// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Querier interface {
	AppendRoleToDomain(ctx context.Context, arg AppendRoleToDomainParams) (DomainRole, error)
	CreateFilePlaceholder(ctx context.Context, arg CreateFilePlaceholderParams) (int64, error)
	CreateImage(ctx context.Context, arg CreateImageParams) (int64, error)
	CreateImagePlaceholder(ctx context.Context, arg CreateImagePlaceholderParams) (int64, error)
	EnsurePage(ctx context.Context, path string) error
	GetArticle(ctx context.Context, arcID sql.NullString) (Article, error)
	GetArticleByDBID(ctx context.Context, id int32) (Article, error)
	GetArticleIDFromSlug(ctx context.Context, slug string) (string, error)
	GetImageBySourceURL(ctx context.Context, srcUrl string) (Image, error)
	GetPage(ctx context.Context, path string) (Page, error)
	GetRolesForAddress(ctx context.Context, emailAddress string) ([]string, error)
	GetRolesForDomain(ctx context.Context, domain string) ([]string, error)
	GetSiteData(ctx context.Context, key string) (json.RawMessage, error)
	ListAddressesWithRole(ctx context.Context, role string) ([]string, error)
	ListAllArcArticles(ctx context.Context, arg ListAllArcArticlesParams) ([]Article, error)
	ListAllSeries(ctx context.Context) ([]string, error)
	ListAllTopics(ctx context.Context) ([]string, error)
	ListAvailableArticles(ctx context.Context, arg ListAvailableArticlesParams) ([]Article, error)
	ListDomainsWithRole(ctx context.Context, role string) ([]string, error)
	ListFiles(ctx context.Context, arg ListFilesParams) ([]File, error)
	ListImages(ctx context.Context, arg ListImagesParams) ([]Image, error)
	ListNewsletters(ctx context.Context, arg ListNewslettersParams) ([]Newsletter, error)
	ListPages(ctx context.Context, arg ListPagesParams) ([]ListPagesRow, error)
	ListSpotlightPAArticles(ctx context.Context) ([]ListSpotlightPAArticlesRow, error)
	ListUnpublishedNewsletters(ctx context.Context, arg ListUnpublishedNewslettersParams) ([]Newsletter, error)
	ListUpcoming(ctx context.Context) ([]Article, error)
	PopScheduled(ctx context.Context) ([]Article, error)
	PopScheduledPages(ctx context.Context) ([]Page, error)
	SetRolesForAddress(ctx context.Context, arg SetRolesForAddressParams) (AddressRole, error)
	SetRolesForDomain(ctx context.Context, arg SetRolesForDomainParams) (DomainRole, error)
	SetSiteData(ctx context.Context, arg SetSiteDataParams) error
	UpdateAlmanackArticle(ctx context.Context, arg UpdateAlmanackArticleParams) (Article, error)
	UpdateArcArticles(ctx context.Context, arcItems json.RawMessage) error
	UpdateFile(ctx context.Context, arg UpdateFileParams) (File, error)
	UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error)
	UpdateNewsletter(ctx context.Context, arg UpdateNewsletterParams) (Newsletter, error)
	UpdateNewsletterArchives(ctx context.Context, arg UpdateNewsletterArchivesParams) (int64, error)
	UpdatePage(ctx context.Context, arg UpdatePageParams) (Page, error)
	UpdateSpotlightPAArticle(ctx context.Context, arg UpdateSpotlightPAArticleParams) (sql.NullTime, error)
	UpdateSpotlightPAArticleLastPublished(ctx context.Context, arcID string) (sql.NullTime, error)
}

var _ Querier = (*Queries)(nil)
