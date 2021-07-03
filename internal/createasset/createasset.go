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

/* Createasset means minting assets on Algorand
The implementations here are going to enable
the creation/minting of assets on Algorand
and on Exitor simultaneously, with the details
of the asset creation transaction being stored
both on a local database and on the public
Algorand Blockchain */

const (
	// The database table for created assets
	CreatedAssetTableName = "CreatedAsset"
)

var (
	// ErrNotFound abstracts the postgres not found error
	ErrNotFound = errors.New("Entity not found")

	// ErrForbidden occurs when a user tries to do sth that is 
	// forbidden to them according to Exitor's access control
	// policies
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// CanReadAsset determines if claims has the authority to access the specified asset by id.
func (repo *Repository) CanReadCreatedAsset(ctx context.Context, claims auth.Claims, id string) error {

		// if the request has claims from a specific created asset, ensure
		// that the claim has the correct access to the asset
		if claims.Audience != "" {
			// select id from CreatedAsset where account_id = [accountID]
			query := sqlbuilder.NewSelectBuilder().Select("id").From(CreatedAssetTableName)
			query.Where(query.And(
					query.Equal("account_id", claims.Audience),
					query.Equal("ID",id),
					//query.Equal("assetCreator") OR query.Equal("assetManager")
			))

			queryStr, args := query.Build()
			queryStr = repo.DbConn.Rebind(queryStr)
			var id string 
			err := repo.DbConn.QueryRowContext(ctx, queryStr, args...).Scan(&id)
			if err != nil && err != sql.ErrNoRows {
				err = errors.Wrapf(err, "query - %s", query.String())
				return err
			}

			// When there is no id returned, then the current claim user
			// does not have access to the specified created asset
			if id == "" {
				return errors.WithStack(ErrForbidden)
			} 
		}

		return nil
}

func (repo *Repository) CanModifyCreatedAsset(ctx context.Context, claims auth.Claims, id string) error {
	    err := repo.CanReadCreatedAsset(ctx, claims, id)
		if err != nil {
			return err
		}

		// Admin user can update an asset they have access to
		if !claims.HasRole(auth.RoleAdmin) {
			return errors.WithStack(ErrForbidden)
		}

		return nil
}
// Clawback function will be implemented when 
// an asset already recorded on the chain is modified


// applyClaimsSelect applies a sub-query to the provided query to enforce ACL based on the
// claims provided.
// 1. No claims, request is internal, no ACL applied
// 2. All role types can access their user ID
func applyClaimsSelect(ctx context.Context, claims auth.Claims, query * sqlbuilder.SelectBuilder) error {
	// if claims are empty, don't apply any ACL
	if claims.Audience == "" {
		return nil
	}

	query.Where(query.Equal("account_id", claims.Audience))
	return nil
}


// TODO Function that claws back and asset created
// The function must call the required parameters as per Algorand's
// asset parameters


// createdassetsMapColumns is the list of columns needed for find
var createdassetsMapColumns = "id,account_id,assetname,walletaddress,status,created_at,updated_at,archived_at"
// Will change above to fit the needed params

func selectQuery() *sqlbuilder.SelectBuilder {
	query := sqlbuilder.NewSelectBuilder()
	query.Select(createdassetsMapColumns)
	query.From(CreatedAssetTableName)
	return query
}


// findRequestQuery generates the select query for the given find request
// TODO: Need to figure out why we cannot parse the args when appending the where to
// the query
func findRequestQuery(req CreatedAssetFindRequest) (*sqlbuilder.SelectBuilder, []interface{}) {
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

// Find() gets all the createdassets from the database based
// on the request params
func (repo *Repository) Find(ctx context.Context, claims auth.Claims, req CreatedAssetFindRequest) (CreatedAsset, error) {
	query, args := findRequestQuery(req)
	return find(ctx, claims, repo.DbConn, query, args, req.IncludeArchived)
}


// this is find, an internal method for getting all the created assets from the 
// database using a select query
func find(ctx context.Context, claims auth.Claims, dbConn *sqlx.DB, query *sqlbuilder.SelectBuilder, args []interface{}, includedArchived bool) (CreatedAsset, error) {
		span, ctx := tracer.StartSpanFromContext(ctx, "internal.createdasset.Find")
		defer span.Finish()

		query.Select(createdassetsMapColumns)
		query.From(CreatedAssetTableName)
		if !includedArchived {
				query.Where(query.IsNull("archived_at"))
		}

		// Check to see if a sub query needs to be applied for the claims
		err := applyClaimsSeelect(ctx, claims, query)
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
				err = errors.WithMessage(err, "find created assets failed")
				return nil, err
		}
		defer rows.Close()


		// Iterate over each row
		resp := []*CreatedAsset{}
		for rows.Next() {
				var (
						m CreatedAsset
						err error
				)
				err = rows.Scan(&m.ID, &m.AccountID, &m.Name, &m.Status, &m.CreatedAt,&m.UpdatedAt, &m.ArchivedAt)
				if err != nil {
						err = errors.Wrapf(err, "query - %s", query.String())
						return nil, err
				}

				resp = append(resp, &m)
		}

		err = rows.Err()
		if err != nil {
				err = errors.Wrapf(err, "query - %s", query.String())
				err = errors.WithMessage(err, "find created assets failed")
				return nil, err
		}

		return resp, nil
}


// ReadByID gets the specified created asset by ID from the database
func (repo *Repository) ReadByID(ctx context.Context, claims auth.Claims, id string) (*CreatedAsset, error) {
	return repo.Read(ctx, claims, CreatedAssetReadRequest{
		ID:              id,
		IncludeArchived: false,
	})
}

// Read gets the specified created assets from the database
func (repo *Repository) Read(ctx context.Context, claims auth.Claims, req CreatedAssetReadRequest) (*CreatedAsset, error) {
		span, ctx := tracer.StartSpanFromContext(ctx, "internal.createdasset.Read")
		defer span.Finish()

		// Validate the request
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
				err = errors.WithMessagef(ErrNotFound, "created asset %s not found", req.ID)
				return nil, err
		}

		u := res[0]
		return u, nil
}

