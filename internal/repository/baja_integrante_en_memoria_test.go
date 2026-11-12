package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// HU-16: la baja lógica en el repositorio en memoria.
func TestProjectRepositoryEnMemoria_DarDeBajaIntegrante(t *testing.T) {
	ctx := context.Background()
	fecha := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	preparar := func(t *testing.T) (*ProjectRepositoryEnMemoria, int64, domain.User, domain.User, domain.User) {
		t.Helper()
		repo := NewProjectRepositoryEnMemoria()
		ana, _ := repo.GuardarUsuario(ctx, domain.User{Nombre: "Ana"})
		beto, _ := repo.GuardarUsuario(ctx, domain.User{Nombre: "Beto"})
		ajeno, _ := repo.GuardarUsuario(ctx, domain.User{Nombre: "Carla"})
		p, _ := repo.GuardarConScrumMaster(ctx, domain.Project{Nombre: "Uno"}, ana.ID)
		if _, err := repo.AgregarIntegrante(ctx, domain.Membership{ProjectID: p.ID, UserID: beto.ID, Role: domain.RolProductBuilder}); err != nil {
			t.Fatalf("no se pudo preparar el integrante: %v", err)
		}
		return repo, p.ID, ana, beto, ajeno
	}

	t.Run("Product Builder activo", func(t *testing.T) {
		repo, p, _, beto, _ := preparar(t)
		if err := repo.DarDeBajaIntegrante(ctx, p, beto.ID, fecha); err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		activos, _ := repo.ListarIntegrantes(ctx, p)
		if got := nombresDe(activos); len(got) != 1 || got[0] != "Ana:scrum_master" {
			t.Errorf("ListarIntegrantes: se esperaba solo a Ana, se obtuvo %v", got)
		}
		historicos, _ := repo.ListarIntegrantesHistoricos(ctx, p)
		if got := nombresDe(historicos); len(got) != 2 || got[1] != "Beto:product_builder" {
			t.Errorf("ListarIntegrantesHistoricos: se esperaba que Beto siga, se obtuvo %v", got)
		}
	})

	t.Run("segunda baja del mismo", func(t *testing.T) {
		repo, p, _, beto, _ := preparar(t)
		if err := repo.DarDeBajaIntegrante(ctx, p, beto.ID, fecha); err != nil {
			t.Fatalf("primera baja: %v", err)
		}
		if err := repo.DarDeBajaIntegrante(ctx, p, beto.ID, fecha); !errors.Is(err, domain.ErrIntegranteNoEncontrado) {
			t.Errorf("se esperaba ErrIntegranteNoEncontrado, se obtuvo %v", err)
		}
	})

	t.Run("usuario que nunca integró", func(t *testing.T) {
		repo, p, _, _, ajeno := preparar(t)
		if err := repo.DarDeBajaIntegrante(ctx, p, ajeno.ID, fecha); !errors.Is(err, domain.ErrIntegranteNoEncontrado) {
			t.Errorf("se esperaba ErrIntegranteNoEncontrado, se obtuvo %v", err)
		}
	})

	t.Run("Scrum Master", func(t *testing.T) {
		repo, p, ana, _, _ := preparar(t)
		if err := repo.DarDeBajaIntegrante(ctx, p, ana.ID, fecha); !errors.Is(err, domain.ErrScrumMasterNoSeQuita) {
			t.Errorf("se esperaba ErrScrumMasterNoSeQuita, se obtuvo %v", err)
		}
		activos, _ := repo.ListarIntegrantes(ctx, p)
		if len(activos) != 2 {
			t.Errorf("el equipo no debe cambiar, se obtuvo %v", nombresDe(activos))
		}
	})
}
