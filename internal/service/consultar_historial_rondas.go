package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// HistorialRondas son las rondas reveladas de una historia, ordenadas por número.
type HistorialRondas struct {
	HistoriaID int64
	Rondas     []RondaDelHistorial
}

// RondaDelHistorial es una ronda revelada con el voto de cada integrante.
type RondaDelHistorial struct {
	Numero     int
	ReveladaEn time.Time
	Votos      []VotoRevelado
}

// VotoRevelado es el voto de un integrante en una ronda revelada.
type VotoRevelado struct {
	UsuarioID int64
	Nombre    string
	Rol       domain.Role
	Valor     int
}

// ConsultarHistorialRondas arma el historial de rondas de una historia. No pide
// quién consulta: esos votos ya son visibles para todo el equipo.
type ConsultarHistorialRondas struct {
	historias  repository.HistoriaRepository
	proyectos  repository.ProjectRepository
	votaciones repository.VotacionRepository
}

func NewConsultarHistorialRondas(historias repository.HistoriaRepository, proyectos repository.ProjectRepository, votaciones repository.VotacionRepository) *ConsultarHistorialRondas {
	return &ConsultarHistorialRondas{historias: historias, proyectos: proyectos, votaciones: votaciones}
}

func (s *ConsultarHistorialRondas) Ejecutar(ctx context.Context, historiaID int64) (HistorialRondas, error) {
	historia, err := s.historias.ObtenerPorID(ctx, historiaID)
	if err != nil {
		return HistorialRondas{}, err
	}
	rondas, err := s.votaciones.ListarRondasReveladas(ctx, historia.ID)
	if err != nil {
		return HistorialRondas{}, err
	}
	integrantes, err := s.proyectos.ListarIntegrantes(ctx, historia.ProyectoID)
	if err != nil {
		return HistorialRondas{}, err
	}
	porUsuario := make(map[int64]domain.Member, len(integrantes))
	for _, m := range integrantes {
		porUsuario[m.UserID] = m
	}

	historial := HistorialRondas{HistoriaID: historia.ID, Rondas: make([]RondaDelHistorial, 0, len(rondas))}
	for _, r := range rondas {
		ronda := RondaDelHistorial{Numero: r.Numero, ReveladaEn: *r.ReveladaEn, Votos: make([]VotoRevelado, 0, len(r.Votos))}
		for _, v := range r.Votos {
			m := porUsuario[v.UsuarioID]
			ronda.Votos = append(ronda.Votos, VotoRevelado{UsuarioID: v.UsuarioID, Nombre: m.Nombre, Rol: m.Role, Valor: v.Valor})
		}
		historial.Rondas = append(historial.Rondas, ronda)
	}
	return historial, nil
}
