package demo

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"sync"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

func baseEnArchivo(t *testing.T) string {
	t.Helper()
	ruta := filepath.Join(t.TempDir(), "demo.db")
	db, err := repository.AbrirSQLite(ruta)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	return ruta
}

func cargarEnParalelo(t *testing.T, dbs []*sql.DB, cargas int) (exitosas, repetidas int) {
	t.Helper()
	var espera sync.WaitGroup
	var cerrojo sync.Mutex
	for i := 0; i < cargas; i++ {
		espera.Add(1)
		go func(db *sql.DB) {
			defer espera.Done()
			_, err := cargador(db).Cargar(context.Background())
			cerrojo.Lock()
			defer cerrojo.Unlock()
			switch {
			case err == nil:
				exitosas++
			case errors.Is(err, ErrYaCargado):
				repetidas++
			default:
				t.Errorf("error inesperado: %v", err)
			}
		}(dbs[i%len(dbs)])
	}
	espera.Wait()
	return exitosas, repetidas
}

func verificarUnaSolaCarga(t *testing.T, db *sql.DB) {
	t.Helper()
	esperado := map[string]int{"usuarios": 4, "proyectos": 1, "integrantes": 3, "sprints": 3, "historias": 10,
		"historias_arrastradas": 1, "registros_horas": 15, "rondas_votacion": 1, "votos": 3}
	for tabla, filas := range fotoDeLaBase(t, db) {
		if filas != esperado[tabla] {
			t.Errorf("%s: %d filas, se esperaban %d", tabla, filas, esperado[tabla])
		}
	}
}

func TestCargar_OchoCargasSimultaneasConUnaConexionDejanUnaSola(t *testing.T) {
	ruta := baseEnArchivo(t)
	db, err := repository.AbrirSQLite(ruta)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	exitosas, repetidas := cargarEnParalelo(t, []*sql.DB{db}, 8)

	if exitosas != 1 || repetidas != 7 {
		t.Errorf("exitosas %d, repetidas %d; se esperaba 1 y 7", exitosas, repetidas)
	}
	verificarUnaSolaCarga(t, db)
}

func TestCargar_DosConexionesSobreElMismoArchivoDejanUnaSola(t *testing.T) {
	ruta := baseEnArchivo(t)
	var dbs []*sql.DB
	for i := 0; i < 2; i++ {
		db, err := repository.AbrirSQLite(ruta)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		dbs = append(dbs, db)
	}

	exitosas, repetidas := cargarEnParalelo(t, dbs, 2)

	if exitosas != 1 || repetidas != 1 {
		t.Errorf("exitosas %d, repetidas %d; se esperaba 1 y 1", exitosas, repetidas)
	}
	verificarUnaSolaCarga(t, dbs[0])
}
