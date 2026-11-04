package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// SQLiteSesionRepository implementa CreadorSesiones, BuscadorSesiones y CerradorSesiones sobre SQLite.
type SQLiteSesionRepository struct {
	db *sql.DB
}

func NewSQLiteSesionRepository(db *sql.DB) *SQLiteSesionRepository {
	return &SQLiteSesionRepository{db: db}
}

func (r *SQLiteSesionRepository) Crear(ctx context.Context, huella string, s domain.Sesion) (domain.Sesion, error) {
	return domain.Sesion{}, nil
}

func (r *SQLiteSesionRepository) ObtenerPorHuella(ctx context.Context, huella string) (domain.SesionConUsuario, error) {
	return domain.SesionConUsuario{}, nil
}

func (r *SQLiteSesionRepository) Cerrar(ctx context.Context, sesionID int64, ahora time.Time) error {
	return nil
}
