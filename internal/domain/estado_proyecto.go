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

// CalcularEstadoProyecto deriva el estado del proyecto. Un proyecto cerrado por
// el Scrum Master es "Finalizado" sin mirar nada más (#300). Si está abierto,
// los Sprints tienen prioridad sobre las fechas: un Sprint Activo mantiene el
// proyecto "En curso"; un Finalizado junto a un Pendiente también (FR-016).
// Con todos los Sprints Finalizado, el proyecto recién es "Finalizado" cuando
// su fecha de fin ya pasó (#173); sin fecha de fin o con fecha futura sigue
// "En curso". Sin ejecución iniciada, deciden las fechas del proyecto,
// comparadas por día calendario (los bordes [inicio, fin] son inclusivos).
func CalcularEstadoProyecto(p Project, sprints []Sprint, ahora time.Time) EstadoProyecto {
	if p.FechaCierre != nil {
		return ProyectoFinalizado
	}
	ahora = inicioDelDiaCivil(ahora)
	for _, s := range sprints {
		if s.Estado == SprintActivo {
			return ProyectoEnCurso
		}
	}
	if len(sprints) > 0 && todosFinalizados(sprints) {
		if p.FechaFin != nil && ahora.After(*p.FechaFin) {
			return ProyectoFinalizado
		}
		return ProyectoEnCurso
	}
	if hayFinalizado(sprints) {
		return ProyectoEnCurso
	}
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

func hayFinalizado(sprints []Sprint) bool {
	for _, s := range sprints {
		if s.Estado == SprintFinalizado {
			return true
		}
	}
	return false
}
