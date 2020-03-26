// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Querier interface {
	AppendRoleToDomain(ctx context.Context, arg AppendRoleToDomainParams) (DomainRole, error)
	GetArticle(ctx context.Context, arcID sql.NullString) (Article, error)
	GetArticleByDBID(ctx context.Context, id int32) (Article, error)
	GetRolesForDomain(ctx context.Context, domain string) ([]string, error)
	ListAllArticles(ctx context.Context) ([]Article, error)
	ListAvailableArticles(ctx context.Context) ([]Article, error)
	ListDomainsWithRole(ctx context.Context, role string) ([]string, error)
	ListSpotlightPAArticles(ctx context.Context) ([]ListSpotlightPAArticlesRow, error)
	ListUpcoming(ctx context.Context) ([]Article, error)
	PopScheduled(ctx context.Context) ([]Article, error)
	SetRolesForDomain(ctx context.Context, arg SetRolesForDomainParams) (DomainRole, error)
	UpdateAlmanackArticle(ctx context.Context, arg UpdateAlmanackArticleParams) (Article, error)
	UpdateArcArticles(ctx context.Context, arcItems json.RawMessage) error
	UpdateSpotlightPAArticle(ctx context.Context, arg UpdateSpotlightPAArticleParams) (Article, error)
}

var _ Querier = (*Queries)(nil)
