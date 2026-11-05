package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

func punteroAEntero(v int) *int { return &v }

func punteroADecimal(v float64) *float64 { return &v }

func fechaDeReporte(dia int) time.Time {
	return time.Date(2026, 10, dia, 0, 0, 0, 0, time.UTC)
}

type fakeDatosMetricasRepository struct {
	datos     []domain.DatosMetricasSprint
	err       error
	llamadas  int
}

func (f *fakeDatosMetricasRepository) ListarDatosSprintsFinalizados(_ context.Context, _ int64) ([]domain.DatosMetricasSprint, error) {
	f.llamadas++
	if f.err != nil {
		return nil, f.err
	}
	return f.datos, nil
}

type fakeArrastradasRepository struct {
	porSprint map[int64][]domain.BacklogItem
	err       error
	llamadas  int
}

func (f *fakeArrastradasRepository) ListarArrastradasPorSprint(_ context.Context, sprintID int64) ([]domain.BacklogItem, error) {
	f.llamadas++
	if f.err != nil {
		return nil, f.err
	}
	historias := f.porSprint[sprintID]
	if historias == nil {
		return []domain.BacklogItem{}, nil
	}
	return historias, nil
}

func sprintDeReporte(t *testing.T, repo *repository.SprintRepositoryEnMemoria, proyectoID int64, goal string, estado domain.EstadoSprint) domain.Sprint {
	t.Helper()
	ctx := context.Background()
	s, err := domain.NewSprint(proyectoID)
	if err != nil {
		t.Fatalf("no se pudo construir el Sprint: %v", err)
	}
	if s, err = repo.Guardar(ctx, s); err != nil {
		t.Fatalf("no se pudo guardar el Sprint: %v", err)
	}
	if estado == domain.SprintPendiente {
		return s
	}
	if err := s.Iniciar(goal, fechaDeReporte(1), fechaDeReporte(15)); err != nil {
		t.Fatalf("no se pudo iniciar el Sprint: %v", err)
	}
	if estado == domain.SprintFinalizado {
		if err := s.Cerrar(); err != nil {
			t.Fatalf("no se pudo cerrar el Sprint: %v", err)
		}
	}
	if err := repo.Actualizar(ctx, s); err != nil {
		t.Fatalf("no se pudo actualizar el Sprint: %v", err)
	}
	return s
}

func historiaDeReporte(id int64, proyectoID int64, sprintID *int64, titulo string, prioridad domain.Prioridad, estado domain.Estado, sp int) domain.BacklogItem {
	return domain.BacklogItem{
		ID: id, ProyectoID: proyectoID, Titulo: titulo, Descripcion: "Descripción", Prioridad: prioridad,
		Estado: estado, EstimacionSP: punteroAEntero(sp), SprintID: sprintID,
	}
}

func defectoDeReporte(id, proyectoID, historiaID, deteccion int64, resolucion *int64, descripcion string) domain.Defecto {
	return domain.Defecto{
		ID: id, ProyectoID: proyectoID, HistoriaID: historiaID, SprintDeteccionID: deteccion,
		Descripcion: descripcion, Severidad: domain.SeveridadAlta, Estado: domain.DefectoPendiente,
		SprintResolucionID: resolucion,
	}
}

func nuevoGenerarReporteDePrueba(
	proyectos *fakeProjectRepository,
	sprints *repository.SprintRepositoryEnMemoria,
	backlog *fakeProductBacklogRepository,
	historias *fakeHistoriasDeSprintRepository,
	arrastradas *fakeArrastradasRepository,
	datos *fakeDatosMetricasRepository,
	defectos *fakeDefectoRepository,
) *GenerarReporte {
	return NewGenerarReporte(proyectos, sprints, backlog, historias, arrastradas, datos, defectos)
}

