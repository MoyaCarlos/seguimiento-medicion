package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// IniciarSprintInput son los datos de entrada del caso de uso.
type IniciarSprintInput struct {
	SprintID    int64
	SprintGoal  string
	FechaInicio time.Time
	FechaFin    time.Time
}

// IniciarSprint pasa un Sprint Pendiente a Activo, garantizando que el
// proyecto no tenga otro Sprint Activo.
type IniciarSprint struct {
	sprints repository.SprintRepository
}

func NewIniciarSprint(sprints repository.SprintRepository) *IniciarSprint {
	return &IniciarSprint{sprints: sprints}
}

func (s *IniciarSprint) Ejecutar(ctx context.Context, input IniciarSprintInput) (domain.Sprint, error) {
	sprint, err := s.sprints.ObtenerPorID(ctx, input.SprintID)
	if err != nil {
		return domain.Sprint{}, err
	}

	// ponytail: chequeo y escritura sin transacción; dos inicios simultáneos
	// podrían pasar los dos. Alcanza para un equipo; si hiciera falta, índice
	// único parcial en SQLite sobre (proyecto_id) WHERE estado = 'Activo'.
	delProyecto, err := s.sprints.ListarPorProyecto(ctx, sprint.ProyectoID)
	if err != nil {
		return domain.Sprint{}, err
	}
	for _, otro := range delProyecto {
		if otro.ID != sprint.ID && otro.Estado == domain.SprintActivo {
			return domain.Sprint{}, domain.ValidationError{Campo: "estado", Mensaje: "el proyecto ya tiene un Sprint activo"}
		}
	}

	if err := sprint.Iniciar(input.SprintGoal, input.FechaInicio, input.FechaFin); err != nil {
		return domain.Sprint{}, err
	}
	if err := s.sprints.Actualizar(ctx, sprint); err != nil {
		return domain.Sprint{}, err
	}
	return sprint, nil
}
