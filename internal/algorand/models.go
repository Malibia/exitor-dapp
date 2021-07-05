package algorand

import (
	"fmt"
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