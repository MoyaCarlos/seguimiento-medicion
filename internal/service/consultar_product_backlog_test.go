package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func proyectosConUno() *fakeProjectRepository {
	return &fakeProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1, Nombre: "Proyecto"}}}
}

func TestConsultarProductBacklog_ProyectoInexistente(t *testing.T) {
	backlog := &fakeProductBacklogRepository{}
	servicio := NewConsultarProductBacklog(&fakeProjectRepository{}, backlog)

	_, err := servicio.Ejecutar(context.Background(), 999)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
	if backlog.llamadas != 0 {
		t.Errorf("no se esperaba consultar el backlog de un proyecto inexistente")
	}
}

func TestConsultarProductBacklog_DevuelveLasHistoriasDelRepositorio(t *testing.T) {
	esperadas := []domain.BacklogItem{
		{ID: 2, ProyectoID: 1, Titulo: "Iniciar sesión", Prioridad: domain.PrioridadMust},
		{ID: 1, ProyectoID: 1, Titulo: "Cerrar sesión", Prioridad: domain.PrioridadShould},
	}
	servicio := NewConsultarProductBacklog(proyectosConUno(), &fakeProductBacklogRepository{historias: esperadas})

	historias, err := servicio.Ejecutar(context.Background(), 1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(historias) != len(esperadas) {
		t.Fatalf("se esperaban %d historias, se obtuvieron %d", len(esperadas), len(historias))
	}
	for i := range esperadas {
		if historias[i].ID != esperadas[i].ID {
			t.Fatalf("se esperaba el orden del repositorio %v, se obtuvo %v", esperadas, historias)
		}
	}
}

func TestConsultarProductBacklog_SinHistoriasDevuelveVacio(t *testing.T) {
	servicio := NewConsultarProductBacklog(proyectosConUno(), &fakeProductBacklogRepository{})

	historias, err := servicio.Ejecutar(context.Background(), 1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if historias == nil || len(historias) != 0 {
		t.Errorf("se esperaba un slice vacío no nil, se obtuvo %#v", historias)
	}
}

func TestConsultarProductBacklog_PropagaErrorDelRepositorio(t *testing.T) {
	errPersistencia := errors.New("falla de persistencia")
	servicio := NewConsultarProductBacklog(proyectosConUno(), &fakeProductBacklogRepository{err: errPersistencia})

	_, err := servicio.Ejecutar(context.Background(), 1)

	if !errors.Is(err, errPersistencia) {
		t.Fatalf("se esperaba el error del repositorio, se obtuvo %v", err)
	}
}