func TestGenerarReporte_EjecutarSprint_DevuelveTodasLasSecciones(t *testing.T) {
	ctx := context.Background()
	proyecto := domain.Project{ID: 1, Nombre: "Proyecto", HorasPorStoryPoint: punteroADecimal(4)}
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: proyecto}}
	sprints := repository.NewSprintRepositoryEnMemoria()
	finalizado := sprintDeReporte(t, sprints, 1, "Cerrar el carrito", domain.SprintFinalizado)

	historias := &fakeHistoriasDeSprintRepository{porSprint: map[int64][]domain.BacklogItem{
		finalizado.ID: {historiaDeReporte(10, 1, &finalizado.ID, "Carrito de compras", domain.PrioridadMust, domain.EstadoCompletada, 5)},
	}}
	arrastradas := &fakeArrastradasRepository{porSprint: map[int64][]domain.BacklogItem{
		finalizado.ID: {historiaDeReporte(11, 1, nil, "Pago", domain.PrioridadCould, domain.EstadoNueva, 3)},
	}}
	datos := &fakeDatosMetricasRepository{datos: []domain.DatosMetricasSprint{{
		SprintID: finalizado.ID, Numero: 1, HistoriasCompletadas: 1, HistoriasArrastradas: 1,
		SPCompletados: 5, SPArrastrados: 3, MinutosReales: 360, DefectosDetectados: 1, DefectosResueltos: 1,
	}}}
	idResolucion := finalizado.ID
	defectos := &fakeDefectoRepository{
		defectos: []domain.Defecto{
			defectoDeReporte(1, 1, 10, finalizado.ID, &idResolucion, "El total no se actualiza"),
		},
		titulos: map[int64]string{10: "Carrito de compras"},
	}

	sut := nuevoGenerarReporteDePrueba(proyectos, sprints, &fakeProductBacklogRepository{}, historias, arrastradas, datos, defectos)
	reporte, err := sut.EjecutarSprint(ctx, finalizado.ID, fechaDeReporte(20))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if reporte.Sprint.Sprint.ID != finalizado.ID {
		t.Errorf("se esperaba el Sprint %d, se obtuvo %d", finalizado.ID, reporte.Sprint.Sprint.ID)
	}
	if len(reporte.Sprint.Planificadas) != 2 {
		t.Fatalf("se esperaban 2 historias planificadas, se obtuvieron %d", len(reporte.Sprint.Planificadas))
	}
	if reporte.Sprint.Planificadas[0].Titulo != "Carrito de compras" || !reporte.Sprint.Planificadas[0].Completada {
		t.Errorf("se esperaba la historia completada primero, se obtuvo %+v", reporte.Sprint.Planificadas[0])
	}
	if reporte.Sprint.Planificadas[1].Titulo != "Pago" || reporte.Sprint.Planificadas[1].Completada {
		t.Errorf("se esperaba la historia arrastrada como no completada, se obtuvo %+v", reporte.Sprint.Planificadas[1])
	}
	if reporte.Sprint.Planificadas[1].EstimacionSP == nil || *reporte.Sprint.Planificadas[1].EstimacionSP != 3 {
		t.Errorf("se esperaba la estimación congelada de 3, se obtuvo %v", reporte.Sprint.Planificadas[1].EstimacionSP)
	}
	if reporte.Sprint.Metricas.HorasReales != 6 {
		t.Errorf("se esperaban 6 horas reales, se obtuvieron %v", reporte.Sprint.Metricas.HorasReales)
	}
	if reporte.Sprint.Metricas.SPPlanificados != 8 || reporte.Sprint.Metricas.SPCompletados != 5 {
		t.Errorf("métricas inesperadas: %+v", reporte.Sprint.Metricas)
	}
	if len(reporte.Sprint.DefectosDetectados) != 1 || len(reporte.Sprint.DefectosResueltos) != 1 {
		t.Errorf("se esperaban un defecto detectado y uno resuelto, se obtuvo %d y %d",
			len(reporte.Sprint.DefectosDetectados), len(reporte.Sprint.DefectosResueltos))
	}
	if reporte.Sprint.DefectosDetectados[0].Historia != "Carrito de compras" {
		t.Errorf("se esperaba el título de la historia, se obtuvo %q", reporte.Sprint.DefectosDetectados[0].Historia)
	}
}

func TestGenerarReporte_EjecutarSprint_Inexistente(t *testing.T) {
	sut := nuevoGenerarReporteDePrueba(
		&fakeProjectRepository{},
		repository.NewSprintRepositoryEnMemoria(),
		&fakeProductBacklogRepository{},
		&fakeHistoriasDeSprintRepository{},
		&fakeArrastradasRepository{},
		&fakeDatosMetricasRepository{},
		&fakeDefectoRepository{},
	)

	_, err := sut.EjecutarSprint(context.Background(), 999, fechaDeReporte(20))

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("se esperaba ErrSprintNoEncontrado, se obtuvo %v", err)
	}
}

func TestGenerarReporte_EjecutarSprint_NoFinalizado(t *testing.T) {
	for _, estado := range []domain.EstadoSprint{domain.SprintPendiente, domain.SprintActivo} {
		t.Run(string(estado), func(t *testing.T) {
			proyecto := domain.Project{ID: 1, Nombre: "Proyecto"}
			proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: proyecto}}
			sprints := repository.NewSprintRepositoryEnMemoria()
			sprint := sprintDeReporte(t, sprints, 1, "Objetivo", estado)
			historias := &fakeHistoriasDeSprintRepository{}
			arrastradas := &fakeArrastradasRepository{}
			datos := &fakeDatosMetricasRepository{}

			sut := nuevoGenerarReporteDePrueba(proyectos, sprints, &fakeProductBacklogRepository{}, historias, arrastradas, datos, &fakeDefectoRepository{})
			_, err := sut.EjecutarSprint(context.Background(), sprint.ID, fechaDeReporte(20))

			if !errors.Is(err, domain.ErrSprintNoFinalizado) {
				t.Fatalf("se esperaba ErrSprintNoFinalizado, se obtuvo %v", err)
			}
			if historias.llamadas != 0 || arrastradas.llamadas != 0 || datos.llamadas != 0 {
				t.Errorf("no se debía leer nada más: historias=%d arrastradas=%d datos=%d",
					historias.llamadas, arrastradas.llamadas, datos.llamadas)
			}
		})
	}
}

