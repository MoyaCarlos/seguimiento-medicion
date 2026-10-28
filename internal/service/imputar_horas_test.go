package service

import (
	"context"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

type escenarioImputar struct {
	historias   *fakeHistoriaRepository
	sprints     *repository.SprintRepositoryEnMemoria
	integrantes *fakeIntegrantesRepository
	horas       *fakeRegistroHorasRepository
	sprint      domain.Sprint
	ahora       time.Time
}

// nuevoEscenarioImputar siembra el proyecto 1 con Ana (usuario 2), un Sprint
// Activo del 2026-10-05 al 2026-10-19 y la historia 3 asignada a él. "Hoy" es
// el 2026-10-10.
func nuevoEscenarioImputar(t *testing.T) escenarioImputar {
	t.Helper()
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint, err := sprints.Guardar(context.Background(), domain.Sprint{
		ProyectoID: 1, SprintGoal: "g", Estado: domain.SprintActivo,
		FechaInicio: time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC),
		FechaFin:    time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
	return escenarioImputar{
		historias: &fakeHistoriaRepository{historias: map[int64]domain.BacklogItem{
			3: {ID: 3, ProyectoID: 1, Estado: domain.EstadoNueva, SprintID: &sprint.ID},
		}},
		sprints: sprints,
		integrantes: &fakeIntegrantesRepository{miembros: []domain.Member{
			{UserID: 2, ProjectID: 1, Nombre: "Ana"},
			{UserID: 5, ProjectID: 1, Nombre: "Jimena"},
		}},
		horas:  &fakeRegistroHorasRepository{},
		sprint: sprint,
		ahora:  time.Date(2026, 10, 10, 12, 0, 0, 0, time.Local),
	}
}

func (e escenarioImputar) service() *ImputarHoras {
	return NewImputarHoras(e.historias, e.sprints, e.integrantes, e.horas)
}

func entradaValida(integranteID int64, horas string) ImputarHorasInput {
	return ImputarHorasInput{HistoriaID: 3, IntegranteID: integranteID, Fecha: "2026-10-10", Actividad: "Implementar", Horas: horas}
}

func TestImputarHoras_RegistraUnaCargaValida(t *testing.T) {
	e := nuevoEscenarioImputar(t)

	detalle, err := e.service().Ejecutar(context.Background(), entradaValida(2, "4"), e.ahora)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(e.horas.registros) != 1 {
		t.Fatalf("se esperaba 1 carga registrada, hay %d", len(e.horas.registros))
	}
	registrada := e.horas.registros[0]
	if registrada.SprintID != e.sprint.ID || registrada.HistoriaID != 3 || registrada.ProyectoID != 1 ||
		registrada.UsuarioID != 2 || registrada.Minutos != 240 || registrada.Actividad != "Implementar" {
		t.Errorf("carga registrada inesperada: %+v", registrada)
	}
	if detalle.ID != registrada.ID || detalle.ID == 0 {
		t.Errorf("se esperaba el ID asignado %d, se obtuvo %d", registrada.ID, detalle.ID)
	}
	if detalle.NombreIntegrante != "Ana" {
		t.Errorf("se esperaba el nombre Ana, se obtuvo %q", detalle.NombreIntegrante)
	}
}

func TestImputarHoras_DosCargasQuedanRegistradas(t *testing.T) {
	e := nuevoEscenarioImputar(t)
	ctx := context.Background()

	_, err1 := e.service().Ejecutar(ctx, entradaValida(2, "1.5"), e.ahora)
	_, err2 := e.service().Ejecutar(ctx, entradaValida(5, "2.5"), e.ahora)

	if err1 != nil || err2 != nil {
		t.Fatalf("no se esperaban errores: %v, %v", err1, err2)
	}
	if len(e.horas.registros) != 2 {
		t.Fatalf("se esperaban 2 cargas, hay %d", len(e.horas.registros))
	}
	if e.horas.registros[0].Minutos+e.horas.registros[1].Minutos != 240 {
		t.Errorf("se esperaban 240 minutos en total, se obtuvo %d", e.horas.registros[0].Minutos+e.horas.registros[1].Minutos)
	}
}
