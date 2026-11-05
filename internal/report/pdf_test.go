package report

import (
	"bytes"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func fechaReporte(dia int) *time.Time {
	fecha := time.Date(2026, 10, dia, 0, 0, 0, 0, time.UTC)
	return &fecha
}

func reporteDeSprintDePrueba() domain.ReporteSprint {
	proyecto := domain.Project{ID: 1, Nombre: "Proyecto de prueba", Descripcion: "Descripción"}
	sprint := domain.Sprint{
		ID: 1, ProyectoID: 1, Numero: 1, SprintGoal: "Cerrar el carrito",
		FechaInicio: *fechaReporte(1), FechaFin: *fechaReporte(15), Estado: domain.SprintFinalizado,
	}
	resolucion := 1
	return domain.ReporteSprint{
		Proyecto: proyecto,
		Sprint: domain.SprintReporte{
			Sprint: sprint,
			Planificadas: []domain.HistoriaReporte{
				{ID: 10, Titulo: "Carrito de compras", Prioridad: domain.PrioridadMust,
					EstimacionSP: entero(5), Estado: domain.EstadoCompletada, Completada: true},
				{ID: 11, Titulo: "Pago", Prioridad: domain.PrioridadCould,
					EstimacionSP: entero(3), Estado: domain.EstadoNueva, Completada: false},
			},
			Metricas: domain.Metricas{
				SPPlanificados: 8, SPCompletados: 5, Velocidad: 5, PorcentajeCompletadas: 50,
				HorasEstimadas: decimal(32), HorasReales: 30, DesviacionHoras: decimal(-2),
				DesviacionPorcentaje: decimal(-6.25), DefectosDetectados: 1, DefectosResueltos: 1,
			},
			DefectosDetectados: []domain.DefectoReporte{{
				Descripcion: "El total no se actualiza", Severidad: domain.SeveridadAlta,
				Estado: domain.DefectoResuelto, Historia: "Carrito de compras",
				SprintDeteccion: 1, SprintResolucion: &resolucion,
			}},
			DefectosResueltos: []domain.DefectoReporte{{
				Descripcion: "El total no se actualiza", Severidad: domain.SeveridadAlta,
				Estado: domain.DefectoResuelto, Historia: "Carrito de compras",
				SprintDeteccion: 1, SprintResolucion: &resolucion,
			}},
		},
	}
}

func reporteDeProyectoDePrueba(sinSprintsFinalizados bool) domain.ReporteProyecto {
	sprint := reporteDeSprintDePrueba()
	activo := domain.Sprint{
		ID: 2, ProyectoID: 1, Numero: 2, SprintGoal: "Pulir la entrega",
		FechaInicio: *fechaReporte(16), FechaFin: *fechaReporte(30), Estado: domain.SprintActivo,
	}
	reporte := domain.ReporteProyecto{
		Proyecto:       domain.Project{ID: 1, Nombre: "Proyecto de prueba", Descripcion: "Descripción"},
		Estado:         domain.ProyectoEnCurso,
		Integrantes:    []domain.Member{{Nombre: "Ana", Role: domain.RolProductBuilder}},
		ProductBacklog: []domain.HistoriaReporte{{ID: 20, Titulo: "Mejoras futuras", Prioridad: domain.PrioridadShould, EstimacionSP: entero(8)}},
		Sprints:        []domain.SprintReporte{sprint.Sprint},
		SprintActivo: &domain.SprintActivoReporte{
			Sprint: activo,
			Historias: []domain.HistoriaReporte{{
				ID: 30, Titulo: "Ajustes finales", Prioridad: domain.PrioridadMust, EstimacionSP: entero(2), Estado: domain.EstadoEnProgreso,
			}},
		},
		Metricas:            sprint.Sprint.Metricas,
		HorasPorStoryPoint:  decimal(4),
		SinSprintsFinalizados: sinSprintsFinalizados,
	}
	if sinSprintsFinalizados {
		reporte.Sprints = []domain.SprintReporte{}
		reporte.Metricas = domain.Metricas{}
	}
	return reporte
}

func entero(v int) *int { return &v }

func decimal(v float64) *float64 { return &v }

func TestGenerarSprint_DevuelveUnPDFValido(t *testing.T) {
	contenido, err := GenerarSprint(reporteDeSprintDePrueba())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(contenido) == 0 {
		t.Fatal("se esperaba un PDF no vacío")
	}
	if !bytes.HasPrefix(contenido, []byte("%PDF-")) {
		t.Errorf("se esperaba un PDF, se obtuvo %q", contenido[:min(8, len(contenido))])
	}
	if !bytes.Contains(contenido, []byte("%%EOF")) {
		t.Error("el PDF no tiene el cierre EOF")
	}
}

func TestGenerarSprint_ConMuchasHistoriasPagina(t *testing.T) {
	reporte := reporteDeSprintDePrueba()
	for i := 0; i < 60; i++ {
		reporte.Sprint.Planificadas = append(reporte.Sprint.Planificadas, domain.HistoriaReporte{
			ID: int64(100 + i), Titulo: "Historia larga con un título que puede ocupar más de una línea en el reporte",
			Prioridad: domain.PrioridadShould, EstimacionSP: entero(3), Estado: domain.EstadoCompletada, Completada: true,
		})
	}

	contenido, err := GenerarSprint(reporte)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if !bytes.HasPrefix(contenido, []byte("%PDF-")) {
		t.Error("se esperaba un PDF válido")
	}
}

func TestGenerarProyecto_DevuelveUnPDFValido(t *testing.T) {
	contenido, err := GenerarProyecto(reporteDeProyectoDePrueba(false))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(contenido) == 0 || !bytes.HasPrefix(contenido, []byte("%PDF-")) || !bytes.Contains(contenido, []byte("%%EOF")) {
		t.Error("se esperaba un PDF válido y no vacío")
	}
}

func TestGenerarProyecto_SinSprintsFinalizados(t *testing.T) {
	contenido, err := GenerarProyecto(reporteDeProyectoDePrueba(true))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(contenido) == 0 || !bytes.HasPrefix(contenido, []byte("%PDF-")) {
		t.Error("se esperaba un PDF válido y no vacío")
	}
}
