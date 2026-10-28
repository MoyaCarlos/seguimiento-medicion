package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// VistaVotacion es la ronda actual de una historia tal como la ve quien consulta.
type VistaVotacion struct {
	HistoriaID  int64
	Numero      int
	ReveladaEn  *time.Time
	Integrantes []VotoDeIntegrante
	Promedio    *domain.Promedio
}

// VotoDeIntegrante indica si un integrante ya votó. Valor es nil mientras la
// ronda está oculta, salvo para quien consulta, que ve su propio voto.
type VotoDeIntegrante struct {
	UsuarioID int64
	Nombre    string
	Rol       domain.Role
	YaVoto    bool
	Valor     *int
}

// ConsultarVotacion arma la vista de la ronda actual de una historia.
type ConsultarVotacion struct {
	historias  repository.HistoriaRepository
	proyectos  repository.ProjectRepository
	votaciones repository.VotacionRepository
}

func NewConsultarVotacion(historias repository.HistoriaRepository, proyectos repository.ProjectRepository, votaciones repository.VotacionRepository) *ConsultarVotacion {
	return &ConsultarVotacion{historias: historias, proyectos: proyectos, votaciones: votaciones}
}

func (s *ConsultarVotacion) Ejecutar(ctx context.Context, historiaID, consultanteID int64) (VistaVotacion, error) {
	return VistaVotacion{}, nil
}
