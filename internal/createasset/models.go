package createasset

import (
	"context"
	"time"

	"database/sql/driver"
	"exitor-dapp/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

// Repository defines the required dependencies for CreatedAsset.
type Repository struct {
	DbConn *sqlx.DB
}

// NewRepository creates a new Repository that defines dependencies for CreatedAsset.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DbConn: db,
	}
}

// CreatedAsset represents the required params
// and wokflow to mint an asset successfully on Exitor
// refernce Algorand Asset Creation Params: 
type CreatedAsset struct {
		ID 			 	string       		`json:"id" validate:"required,uuid" example:"928981298-162y2-612y2u12"`
		AccountID    	string       	`json:"account_id" validate:"required,uuid" truss:"api-create"`
		WalletAddress 	string      		`json:"wallet_address" validate:"required,uuid" truss:"api-create"`
		Total           uint64          `json:"total_assetissuance" validate:"required,uuid" example:"the total base units of the asset being created" truss:"api-create"`
		AssetName    	string       		`json:"assetName validate:"required" example:"Kwa Jeff Limited"`
		Decimals        uint32           `json:assetDecimalsDenomination validate:"required,uuid" truss:"api-create"`
		DefaultFrozen   bool             `json:defaultAssetsFrozen validate:"required,uuid" truss:"api-create"`
		URL             string          `json:assetUrl validate:"required" truss:"api-create"`
		Status       AssetCreationStatus  `json:"status" validate:"omitempty,oneof=active disabled"  enums:"active, disabled" swaggertype:"string" example:"active"`
		CreatedAt  time.Time       `json:"created_at" truss:"api-read"`
		UpdatedAt  time.Time       `json:"updated_at" truss:"api-read"`
		ArchivedAt *pq.NullTime    `json:"archived_at,omitempty" truss:"api-hide"`
}

// CreatedAssetResponse is the workflow/params that is returned for display once 
// the asset is created
type CreatedAssetResponse struct {
	ID 			 	string       		`json:"id" validate:"required,uuid" example:"928981298-162y2-612y2u12"`
	AccountID    	string       	`json:"account_id" validate:"required,uuid" truss:"api-create"`
	WalletAddress 	string      		`json:"wallet_address" validate:"required,uuid" truss:"api-create"`
	Total           uint64          `json:"total_assetissuance" validate:"required,uuid" example:"the total base units of the asset being created" truss:"api-create"`
	AssetName    	string       		`json:"assetName validate:"required" example:"Kwa Jeff Limited"`
	Decimals        uint32           `json:assetDecimalsDenomination validate:"required,uuid" truss:"api-create"`
	DefaultFrozen   bool             `json:defaultAssetsFrozen validate:"required,uuid" truss:"api-create"`
	AssetURL             string          `json:assetUrl validate:"required" truss:"api-create"`
	AssetManagerAddress     string // Address for Exitor Org 
	Status     web.EnumResponse  `json:"status"`                // Status is enum with values [active, disabled].
	CreatedAt  web.TimeResponse  `json:"created_at"`            // CreatedAt contains multiple format options for display.
	UpdatedAt  web.TimeResponse  `json:"updated_at"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt *web.TimeResponse `json:"archived_at,omitempty"` // ArchivedAt contains multiple format options for display.
}


// Function to transform CreatedAsset and CreatedAssetResponse for display
func (m *CreatedAsset) Response(ctx context.Context) *CreatedAssetResponse {
	if m == nil {
		return nil
	}

	r := &CreatedAssetResponse{
		ID: 	m.ID,
		AccountID: m.AccountID,
		AssetName:   m.AssetName,
		Status:   web.NewEnumResponse(ctx, m.Status, CreatedAssetStatus_ValuesInterface()...),
		CreatedAt web.NewTimeResponse(ctx, m.CreatedAt),
		UpdatedAt web.NewTimeResponse(ctx, m.UpdatedAt),
	}

	if m.ArchivedAt != nil && !m.ArchiveAt.Time.IsZero() {
			at := web.NeWTimeResponse(ctx, m.ArchiveAt.Time)
			r.ArchiveAt = &at
	}

	return r
}

// a list of created assets
type CreatedAssets []*CreatedAsset

// The Response Function transforms a list of Created Assets(under CreatedAssets) to a list of CreatedAssetResponses
func (m *CreatedAssets) Response(ctx contex.Context) []*CreatedAssetResponse {
		var l []*CreatedAssetResponse 
		if m != nil && len(*m) > 0 {
				for _, n := range *m {
					l = append(l, n.Response(ctx))
				}
		}

		return l
}

// CreatedAssetCreateRequest contains information needed to create a new Asset
type CreatedAssetCreateRequest struct {
		ID 			 	string       		`json:"id" validate:"required,uuid" example:"928981298-162y2-612y2u12"`
		AccountID    	string       	`json:"account_id" validate:"required,uuid" truss:"api-create"`
		WalletAddress 	string      		`json:"algorand_wallet_address" validate:"required,uuid" truss:"api-create"`
		Total           uint64          `json:"total_assetIssuance" validate:"required,uuid" example:"the total base units of the asset being created" truss:"api-create"`
		AssetName    	string       		`json:"assetName validate:"required" example:"Kwa Jeff Limited"`
		Decimals        uint32           `json:assetDecimalsDenomination validate:"required,uuid" truss:"api-create"`
		DefaultFrozen   bool             `json:defaultAssetsFrozen validate:"required,uuid" truss:"api-create"`
		URL             url          `json:assetUrl validate:"required" truss:"api-create"`
		// converted url from string to url
		Status       AssetCreationStatus  `json:"status" validate:"omitempty,oneof=active disabled"  enums:"active, disabled" swaggertype:"string" example:"active"`
		CreatedAt  time.Time       `json:"created_at" truss:"api-read"`
		UpdatedAt  time.Time       `json:"updated_at" truss:"api-read"`
		ArchivedAt *pq.NullTime    `json:"archived_at,omitempty" truss:"api-hide"`
}


