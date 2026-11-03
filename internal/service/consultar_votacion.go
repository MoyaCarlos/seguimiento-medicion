package service

import (
	"cmp"
	"context"
	"slices"
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

// Ejecutar devuelve a todos los integrantes del proyecto ordenados por id, así
// quien se sumó con la ronda abierta figura entre los que faltan votar.
func (s *ConsultarVotacion) Ejecutar(ctx context.Context, historiaID, consultanteID int64) (VistaVotacion, error) {
	historia, err := s.historias.ObtenerPorID(ctx, historiaID)
	if err != nil {
		return VistaVotacion{}, err
	}
	ronda, err := s.votaciones.ObtenerRondaActual(ctx, historia.ID)
	if err != nil {
		return VistaVotacion{}, err
	}
	integrantes, err := s.proyectos.ListarIntegrantes(ctx, historia.ProyectoID)
	if err != nil {
		return VistaVotacion{}, err
	}

	vista := VistaVotacion{HistoriaID: historia.ID, Numero: ronda.Numero, ReveladaEn: ronda.ReveladaEn}
	for _, m := range slices.SortedFunc(slices.Values(integrantes), func(a, b domain.Member) int { return cmp.Compare(a.UserID, b.UserID) }) {
		voto := VotoDeIntegrante{UsuarioID: m.UserID, Nombre: m.Nombre, Rol: m.Role}
		if i := slices.IndexFunc(ronda.Votos, func(v domain.Voto) bool { return v.UsuarioID == m.UserID }); i >= 0 {
			voto.YaVoto = true
			if ronda.Revelada() || m.UserID == consultanteID {
				valor := ronda.Votos[i].Valor
				voto.Valor = &valor
			}
		}
		vista.Integrantes = append(vista.Integrantes, voto)
	}
	if promedio, ok := ronda.Promedio(); ok {
		vista.Promedio = &promedio
	}
	return vista, nil
}
