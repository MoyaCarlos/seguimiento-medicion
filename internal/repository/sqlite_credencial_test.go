package repository

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func nuevoUsuario(t *testing.T, nombre string) domain.User {
	t.Helper()
	u, err := domain.NewUser(nombre)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	return u
}

func TestSQLiteCredencialRepository_RegistrarGuardaElHash(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteCredencialRepository(nuevaBDMigrada(t))

	registrado, err := repo.Registrar(ctx, nuevoUsuario(t, "ana.perez"), "hash-de-ana")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if registrado.ID == 0 || registrado.CreadoEn.IsZero() {
		t.Errorf("se esperaban ID y CreadoEn asignados: %+v", registrado)
	}

	u, hash, err := repo.ObtenerPorNombreNormalizado(ctx, "ana.perez")
	if err != nil {
		t.Fatalf("no se esperaba error al buscar: %v", err)
	}
	if u.ID != registrado.ID || u.Nombre != "ana.perez" || hash != "hash-de-ana" {
		t.Errorf("usuario o hash inesperado: %+v, %q", u, hash)
	}
}

func TestSQLiteCredencialRepository_NombreRepetidoEsDuplicado(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteCredencialRepository(nuevaBDMigrada(t))

	if _, err := repo.Registrar(ctx, nuevoUsuario(t, "ana.perez"), "h1"); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if _, err := repo.Registrar(ctx, nuevoUsuario(t, "Ana.Perez"), "h2"); !errors.Is(err, ErrUsuarioDuplicado) {
		t.Fatalf("se esperaba ErrUsuarioDuplicado, se obtuvo %v", err)
	}
}

func TestSQLiteCredencialRepository_DuplicadoContraUsuarioSinContrasena(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	if _, err := NewSQLiteUserRepository(db).Guardar(ctx, nuevoUsuario(t, "luis")); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	repo := NewSQLiteCredencialRepository(db)

	if _, err := repo.Registrar(ctx, nuevoUsuario(t, "Luis"), "h"); !errors.Is(err, ErrUsuarioDuplicado) {
		t.Fatalf("se esperaba ErrUsuarioDuplicado, se obtuvo %v", err)
	}
}

func TestSQLiteCredencialRepository_UsuarioSinContrasenaDevuelveHashVacio(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	previo, err := NewSQLiteUserRepository(db).Guardar(ctx, nuevoUsuario(t, "luis"))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	repo := NewSQLiteCredencialRepository(db)

	u, hash, err := repo.ObtenerPorNombreNormalizado(ctx, "luis")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if u.ID != previo.ID || hash != "" {
		t.Errorf("se esperaba el usuario previo con hash vacío: %+v, %q", u, hash)
	}
}

func TestSQLiteCredencialRepository_InexistenteEsNoEncontrado(t *testing.T) {
	repo := NewSQLiteCredencialRepository(nuevaBDMigrada(t))

	_, _, err := repo.ObtenerPorNombreNormalizado(context.Background(), "nadie")
	if !errors.Is(err, domain.ErrNoEncontrado) {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestSQLiteCredencialRepository_RegistroConcurrenteDelMismoNombre(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteCredencialRepository(nuevaBDMigrada(t))

	const intentos = 2
	u := nuevoUsuario(t, "ana.perez")
	errs := make([]error, intentos)
	var listos, arranque sync.WaitGroup
	listos.Add(intentos)
	arranque.Add(1)
	for i := 0; i < intentos; i++ {
		go func() {
			defer listos.Done()
			arranque.Wait()
			_, errs[i] = repo.Registrar(ctx, u, "hash")
		}()
	}
	arranque.Done()
	listos.Wait()

	var creados, duplicados int
	for _, err := range errs {
		switch {
		case err == nil:
			creados++
		case errors.Is(err, ErrUsuarioDuplicado):
			duplicados++
		default:
			t.Errorf("error inesperado: %v", err)
		}
	}
	if creados != 1 || duplicados != 1 {
		t.Errorf("se esperaba 1 creado y 1 duplicado, se obtuvo %d y %d", creados, duplicados)
	}
}
