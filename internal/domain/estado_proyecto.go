package domain

import "time"

// EstadoProyecto representa la etapa del ciclo de vida de un proyecto.
// Es un valor derivado: se calcula on-demand y no se persiste.
type EstadoProyecto string

const (
	ProyectoPlanificado EstadoProyecto = "Planificado"
	ProyectoEnCurso     EstadoProyecto = "En curso"
	ProyectoFinalizado  EstadoProyecto = "Finalizado"
)

// CalcularEstadoProyecto deriva el estado del proyecto dando prioridad a los
// Sprints sobre las fechas: un Sprint Activo mantiene el proyecto "En curso";
// sin ejecución iniciada, deciden las fechas del proyecto, comparadas por día
// calendario (los bordes [inicio, fin] son inclusivos).
func CalcularEstadoProyecto(p Project, sprints []Sprint, ahora time.Time) EstadoProyecto {
	for _, s := range sprints {
		if s.Estado == SprintActivo {
			return ProyectoEnCurso
		}
	}
	if len(sprints) > 0 && todosFinalizados(sprints) {
		return ProyectoFinalizado
	}
	ahora = inicioDelDiaCivil(ahora)
	if p.FechaInicio == nil || ahora.Before(*p.FechaInicio) {
		return ProyectoPlanificado
	}
	if p.FechaFin != nil && ahora.After(*p.FechaFin) {
		return ProyectoFinalizado
	}
	return ProyectoEnCurso
}

// inicioDelDiaCivil reduce un instante a la medianoche UTC de su día civil en
// la zona horaria local del servidor, para comparar por día calendario contra
// las fechas del proyecto (que HU-04 guarda a las 00:00 UTC).
func inicioDelDiaCivil(ahora time.Time) time.Time {
	local := ahora.Local()
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.UTC)
}

func todosFinalizados(sprints []Sprint) bool {
	for _, s := range sprints {
		if s.Estado != SprintFinalizado {
			return false
		}
	}
	return true
}
