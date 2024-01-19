package handlers

import (
	"context"
	"database/sql"

	"github.com/markmumba/chasebank/internal/database"
)

type Applicaton struct {
	SDB *sql.DB
	DB  *database.Queries
	Ctx context.Context
}
