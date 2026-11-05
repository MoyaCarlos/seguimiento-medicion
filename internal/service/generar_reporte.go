package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// GenerarReporte arma el contenido de los reportes PDF de un Sprint finalizado
// y de un proyecto (HU-10). Es solo lectura: no escribe nada.
type GenerarReporte struct {
	proyectos   repository.ProjectRepository
	sprints     repository.SprintRepository
	backlog     repository.ProductBacklogRepository
	historias   repository.HistoriasDeSprintRepository
	arrastradas repository.HistoriasArrastradasRepository
	datos       repository.DatosMetricasRepository
	defectos    repository.DefectoRepository
}

func NewGenerarReporte(
	proyectos repository.ProjectRepository,
	sprints repository.SprintRepository,
	backlog repository.ProductBacklogRepository,
	historias repository.HistoriasDeSprintRepository,
	arrastradas repository.HistoriasArrastradasRepository,
	datos repository.DatosMetricasRepository,
	defectos repository.DefectoRepository,
) *GenerarReporte {
	return &GenerarReporte{
		proyectos: proyectos, sprints: sprints, backlog: backlog,
		historias: historias, arrastradas: arrastradas, datos: datos, defectos: defectos,
	}
}

// EjecutarSprint arma el reporte de un Sprint finalizado.
func (s *GenerarReporte) EjecutarSprint(ctx context.Context, sprintID int64, _ time.Time) (domain.ReporteSprint, error) {
	sprint, err := s.sprints.ObtenerPorID(ctx, sprintID)
	if err != nil {
		return domain.ReporteSprint{}, err
	}
	if sprint.Estado != domain.SprintFinalizado {
		return domain.ReporteSprint{}, domain.ErrSprintNoFinalizado
	}
	proyecto, err := s.proyectos.ObtenerPorID(ctx, sprint.ProyectoID)
	if err != nil {
		return domain.ReporteSprint{}, err
	}
	sprints, err := s.sprints.ListarPorProyecto(ctx, sprint.ProyectoID)
	if err != nil {
		return domain.ReporteSprint{}, err
	}
	datos, err := s.datos.ListarDatosSprintsFinalizados(ctx, sprint.ProyectoID)
	if err != nil {
		return domain.ReporteSprint{}, err
	}
	defectos, err := s.defectos.ListarPorProyecto(ctx, sprint.ProyectoID)
	if err != nil {
		return domain.ReporteSprint{}, err
	}

	metricas := domain.CalcularMetricas(datos, proyecto.HorasPorStoryPoint)
	delSprint := buscarMetricasDeSprint(metricas, sprint.ID)
	reporte, err := s.armarSprint(ctx, sprint, numerosDeSprint(sprints), delSprint, defectos)
	if err != nil {
		return domain.ReporteSprint{}, err
	}
	return domain.ReporteSprint{Proyecto: proyecto, Sprint: reporte}, nil
}

func (s *GenerarReporte) armarSprint(
	ctx context.Context,
	sprint domain.Sprint,
	numeros map[int64]int,
	metricas domain.Metricas,
	defectos []domain.DefectoDetalle,
) (domain.SprintReporte, error) {
	completadas, err := s.historias.ListarPorSprint(ctx, sprint.ID)
	if err != nil {
		return domain.SprintReporte{}, err
	}
	arrastradas, err := s.arrastradas.ListarArrastradasPorSprint(ctx, sprint.ID)
	if err != nil {
		return domain.SprintReporte{}, err
	}

	planificadas := make([]domain.HistoriaReporte, 0, len(completadas)+len(arrastradas))
	for _, historia := range completadas {
		planificadas = append(planificadas, aHistoriaReporte(historia, historia.Estado == domain.EstadoCompletada))
	}
	for _, historia := range arrastradas {
		planificadas = append(planificadas, aHistoriaReporte(historia, false))
	}

	detectados, resueltos := defectosDelSprint(defectos, sprint.ID, numeros)
	return domain.SprintReporte{
		Sprint:             sprint,
		Planificadas:       planificadas,
		Metricas:           metricas,
		DefectosDetectados: detectados,
		DefectosResueltos:  resueltos,
	}, nil
}

func aHistoriaReporte(historia domain.BacklogItem, completada bool) domain.HistoriaReporte {
	return domain.HistoriaReporte{
		ID: historia.ID, Titulo: historia.Titulo, Prioridad: historia.Prioridad,
		EstimacionSP: historia.EstimacionSP, Estado: historia.Estado, Completada: completada,
	}
}

func numerosDeSprint(sprints []domain.Sprint) map[int64]int {
	numeros := make(map[int64]int, len(sprints))
	for _, sprint := range sprints {
		numeros[sprint.ID] = sprint.Numero
	}
	return numeros
}

func buscarMetricasDeSprint(metricas domain.MetricasProyecto, sprintID int64) domain.Metricas {
	for _, sprint := range metricas.Sprints {
		if sprint.SprintID == sprintID {
			return sprint.Metricas
		}
	}
	return domain.Metricas{}
}

func defectosDelSprint(defectos []domain.DefectoDetalle, sprintID int64, numeros map[int64]int) ([]domain.DefectoReporte, []domain.DefectoReporte) {
	detectados := make([]domain.DefectoReporte, 0)
	resueltos := make([]domain.DefectoReporte, 0)
	for _, defecto := range defectos {
		if defecto.SprintDeteccionID == sprintID {
			detectados = append(detectados, aDefectoReporte(defecto, numeros))
		}
		if defecto.SprintResolucionID != nil && *defecto.SprintResolucionID == sprintID {
			resueltos = append(resueltos, aDefectoReporte(defecto, numeros))
		}
	}
	return detectados, resueltos
}

func aDefectoReporte(defecto domain.DefectoDetalle, numeros map[int64]int) domain.DefectoReporte {
	var resolucion *int
	if defecto.SprintResolucionID != nil {
		numero := numeros[*defecto.SprintResolucionID]
		resolucion = &numero
	}
	return domain.DefectoReporte{
		Descripcion:      defecto.Descripcion,
		Severidad:        defecto.Severidad,
		Estado:           defecto.Estado,
		Historia:         defecto.TituloHistoria,
		SprintDeteccion:  numeros[defecto.SprintDeteccionID],
		SprintResolucion: resolucion,
	}
}

// EjecutarProyecto arma el reporte completo de un proyecto.
func (s *GenerarReporte) EjecutarProyecto(_ context.Context, _ int64, _ time.Time) (domain.ReporteProyecto, error) {
	return domain.ReporteProyecto{}, nil
}