// Create inserts a new created asset into the database
func (repo *Repository) Create(ctx context.Context, claims auth.Claims, req CreatedAssetCreateRequest, now time.Time) (*CreatedAsset, error) {
		span, ctx := tracer.StartSpanFromContext(ctx, "internal.createdasset.Create")
		defer span.Finish()
		if claims.Audience != "" {
			// Admin users can update created assets they have access to
			if !claims.HasRole(auth.RoleAdmin) {
					return nil, errors.WithStack(ErrForbidden)
			}

			if req.AccountID != "" {
				// Request accountId must match claims
				if req.AccountID != claims.Audience {
						return nil, errors.WithStack(ErrForbidden)
				}
		
			} else {
					// Set the accountId from claims
					req.AccountID = claims.Audience
			}
		}

		// Validate the request
		v := webcontext.Validator()
		err := v.Struct(req)
		if err != nil {
				return nil, err
		}

		// If now empty set it to the current time
		if now.IsZero() {
				now = time.Now()
		}

		// Always store the time as UTC
		now = now.UTC()
		// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)
	m := CreatedAsset{
		ID:        uuid.NewRandom().String(),
		AccountID: req.AccountID,
		WalletAddress: req.WalletAddress,
		Total:		req.Total,
		AssetName:      req.AssetName,
		Decimals: 	req.Decimals,
		DefaultFrozen: req.DefaultFrozen,
		URL:			req.URL,
		Status:    CreatedAssetStatus_Active,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if req.Status != nil {
			m.Status = *req.Status
	}

	// Build the insert SQL statement.
	query := sqlbuilder.NewInsertBuilder()
	query.InsertInto(CreatedAssetTableName)
	query.Cols(
		"id",
		"account_id",
		"algorand_wallet_address",
		"total_assetIssuance",
		"assetName",
		"assetDecimalsDenomination",
		"defaultAssetsFrozen",
		"assetUrl",
		"status",
		"created_at",
		"updated_at",
		"archived_at",
	)

	query.Values(
		m.ID,
		m.AccountID,
		m.WalletAddress,
		m.Total,
		m.AssetName,
		m.Decimals,
		m.DefaultFrozen,
		m.URL,
		m.Status,
		m.CreatedAt,
		m.UpdatedAt,
		m.ArchivedAt,
	)

	// Execute the query with the provided context
	sql, args := query.Build()
	sql = repo.DbConn.Rebind(sql)
	_, err = repo.DbConn.ExecContext(ctx, sql, args...)
	if err != nil {
			err = errors.Wrapf(err, "query - %s", query.String())
			err = errors.WithMessage(err, "create asset failed")
			return nil, err
	}

	return &m, nil
}

// Update replaces a created asset in the database
func (repo *Repository) Update(ctx context.Context, claims auth.Claims, req CreatedAssetUpdateRequest, now time.Time) error {
		span, ctx := tracer.StartSpanFromContext(ctx, "internal.createdasset.Update")
		defer span.Finish()

		// Validate the request.
		v := webcontext.Validator()
		err := v.Struct(req)
		if err != nil {
				return err
		}

		// Ensure the claims can modify the created asset specified in the request.
		err = repo.CanModifyCreatedAsset((ctx, claims, req.ID))
		if err != nil {
				return err
		}

		// If now empty set it to the current time
		if now.IsZero() {
				now = time.Now()
		}

		// Always store the time as UTC
		now = now.UTC()
		// Postgres truncates times to milliseconds when storing. We and do the same
		// here so the value we return is consistent with what we store.
		now = now.Truncate(time.Millisecond)
		// Build the update SQL statement.
		query := sqlbuilder.NewUpdateBuilder()
		query.Update(createdCreatedAssetTableName)
		var fields []string
		if req.Name != nil {
			fields = append(fields, query.Assign("Asset Name", req.AssetName))
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
			err = errors.WithMessagef(err, "update created asset %s failed", req.ID)
			return err
		}

		return nil
	}

}