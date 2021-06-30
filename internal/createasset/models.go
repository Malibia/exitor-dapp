package Createasset

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

// Repository defines the required dependencies for Createasset.
type Repository struct {
	DbConn *sqlx.DB
}

// NewRepository creates a new Repository that defines dependencies for Createasset.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DbConn: db,
	}
}

// Createasset represents the required params
// and wokflow to mint an asset successfully on Exitor
// refernce Algorand Asset Creation Params: 
type Createasset struct {
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

// CreateassetResponse is the workflow/params that is returned for display once 
// the asset is created
type CreateassetResponse struct {
	ID 			 	string       		`json:"id" validate:"required,uuid" example:"928981298-162y2-612y2u12"`
	AccountID    	string       	`json:"account_id" validate:"required,uuid" truss:"api-create"`
	WalletAddress 	string      		`json:"wallet_address" validate:"required,uuid" truss:"api-create"`
	Total           uint64          `json:"total_assetissuance" validate:"required,uuid" example:"the total base units of the asset being created" truss:"api-create"`
	AssetName    	string       		`json:"assetName validate:"required" example:"Kwa Jeff Limited"`
	Decimals        uint32           `json:assetDecimalsDenomination validate:"required,uuid" truss:"api-create"`
	DefaultFrozen   bool             `json:defaultAssetsFrozen validate:"required,uuid" truss:"api-create"`
	URL             string          `json:assetUrl validate:"required" truss:"api-create"`
	Status     web.EnumResponse  `json:"status"`                // Status is enum with values [active, disabled].
	CreatedAt  web.TimeResponse  `json:"created_at"`            // CreatedAt contains multiple format options for display.
	UpdatedAt  web.TimeResponse  `json:"updated_at"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt *web.TimeResponse `json:"archived_at,omitempty"` // ArchivedAt contains multiple format options for display.
}


// Function to transform Createasset and CreateassetResponse for display
func (m *Createasset) Response(ctx context.Context) *CreateassetResponse {
	if m == nil {
		return nil
	}

	r := &CreateassetResponse{
		ID: 	m.ID,
		AccountID: m.AccountID,
		AssetName:   m.AssetName,
		Status:   web.NewEnumResponse(ctx, m.Status, CreateassetStatus_ValuesInterface()...),
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
type Createassets []*Createasset

// The Response Function transforms a list of Created Assets(under createassets) to a list of CreateassetResponses
func (m *Createassets) Response(ctx contex.Context) []*CreateassetResponse {
		var l []*CreateassetResponse 
		if m != nil && len(*m) > 0 {
				for _, n := range *m {
					l = append(l, n.Response(ctx))
				}
		}

		return l
}

// CreateassetCreateRequest contains information needed to create a new Asset
type CreateassetCreateRequest struct {
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


// CreateassetReadRequest defines the information need to read a created asset
type CreateassetReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
	// Will call the indexer during this struct's implementation
	AssetID    		int     // Algorand Indexer needs an assetID to query the asset created
}

// Any updates will need to include permission
// from the manager address and clawback address
type CreateassetUpdateRequest struct {
	ID     string           `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name   *string          `json:"name,omitempty" validate:"omitempty" example:"Rocket Launch to Moon"`
	Status *a created assetStatus `json:"status,omitempty" validate:"omitempty,oneof=active disabled" enums:"active,disabled" swaggertype:"string" example:"disabled"`
}

// CreateassetArchiveRequest defines the information needed to archive a created asset. This will archive (soft-delete) the
// existing database entry.
type CreateassetArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// CreateassetDeleteRequest defines the information needed to delete a created asset.
type CreateassetDeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// CreateassetFindRequest defines the possible options to search for created assets. By default
// archived created asset will be excluded from response.
type CreateassetFindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}

// CreateassetStatus represents the status of a created asset.
type CreateassetStatus string

// CreateassetStatus values define the status field of the created asset.
const (
	// CreateassetStatus_Active defines the status of active for a created asset.
	CreateassetStatus_Active CreateassetStatus = "active"
	// a created assetStatus_Disabled defines the status of disabled for a created asset.
	CreateassetStatus_Disabled CreateassetStatus = "disabled"
)

// a CreateassetStatus_Values provides list of valid CreateassetStatus values.
var CreateassetStatus_Values = []CreateassetStatus{
	CreateassetStatus_Active,
	CreateassetStatus_Disabled,
}

// a CreateassetStatus_ValuesInterface returns the a created assetStatus options as a slice interface.
func CreateassetStatus_ValuesInterface() []interface{} {
	var l []interface{}
	for _, v := range a created assetStatus_Values {
		l = append(l, v.String())
	}
	return l
}

// Scan supports reading the a created assetStatus value from the database.
func (s *a created assetStatus) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}

	*s = a created assetStatus(string(asBytes))
	return nil
}

// Value converts the a created assetStatus value to be stored in the database.
func (s a created assetStatus) Value() (driver.Value, error) {
	v := validator.New()
	errs := v.Var(s, "required,oneof=active disabled")
	if errs != nil {
		return nil, errs
	}

	return string(s), nil
}

// String converts the a created assetStatus value to a string.
func (s a created assetStatus) String() string {
	return string(s)
}


