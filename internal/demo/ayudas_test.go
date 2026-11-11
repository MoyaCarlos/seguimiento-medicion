package demo

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

const contrasenaDemo = "demo1234"

// hasheadorBarato usa una iteración: las pruebas no deben pagar 600 000.
var hasheadorBarato = service.NewServicioContrasenas(1)

type hasheadorContador struct{ llamadas int }

func (h *hasheadorContador) Hashear(contrasena string) (string, error) {
	h.llamadas++
	return hasheadorBarato.Hashear(contrasena)
}

func hashDePrueba(t *testing.T) string {
	t.Helper()
	hash, err := hasheadorBarato.Hashear(contrasenaDemo)
	if err != nil {
		t.Fatalf("hashear: %v", err)
	}
	return hash
}

var ahoraFijo = time.Date(2026, 10, 30, 15, 0, 0, 0, time.UTC)

// ahoraOtroDia solo lo usa la prueba de fechas relativas: ahoraFijo sigue
// siendo la fecha de las demás pruebas.
var ahoraOtroDia = time.Date(2026, 11, 9, 18, 30, 0, 0, time.UTC)

func baseMigrada(t *testing.T) *sql.DB {
	t.Helper()
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		t.Fatalf("abrir base: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		t.Fatalf("migrar: %v", err)
	}
	return db
}

func ejecutar(t *testing.T, db *sql.DB, bloques ...bloque) *carga {
	t.Helper()
	c, err := ejecutarBloques(context.Background(), db, ahoraFijo, hashDePrueba(t), bloques)
	if err != nil {
		t.Fatalf("ejecutar bloques: %v", err)
	}
	return c
}

func contar(t *testing.T, db *sql.DB, consulta string, args ...any) int {
	t.Helper()
	var n int
	if err := db.QueryRow(consulta, args...).Scan(&n); err != nil {
		t.Fatalf("%s: %v", consulta, err)
	}
	return n
}
