package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// Con SQLite real, agregar como Scrum Master a quien ya lo es viola la clave
// primaria y el índice a la vez; el service debe devolver el duplicado (400).
func TestAsignarIntegrante_SQLiteReal_ScrumMasterQueYaEsScrumMasterEsDuplicado(t *testing.T) {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	ctx := context.Background()
	if err := repository.AplicarMigraciones(ctx, db); err != nil {
		t.Fatalf("no se pudieron aplicar las migraciones: %v", err)
	}
	proyectos := repository.NewSQLiteProjectRepository(db)
	usuarios := repository.NewSQLiteUserRepository(db)
	ana, err := usuarios.Guardar(ctx, domain.User{Nombre: "Ana", NombreNormalizado: "ana"})
	if err != nil {
		t.Fatalf("no se pudo preparar a Ana: %v", err)
	}
	proyecto, err := NewCrearProyecto(proyectos).Ejecutar(ctx, CrearProyectoInput{Nombre: "Proyecto", CreadorID: ana.ID})
	if err != nil {
		t.Fatalf("no se pudo preparar el proyecto: %v", err)
	}

	_, err = NewAsignarIntegrante(proyectos, usuarios).Ejecutar(ctx, AsignarIntegranteInput{
		ProyectoID: proyecto.ID, Nombre: "Ana", Rol: domain.RolScrumMaster,
	})

	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "integrante" {
		t.Fatalf("se esperaba el duplicado de integrante (400, campo integrante), se obtuvo %v", err)
	}
	if errors.Is(err, domain.ErrScrumMasterYaExiste) {
		t.Errorf("no debe devolver ErrScrumMasterYaExiste para quien ya es integrante")
	}
}
