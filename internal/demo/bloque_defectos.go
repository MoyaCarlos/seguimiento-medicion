package demo

import (
	"context"
	"database/sql"
)

// bloqueDefectos todavía no crea nada.
func bloqueDefectos(_ context.Context, _ *sql.Tx, _ *carga) error { return nil }
