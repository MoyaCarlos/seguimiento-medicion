package service

import (
	"context"
	"errors"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// RegistrarDefecto registra un defecto Pendiente con el Sprint activo del
// proyecto como Sprint de detección (HU-08).
type RegistrarDefecto struct {
	sprints   repository.SprintRepository
	historias repository.HistoriaRepository
	defectos  repository.DefectoRepository
}

// RegistrarDefectoInput son los datos de un defecto. El estado y los Sprints
// los asigna el sistema.
type RegistrarDefectoInput struct {
	ProyectoID  int64
	HistoriaID  int64
	Descripcion string
	Severidad   string
}

func NewRegistrarDefecto(sprints repository.SprintRepository, historias repository.HistoriaRepository, defectos repository.DefectoRepository) *RegistrarDefecto {
	return &RegistrarDefecto{sprints: sprints, historias: historias, defectos: defectos}
}

// Ejecutar valida en el orden de la spec: Sprint activo, historia relacionada
// y, en el dominio, descripción y severidad. La existencia del proyecto la
// verifica el handler antes de leer el cuerpo.
func (s *RegistrarDefecto) Ejecutar(ctx context.Context, in RegistrarDefectoInput) (domain.DefectoDetalle, error) {
	sprints, err := s.sprints.ListarPorProyecto(ctx, in.ProyectoID)
	if err != nil {
		return domain.DefectoDetalle{}, err
	}
	activo, hay := domain.SprintActivoDe(sprints)
	if !hay {
		return domain.DefectoDetalle{}, domain.ErrSinSprintActivoParaDetectar
	}
	historia, err := s.historiaRelacionada(ctx, in.ProyectoID, in.HistoriaID)
	if err != nil {
		return domain.DefectoDetalle{}, err
	}
	defecto, err := domain.NewDefecto(in.ProyectoID, historia.ID, activo, in.Descripcion, in.Severidad)
	if err != nil {
		return domain.DefectoDetalle{}, err
	}
	defecto, err = s.defectos.Registrar(ctx, defecto)
	if err != nil {
		return domain.DefectoDetalle{}, err
	}
	return domain.DefectoDetalle{Defecto: defecto, TituloHistoria: historia.Titulo}, nil
}

// historiaRelacionada devuelve domain.ErrHistoriaRelacionadaAjena si la
// historia no existe o es de otro proyecto. Su Sprint y su estado no importan.
func (s *RegistrarDefecto) historiaRelacionada(ctx context.Context, proyectoID, historiaID int64) (domain.BacklogItem, error) {
	if historiaID == 0 {
		return domain.BacklogItem{}, domain.ErrHistoriaRelacionadaObligatoria
	}
	historia, err := s.historias.ObtenerPorID(ctx, historiaID)
	if errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		return domain.BacklogItem{}, domain.ErrHistoriaRelacionadaAjena
	}
	if err != nil {
		return domain.BacklogItem{}, err
	}
	if historia.ProyectoID != proyectoID {
		return domain.BacklogItem{}, domain.ErrHistoriaRelacionadaAjena
	}
	return historia, nil
}
