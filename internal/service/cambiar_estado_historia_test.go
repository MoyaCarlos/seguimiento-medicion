package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

func relojFijo() time.Time { return time.Date(2026, 10, 29, 15, 30, 0, 0, time.UTC) }

func historiaEnSprint(id, sprintID int64, estado domain.Estado) domain.BacklogItem {
	h := historiaSinSprint(id, 1, estado)
	h.SprintID = &sprintID
	return h
}

func fakeCambioConHistorias(items ...domain.BacklogItem) *fakeCambioEstadoRepository {
	f := &fakeCambioEstadoRepository{historias: map[int64]domain.BacklogItem{}}
	for _, h := range items {
		f.historias[h.ID] = h
	}
	return f
}

func TestCambiarEstadoHistoria_NuevaAEnProgreso(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	activo := sembrarSprintActivo(t, sprints, 1)
	historias := fakeCambioConHistorias(historiaEnSprint(1, activo.ID, domain.EstadoNueva))

	h, err := NewCambiarEstadoHistoria(historias, sprints, relojFijo).Ejecutar(context.Background(), 1, domain.EstadoEnProgreso)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if h.Estado != domain.EstadoEnProgreso {
		t.Errorf("Estado = %q, quería En progreso", h.Estado)
	}
	if h.FechaFinalizacion != nil {
		t.Errorf("FechaFinalizacion = %v, quería nil", h.FechaFinalizacion)
	}
	llamadas := historias.registroDeLlamadas()
	if len(llamadas) != 1 {
		t.Fatalf("CambiarEstado se llamó %d veces, quería 1", len(llamadas))
	}
	if l := llamadas[0]; l.desde != domain.EstadoNueva || l.hacia != domain.EstadoEnProgreso || l.fecha != nil {
		t.Errorf("llamada = %+v, quería desde Nueva, hacia En progreso, fecha nil", l)
	}
}

func TestCambiarEstadoHistoria_EnProgresoACompletada(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	activo := sembrarSprintActivo(t, sprints, 1)
	historias := fakeCambioConHistorias(historiaEnSprint(1, activo.ID, domain.EstadoEnProgreso))

	h, err := NewCambiarEstadoHistoria(historias, sprints, relojFijo).Ejecutar(context.Background(), 1, domain.EstadoCompletada)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if h.Estado != domain.EstadoCompletada {
		t.Errorf("Estado = %q, quería Completada", h.Estado)
	}
	if h.FechaFinalizacion == nil || !h.FechaFinalizacion.Equal(relojFijo()) {
		t.Errorf("FechaFinalizacion = %v, quería %v", h.FechaFinalizacion, relojFijo())
	}
	llamadas := historias.registroDeLlamadas()
	if len(llamadas) != 1 {
		t.Fatalf("CambiarEstado se llamó %d veces, quería 1", len(llamadas))
	}
	if l := llamadas[0]; l.desde != domain.EstadoEnProgreso || l.hacia != domain.EstadoCompletada || l.fecha == nil {
		t.Errorf("llamada = %+v, quería desde En progreso, hacia Completada, fecha no nula", l)
	}
}

func TestCambiarEstadoHistoria_HistoriaInexistente(t *testing.T) {
	historias := fakeCambioConHistorias()
	caso := NewCambiarEstadoHistoria(historias, repository.NewSprintRepositoryEnMemoria(), relojFijo)

	_, err := caso.Ejecutar(context.Background(), 99, domain.EstadoEnProgreso)

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Errorf("err = %v, quería ErrHistoriaNoEncontrada", err)
	}
	if n := len(historias.registroDeLlamadas()); n != 0 {
		t.Errorf("CambiarEstado se llamó %d veces, no debía llamarse", n)
	}
}

