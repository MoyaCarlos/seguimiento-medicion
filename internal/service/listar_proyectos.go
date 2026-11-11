package service

import (
	"context"
	"sort"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type listadorProyectos interface {
	ListarPorIntegrante(ctx context.Context, usuarioID int64) ([]domain.ProyectoConRol, error)
}

// ListarProyectos es «Mis proyectos» (HU-21): los proyectos que integra el
// usuario, con su rol en cada uno, ordenados por id.
type ListarProyectos struct {
	proyectos listadorProyectos
}

func NewListarProyectos(proyectos listadorProyectos) *ListarProyectos {
	return &ListarProyectos{proyectos: proyectos}
}

// Ejecutar devuelve los proyectos de usuarioID; un arreglo vacío (nunca un
// error de permisos) si no integra ninguno.
func (s *ListarProyectos) Ejecutar(ctx context.Context, usuarioID int64) ([]domain.ProyectoConRol, error) {
	proyectos, err := s.proyectos.ListarPorIntegrante(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	if proyectos == nil {
		proyectos = []domain.ProyectoConRol{}
	}
	sort.Slice(proyectos, func(i, j int) bool { return proyectos[i].Proyecto.ID < proyectos[j].Proyecto.ID })
	return proyectos, nil
}
