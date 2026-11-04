package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

// Crear persiste la sesión con la huella de su token y devuelve una copia con
// el ID asignado. Un usuario inexistente devuelve ErrIntegridadReferencial.
func (r *SQLiteSesionRepository) Crear(ctx context.Context, huella string, s domain.Sesion) (domain.Sesion, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO sesiones (usuario_id, token_hash, creado_en, vence_en) VALUES (?, ?, ?, ?)`,
		s.UsuarioID, huella, s.CreadoEn.UTC().Format(time.RFC3339), s.VenceEn.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return domain.Sesion{}, fmt.Errorf("crear sesión: %w", errorDeEscritura(err))
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Sesion{}, fmt.Errorf("obtener id de la sesión: %w", err)
	}
	s.ID = id
	return s, nil
}

// ObtenerPorHuella devuelve la sesión con su usuario, también si está vencida o cerrada.
func (r *SQLiteSesionRepository) ObtenerPorHuella(ctx context.Context, huella string) (domain.SesionConUsuario, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT s.id, s.usuario_id, s.creado_en, s.vence_en, s.cerrada_en, u.nombre, u.nombre_normalizado, u.creado_en
		FROM sesiones s JOIN usuarios u ON u.id = s.usuario_id
		WHERE s.token_hash = ?`,
		huella,
	)
	var (
		sc                               domain.SesionConUsuario
		creadaEn, venceEn, usuarioCreado string
		cerradaEn                        sql.NullString
	)
	err := row.Scan(&sc.ID, &sc.UsuarioID, &creadaEn, &venceEn, &cerradaEn,
		&sc.Usuario.Nombre, &sc.Usuario.NombreNormalizado, &usuarioCreado)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.SesionConUsuario{}, domain.ErrNoEncontrado
		}
		return domain.SesionConUsuario{}, fmt.Errorf("buscar sesión: %w", err)
	}
	sc.CreadoEn = parsearTiempo(creadaEn)
	sc.VenceEn = parsearTiempo(venceEn)
	if cerradaEn.Valid {
		cierre := parsearTiempo(cerradaEn.String)
		sc.CerradaEn = &cierre
	}
	sc.Usuario.ID = sc.UsuarioID
	sc.Usuario.CreadoEn = parsearTiempo(usuarioCreado)
	return sc, nil
}

// Cerrar fija cerrada_en; la fila no se borra.
func (r *SQLiteSesionRepository) Cerrar(ctx context.Context, sesionID int64, ahora time.Time) error {
	if _, err := r.db.ExecContext(ctx,
		`UPDATE sesiones SET cerrada_en = ? WHERE id = ?`,
		ahora.UTC().Format(time.RFC3339), sesionID,
	); err != nil {
		return fmt.Errorf("cerrar sesión: %w", err)
	}
	return nil
}
