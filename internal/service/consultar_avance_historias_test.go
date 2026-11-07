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

func historiaEnSprintParaAvance(id int64, estado domain.Estado, sprintID int64) domain.BacklogItem {
	h := historiaParaAvance(id, estado, puntos(3))
	h.SprintID = &sprintID
	return h
}

func TestConsultarAvanceHistorias_InformaElSprintDeCadaHistoria(t *testing.T) {
	historias := &fakeAvanceHistoriasRepository{historias: []domain.BacklogItem{
		historiaParaAvance(1, domain.EstadoNueva, puntos(3)),
		historiaEnSprintParaAvance(2, domain.EstadoEnProgreso, 20),
		historiaEnSprintParaAvance(3, domain.EstadoCompletada, 10),
	}}
	sprints := &fakeListadorSprints{sprints: []domain.Sprint{
		{ID: 10, ProyectoID: 1, Numero: 1, Estado: domain.SprintFinalizado},
		{ID: 20, ProyectoID: 1, Numero: 2, Estado: domain.SprintActivo},
	}}
	servicio := NewConsultarAvanceHistorias(proyectosConUno(), historias, sprints)

	avance, err := servicio.Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(avance.Grupos) != 3 {
		t.Fatalf("se esperaban 3 grupos, se obtuvieron %d", len(avance.Grupos))
	}
	if got := avance.Grupos[0].Historias[0].Sprint; got != nil {
		t.Errorf("la historia sin Sprint debe quedar en el Product Backlog (nil), se obtuvo %+v", *got)
	}
	if got := avance.Grupos[1].Historias[0].Sprint; got == nil || *got != (domain.SprintRef{ID: 20, Numero: 2}) {
		t.Errorf("se esperaba el Sprint {20 2} en la historia En progreso, se obtuvo %v", got)
	}
	if got := avance.Grupos[2].Historias[0].Sprint; got == nil || *got != (domain.SprintRef{ID: 10, Numero: 1}) {
		t.Errorf("se esperaba el Sprint {10 1} en la historia Completada, se obtuvo %v", got)
	}
	if sprints.llamadas != 1 {
		t.Errorf("se esperaba un solo listado de Sprints, hubo %d", sprints.llamadas)
	}
}

func TestConsultarAvanceHistorias_ErrorAlListarLosSprintsSePropagaSinDatosParciales(t *testing.T) {
	falla := errors.New("falló el listado de Sprints")
	servicio := NewConsultarAvanceHistorias(proyectosConUno(), &fakeAvanceHistoriasRepository{},
		&fakeListadorSprints{err: falla})

	avance, err := servicio.Ejecutar(context.Background(), 1)

	if !errors.Is(err, falla) {
		t.Fatalf("se esperaba el error del listado, se obtuvo %v", err)
	}
	if len(avance.Grupos) != 0 {
		t.Errorf("no se esperaban datos parciales, se obtuvo %+v", avance)
	}
}

func TestConsultarAvanceHistorias_SprintAusenteDelListadoEsErrorInternoYNoNoEncontrado(t *testing.T) {
	historias := &fakeAvanceHistoriasRepository{historias: []domain.BacklogItem{
		historiaEnSprintParaAvance(1, domain.EstadoEnProgreso, 99),
	}}
	sprints := &fakeListadorSprints{sprints: []domain.Sprint{{ID: 10, ProyectoID: 1, Numero: 1}}}
	servicio := NewConsultarAvanceHistorias(proyectosConUno(), historias, sprints)

	avance, err := servicio.Ejecutar(context.Background(), 1)

	if err == nil {
		t.Fatal("se esperaba un error por el Sprint ausente del listado")
	}
	if errors.Is(err, domain.ErrSprintNoEncontrado) || errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("las fuentes inconsistentes son un error interno, no un 404: %v", err)
	}
	if len(avance.Grupos) != 0 {
		t.Errorf("no se esperaban datos parciales, se obtuvo %+v", avance)
	}
	if sprints.llamadas != 1 {
		t.Errorf("se esperaba un solo listado de Sprints, hubo %d", sprints.llamadas)
	}
}

func TestConsultarAvanceHistorias_ProyectoSinHistoriasTieneLosTresGruposYLosTotalesEnCero(t *testing.T) {
	servicio := NewConsultarAvanceHistorias(proyectosConUno(), &fakeAvanceHistoriasRepository{}, &fakeListadorSprints{})

	avance, err := servicio.Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(avance.Grupos) != 3 {
		t.Fatalf("se esperaban 3 grupos, se obtuvieron %d", len(avance.Grupos))
	}
	for _, g := range avance.Grupos {
		if g.Cantidad != 0 || g.StoryPoints != 0 || g.Historias == nil || len(g.Historias) != 0 {
			t.Errorf("grupo %q: se esperaba vacío con Historias no nil, se obtuvo %+v", g.Estado, g)
		}
	}
	if avance.Total != (domain.TotalesAvance{}) {
		t.Errorf("se esperaban totales en cero, se obtuvo %+v", avance.Total)
	}
}
