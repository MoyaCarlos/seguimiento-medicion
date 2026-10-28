package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func sprintActivoEnMemoria(t *testing.T, repo *SprintRepositoryEnMemoria, proyectoID int64) domain.Sprint {
	t.Helper()
	s, err := domain.NewSprint(proyectoID)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	inicio := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	if err := s.Iniciar("Objetivo", inicio, inicio.AddDate(0, 0, 14)); err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, _ := repo.Guardar(context.Background(), s)
	return guardado
}

func historiaDe(id, proyectoID, sprintID int64, estimacion *int) domain.BacklogItem {
	return domain.BacklogItem{ID: id, ProyectoID: proyectoID, Estado: domain.EstadoNueva, EstimacionSP: estimacion, SprintID: &sprintID}
}

func TestSprintRepositoryEnMemoria_CerrarConArrastre_RegistraFoto(t *testing.T) {
	repo := NewSprintRepositoryEnMemoria()
	activo := sprintActivoEnMemoria(t, repo, 1)
	estimada := historiaDe(1, 1, activo.ID, sp(8))

	err := repo.CerrarConArrastre(context.Background(), cerrado(t, activo), []domain.BacklogItem{estimada, historiaDe(2, 1, activo.ID, nil)})
	if err != nil {
		t.Fatalf("no se esperaba error al cerrar: %v", err)
	}
	*estimada.EstimacionSP = 13

	if s, _ := repo.ObtenerPorID(context.Background(), activo.ID); s.Estado != domain.SprintFinalizado {
		t.Errorf("el Sprint debía quedar Finalizado, quedó %q", s.Estado)
	}
	arrastradas := repo.Arrastradas(activo.ID)
	if len(arrastradas) != 2 || arrastradas[0].ID != 1 || arrastradas[1].ID != 2 {
		t.Fatalf("se esperaban las historias 1 y 2 arrastradas, se obtuvo %+v", arrastradas)
	}
	if arrastradas[0].EstimacionSP == nil || *arrastradas[0].EstimacionSP != 8 {
		t.Errorf("la foto de Story Points debía ser 8 aunque se edite la historia, quedó %v", arrastradas[0].EstimacionSP)
	}
	if arrastradas[1].EstimacionSP != nil {
		t.Errorf("la historia sin estimar debía quedar sin Story Points, quedó %d", *arrastradas[1].EstimacionSP)
	}
}

func TestSprintRepositoryEnMemoria_CerrarConArrastre_FalloNoAplicaNada(t *testing.T) {
	casos := map[string]struct {
		invalida func(sprintID int64) domain.BacklogItem
		esperado error
	}{
		"historia de otro proyecto": {func(sprintID int64) domain.BacklogItem { return historiaDe(2, 2, sprintID, nil) }, ErrIntegridadReferencial},
		"historia de otro Sprint":   {func(sprintID int64) domain.BacklogItem { return historiaDe(2, 1, sprintID+1, nil) }, ErrHistoriaFueraDelSprint},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			repo := NewSprintRepositoryEnMemoria()
			activo := sprintActivoEnMemoria(t, repo, 1)

			err := repo.CerrarConArrastre(context.Background(), cerrado(t, activo),
				[]domain.BacklogItem{historiaDe(1, 1, activo.ID, sp(5)), c.invalida(activo.ID)})

			if !errors.Is(err, c.esperado) {
				t.Fatalf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if s, _ := repo.ObtenerPorID(context.Background(), activo.ID); s.Estado != domain.SprintActivo {
				t.Errorf("el Sprint debía seguir Activo, quedó %q", s.Estado)
			}
			if got := repo.Arrastradas(activo.ID); len(got) != 0 {
				t.Errorf("no debía quedar ninguna historia arrastrada, quedaron %+v", got)
			}
		})
	}
}

func TestSprintRepositoryEnMemoria_CerrarConArrastre_SprintInexistente(t *testing.T) {
	s := domain.Sprint{ID: 999, ProyectoID: 1, Estado: domain.SprintFinalizado}

	err := NewSprintRepositoryEnMemoria().CerrarConArrastre(context.Background(), s, nil)

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("se esperaba ErrSprintNoEncontrado, se obtuvo %v", err)
	}
}
