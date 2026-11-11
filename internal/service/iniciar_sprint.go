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
	// UsuarioID es quien inicia (la sesión, HU-21); debe ser el Scrum Master.
	UsuarioID int64
}

// IniciarSprint pasa un Sprint Pendiente a Activo, garantizando que el
// proyecto no tenga otro Sprint Activo.
type IniciarSprint struct {
	sprints     repository.SprintRepository
	integrantes repository.IntegrantesRepository
}

func NewIniciarSprint(sprints repository.SprintRepository, integrantes repository.IntegrantesRepository) *IniciarSprint {
	return &IniciarSprint{sprints: sprints, integrantes: integrantes}
}

// Autorizar devuelve el error de ErrSprintNoEncontrado (404) o de permisos
// (403) sin tocar nada: el handler lo llama antes de leer el cuerpo, así un
// Product Builder recibe 403 aunque mande fechas inválidas (HU-21, FR-006).
func (s *IniciarSprint) Autorizar(ctx context.Context, sprintID, usuarioID int64) error {
	_, err := s.sprintDelScrumMaster(ctx, sprintID, usuarioID)
	return err
}

func (s *IniciarSprint) sprintDelScrumMaster(ctx context.Context, sprintID, usuarioID int64) (domain.Sprint, error) {
	sprint, err := s.sprints.ObtenerPorID(ctx, sprintID)
	if err != nil {
		return domain.Sprint{}, err
	}
	if err := exigirScrumMaster(ctx, s.integrantes, sprint.ProyectoID, usuarioID); err != nil {
		return domain.Sprint{}, err
	}
	return sprint, nil
}

func (s *IniciarSprint) Ejecutar(ctx context.Context, input IniciarSprintInput) (domain.Sprint, error) {
	sprint, err := s.sprintDelScrumMaster(ctx, input.SprintID, input.UsuarioID)
	if err != nil {
		return domain.Sprint{}, err
	}

	// Estos chequeos dan el mensaje; la carrera la cubre la base: el índice
	// único parcial impide un segundo Activo y Actualizar solo inicia un
	// Sprint que sigue Pendiente (#158).
	delProyecto, err := s.sprints.ListarPorProyecto(ctx, sprint.ProyectoID)
	if err != nil {
		return domain.Sprint{}, err
	}
	for _, otro := range delProyecto {
		if otro.ID != sprint.ID && otro.Estado == domain.SprintActivo {
			return domain.Sprint{}, domain.ErrYaHaySprintActivo
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
