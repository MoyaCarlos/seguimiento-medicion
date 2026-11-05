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
func (s *GenerarReporte) EjecutarSprint(_ context.Context, _ int64, _ time.Time) (domain.ReporteSprint, error) {
	return domain.ReporteSprint{}, nil
}

// EjecutarProyecto arma el reporte completo de un proyecto.
func (s *GenerarReporte) EjecutarProyecto(_ context.Context, _ int64, _ time.Time) (domain.ReporteProyecto, error) {
	return domain.ReporteProyecto{}, nil
}
