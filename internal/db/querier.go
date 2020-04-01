// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Querier interface {
	AppendRoleToDomain(ctx context.Context, arg AppendRoleToDomainParams) (DomainRole, error)
	CreateImagePlaceholder(ctx context.Context, arg CreateImagePlaceholderParams) (int64, error)
	GetArticle(ctx context.Context, arcID sql.NullString) (Article, error)
	GetArticleByDBID(ctx context.Context, id int32) (Article, error)
	GetImageBySourceURL(ctx context.Context, srcUrl string) (Image, error)
	GetRolesForDomain(ctx context.Context, domain string) ([]string, error)
	ListAllArticles(ctx context.Context) ([]Article, error)
	ListAvailableArticles(ctx context.Context) ([]Article, error)
	ListDomainsWithRole(ctx context.Context, role string) ([]string, error)
	ListImages(ctx context.Context, arg ListImagesParams) ([]Image, error)
	ListSpotlightPAArticles(ctx context.Context) ([]ListSpotlightPAArticlesRow, error)
	ListUpcoming(ctx context.Context) ([]Article, error)
	PopScheduled(ctx context.Context) ([]Article, error)
	SetRolesForDomain(ctx context.Context, arg SetRolesForDomainParams) (DomainRole, error)
	UpdateAlmanackArticle(ctx context.Context, arg UpdateAlmanackArticleParams) (Article, error)
	UpdateArcArticles(ctx context.Context, arcItems json.RawMessage) error
	UpdateImage(ctx context.Context, arg UpdateImageParams) (Image, error)
	UpdateSpotlightPAArticle(ctx context.Context, arg UpdateSpotlightPAArticleParams) (Article, error)
}

var _ Querier = (*Queries)(nil)
