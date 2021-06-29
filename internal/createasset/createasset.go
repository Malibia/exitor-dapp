package createasset

import (
	"context"
	"database/sql"
	"time"

	"exitor-dapp/internal/platform/auth"
	"exitor-dapp/internal/platform/web/webcontext"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	// The database table for createasset
	createassetTableName = "createassets"
)

var (
	// ErrNotFound abstracts the postgres not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// CanReadcreateasset determines if claims has the authority to access the specified createasset by id.
func (repo *Repository) CanReadcreateasset(ctx context.Context, claims auth.Claims, id string) error {

	// If the request has claims from a specific createasset, ensure that the claims
	// has the correct access to the createasset.
	if claims.Audience != "" {
		// select id from createassets where account_id = [accountID]
		query := sqlbuilder.NewSelectBuilder().Select("id").From(createassetTableName)
		query.Where(query.And(
			query.Equal("account_id", claims.Audience),
			query.Equal("ID", id),
		))

		queryStr, args := query.Build()
		queryStr = repo.DbConn.Rebind(queryStr)
		var id string
		err := repo.DbConn.QueryRowContext(ctx, queryStr, args...).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			err = errors.Wrapf(err, "query - %s", query.String())
			return err
		}

		// When there is no id returned, then the current claim user does not have access
		// to the specified createasset.
		if id == "" {
			return errors.WithStack(ErrForbidden)
		}

	}

	return nil
}

// CanModifycreateasset determines if claims has the authority to modify the specified createasset by id.
func (repo *Repository) CanModifycreateasset(ctx context.Context, claims auth.Claims, id string) error {
	err := repo.CanReadcreateasset(ctx, claims, id)
	if err != nil {
		return err
	}

	// Admin users can update createassets they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	return nil
}

// applyClaimsSelect applies a sub-query to the provided query to enforce ACL based on the claims provided.
// 	1. No claims, request is internal, no ACL applied
// 	2. All role types can access their user ID
func applyClaimsSelect(ctx context.Context, claims auth.Claims, query *sqlbuilder.SelectBuilder) error {
	// Claims are empty, don't apply any ACL
	if claims.Audience == "" {
		return nil
	}

	query.Where(query.Equal("account_id", claims.Audience))
	return nil
}

// createassetMapColumns is the list of columns needed for find.
var createassetMapColumns = "id,account_id,name,status,created_at,updated_at,archived_at"

// selectQuery constructs a base select query for createasset.
func selectQuery() *sqlbuilder.SelectBuilder {
	query := sqlbuilder.NewSelectBuilder()
	query.Select(createassetMapColumns)
	query.From(createassetTableName)
	return query
}

// findRequestQuery generates the select query for the given find request.
// TODO: Need to figure out why can't parse the args when appending the where
// 			to the query.
func findRequestQuery(req createassetFindRequest) (*sqlbuilder.SelectBuilder, []interface{}) {
	query := selectQuery()

	if req.Where != "" {
		query.Where(query.And(req.Where))
	}

	if len(req.Order) > 0 {
		query.OrderBy(req.Order...)
	}

	if req.Limit != nil {
		query.Limit(int(*req.Limit))
	}

	if req.Offset != nil {
		query.Offset(int(*req.Offset))
	}

	return query, req.Args
}

// Find gets all the createassets from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, claims auth.Claims, req createassetFindRequest) (createassets, error) {
	query, args := findRequestQuery(req)
	return find(ctx, claims, repo.DbConn, query, args, req.IncludeArchived)
}

