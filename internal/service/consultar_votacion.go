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
	Diferencias *domain.Diferencias
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
	regla      domain.ReglaConsenso
}

func NewConsultarVotacion(historias repository.HistoriaRepository, proyectos repository.ProjectRepository, votaciones repository.VotacionRepository, regla domain.ReglaConsenso) *ConsultarVotacion {
	return &ConsultarVotacion{historias: historias, proyectos: proyectos, votaciones: votaciones, regla: regla}
}

// Ejecutar devuelve a los integrantes activos del proyecto, más quien tenga un
// voto en la ronda actual aunque después lo hayan dado de baja (HU-16), ordenados
// por id: quien se sumó con la ronda abierta figura entre los que faltan votar,
// y la lista mostrada coincide con los votos del promedio.
func (s *ConsultarVotacion) Ejecutar(ctx context.Context, historiaID, consultanteID int64) (VistaVotacion, error) {
	historia, err := s.historias.ObtenerPorID(ctx, historiaID)
	if err != nil {
		return VistaVotacion{}, err
	}
	if err := verificarIntegrante(ctx, s.proyectos, historia.ProyectoID, consultanteID); err != nil {
		return VistaVotacion{}, err
	}
	ronda, err := s.votaciones.ObtenerRondaActual(ctx, historia.ID)
	if err != nil {
		return VistaVotacion{}, err
	}
	integrantes, err := s.participantes(ctx, historia.ProyectoID, ronda)
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
	if diferencias, ok := ronda.Diferencias(s.regla); ok {
		vista.Diferencias = &diferencias
	}
	return vista, nil
}

// participantes son los integrantes activos y quienes votaron en la ronda, con
// nombre y rol tomados del equipo histórico.
func (s *ConsultarVotacion) participantes(ctx context.Context, proyectoID int64, ronda domain.RondaVotacion) ([]domain.Member, error) {
	activos, err := s.proyectos.ListarIntegrantes(ctx, proyectoID)
	if err != nil {
		return nil, err
	}
	historicos, err := s.proyectos.ListarIntegrantesHistoricos(ctx, proyectoID)
	if err != nil {
		return nil, err
	}
	participantes := activos
	for _, m := range historicos {
		votoEnLaRonda := slices.ContainsFunc(ronda.Votos, func(v domain.Voto) bool { return v.UsuarioID == m.UserID })
		yaListado := slices.ContainsFunc(participantes, func(a domain.Member) bool { return a.UserID == m.UserID })
		if votoEnLaRonda && !yaListado {
			participantes = append(participantes, m)
		}
	}
	return participantes, nil
}
