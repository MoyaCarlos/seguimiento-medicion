package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// jimenaDadaDeBaja suma a Jimena (id 2) como Product Builder y la da de baja (HU-16).
func jimenaDadaDeBaja(t *testing.T) (*fakeProjectRepository, *fakeCredenciales, int64, domain.Member) {
	t.Helper()
	proyectos, credenciales, proyectoID := equipoDePrueba(t)
	alta, err := agregar(proyectos, credenciales, proyectoID, "Jimena", domain.RolProductBuilder)
	if err != nil {
		t.Fatalf("no se pudo sumar a Jimena: %v", err)
	}
	if err := proyectos.DarDeBajaIntegrante(context.Background(), proyectoID, alta.UserID, time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("no se pudo dar de baja a Jimena: %v", err)
	}
	return proyectos, credenciales, proyectoID, alta
}

// HU-16, Escenario 3: agregar de nuevo a quien fue dado de baja lo reactiva como
// Product Builder, una sola vez y con su alta original.
func TestAsignarIntegrante_ReactivaAUnDadoDeBaja(t *testing.T) {
	proyectos, credenciales, proyectoID, alta := jimenaDadaDeBaja(t)

	miembro, err := agregar(proyectos, credenciales, proyectoID, "jimena", "")

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if miembro.UserID != alta.UserID || miembro.Role != domain.RolProductBuilder || !miembro.CreadoEn.Equal(alta.CreadoEn) {
		t.Errorf("se esperaba a Jimena reactivada como product_builder con su alta original, se obtuvo %+v", miembro)
	}
	activos, _ := proyectos.ListarIntegrantes(context.Background(), proyectoID)
	if len(activos) != 2 || len(proyectos.miembros) != 2 {
		t.Errorf("se esperaba a Jimena una sola vez: activos %+v, vínculos %+v", activos, proyectos.miembros)
	}
}

// HU-16 (#236): un dado de baja no vuelve como Scrum Master, porque el proyecto
// ya tiene uno.
func TestAsignarIntegrante_DadoDeBajaComoScrumMasterEsConflicto(t *testing.T) {
	proyectos, credenciales, proyectoID, _ := jimenaDadaDeBaja(t)

	_, err := agregar(proyectos, credenciales, proyectoID, "Jimena", domain.RolScrumMaster)

	if !errors.Is(err, domain.ErrScrumMasterYaExiste) {
		t.Fatalf("se esperaba ErrScrumMasterYaExiste, se obtuvo %v", err)
	}
	if activos, _ := proyectos.ListarIntegrantes(context.Background(), proyectoID); len(activos) != 1 {
		t.Errorf("Jimena no debía reactivarse: %+v", activos)
	}
}

// HU-16, Escenario 4 (Clarification Q2): al Scrum Master no se le cambia el rol;
// agregarlo de nuevo, con cualquier rol, es un integrante duplicado.
func TestAsignarIntegrante_ElScrumMasterNoCambiaDeRol(t *testing.T) {
	for _, rol := range []domain.Role{domain.RolProductBuilder, domain.RolScrumMaster} {
		proyectos, credenciales, proyectoID := equipoDePrueba(t)

		_, err := agregar(proyectos, credenciales, proyectoID, "Ana", rol)

		var verr domain.ValidationError
		if !errors.As(err, &verr) || verr.Mensaje != "el integrante ya pertenece al proyecto" {
			t.Errorf("rol %q: se esperaba el integrante duplicado, se obtuvo %v", rol, err)
		}
		activos, _ := proyectos.ListarIntegrantes(context.Background(), proyectoID)
		if len(activos) != 1 || activos[0].Role != domain.RolScrumMaster {
			t.Errorf("rol %q: el equipo no debía cambiar: %+v", rol, activos)
		}
	}
}
