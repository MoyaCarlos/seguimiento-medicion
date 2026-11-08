package service

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func sprintParaHistorial(id int64, numero int, estado domain.EstadoSprint) domain.Sprint {
	return domain.Sprint{ID: id, ProyectoID: 1, Numero: numero, SprintGoal: "Goal", Estado: estado}
}

func metricasDeSprint(sprintID int64, numero, spCompletados int) domain.MetricasDeSprint {
	return domain.MetricasDeSprint{SprintID: sprintID, Numero: numero, Metricas: domain.Metricas{SPCompletados: spCompletados}}
}

func TestConsultarHistorialSprints_UneCadaSprintFinalizadoConSusMetricas(t *testing.T) {
	sprints := &fakeListadorSprints{sprints: []domain.Sprint{
		sprintParaHistorial(20, 2, domain.SprintFinalizado),
		sprintParaHistorial(30, 3, domain.SprintActivo),
		sprintParaHistorial(10, 1, domain.SprintFinalizado),
		sprintParaHistorial(40, 4, domain.SprintPendiente),
	}}
	metricas := &fakeCalculadorMetricas{resultado: domain.MetricasProyecto{Sprints: []domain.MetricasDeSprint{
		metricasDeSprint(10, 1, 5),
		metricasDeSprint(20, 2, 8),
	}}}
	servicio := NewConsultarHistorialSprints(metricas, sprints)

	historial, err := servicio.Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	var ids []int64
	var completados []int
	for _, h := range historial {
		ids = append(ids, h.Sprint.ID)
		completados = append(completados, h.Metricas.SPCompletados)
	}
	if !slices.Equal(ids, []int64{10, 20}) {
		t.Errorf("se esperaban los Sprints [10 20] en orden de número, se obtuvieron %v", ids)
	}
	if !slices.Equal(completados, []int{5, 8}) {
		t.Errorf("cada Sprint debía llevar sus métricas: se esperaban [5 8], se obtuvieron %v", completados)
	}
}

func TestConsultarHistorialSprints_SinSprintsFinalizadosDevuelveSliceVacio(t *testing.T) {
	sprints := &fakeListadorSprints{sprints: []domain.Sprint{sprintParaHistorial(10, 1, domain.SprintActivo)}}
	servicio := NewConsultarHistorialSprints(&fakeCalculadorMetricas{}, sprints)

	historial, err := servicio.Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if historial == nil || len(historial) != 0 {
		t.Fatalf("se esperaba un slice vacío no nil, se obtuvo %#v", historial)
	}
}

func TestConsultarHistorialSprints_ProyectoInexistenteNoListaSprints(t *testing.T) {
	sprints := &fakeListadorSprints{}
	servicio := NewConsultarHistorialSprints(&fakeCalculadorMetricas{err: domain.ErrProyectoNoEncontrado}, sprints)

	_, err := servicio.Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
	if sprints.llamadas != 0 {
		t.Errorf("no se debían listar los Sprints de un proyecto inexistente; llamadas: %d", sprints.llamadas)
	}
}

func TestConsultarHistorialSprints_SprintFinalizadoSinMetricasEsErrorInterno(t *testing.T) {
	sprints := &fakeListadorSprints{sprints: []domain.Sprint{sprintParaHistorial(10, 1, domain.SprintFinalizado)}}
	servicio := NewConsultarHistorialSprints(&fakeCalculadorMetricas{}, sprints)

	historial, err := servicio.Ejecutar(context.Background(), 1)

	if err == nil {
		t.Fatal("se esperaba un error por la inconsistencia entre fuentes")
	}
	var verr domain.ValidationError
	if errors.As(err, &verr) || errors.Is(err, domain.ErrProyectoNoEncontrado) || errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Errorf("se esperaba un error interno, se obtuvo %v", err)
	}
	if historial != nil {
		t.Errorf("no se debían devolver datos parciales, se obtuvo %#v", historial)
	}
}

func TestConsultarHistorialSprints_PropagaErroresDeLectura(t *testing.T) {
	falla := errors.New("falla de lectura")
	casos := []struct {
		nombre   string
		metricas *fakeCalculadorMetricas
		sprints  *fakeListadorSprints
	}{
		{"del cálculo de métricas", &fakeCalculadorMetricas{err: falla}, &fakeListadorSprints{}},
		{"del listado de Sprints", &fakeCalculadorMetricas{}, &fakeListadorSprints{err: falla}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			historial, err := NewConsultarHistorialSprints(c.metricas, c.sprints).Ejecutar(context.Background(), 1)

			if !errors.Is(err, falla) {
				t.Fatalf("se esperaba la falla de lectura, se obtuvo %v", err)
			}
			if historial != nil {
				t.Errorf("no se debían devolver datos parciales, se obtuvo %#v", historial)
			}
		})
	}
}
