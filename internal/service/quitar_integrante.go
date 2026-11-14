package service

import (
	"context"
	"slices"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// QuitarIntegranteInput son los datos para quitar a un integrante del equipo.
type QuitarIntegranteInput struct {
	ProyectoID int64
	// UsuarioID es la persona que se quita.
	UsuarioID int64
	// ActorID es quien pide la baja (la sesión); debe ser el Scrum Master.
	ActorID int64
}

// QuitarIntegrante da de baja a un integrante del equipo conservando sus
// aportes (HU-16).
type QuitarIntegrante struct {
	proyectos repository.ProjectRepository
	reloj     func() time.Time
}

func NewQuitarIntegrante(proyectos repository.ProjectRepository, reloj func() time.Time) *QuitarIntegrante {
	return &QuitarIntegrante{proyectos: proyectos, reloj: reloj}
}

// Ejecutar verifica, en el orden de FR-013, que el proyecto exista (404), que
// quien actúa sea su Scrum Master (403) con el proyecto abierto (409), que la
// persona sea integrante activo (404) y que no sea el Scrum Master (409). El
// repositorio vuelve a verificar las dos últimas en la misma escritura.
func (s *QuitarIntegrante) Ejecutar(ctx context.Context, input QuitarIntegranteInput) error {
	if _, err := s.proyectos.ObtenerPorID(ctx, input.ProyectoID); err != nil {
		return err
	}
	if err := exigirScrumMaster(ctx, s.proyectos, input.ProyectoID, input.ActorID); err != nil {
		return err
	}
	integrantes, err := s.proyectos.ListarIntegrantes(ctx, input.ProyectoID)
	if err != nil {
		return err
	}
	i := slices.IndexFunc(integrantes, func(m domain.Member) bool { return m.UserID == input.UsuarioID })
	if i < 0 {
		return domain.ErrIntegranteNoEncontrado
	}
	if err := integrantes[i].ValidarBaja(); err != nil {
		return err
	}
	return s.proyectos.DarDeBajaIntegrante(ctx, input.ProyectoID, input.UsuarioID, s.reloj())
}
