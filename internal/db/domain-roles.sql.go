// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: domain-roles.sql

package db

import (
	"context"
)

const appendRoleToDomain = `-- name: AppendRoleToDomain :one
INSERT INTO domain_roles ("domain", roles)
  VALUES ($1, ARRAY[$2::text])
ON CONFLICT (lower("domain"))
  DO UPDATE SET
    roles = CASE WHEN NOT (domain_roles.roles::text[] @> ARRAY[$2]) THEN
      domain_roles.roles::text[] || ARRAY[$2]
    ELSE
      domain_roles.roles
    END
  RETURNING
    id, domain, roles, created_at, updated_at
`

type AppendRoleToDomainParams struct {
	Domain string `json:"domain"`
	Role   string `json:"role"`
}

func (q *Queries) AppendRoleToDomain(ctx context.Context, arg AppendRoleToDomainParams) (DomainRole, error) {
	row := q.db.QueryRow(ctx, appendRoleToDomain, arg.Domain, arg.Role)
	var i DomainRole
	err := row.Scan(
		&i.ID,
		&i.Domain,
		&i.Roles,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRolesForDomain = `-- name: GetRolesForDomain :one
SELECT
  roles
FROM
  domain_roles
WHERE
  "domain" ILIKE $1
`

func (q *Queries) GetRolesForDomain(ctx context.Context, domain string) ([]string, error) {
	row := q.db.QueryRow(ctx, getRolesForDomain, domain)
	var roles []string
	err := row.Scan(&roles)
	return roles, err
}

const listDomainsWithRole = `-- name: ListDomainsWithRole :many
SELECT
  "domain"
FROM
  "domain_roles"
WHERE
  "roles" @> ARRAY[$1::text]
ORDER BY
  "domain" ASC
`

func (q *Queries) ListDomainsWithRole(ctx context.Context, role string) ([]string, error) {
	rows, err := q.db.Query(ctx, listDomainsWithRole, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			return nil, err
		}
		items = append(items, domain)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setRolesForDomain = `-- name: SetRolesForDomain :one
INSERT INTO domain_roles ("domain", roles)
  VALUES ($1, $2)
ON CONFLICT (lower("domain"))
  DO UPDATE SET
    roles = $2
  RETURNING
    id, domain, roles, created_at, updated_at
`

type SetRolesForDomainParams struct {
	Domain string   `json:"domain"`
	Roles  []string `json:"roles"`
}

func (q *Queries) SetRolesForDomain(ctx context.Context, arg SetRolesForDomainParams) (DomainRole, error) {
	row := q.db.QueryRow(ctx, setRolesForDomain, arg.Domain, arg.Roles)
	var i DomainRole
	err := row.Scan(
		&i.ID,
		&i.Domain,
		&i.Roles,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
