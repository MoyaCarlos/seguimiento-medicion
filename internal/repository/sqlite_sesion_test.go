package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func registrarUsuarioConSesiones(t *testing.T, repo *SQLiteCredencialRepository, nombre string) domain.User {
	t.Helper()
	u, err := repo.Registrar(context.Background(), nuevoUsuario(t, nombre), "hash")
	if err != nil {
		t.Fatalf("no se esperaba error al registrar: %v", err)
	}
	return u
}

func sesionDe(usuarioID int64, creadaEn time.Time) domain.Sesion {
	return domain.Sesion{UsuarioID: usuarioID, CreadoEn: creadaEn, VenceEn: creadaEn.Add(domain.DuracionSesion)}
}

func TestSQLiteSesionRepository_CrearYObtenerPorHuellaDevuelveLaSesionConUsuario(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	usuario := registrarUsuarioConSesiones(t, NewSQLiteCredencialRepository(db), "ana.perez")
	repo := NewSQLiteSesionRepository(db)
	creada := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)

	guardada, err := repo.Crear(ctx, "huella-1", sesionDe(usuario.ID, creada))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if guardada.ID == 0 {
		t.Errorf("se esperaba ID asignado: %+v", guardada)
	}

	obtenida, err := repo.ObtenerPorHuella(ctx, "huella-1")
	if err != nil {
		t.Fatalf("no se esperaba error al buscar: %v", err)
	}
	if obtenida.ID != guardada.ID || obtenida.UsuarioID != usuario.ID {
		t.Errorf("sesión inesperada: %+v", obtenida)
	}
	if obtenida.Usuario.ID != usuario.ID || obtenida.Usuario.Nombre != "ana.perez" {
		t.Errorf("usuario inesperado: %+v", obtenida.Usuario)
	}
	if !obtenida.CreadoEn.Equal(creada) || !obtenida.VenceEn.Equal(creada.Add(domain.DuracionSesion)) {
		t.Errorf("instantes inesperados: %+v", obtenida.Sesion)
	}
	if obtenida.CerradaEn != nil {
		t.Errorf("la sesión nueva no debería estar cerrada: %v", obtenida.CerradaEn)
	}
}

func TestSQLiteSesionRepository_HuellaInexistenteEsNoEncontrado(t *testing.T) {
	repo := NewSQLiteSesionRepository(nuevaBDMigrada(t))

	if _, err := repo.ObtenerPorHuella(context.Background(), "no-existe"); !errors.Is(err, domain.ErrNoEncontrado) {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestSQLiteSesionRepository_HuellaRepetidaEsError(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	usuario := registrarUsuarioConSesiones(t, NewSQLiteCredencialRepository(db), "ana.perez")
	repo := NewSQLiteSesionRepository(db)
	creada := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)

	if _, err := repo.Crear(ctx, "huella-1", sesionDe(usuario.ID, creada)); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	_, err := repo.Crear(ctx, "huella-1", sesionDe(usuario.ID, creada))
	if err == nil {
		t.Fatal("se esperaba error de unicidad con la huella repetida")
	}
	if errors.Is(err, ErrIntegridadReferencial) {
		t.Errorf("la huella repetida no es una violación de llave foránea: %v", err)
	}
}

func TestSQLiteSesionRepository_UsuarioInexistenteEsIntegridadReferencial(t *testing.T) {
	repo := NewSQLiteSesionRepository(nuevaBDMigrada(t))
	creada := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)

	_, err := repo.Crear(context.Background(), "huella-1", sesionDe(999, creada))
	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Fatalf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
}

func TestSQLiteSesionRepository_CerrarFijaCerradaEnYNoBorraLaFila(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	usuario := registrarUsuarioConSesiones(t, NewSQLiteCredencialRepository(db), "ana.perez")
	repo := NewSQLiteSesionRepository(db)
	creada := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	cierre := creada.Add(time.Hour)

	guardada, err := repo.Crear(ctx, "huella-1", sesionDe(usuario.ID, creada))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if err := repo.Cerrar(ctx, guardada.ID, cierre); err != nil {
		t.Fatalf("no se esperaba error al cerrar: %v", err)
	}

	obtenida, err := repo.ObtenerPorHuella(ctx, "huella-1")
	if err != nil {
		t.Fatalf("la fila debería seguir existiendo: %v", err)
	}
	if obtenida.CerradaEn == nil || !obtenida.CerradaEn.Equal(cierre) {
		t.Errorf("se esperaba CerradaEn = %v, se obtuvo %v", cierre, obtenida.CerradaEn)
	}
}

func TestSQLiteSesionRepository_CerrarUnaSesionNoTocaOtraDelMismoUsuario(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	usuario := registrarUsuarioConSesiones(t, NewSQLiteCredencialRepository(db), "ana.perez")
	repo := NewSQLiteSesionRepository(db)
	creada := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)

	primera, err := repo.Crear(ctx, "huella-1", sesionDe(usuario.ID, creada))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if _, err := repo.Crear(ctx, "huella-2", sesionDe(usuario.ID, creada)); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if err := repo.Cerrar(ctx, primera.ID, creada.Add(time.Hour)); err != nil {
		t.Fatalf("no se esperaba error al cerrar: %v", err)
	}

	otra, err := repo.ObtenerPorHuella(ctx, "huella-2")
	if err != nil {
		t.Fatalf("no se esperaba error al buscar: %v", err)
	}
	if otra.CerradaEn != nil {
		t.Errorf("la otra sesión no debería estar cerrada: %v", otra.CerradaEn)
	}
}