// CreatedAssetReadRequest defines the information need to read a created asset
type CreatedAssetReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
	// Will call the indexer during this struct's implementation
	AssetID    		int     // Algorand Indexer needs an assetID to query the asset created
}

// Any updates will need to include permission
// from the manager address and clawback address
type CreatedAssetUpdateRequest struct {
	ID     string           `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name   *string          `json:"name,omitempty" validate:"omitempty" example:"Rocket Launch to Moon"`
	Status *a created assetStatus `json:"status,omitempty" validate:"omitempty,oneof=active disabled" enums:"active,disabled" swaggertype:"string" example:"disabled"`
}

// CreatedAssetArchiveRequest defines the information needed to archive a created asset. This will archive (soft-delete) the
// existing database entry.
type CreatedAssetArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// CreatedAssetDeleteRequest defines the information needed to delete a created asset.
type CreatedAssetDeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// CreatedAssetFindRequest defines the possible options to search for created assets. By default
// archived created asset will be excluded from response.
type CreatedAssetFindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}

// CreatedAssetStatus represents the status of a created asset.
type CreatedAssetStatus string

// CreatedAssetStatus values define the status field of the created asset.
const (
	// CreatedAssetStatus_Active defines the status of active for a created asset.
	CreatedAssetStatus_Active CreatedAssetStatus = "active"
	// a created assetStatus_Disabled defines the status of disabled for a created asset.
	CreatedAssetStatus_Disabled CreatedAssetStatus = "disabled"
)

// a CreatedAssetStatus_Values provides list of valid CreatedAssetStatus values.
var CreatedAssetStatus_Values = []CreatedAssetStatus{
	CreatedAssetStatus_Active,
	CreatedAssetStatus_Disabled,
}

// CreatedAssetStatus_ValuesInterface returns the CreatedAssetStatus options as a slice interface.
func CreatedAssetStatus_ValuesInterface() []interface{} {
	var l []interface{}
	for _, v := range CreatedAssetStatus_Values {
		l = append(l, v.String())
	}
	return l
}

// Scan supports reading the CreatedAssetStatus value from the database.
func (s *CreatedAssetStatus) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}

	*s = CreatedAssetStatus(string(asBytes))
	return nil
}

// Value converts the CreatedAssetStatus value to be stored in the database.
func (s CreatedAssetStatus) Value() (driver.Value, error) {
	v := validator.New()
	errs := v.Var(s, "required,oneof=active disabled")
	if errs != nil {
		return nil, errs
	}

	return string(s), nil
}

// String converts the CreatedAssetStatus value to a string.
func (s CreatedAssetStatus) String() string {
	return string(s)
}


