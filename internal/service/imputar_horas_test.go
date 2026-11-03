package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

type escenarioImputar struct {
	historias   *fakeHistoriaRepository
	sprints     *fakeSprintsQueCuentan
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
	sprints := &fakeSprintsQueCuentan{SprintRepositoryEnMemoria: repository.NewSprintRepositoryEnMemoria()}
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

func (e escenarioImputar) moverHistoriaA(t *testing.T, estado domain.EstadoSprint) {
	t.Helper()
	otro, err := e.sprints.Guardar(context.Background(), domain.Sprint{ProyectoID: 1, SprintGoal: "g", Estado: estado})
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
	h := e.historias.historias[3]
	h.SprintID = &otro.ID
	e.historias.historias[3] = h
}

func (e escenarioImputar) sacarHistoriaDelSprint() {
	h := e.historias.historias[3]
	h.SprintID = nil
	e.historias.historias[3] = h
}

func TestImputarHoras_Rechazos(t *testing.T) {
	casos := []struct {
		nombre   string
		preparar func(t *testing.T, e escenarioImputar)
		entrada  ImputarHorasInput
		esperado error
	}{
		{"historia inexistente", nil,
			ImputarHorasInput{HistoriaID: 99, IntegranteID: 2, Fecha: "2026-10-10", Actividad: "x", Horas: "1"}, domain.ErrHistoriaNoEncontrada},
		{"historia sin Sprint", func(_ *testing.T, e escenarioImputar) { e.sacarHistoriaDelSprint() },
			entradaValida(2, "1"), domain.ErrHistoriaFueraDelSprintActivo},
		{"historia en un Sprint Pendiente", func(t *testing.T, e escenarioImputar) { e.moverHistoriaA(t, domain.SprintPendiente) },
			entradaValida(2, "1"), domain.ErrHistoriaFueraDelSprintActivo},
		{"historia en un Sprint Finalizado", func(t *testing.T, e escenarioImputar) { e.moverHistoriaA(t, domain.SprintFinalizado) },
			entradaValida(2, "1"), domain.ErrHistoriaFueraDelSprintActivo},
		{"integrante de otro proyecto", nil, entradaValida(9, "1"), domain.ErrIntegranteAjeno},
		{"integrante inexistente", nil, entradaValida(99, "1"), domain.ErrIntegranteAjeno},
		{"horas inválidas", nil, entradaValida(2, "0"),
			domain.ValidationError{Campo: "horas", Mensaje: "las horas deben ser mayores que cero"}},
		{"fuera del Sprint, integrante ajeno y horas inválidas", func(_ *testing.T, e escenarioImputar) { e.sacarHistoriaDelSprint() },
			entradaValida(99, "0"), domain.ErrHistoriaFueraDelSprintActivo},
		{"integrante ajeno y horas inválidas", nil, entradaValida(99, "0"), domain.ErrIntegranteAjeno},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := nuevoEscenarioImputar(t)
			e.integrantes.miembros = append(e.integrantes.miembros, domain.Member{UserID: 9, ProjectID: 2, Nombre: "Pedro"})
			if c.preparar != nil {
				c.preparar(t, e)
			}

			_, err := e.service().Ejecutar(context.Background(), c.entrada, e.ahora)

			if !errors.Is(err, c.esperado) {
				t.Errorf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if len(e.horas.registros) != 0 {
				t.Errorf("no se debía registrar ninguna carga, hay %d", len(e.horas.registros))
			}
		})
	}
}

func TestImputarHoras_HistoriaSinSprintNoConsultaElSprint(t *testing.T) {
	e := nuevoEscenarioImputar(t)
	e.sacarHistoriaDelSprint()

	_, _ = e.service().Ejecutar(context.Background(), entradaValida(2, "1"), e.ahora)

	if e.sprints.lecturas != 0 {
		t.Errorf("no se debía consultar el Sprint, se consultó %d veces", e.sprints.lecturas)
	}
}

func TestImputarHoras_PropagaElCierreDelSprintAlRegistrar(t *testing.T) {
	e := nuevoEscenarioImputar(t)
	e.horas.errRegistrar = domain.ErrHistoriaFueraDelSprintActivo

	_, err := e.service().Ejecutar(context.Background(), entradaValida(2, "1"), e.ahora)

	if !errors.Is(err, domain.ErrHistoriaFueraDelSprintActivo) {
		t.Errorf("se esperaba ErrHistoriaFueraDelSprintActivo, se obtuvo %v", err)
	}
}