func TestCambiarEstadoHistoria_HistoriaSinSprint(t *testing.T) {
	historias := fakeCambioConHistorias(historiaSinSprint(1, 1, domain.EstadoNueva))
	caso := NewCambiarEstadoHistoria(historias, repository.NewSprintRepositoryEnMemoria(), relojFijo)

	_, err := caso.Ejecutar(context.Background(), 1, domain.EstadoEnProgreso)

	if !errors.Is(err, domain.ErrHistoriaSinSprint) {
		t.Errorf("err = %v, quería ErrHistoriaSinSprint", err)
	}
	if n := len(historias.registroDeLlamadas()); n != 0 {
		t.Errorf("CambiarEstado se llamó %d veces, no debía llamarse", n)
	}
}

func TestCambiarEstadoHistoria_SprintNoActivo(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	pendiente := sembrarSprintPendiente(t, sprints, 1)
	finalizado := sembrarSprintFinalizado(t, sprints, 1)

	for _, c := range []struct {
		nombre   string
		sprintID int64
	}{{"Sprint pendiente", pendiente.ID}, {"Sprint finalizado", finalizado.ID}} {
		t.Run(c.nombre, func(t *testing.T) {
			historias := fakeCambioConHistorias(historiaEnSprint(1, c.sprintID, domain.EstadoNueva))

			_, err := NewCambiarEstadoHistoria(historias, sprints, relojFijo).Ejecutar(context.Background(), 1, domain.EstadoEnProgreso)

			if !errors.Is(err, domain.ErrSprintNoActivoParaEstado) {
				t.Errorf("err = %v, quería ErrSprintNoActivoParaEstado", err)
			}
			if n := len(historias.registroDeLlamadas()); n != 0 {
				t.Errorf("CambiarEstado se llamó %d veces, no debía llamarse", n)
			}
		})
	}
}

func TestCambiarEstadoHistoria_ErroresDeReglaSePropaganTalCual(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	activo := sembrarSprintActivo(t, sprints, 1)

	casos := []struct {
		nombre   string
		actual   domain.Estado
		nuevo    domain.Estado
		esperado error
	}{
		{"salto directo", domain.EstadoNueva, domain.EstadoCompletada, domain.ErrSaltoDirecto},
		{"retroceso", domain.EstadoEnProgreso, domain.EstadoNueva, domain.ErrRetrocesoEstado},
		{"estado repetido", domain.EstadoEnProgreso, domain.EstadoEnProgreso, domain.ErrEstadoRepetido},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			historias := fakeCambioConHistorias(historiaEnSprint(1, activo.ID, c.actual))

			_, err := NewCambiarEstadoHistoria(historias, sprints, relojFijo).Ejecutar(context.Background(), 1, c.nuevo)

			if !errors.Is(err, c.esperado) {
				t.Errorf("err = %v, quería %v", err, c.esperado)
			}
			if n := len(historias.registroDeLlamadas()); n != 0 {
				t.Errorf("CambiarEstado se llamó %d veces, no debía llamarse", n)
			}
		})
	}
}

// ganaOtroCambioDeEstado lee la historia en su estado original, pero rechaza la
// escritura: simula que otro cambio ganó entre la lectura y el UPDATE.
type ganaOtroCambioDeEstado struct {
	*fakeCambioEstadoRepository
}

func (ganaOtroCambioDeEstado) CambiarEstado(context.Context, int64, domain.Estado, domain.Estado, *time.Time) error {
	return domain.ErrConflictoEstadoHistoria
}

func TestCambiarEstadoHistoria_PropagaLaCarreraDelRepositorio(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	activo := sembrarSprintActivo(t, sprints, 1)
	historias := ganaOtroCambioDeEstado{fakeCambioConHistorias(historiaEnSprint(1, activo.ID, domain.EstadoNueva))}

	_, err := NewCambiarEstadoHistoria(historias, sprints, relojFijo).Ejecutar(context.Background(), 1, domain.EstadoEnProgreso)

	if !errors.Is(err, domain.ErrConflictoEstadoHistoria) {
		t.Errorf("err = %v, quería ErrConflictoEstadoHistoria", err)
	}
}
