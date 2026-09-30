package domain

import "time"

// Project es el contenedor de todo lo demás: Sprints, historias, esfuerzo.
// Los campos de fechas/estado (HU-13) se agregan sobre esta misma struct,
// no se cambia el nombre ni se remueven los que ya hay.
type Project struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}
