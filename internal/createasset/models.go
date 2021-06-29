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

// Createasset represents a workflow.
type Createasset struct {
	ID         string          `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	AccountID  string          `json:"account_id" validate:"required,uuid" truss:"api-create"`
	Name       string          `json:"name"  validate:"required" example:"Rocket Launch"`
	Status     CreateassetStatus `json:"status" validate:"omitempty,oneof=active disabled" enums:"active,disabled" swaggertype:"string" example:"active"`
	CreatedAt  time.Time       `json:"created_at" truss:"api-read"`
	UpdatedAt  time.Time       `json:"updated_at" truss:"api-read"`
	ArchivedAt *pq.NullTime    `json:"archived_at,omitempty" truss:"api-hide"`
}

// CreateassetResponse represents a workflow that is returned for display.
type CreateassetResponse struct {
	ID         string            `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	AccountID  string            `json:"account_id" validate:"required,uuid" truss:"api-create" example:"c4653bf9-5978-48b7-89c5-95704aebb7e2"`
	Name       string            `json:"name"  validate:"required" example:"Rocket Launch"`
	Status     web.EnumResponse  `json:"status"`                // Status is enum with values [active, disabled].
	CreatedAt  web.TimeResponse  `json:"created_at"`            // CreatedAt contains multiple format options for display.
	UpdatedAt  web.TimeResponse  `json:"updated_at"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt *web.TimeResponse `json:"archived_at,omitempty"` // ArchivedAt contains multiple format options for display.
}

// Response transforms Createasset and CreateassetResponse that is used for display.
// Additional filtering by context values or translations could be applied.
func (m *Createasset) Response(ctx context.Context) *CreateassetResponse {
	if m == nil {
		return nil
	}

	r := &CreateassetResponse{
		ID:        m.ID,
		AccountID: m.AccountID,
		Name:      m.Name,
		Status:    web.NewEnumResponse(ctx, m.Status, CreateassetStatus_ValuesInterface()...),
		CreatedAt: web.NewTimeResponse(ctx, m.CreatedAt),
		UpdatedAt: web.NewTimeResponse(ctx, m.UpdatedAt),
	}

	if m.ArchivedAt != nil && !m.ArchivedAt.Time.IsZero() {
		at := web.NewTimeResponse(ctx, m.ArchivedAt.Time)
		r.ArchivedAt = &at
	}

	return r
}

// Createassets a list of Createassets.
type Createassets []*Createasset

// Response transforms a list of Createassets to a list of CreateassetResponses.
func (m *Createassets) Response(ctx context.Context) []*CreateassetResponse {
	var l []*CreateassetResponse
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// CreateassetCreateRequest contains information needed to create a new Createasset.
type CreateassetCreateRequest struct {
	AccountID string           `json:"account_id" validate:"required,uuid"  example:"c4653bf9-5978-48b7-89c5-95704aebb7e2"`
	Name      string           `json:"name" validate:"required"  example:"Rocket Launch"`
	Status    *CreateassetStatus `json:"status,omitempty" validate:"omitempty,oneof=active disabled" enums:"active,disabled" swaggertype:"string" example:"active"`
}

// CreateassetReadRequest defines the information needed to read a Createasset.
type CreateassetReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// CreateassetUpdateRequest defines what information may be provided to modify an existing
// Createasset. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank.
type CreateassetUpdateRequest struct {
	ID     string           `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name   *string          `json:"name,omitempty" validate:"omitempty" example:"Rocket Launch to Moon"`
	Status *CreateassetStatus `json:"status,omitempty" validate:"omitempty,oneof=active disabled" enums:"active,disabled" swaggertype:"string" example:"disabled"`
}

// CreateassetArchiveRequest defines the information needed to archive a Createasset. This will archive (soft-delete) the
// existing database entry.
type CreateassetArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// CreateassetDeleteRequest defines the information needed to delete a Createasset.
type CreateassetDeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// CreateassetFindRequest defines the possible options to search for Createassets. By default
// archived Createasset will be excluded from response.
type CreateassetFindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}

// CreateassetStatus represents the status of Createasset.
type CreateassetStatus string

// CreateassetStatus values define the status field of Createasset.
const (
	// CreateassetStatus_Active defines the status of active for Createasset.
	CreateassetStatus_Active CreateassetStatus = "active"
	// CreateassetStatus_Disabled defines the status of disabled for Createasset.
	CreateassetStatus_Disabled CreateassetStatus = "disabled"
)

// CreateassetStatus_Values provides list of valid CreateassetStatus values.
var CreateassetStatus_Values = []CreateassetStatus{
	CreateassetStatus_Active,
	CreateassetStatus_Disabled,
}

// CreateassetStatus_ValuesInterface returns the CreateassetStatus options as a slice interface.
func CreateassetStatus_ValuesInterface() []interface{} {
	var l []interface{}
	for _, v := range CreateassetStatus_Values {
		l = append(l, v.String())
	}
	return l
}

// Scan supports reading the CreateassetStatus value from the database.
func (s *CreateassetStatus) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}

	*s = CreateassetStatus(string(asBytes))
	return nil
}

// Value converts the CreateassetStatus value to be stored in the database.
func (s CreateassetStatus) Value() (driver.Value, error) {
	v := validator.New()
	errs := v.Var(s, "required,oneof=active disabled")
	if errs != nil {
		return nil, errs
	}

	return string(s), nil
}

// String converts the CreateassetStatus value to a string.
func (s CreateassetStatus) String() string {
	return string(s)
}
