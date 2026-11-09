package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestProjectRepositoryEnMemoria_AgregarIntegranteProyectoInexistente(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	_, err := repo.AgregarIntegrante(context.Background(), domain.Membership{
		ProjectID: 1, UserID: 1, Role: domain.RolScrumMaster,
	})
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
}

func TestProjectRepositoryEnMemoria_AgregarIntegranteNoValidaUsuario(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	guardado, err := repo.Guardar(context.Background(), domain.Project{Nombre: "Proyecto"})
	if err != nil {
		t.Fatalf("no se esperaba error al guardar: %v", err)
	}
	_, err = repo.AgregarIntegrante(context.Background(), domain.Membership{
		ProjectID: guardado.ID, UserID: 42, Role: domain.RolScrumMaster,
	})
	if err != nil {
		t.Fatalf("el fake no debe validar usuarios en otro repositorio: %v", err)
	}
}

func TestProjectRepositoryEnMemoria_GuardarConScrumMaster_RollbackSiFallaMembresia(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	repo.miembros = []domain.Membership{{ProjectID: 1, UserID: 42, Role: domain.RolScrumMaster}}

	_, err := repo.GuardarConScrumMaster(context.Background(), domain.Project{Nombre: "Proyecto"}, 42)
	if !errors.Is(err, ErrMiembroDuplicado) {
		t.Fatalf("se esperaba ErrMiembroDuplicado, se obtuvo %v", err)
	}
	if _, err := repo.ObtenerPorID(context.Background(), 1); !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("el proyecto debía revertirse, se obtuvo %v", err)
	}
}

func TestProjectRepositoryEnMemoria_AgregarIntegranteDuplicado(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	guardado, _ := repo.Guardar(context.Background(), domain.Project{Nombre: "Proyecto"})
	usuario, _ := repo.GuardarUsuario(context.Background(), domain.User{Nombre: "Ana", NombreNormalizado: "ana"})
	m := domain.Membership{ProjectID: guardado.ID, UserID: usuario.ID, Role: domain.RolScrumMaster}

	if _, err := repo.AgregarIntegrante(context.Background(), m); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if _, err := repo.AgregarIntegrante(context.Background(), m); !errors.Is(err, ErrMiembroDuplicado) {
		t.Fatalf("se esperaba ErrMiembroDuplicado, se obtuvo %v", err)
	}
}

func TestProjectRepositoryEnMemoria_ListarIntegrantesDevuelveNombre(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	guardado, _ := repo.Guardar(context.Background(), domain.Project{Nombre: "Proyecto"})
	usuario, _ := repo.GuardarUsuario(context.Background(), domain.User{Nombre: "Candela", NombreNormalizado: "candela"})

	if _, err := repo.AgregarIntegrante(context.Background(), domain.Membership{
		ProjectID: guardado.ID, UserID: usuario.ID, Role: domain.RolProductBuilder,
	}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	miembros, err := repo.ListarIntegrantes(context.Background(), guardado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(miembros) != 1 || miembros[0].Nombre != "Candela" {
		t.Fatalf("se esperaba el integrante con su nombre, se obtuvo %+v", miembros)
	}
}

func TestProjectRepositoryEnMemoria_AgregarIntegrante_SegundoScrumMasterEsErrorDeDominio(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	guardado, _ := repo.Guardar(context.Background(), domain.Project{Nombre: "Proyecto"})
	ana, _ := repo.GuardarUsuario(context.Background(), domain.User{Nombre: "Ana", NombreNormalizado: "ana"})
	beto, _ := repo.GuardarUsuario(context.Background(), domain.User{Nombre: "Beto", NombreNormalizado: "beto"})
	if _, err := repo.AgregarIntegrante(context.Background(), domain.Membership{ProjectID: guardado.ID, UserID: ana.ID, Role: domain.RolScrumMaster}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	_, err := repo.AgregarIntegrante(context.Background(), domain.Membership{ProjectID: guardado.ID, UserID: beto.ID, Role: domain.RolScrumMaster})

	if !errors.Is(err, domain.ErrScrumMasterYaExiste) {
		t.Fatalf("se esperaba domain.ErrScrumMasterYaExiste, se obtuvo %v", err)
	}
}

func TestProjectRepositoryEnMemoria_AgregarIntegrante_ProductBuildersYScrumMasterEnOtroProyecto(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	uno, _ := repo.Guardar(context.Background(), domain.Project{Nombre: "Uno"})
	dos, _ := repo.Guardar(context.Background(), domain.Project{Nombre: "Dos"})
	altas := []domain.Membership{
		{ProjectID: uno.ID, UserID: 1, Role: domain.RolScrumMaster},
		{ProjectID: uno.ID, UserID: 2, Role: domain.RolProductBuilder},
		{ProjectID: uno.ID, UserID: 3, Role: domain.RolProductBuilder},
		{ProjectID: dos.ID, UserID: 1, Role: domain.RolScrumMaster},
	}
	for _, m := range altas {
		if _, err := repo.AgregarIntegrante(context.Background(), m); err != nil {
			t.Fatalf("no se esperaba error para %+v: %v", m, err)
		}
	}
}
