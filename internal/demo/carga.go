// Package demo carga el proyecto de demostración (issue #144) en una sola
// transacción, armada con bloques independientes. Cada historia futura que
// agregue datos suma su propio bloque a la lista (ver specs/021-datos-de-prueba).
package demo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var (
	// ErrYaCargado indica que el proyecto de demostración ya existe.
	ErrYaCargado = errors.New("los datos de prueba ya están cargados")
	// ErrIdOcupado indica que el id del proyecto de demostración lo usa otro proyecto.
	ErrIdOcupado = errors.New("el id 25 ya lo usa otro proyecto")
)

// Resumen cuenta lo que creó la carga.
type Resumen struct {
	ProyectoID     int64 `json:"proyecto_id"`
	Usuarios       int   `json:"usuarios"`
	Integrantes    int   `json:"integrantes"`
	Sprints        int   `json:"sprints"`
	Historias      int   `json:"historias"`
	RegistrosHoras int   `json:"registros_horas"`
	Rondas         int   `json:"rondas_votacion"`
	Votos          int   `json:"votos"`
}

// carga es el estado compartido entre los bloques de una misma transacción.
type carga struct {
	ahora     time.Time
	hoy       time.Time
	usuarios  map[string]int64
	sprints   []int64
	historias map[string]int64
	resumen   Resumen
}

type bloque func(ctx context.Context, tx *sql.Tx, c *carga) error

func nuevaCarga(ahora time.Time) *carga {
	ahora = ahora.UTC()
	return &carga{
		ahora:     ahora,
		hoy:       ahora.Truncate(24 * time.Hour),
		usuarios:  map[string]int64{},
		historias: map[string]int64{},
	}
}

// dia devuelve el día relativo a hoy (dia(-14) es hace 14 días).
func (c *carga) dia(desplazamiento int) time.Time {
	return c.hoy.AddDate(0, 0, desplazamiento)
}

func (c *carga) creadoEn() string {
	return c.ahora.Format(time.RFC3339)
}

// ejecutarBloques corre los bloques en orden dentro de una transacción: si uno
// falla se deshace todo.
func ejecutarBloques(ctx context.Context, db *sql.DB, ahora time.Time, bloques []bloque) (*carga, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("iniciar transacción de la carga: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	c := nuevaCarga(ahora)
	for _, ejecutar := range bloques {
		if err := ejecutar(ctx, tx, c); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("confirmar la carga: %w", err)
	}
	return c, nil
}

// bloquesV1 es el orden de la carga. Una historia futura agrega su bloque acá.
var bloquesV1 = []bloque{
	bloqueUsuarios,
	bloqueProyecto,
	bloqueSprints,
	bloqueHistorias,
	bloqueHoras,
	bloqueVotacion,
}

// Cargador carga el proyecto de demostración en una sola transacción.
type Cargador struct {
	db      *sql.DB
	ahora   func() time.Time
	bloques []bloque
}

func NewCargador(db *sql.DB, ahora func() time.Time) *Cargador {
	return &Cargador{db: db, ahora: ahora, bloques: bloquesV1}
}

// Cargar crea todos los datos o ninguno.
func (l *Cargador) Cargar(ctx context.Context) (Resumen, error) {
	c, err := ejecutarBloques(ctx, l.db, l.ahora(), l.bloques)
	if esViolacionDeClave(err) {
		return Resumen{}, l.clasificarConflicto(ctx, err)
	}
	if err != nil {
		return Resumen{}, err
	}
	return c.resumen, nil
}

// clasificarConflicto se llama con la transacción ya deshecha: distingue la
// carga repetida del id ocupado por el nombre del proyecto 25.
func (l *Cargador) clasificarConflicto(ctx context.Context, original error) error {
	var nombre string
	err := l.db.QueryRowContext(ctx, `SELECT nombre FROM proyectos WHERE id = ?`, proyectoID).Scan(&nombre)
	switch {
	case err != nil:
		return original
	case nombre == nombreProyecto:
		return ErrYaCargado
	default:
		return ErrIdOcupado
	}
}

func esViolacionDeClave(err error) bool {
	var sqliteErr *sqlite.Error
	if !errors.As(err, &sqliteErr) {
		return false
	}
	codigo := sqliteErr.Code()
	return codigo == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY || codigo == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}
