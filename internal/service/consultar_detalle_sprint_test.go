package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type fuentesDelDetalle struct {
	proyectos stubBuscadorProyectos
	sprints   *fakeObtenedorSprint
	metricas  *fakeCalculadorMetricas
	historial *fakeHistorialSprintRepository
}

func nuevasFuentesDelDetalle() fuentesDelDetalle {
	return fuentesDelDetalle{
		sprints: &fakeObtenedorSprint{sprints: map[int64]domain.Sprint{
			10: sprintParaHistorial(10, 1, domain.SprintFinalizado),
			20: sprintParaHistorial(20, 2, domain.SprintActivo),
			30: sprintParaHistorial(30, 3, domain.SprintPendiente),
			40: {ID: 40, ProyectoID: 2, Numero: 1, Estado: domain.SprintFinalizado},
		}},
		metricas: &fakeCalculadorMetricas{resultado: domain.MetricasProyecto{Sprints: []domain.MetricasDeSprint{
			metricasDeSprint(10, 1, 5),
		}}},
		historial: &fakeHistorialSprintRepository{
			historias: []domain.HistoriaComprometida{{ID: 1, Titulo: "Login", Completada: true}},
			minutos:   map[int64]int{1: 90},
			defectos:  []domain.DefectoDetalle{{Defecto: domain.Defecto{ID: 7, SprintDeteccionID: 10}}},
		},
	}
}

func (f fuentesDelDetalle) consultar(proyectoID, sprintID int64) (domain.DetalleSprint, error) {
	return NewConsultarDetalleSprint(f.proyectos, f.sprints, f.metricas, f.historial).Ejecutar(context.Background(), proyectoID, sprintID)
}

func TestConsultarDetalleSprint_ArmaElDetalleDeUnSprintFinalizado(t *testing.T) {
	detalle, err := nuevasFuentesDelDetalle().consultar(1, 10)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if detalle.Sprint.ID != 10 || detalle.Metricas.SPCompletados != 5 {
		t.Errorf("se esperaba el Sprint 10 con sus métricas de HU-03, se obtuvo Sprint %d con %d SP completados",
			detalle.Sprint.ID, detalle.Metricas.SPCompletados)
	}
	if len(detalle.Historias) != 1 || detalle.Historias[0].Minutos != 90 || detalle.Historias[0].Horas != 1.5 {
		t.Errorf("se esperaba la historia con 90 minutos y 1.5 horas, se obtuvo %+v", detalle.Historias)
	}
	if len(detalle.Defectos) != 1 || !detalle.Defectos[0].DetectadoEnSprint {
		t.Errorf("se esperaba el defecto marcado como detectado en el Sprint, se obtuvo %+v", detalle.Defectos)
	}
}

func TestConsultarDetalleSprint_ErroresEnOrden(t *testing.T) {
	casos := []struct {
		nombre     string
		preparar   func(*fuentesDelDetalle)
		proyectoID int64
		sprintID   int64
		esperado   error
	}{
		{"proyecto inexistente", func(f *fuentesDelDetalle) { f.proyectos.err = domain.ErrProyectoNoEncontrado }, 99, 10, domain.ErrProyectoNoEncontrado},
		{"Sprint inexistente", nil, 1, 999, domain.ErrSprintNoEncontrado},
		{"Sprint de otro proyecto", nil, 1, 40, domain.ErrSprintNoEncontrado},
		{"Sprint Activo", nil, 1, 20, domain.ErrSprintNoFinalizado},
		{"Sprint Pendiente", nil, 1, 30, domain.ErrSprintNoFinalizado},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			f := nuevasFuentesDelDetalle()
			if c.preparar != nil {
				c.preparar(&f)
			}

			_, err := f.consultar(c.proyectoID, c.sprintID)

			if !errors.Is(err, c.esperado) {
				t.Fatalf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if f.historial.lecturas != 0 {
				t.Errorf("no se debían leer historias, horas ni defectos; lecturas: %d", f.historial.lecturas)
			}
		})
	}
}

func TestConsultarDetalleSprint_SprintFinalizadoSinMetricasEsErrorInterno(t *testing.T) {
	f := nuevasFuentesDelDetalle()
	f.metricas.resultado = domain.MetricasProyecto{}

	_, err := f.consultar(1, 10)

	if err == nil {
		t.Fatal("se esperaba un error por la inconsistencia entre fuentes")
	}
	var verr domain.ValidationError
	if errors.As(err, &verr) || errors.Is(err, domain.ErrSprintNoEncontrado) || errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Errorf("se esperaba un error interno, se obtuvo %v", err)
	}
}

func TestConsultarDetalleSprint_PropagaErroresDeLecturaSinDatosParciales(t *testing.T) {
	falla := errors.New("falla de lectura")
	casos := []struct {
		nombre   string
		preparar func(*fuentesDelDetalle)
	}{
		{"del Sprint", func(f *fuentesDelDetalle) { f.sprints.err = falla }},
		{"de las métricas", func(f *fuentesDelDetalle) { f.metricas.err = falla }},
		{"de las historias", func(f *fuentesDelDetalle) { f.historial.errHistoria = falla }},
		{"de los minutos", func(f *fuentesDelDetalle) { f.historial.errMinutos = falla }},
		{"de los defectos", func(f *fuentesDelDetalle) { f.historial.errDefectos = falla }},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			f := nuevasFuentesDelDetalle()
			c.preparar(&f)

			detalle, err := f.consultar(1, 10)

			if !errors.Is(err, falla) {
				t.Fatalf("se esperaba la falla de lectura, se obtuvo %v", err)
			}
			if detalle.Sprint.ID != 0 || detalle.Historias != nil || detalle.Defectos != nil {
				t.Errorf("no se debían devolver datos parciales, se obtuvo %+v", detalle)
			}
		})
	}
}
