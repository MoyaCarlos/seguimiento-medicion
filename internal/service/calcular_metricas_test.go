package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type fakeDatosMetricas struct {
	datos     []domain.DatosMetricasSprint
	err       error
	consultas []int64
}

func (f *fakeDatosMetricas) ListarDatosSprintsFinalizados(_ context.Context, proyectoID int64) ([]domain.DatosMetricasSprint, error) {
	f.consultas = append(f.consultas, proyectoID)
	return f.datos, f.err
}

func TestCalcularMetricas_UsaLasHorasPorStoryPointDelProyecto(t *testing.T) {
	horas := 4.0
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1, HorasPorStoryPoint: &horas}}}
	datos := &fakeDatosMetricas{datos: []domain.DatosMetricasSprint{{SprintID: 7, Numero: 1, HistoriasCompletadas: 1, SPCompletados: 5}}}

	metricas, err := NewCalcularMetricas(proyectos, datos).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(metricas.Sprints) != 1 || metricas.Sprints[0].SprintID != 7 {
		t.Fatalf("se esperaba el Sprint 7, se obtuvo %+v", metricas.Sprints)
	}
	if estimadas := metricas.Proyecto.HorasEstimadas; estimadas == nil || *estimadas != 20 {
		t.Errorf("se esperaban 20 horas estimadas, se obtuvo %v", estimadas)
	}
	if len(datos.consultas) != 1 || datos.consultas[0] != 1 {
		t.Errorf("se esperaba consultar el proyecto 1, se consultó %v", datos.consultas)
	}
	if proyectos.creates != 0 || proyectos.updates != 0 {
		t.Errorf("la consulta no debe escribir: %d altas y %d ediciones", proyectos.creates, proyectos.updates)
	}
}

func TestCalcularMetricas_ProyectoSinHorasPorStoryPoint(t *testing.T) {
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1}}}
	datos := &fakeDatosMetricas{datos: []domain.DatosMetricasSprint{{SprintID: 7, SPCompletados: 5}}}

	metricas, err := NewCalcularMetricas(proyectos, datos).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if metricas.Proyecto.SPCompletados != 5 {
		t.Errorf("se esperaban 5 SP completados, se obtuvo %d", metricas.Proyecto.SPCompletados)
	}
	if metricas.Proyecto.HorasEstimadas != nil {
		t.Errorf("se esperaba sin datos, se obtuvo %v", *metricas.Proyecto.HorasEstimadas)
	}
}

func TestCalcularMetricas_ProyectoInexistente(t *testing.T) {
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{}}
	datos := &fakeDatosMetricas{}

	_, err := NewCalcularMetricas(proyectos, datos).Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
	if len(datos.consultas) != 0 {
		t.Errorf("no se debían consultar los datos: %v", datos.consultas)
	}
}

func TestCalcularMetricas_PropagaElErrorDelRepositorio(t *testing.T) {
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1}}}
	falla := errors.New("falla de lectura")

	_, err := NewCalcularMetricas(proyectos, &fakeDatosMetricas{err: falla}).Ejecutar(context.Background(), 1)

	if !errors.Is(err, falla) {
		t.Fatalf("se esperaba la falla de lectura, se obtuvo %v", err)
	}
}