func TestGenerarReporte_EjecutarSprint_SinHorasNiDefectos(t *testing.T) {
	proyecto := domain.Project{ID: 1, Nombre: "Proyecto", HorasPorStoryPoint: punteroADecimal(4)}
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: proyecto}}
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sprintDeReporte(t, sprints, 1, "Vacío", domain.SprintFinalizado)

	sut := nuevoGenerarReporteDePrueba(
		proyectos, sprints, &fakeProductBacklogRepository{},
		&fakeHistoriasDeSprintRepository{}, &fakeArrastradasRepository{},
		&fakeDatosMetricasRepository{datos: []domain.DatosMetricasSprint{{SprintID: sprint.ID, Numero: 1}}},
		&fakeDefectoRepository{},
	)

	reporte, err := sut.EjecutarSprint(context.Background(), sprint.ID, fechaDeReporte(20))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(reporte.Sprint.Planificadas) != 0 || reporte.Sprint.Planificadas == nil {
		t.Errorf("se esperaba un slice vacío no nil, se obtuvo %#v", reporte.Sprint.Planificadas)
	}
	if reporte.Sprint.Metricas.HorasReales != 0 || reporte.Sprint.Metricas.SPPlanificados != 0 {
		t.Errorf("se esperaban métricas en cero, se obtuvo %+v", reporte.Sprint.Metricas)
	}
	if len(reporte.Sprint.DefectosDetectados) != 0 || reporte.Sprint.DefectosDetectados == nil {
		t.Errorf("se esperaba un slice vacío no nil de defectos, se obtuvo %#v", reporte.Sprint.DefectosDetectados)
	}
}

func TestGenerarReporte_EjecutarSprint_DefectoDeOtroSprintSoloResuelto(t *testing.T) {
	proyecto := domain.Project{ID: 1, Nombre: "Proyecto"}
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: proyecto}}
	sprints := repository.NewSprintRepositoryEnMemoria()
	anterior := sprintDeReporte(t, sprints, 1, "Anterior", domain.SprintFinalizado)
	actual := sprintDeReporte(t, sprints, 1, "Actual", domain.SprintFinalizado)
	idActual := actual.ID
	defectos := &fakeDefectoRepository{
		defectos: []domain.Defecto{defectoDeReporte(1, 1, 10, anterior.ID, &idActual, "Reaparece")},
		titulos:  map[int64]string{10: "Historia"},
	}

	sut := nuevoGenerarReporteDePrueba(
		proyectos, sprints, &fakeProductBacklogRepository{},
		&fakeHistoriasDeSprintRepository{}, &fakeArrastradasRepository{},
		&fakeDatosMetricasRepository{}, defectos,
	)

	reporte, err := sut.EjecutarSprint(context.Background(), actual.ID, fechaDeReporte(20))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(reporte.Sprint.DefectosDetectados) != 0 {
		t.Errorf("no debía aparecer como detectado en el Sprint actual: %+v", reporte.Sprint.DefectosDetectados)
	}
	if len(reporte.Sprint.DefectosResueltos) != 1 {
		t.Fatalf("se esperaba el defecto como resuelto, se obtuvo %d", len(reporte.Sprint.DefectosResueltos))
	}
	if reporte.Sprint.DefectosResueltos[0].SprintDeteccion != anterior.Numero {
		t.Errorf("se esperaba el Sprint de detección %d, se obtuvo %d",
			anterior.Numero, reporte.Sprint.DefectosResueltos[0].SprintDeteccion)
	}
}

func TestGenerarReporte_EjecutarSprint_PropagaErrores(t *testing.T) {
	proyecto := domain.Project{ID: 1, Nombre: "Proyecto"}
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: proyecto}}
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sprintDeReporte(t, sprints, 1, "Objetivo", domain.SprintFinalizado)
	falla := errors.New("falla de lectura")

	sut := nuevoGenerarReporteDePrueba(
		proyectos, sprints, &fakeProductBacklogRepository{},
		&fakeHistoriasDeSprintRepository{}, &fakeArrastradasRepository{},
		&fakeDatosMetricasRepository{err: falla}, &fakeDefectoRepository{},
	)

	_, err := sut.EjecutarSprint(context.Background(), sprint.ID, fechaDeReporte(20))
	if !errors.Is(err, falla) {
		t.Fatalf("se esperaba el error del repositorio, se obtuvo %v", err)
	}
}
