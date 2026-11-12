package repository

import (
	"context"
	"errors"
	"testing"
	"time"

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

// Spec 044 (#300). El repositorio en memoria no conoce los Sprints: la regla
// del Sprint Activo la da el dominio (Project.Cerrar).
func TestProjectRepositoryEnMemoria_Cerrar(t *testing.T) {
	ctx := context.Background()
	hoy := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)

	t.Run("abierto guarda la fecha", func(t *testing.T) {
		repo := NewProjectRepositoryEnMemoria()
		p, _ := repo.Guardar(ctx, domain.Project{Nombre: "Tienda"})

		if err := repo.Cerrar(ctx, p.ID, hoy); err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		leido, _ := repo.ObtenerPorID(ctx, p.ID)
		if leido.FechaCierre == nil || !leido.FechaCierre.Equal(hoy) {
			t.Errorf("se esperaba la fecha de cierre %v, se obtuvo %v", hoy, leido.FechaCierre)
		}
	})

	t.Run("dos cierres seguidos: el segundo informa que ya está cerrado", func(t *testing.T) {
		repo := NewProjectRepositoryEnMemoria()
		p, _ := repo.Guardar(ctx, domain.Project{Nombre: "Tienda"})
		_ = repo.Cerrar(ctx, p.ID, hoy)

		if err := repo.Cerrar(ctx, p.ID, hoy.AddDate(0, 0, 1)); !errors.Is(err, domain.ErrProyectoYaCerrado) {
			t.Fatalf("se esperaba ErrProyectoYaCerrado, se obtuvo %v", err)
		}
		leido, _ := repo.ObtenerPorID(ctx, p.ID)
		if !leido.FechaCierre.Equal(hoy) {
			t.Errorf("la fecha del primer cierre no debe cambiar, se obtuvo %v", *leido.FechaCierre)
		}
	})

	t.Run("inexistente", func(t *testing.T) {
		if err := NewProjectRepositoryEnMemoria().Cerrar(ctx, 99, hoy); !errors.Is(err, domain.ErrProyectoNoEncontrado) {
			t.Fatalf("se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
		}
	})
}

// Como en SQLite, editar el proyecto no toca la fecha de cierre.
func TestProjectRepositoryEnMemoria_ActualizarNoPisaLaFechaDeCierre(t *testing.T) {
	ctx := context.Background()
	repo := NewProjectRepositoryEnMemoria()
	p, _ := repo.Guardar(ctx, domain.Project{Nombre: "Tienda"})
	hoy := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)
	_ = repo.Cerrar(ctx, p.ID, hoy)

	p.Nombre = "Otro nombre"
	p.FechaCierre = nil
	if err := repo.Actualizar(ctx, p); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	leido, _ := repo.ObtenerPorID(ctx, p.ID)
	if leido.Nombre != "Otro nombre" || leido.FechaCierre == nil {
		t.Errorf("se esperaba el nombre editado y la fecha de cierre intacta, se obtuvo %+v", leido)
	}
}

func TestProjectRepositoryEnMemoria_ProyectoCerrado(t *testing.T) {
	ctx := context.Background()
	repo := NewProjectRepositoryEnMemoria()
	abierto, _ := repo.Guardar(ctx, domain.Project{Nombre: "Abierto"})
	cerrado, _ := repo.Guardar(ctx, domain.Project{Nombre: "Cerrado"})
	_ = repo.Cerrar(ctx, cerrado.ID, time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC))

	for id, esperado := range map[int64]bool{abierto.ID: false, cerrado.ID: true, 99: false} {
		got, err := repo.ProyectoCerrado(ctx, id)
		if err != nil {
			t.Fatalf("proyecto %d: no se esperaba error: %v", id, err)
		}
		if got != esperado {
			t.Errorf("proyecto %d: se esperaba %v, se obtuvo %v", id, esperado, got)
		}
	}
}

func TestProjectRepositoryEnMemoria_Reabrir(t *testing.T) {
	ctx := context.Background()
	repo := NewProjectRepositoryEnMemoria()
	p, _ := repo.Guardar(ctx, domain.Project{Nombre: "Tienda"})

	if err := repo.Reabrir(ctx, p.ID); !errors.Is(err, domain.ErrProyectoNoCerrado) {
		t.Fatalf("abierto: se esperaba ErrProyectoNoCerrado, se obtuvo %v", err)
	}
	_ = repo.Cerrar(ctx, p.ID, time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC))
	if err := repo.Reabrir(ctx, p.ID); err != nil {
		t.Fatalf("cerrado: no se esperaba error: %v", err)
	}
	if leido, _ := repo.ObtenerPorID(ctx, p.ID); leido.FechaCierre != nil {
		t.Errorf("se esperaba sin fecha de cierre, se obtuvo %v", *leido.FechaCierre)
	}
	if err := repo.Reabrir(ctx, p.ID); !errors.Is(err, domain.ErrProyectoNoCerrado) {
		t.Fatalf("segunda reapertura: se esperaba ErrProyectoNoCerrado, se obtuvo %v", err)
	}
	if err := repo.Reabrir(ctx, 99); !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("inexistente: se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
}
