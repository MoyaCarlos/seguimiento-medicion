package domain

import "sort"

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
