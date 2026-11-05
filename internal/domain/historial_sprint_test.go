package domain

import (
	"errors"
	"slices"
	"testing"
)

func sprintDeHistorial(id int64, numero int, estado EstadoSprint) Sprint {
	return Sprint{ID: id, ProyectoID: 1, Numero: numero, Estado: estado}
}

func numerosDe(sprints []Sprint) []int {
	numeros := make([]int, 0, len(sprints))
	for _, s := range sprints {
		numeros = append(numeros, s.Numero)
	}
	return numeros
}

func TestSprintsFinalizados(t *testing.T) {
	casos := []struct {
		nombre   string
		sprints  []Sprint
		esperado []int
	}{
		{
			nombre: "solo los Finalizado, sin Activo ni Pendiente",
			sprints: []Sprint{
				sprintDeHistorial(1, 1, SprintFinalizado),
				sprintDeHistorial(2, 2, SprintActivo),
				sprintDeHistorial(3, 3, SprintPendiente),
			},
			esperado: []int{1},
		},
		{
			nombre: "ordenados por número aunque lleguen desordenados",
			sprints: []Sprint{
				sprintDeHistorial(3, 3, SprintFinalizado),
				sprintDeHistorial(1, 1, SprintFinalizado),
				sprintDeHistorial(2, 2, SprintFinalizado),
			},
			esperado: []int{1, 2, 3},
		},
		{
			nombre:   "sin finalizados",
			sprints:  []Sprint{sprintDeHistorial(1, 1, SprintActivo)},
			esperado: []int{},
		},
		{
			nombre:   "entrada vacía",
			sprints:  nil,
			esperado: []int{},
		},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			obtenidos := SprintsFinalizados(c.sprints)

			if obtenidos == nil {
				t.Fatal("se esperaba un slice vacío, nunca nil")
			}
			if numeros := numerosDe(obtenidos); !slices.Equal(numeros, c.esperado) {
				t.Errorf("se esperaban los Sprints %v, se obtuvieron %v", c.esperado, numeros)
			}
		})
	}
}

func TestValidarSprintFinalizado(t *testing.T) {
	casos := []struct {
		estado   EstadoSprint
		esperado error
	}{
		{SprintFinalizado, nil},
		{SprintActivo, ErrSprintNoFinalizado},
		{SprintPendiente, ErrSprintNoFinalizado},
	}
	for _, c := range casos {
		t.Run(string(c.estado), func(t *testing.T) {
			err := ValidarSprintFinalizado(sprintDeHistorial(1, 1, c.estado))

			if !errors.Is(err, c.esperado) {
				t.Errorf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
		})
	}
}

func idPtr(id int64) *int64 { return &id }

func spPtr(sp int) *int { return &sp }

func defectoDeHistorial(id, detectadoEn int64, resueltoEn *int64) DefectoDetalle {
	return DefectoDetalle{
		Defecto:        Defecto{ID: id, ProyectoID: 1, HistoriaID: 1, SprintDeteccionID: detectadoEn, SprintResolucionID: resueltoEn},
		TituloHistoria: "Login",
	}
}

func TestArmarDetalleSprint_MinutosYHorasPorHistoria(t *testing.T) {
	base := SprintHistorial{Sprint: sprintDeHistorial(10, 1, SprintFinalizado)}
	historias := []HistoriaComprometida{
		{ID: 3, Titulo: "Login", EstimacionSP: spPtr(5), Completada: true},
		{ID: 1, Titulo: "Reportes", EstimacionSP: spPtr(3)},
		{ID: 2, Titulo: "Sin horas"},
	}
	minutos := map[int64]int{3: 90, 1: 50}

	detalle := ArmarDetalleSprint(base, historias, minutos, nil)

	esperadas := []struct {
		id      int64
		minutos int
		horas   float64
	}{{3, 90, 1.5}, {1, 50, 0.83}, {2, 0, 0}}
	if len(detalle.Historias) != len(esperadas) {
		t.Fatalf("se esperaban %d historias, se obtuvieron %d", len(esperadas), len(detalle.Historias))
	}
	for i, e := range esperadas {
		h := detalle.Historias[i]
		if h.ID != e.id || h.Minutos != e.minutos || h.Horas != e.horas {
			t.Errorf("posición %d: se esperaba {id %d, %d min, %v h}, se obtuvo {id %d, %d min, %v h}",
				i, e.id, e.minutos, e.horas, h.ID, h.Minutos, h.Horas)
		}
	}
	if detalle.Historias[0].Titulo != "Login" || !detalle.Historias[0].Completada || *detalle.Historias[0].EstimacionSP != 5 {
		t.Errorf("se debían conservar los datos de la historia, se obtuvo %+v", detalle.Historias[0])
	}
	if detalle.Historias[2].EstimacionSP != nil {
		t.Errorf("una historia sin estimar debía seguir sin estimar, se obtuvo %v", *detalle.Historias[2].EstimacionSP)
	}
	if detalle.Sprint.ID != 10 {
		t.Errorf("se esperaba el Sprint 10 en el detalle, se obtuvo %d", detalle.Sprint.ID)
	}
}

func TestArmarDetalleSprint_MarcasDeLosDefectos(t *testing.T) {
	base := SprintHistorial{Sprint: sprintDeHistorial(10, 1, SprintFinalizado)}
	defectos := []DefectoDetalle{
		defectoDeHistorial(1, 10, nil),       // detectado acá, pendiente
		defectoDeHistorial(2, 10, idPtr(20)), // detectado acá, resuelto en otro
		defectoDeHistorial(3, 5, idPtr(10)),  // detectado en otro, resuelto acá
		defectoDeHistorial(4, 10, idPtr(10)), // detectado y resuelto acá
	}

	detalle := ArmarDetalleSprint(base, nil, nil, defectos)

	esperadas := []struct {
		id                  int64
		detectado, resuelto bool
	}{{1, true, false}, {2, true, false}, {3, false, true}, {4, true, true}}
	if len(detalle.Defectos) != len(esperadas) {
		t.Fatalf("se esperaban %d defectos (cada uno una sola vez), se obtuvieron %d", len(esperadas), len(detalle.Defectos))
	}
	for i, e := range esperadas {
		d := detalle.Defectos[i]
		if d.ID != e.id || d.DetectadoEnSprint != e.detectado || d.ResueltoEnSprint != e.resuelto {
			t.Errorf("posición %d: se esperaba {id %d, detectado %v, resuelto %v}, se obtuvo {id %d, detectado %v, resuelto %v}",
				i, e.id, e.detectado, e.resuelto, d.ID, d.DetectadoEnSprint, d.ResueltoEnSprint)
		}
	}
	if detalle.Defectos[0].TituloHistoria != "Login" {
		t.Errorf("se debía conservar el título de la historia del defecto, se obtuvo %q", detalle.Defectos[0].TituloHistoria)
	}
}

func TestArmarDetalleSprint_SinHistoriasNiDefectosDevuelveSlicesVacios(t *testing.T) {
	detalle := ArmarDetalleSprint(SprintHistorial{Sprint: sprintDeHistorial(10, 1, SprintFinalizado)}, nil, nil, nil)

	if detalle.Historias == nil || len(detalle.Historias) != 0 {
		t.Errorf("se esperaban historias vacías no nil, se obtuvo %#v", detalle.Historias)
	}
	if detalle.Defectos == nil || len(detalle.Defectos) != 0 {
		t.Errorf("se esperaban defectos vacíos no nil, se obtuvo %#v", detalle.Defectos)
	}
}
