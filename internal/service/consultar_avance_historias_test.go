package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func historiaParaAvance(id int64, estado domain.Estado, sp *int) domain.BacklogItem {
	return domain.BacklogItem{ID: id, ProyectoID: 1, Titulo: "Historia", Estado: estado, EstimacionSP: sp}
}

func TestConsultarAvanceHistorias_AgrupaLasHistoriasDelRepositorio(t *testing.T) {
	historias := &fakeAvanceHistoriasRepository{historias: []domain.BacklogItem{
		historiaParaAvance(1, domain.EstadoNueva, puntos(3)),
		historiaParaAvance(2, domain.EstadoEnProgreso, puntos(5)),
		historiaParaAvance(3, domain.EstadoCompletada, puntos(8)),
		historiaParaAvance(4, domain.EstadoCompletada, nil),
	}}
	servicio := NewConsultarAvanceHistorias(proyectosConUno(), historias, &fakeListadorSprints{})

	avance, err := servicio.Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(avance.Grupos) != 3 {
		t.Fatalf("se esperaban 3 grupos, se obtuvieron %d", len(avance.Grupos))
	}
	esperados := []struct{ cantidad, storyPoints int }{{1, 3}, {1, 5}, {2, 8}}
	for i, e := range esperados {
		g := avance.Grupos[i]
		if g.Cantidad != e.cantidad || g.StoryPoints != e.storyPoints {
			t.Errorf("grupo %q: se esperaban %d historias y %d Story Points, se obtuvo %d y %d",
				g.Estado, e.cantidad, e.storyPoints, g.Cantidad, g.StoryPoints)
		}
	}
	if avance.Total != (domain.TotalesAvance{Historias: 4, StoryPoints: 16}) {
		t.Errorf("se esperaban 4 historias y 16 Story Points en total, se obtuvo %+v", avance.Total)
	}
	if historias.llamadas != 1 {
		t.Errorf("se esperaba una sola lectura de historias, hubo %d", historias.llamadas)
	}
}

func TestConsultarAvanceHistorias_ProyectoInexistente(t *testing.T) {
	historias := &fakeAvanceHistoriasRepository{}
	sprints := &fakeListadorSprints{}
	servicio := NewConsultarAvanceHistorias(&fakeProjectRepository{}, historias, sprints)

	avance, err := servicio.Ejecutar(context.Background(), 999)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
	if len(avance.Grupos) != 0 {
		t.Errorf("no se esperaban datos parciales, se obtuvo %+v", avance)
	}
	if historias.llamadas != 0 || sprints.llamadas != 0 {
		t.Errorf("no se esperaba leer historias (%d) ni Sprints (%d) de un proyecto inexistente", historias.llamadas, sprints.llamadas)
	}
}

func TestConsultarAvanceHistorias_PropagaLosErroresSinDatosParciales(t *testing.T) {
	falla := errors.New("falló la lectura")
	casos := map[string]struct {
		proyectos *fakeProjectRepository
		historias *fakeAvanceHistoriasRepository
	}{
		"del repositorio de proyectos": {
			proyectos: &fakeProjectRepository{err: falla},
			historias: &fakeAvanceHistoriasRepository{},
		},
		"del repositorio de historias": {
			proyectos: proyectosConUno(),
			historias: &fakeAvanceHistoriasRepository{err: falla},
		},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			servicio := NewConsultarAvanceHistorias(c.proyectos, c.historias, &fakeListadorSprints{})

			avance, err := servicio.Ejecutar(context.Background(), 1)

			if !errors.Is(err, falla) {
				t.Fatalf("se esperaba el error de la lectura, se obtuvo %v", err)
			}
			if len(avance.Grupos) != 0 {
				t.Errorf("no se esperaban datos parciales, se obtuvo %+v", avance)
			}
		})
	}
}