// find internal method for getting all the createassets from the database using a select query.
func find(ctx context.Context, claims auth.Claims, dbConn *sqlx.DB, query *sqlbuilder.SelectBuilder, args []interface{}, includedArchived bool) (createassets, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.createasset.Find")
	defer span.Finish()

	query.Select(createassetMapColumns)
	query.From(createassetTableName)
	if !includedArchived {
		query.Where(query.IsNull("archived_at"))
	}

	// Check to see if a sub query needs to be applied for the claims.
	err := applyClaimsSelect(ctx, claims, query)
	if err != nil {
		return nil, err
	}

	queryStr, queryArgs := query.Build()
	queryStr = dbConn.Rebind(queryStr)
	args = append(args, queryArgs...)
	// Fetch all entries from the db.
	rows, err := dbConn.QueryContext(ctx, queryStr, args...)
	if err != nil {
		err = errors.Wrapf(err, "query - %s", query.String())
		err = errors.WithMessage(err, "find createassets failed")
		return nil, err
	}
	defer rows.Close()

	// Iterate over each row.
	resp := []*createasset{}
	for rows.Next() {
		var (
			m   createasset
			err error
		)
		err = rows.Scan(&m.ID, &m.AccountID, &m.Name, &m.Status, &m.CreatedAt, &m.UpdatedAt, &m.ArchivedAt)
		if err != nil {
			err = errors.Wrapf(err, "query - %s", query.String())
			return nil, err
		}

		resp = append(resp, &m)
	}

	err = rows.Err()
	if err != nil {
		err = errors.Wrapf(err, "query - %s", query.String())
		err = errors.WithMessage(err, "find createassets failed")
		return nil, err
	}

	return resp, nil
}

// ReadByID gets the specified createasset by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, claims auth.Claims, id string) (*createasset, error) {
	return repo.Read(ctx, claims, createassetReadRequest{
		ID:              id,
		IncludeArchived: false,
	})
}

// Read gets the specified createasset from the database.
func (repo *Repository) Read(ctx context.Context, claims auth.Claims, req createassetReadRequest) (*createasset, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.createasset.Read")
	defer span.Finish()

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return nil, err
	}

	// Filter base select query by id
	query := sqlbuilder.NewSelectBuilder()
	query.Where(query.Equal("id", req.ID))

	res, err := find(ctx, claims, repo.DbConn, query, []interface{}{}, req.IncludeArchived)
	if err != nil {
		return nil, err
	} else if res == nil || len(res) == 0 {
		err = errors.WithMessagef(ErrNotFound, "createasset %s not found", req.ID)
		return nil, err
	}

	u := res[0]
	return u, nil
}

// Create inserts a new createasset into the database.
func (repo *Repository) Create(ctx context.Context, claims auth.Claims, req createassetCreateRequest, now time.Time) (*createasset, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.createasset.Create")
	defer span.Finish()
	if claims.Audience != "" {
		// Admin users can update createassets they have access to.
		if !claims.HasRole(auth.RoleAdmin) {
			return nil, errors.WithStack(ErrForbidden)
		}

		if req.AccountID != "" {
			// Request accountId must match claims.
			if req.AccountID != claims.Audience {
				return nil, errors.WithStack(ErrForbidden)
			}

		} else {
			// Set the accountId from claims.
			req.AccountID = claims.Audience
		}

	}

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return nil, err
	}

	// If now empty set it to the current time.
	if now.IsZero() {
		now = time.Now()
	}

	// Always store the time as UTC.
	now = now.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)
	m := createasset{
		ID:        uuid.NewRandom().String(),
		AccountID: req.AccountID,
		Name:      req.Name,
		Status:    createassetStatus_Active,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if req.Status != nil {
		m.Status = *req.Status
	}

	// Build the insert SQL statement.
	query := sqlbuilder.NewInsertBuilder()
	query.InsertInto(createassetTableName)
	query.Cols(
		"id",
		"account_id",
		"name",
		"status",
		"created_at",
		"updated_at",
		"archived_at",
	)

	query.Values(
		m.ID,
		m.AccountID,
		m.Name,
		m.Status,
		m.CreatedAt,
		m.UpdatedAt,
		m.ArchivedAt,
	)

	// Execute the query with the provided context.
	sql, args := query.Build()
	sql = repo.DbConn.Rebind(sql)
	_, err = repo.DbConn.ExecContext(ctx, sql, args...)
	if err != nil {
		err = errors.Wrapf(err, "query - %s", query.String())
		err = errors.WithMessage(err, "create createasset failed")
		return nil, err
	}

	return &m, nil
}

