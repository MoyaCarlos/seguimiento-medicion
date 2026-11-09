package domain

// HistoriaReporte es una historia tal como aparece en el reporte PDF (HU-10).
type HistoriaReporte struct {
	ID           int64
	Titulo       string
	Prioridad    Prioridad
	EstimacionSP *int
	Estado       Estado
	Completada   bool
}

// DefectoReporte es un defecto con los números de Sprint que muestra el
// reporte. SprintResolucion es nil mientras el defecto esté Pendiente (HU-10).
type DefectoReporte struct {
	Descripcion      string
	Severidad        Severidad
	Estado           EstadoDefecto
	Historia         string
	SprintDeteccion  int
	SprintResolucion *int
}

// SprintReporte reúne las secciones de un Sprint finalizado (HU-10).
type SprintReporte struct {
	Sprint             Sprint
	Planificadas       []HistoriaReporte
	Metricas           Metricas
	DefectosDetectados []DefectoReporte
	DefectosResueltos  []DefectoReporte
}

// SprintActivoReporte es el Sprint Activo con sus historias; no suma a los
// totales del proyecto (HU-10).
type SprintActivoReporte struct {
	Sprint    Sprint
	Historias []HistoriaReporte
}

// ReporteSprint es el contenido del reporte de un Sprint finalizado (HU-10).
type ReporteSprint struct {
	Proyecto Project
	Sprint   SprintReporte
}

// ReporteProyecto es el contenido del reporte de un proyecto (HU-10). No se
// persiste: se arma en cada solicitud.
type ReporteProyecto struct {
	Proyecto              Project
	Estado                EstadoProyecto
	Integrantes           []Member
	ProductBacklog        []HistoriaReporte
	Sprints               []SprintReporte
	SprintActivo          *SprintActivoReporte
	Metricas              Metricas
	HorasPorStoryPoint    *float64
	SinSprintsFinalizados bool
}
