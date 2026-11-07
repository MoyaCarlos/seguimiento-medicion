package domain

import "sort"

// HistoriaComprometida es una historia completada en un Sprint o arrastrada al
// cerrarlo (D2c). En una arrastrada, EstimacionSP es la foto que se guardó al
// cerrar; nil significa "sin estimar". Minutos y Horas son los imputados en ese
// Sprint.
type HistoriaComprometida struct {
	ID           int64
	Titulo       string
	EstimacionSP *int
	Completada   bool
	Minutos      int
	Horas        float64
}

// DefectoDeSprint es un defecto detectado o resuelto en el Sprint consultado.
type DefectoDeSprint struct {
	DefectoDetalle
	DetectadoEnSprint bool
	ResueltoEnSprint  bool
}

// SprintHistorial es un Sprint finalizado con sus métricas de HU-03.
type SprintHistorial struct {
	Sprint   Sprint
	Metricas Metricas
}

// DetalleSprint es un Sprint finalizado con sus historias comprometidas y sus
// defectos.
type DetalleSprint struct {
	SprintHistorial
	Historias []HistoriaComprometida
	Defectos  []DefectoDeSprint
}

// SprintsFinalizados devuelve los Sprints Finalizado ordenados por número.
func SprintsFinalizados(sprints []Sprint) []Sprint {
	finalizados := make([]Sprint, 0, len(sprints))
	for _, s := range sprints {
		if s.Estado == SprintFinalizado {
			finalizados = append(finalizados, s)
		}
	}
	sort.Slice(finalizados, func(i, j int) bool { return finalizados[i].Numero < finalizados[j].Numero })
	return finalizados
}

// ValidarSprintFinalizado rechaza un Sprint que todavía no está Finalizado.
func ValidarSprintFinalizado(sprint Sprint) error {
	if sprint.Estado != SprintFinalizado {
		return ErrSprintNoFinalizado
	}
	return nil
}

// ArmarDetalleSprint une a cada historia comprometida sus minutos y horas en el
// Sprint, y marca si cada defecto se detectó o se resolvió en él.
func ArmarDetalleSprint(base SprintHistorial, historias []HistoriaComprometida, minutos map[int64]int, defectos []DefectoDetalle) DetalleSprint {
	detalle := DetalleSprint{
		SprintHistorial: base,
		Historias:       make([]HistoriaComprometida, 0, len(historias)),
		Defectos:        make([]DefectoDeSprint, 0, len(defectos)),
	}
	for _, h := range historias {
		h.Minutos = minutos[h.ID]
		h.Horas = redondear(minutosAHoras(h.Minutos))
		detalle.Historias = append(detalle.Historias, h)
	}
	sprintID := base.Sprint.ID
	for _, d := range defectos {
		detalle.Defectos = append(detalle.Defectos, DefectoDeSprint{
			DefectoDetalle:    d,
			DetectadoEnSprint: d.SprintDeteccionID == sprintID,
			ResueltoEnSprint:  d.SprintResolucionID != nil && *d.SprintResolucionID == sprintID,
		})
	}
	return detalle
}