// Update replaces an createasset in the database.
func (repo *Repository) Update(ctx context.Context, claims auth.Claims, req createassetUpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.createasset.Update")
	defer span.Finish()

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	// Ensure the claims can modify the createasset specified in the request.
	err = repo.CanModifycreateasset(ctx, claims, req.ID)
	if err != nil {
		return err
	}

	// If now empty set it to the current time.
	if now.IsZero() {
		now = time.Now()
	}

	// Always store the time as UTC.
	now = now.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)
	// Build the update SQL statement.
	query := sqlbuilder.NewUpdateBuilder()
	query.Update(createassetTableName)
	var fields []string
	if req.Name != nil {
		fields = append(fields, query.Assign("name", req.Name))
	}

	if req.Status != nil {
		fields = append(fields, query.Assign("status", req.Status))
	}

	// If there's nothing to update we can quit early.
	if len(fields) == 0 {
		return nil
	}

	// Append the updated_at field
	fields = append(fields, query.Assign("updated_at", now))
	query.Set(fields...)
	query.Where(query.Equal("id", req.ID))
	// Execute the query with the provided context.
	sql, args := query.Build()
	sql = repo.DbConn.Rebind(sql)
	_, err = repo.DbConn.ExecContext(ctx, sql, args...)
	if err != nil {
		err = errors.Wrapf(err, "query - %s", query.String())
		err = errors.WithMessagef(err, "update createasset %s failed", req.ID)
		return err
	}

	return nil
}

// Archive soft deleted the createasset from the database.
func (repo *Repository) Archive(ctx context.Context, claims auth.Claims, req createassetArchiveRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.createasset.Archive")
	defer span.Finish()

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	// Ensure the claims can modify the createasset specified in the request.
	err = repo.CanModifycreateasset(ctx, claims, req.ID)
	if err != nil {
		return err
	}

	// If now empty set it to the current time.
	if now.IsZero() {
		now = time.Now()
	}

	// Always store the time as UTC.
	now = now.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)
	// Build the update SQL statement.
	query := sqlbuilder.NewUpdateBuilder()
	query.Update(createassetTableName)
	query.Set(
		query.Assign("archived_at", now),
	)

	query.Where(query.Equal("id", req.ID))
	// Execute the query with the provided context.
	sql, args := query.Build()
	sql = repo.DbConn.Rebind(sql)
	_, err = repo.DbConn.ExecContext(ctx, sql, args...)
	if err != nil {
		err = errors.Wrapf(err, "query - %s", query.String())
		err = errors.WithMessagef(err, "archive createasset %s failed", req.ID)
		return err
	}

	return nil
}

// Delete removes an createasset from the database.
func (repo *Repository) Delete(ctx context.Context, claims auth.Claims, req createassetDeleteRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.createasset.Delete")
	defer span.Finish()

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	// Ensure the claims can modify the createasset specified in the request.
	err = repo.CanModifycreateasset(ctx, claims, req.ID)
	if err != nil {
		return err
	}

	// Build the delete SQL statement.
	query := sqlbuilder.NewDeleteBuilder()
	query.DeleteFrom(createassetTableName)
	query.Where(query.Equal("id", req.ID))
	// Execute the query with the provided context.
	sql, args := query.Build()
	sql = repo.DbConn.Rebind(sql)
	_, err = repo.DbConn.ExecContext(ctx, sql, args...)
	if err != nil {
		err = errors.Wrapf(err, "query - %s", query.String())
		err = errors.WithMessagef(err, "delete createasset %s failed", req.ID)
		return err
	}

	return nil
}
