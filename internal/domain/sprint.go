package domain

import (
	"strings"
	"time"
)

// Sprint es una iteración de trabajo de un proyecto.
type Sprint struct {
	ID          int64
	ProyectoID  int64
	SprintGoal  string
	FechaInicio time.Time
	FechaFin    time.Time
	Estado      EstadoSprint
}

// NewSprint crea un Sprint en estado Pendiente. El Goal y las fechas se
// definen recién al iniciarlo.
func NewSprint(proyectoID int64) (Sprint, error) {
	if proyectoID <= 0 {
		return Sprint{}, ValidationError{Campo: "proyecto_id", Mensaje: "el proyecto es obligatorio"}
	}
	return Sprint{ProyectoID: proyectoID, Estado: SprintPendiente}, nil
}

// Iniciar pasa el Sprint de Pendiente a Activo con su Goal y plazos.
func (s *Sprint) Iniciar(goal string, inicio, fin time.Time) error {
	if s.Estado != SprintPendiente {
		return ValidationError{Campo: "estado", Mensaje: "solo se puede iniciar un Sprint pendiente"}
	}
	goal = strings.TrimSpace(goal)
	if goal == "" {
		return ValidationError{Campo: "sprint_goal", Mensaje: "el Sprint Goal es obligatorio"}
	}
	if inicio.IsZero() {
		return ValidationError{Campo: "fecha_inicio", Mensaje: "la fecha de inicio es obligatoria"}
	}
	if !fin.After(inicio) {
		return ValidationError{Campo: "fecha_fin", Mensaje: "la fecha de fin debe ser posterior a la de inicio"}
	}
	s.SprintGoal, s.FechaInicio, s.FechaFin, s.Estado = goal, inicio, fin, SprintActivo
	return nil
}

// Cerrar pasa el Sprint de Activo a Finalizado.
func (s *Sprint) Cerrar() error {
	if s.Estado != SprintActivo {
		return ValidationError{Campo: "estado", Mensaje: "solo se puede cerrar un Sprint activo"}
	}
	s.Estado = SprintFinalizado
	return nil
}
