package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type stubListadorProyectos struct {
	proyectos map[int64][]domain.ProyectoConRol
	pedidoDe  int64
	err       error
}

func (s *stubListadorProyectos) ListarPorIntegrante(_ context.Context, usuarioID int64) ([]domain.ProyectoConRol, error) {
	s.pedidoDe = usuarioID
	return s.proyectos[usuarioID], s.err
}

// HU-21, escenario 5: «Mis proyectos» devuelve solo los de quien consulta, con
// su rol en cada uno, ordenados por id.
func TestListarProyectos_SoloLosDelUsuarioConSuRol(t *testing.T) {
	stub := &stubListadorProyectos{proyectos: map[int64][]domain.ProyectoConRol{
		1: {
			{Proyecto: domain.Project{ID: 3, Nombre: "Tres"}, Rol: domain.RolProductBuilder},
			{Proyecto: domain.Project{ID: 1, Nombre: "Uno"}, Rol: domain.RolScrumMaster},
		},
		2: {{Proyecto: domain.Project{ID: 2, Nombre: "Dos"}, Rol: domain.RolScrumMaster}},
	}}

	proyectos, err := NewListarProyectos(stub).Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if stub.pedidoDe != 1 {
		t.Errorf("se debían pedir los proyectos del usuario 1, se pidieron los de %d", stub.pedidoDe)
	}
	if len(proyectos) != 2 {
		t.Fatalf("se esperaban 2 proyectos, se obtuvieron %+v", proyectos)
	}
	if proyectos[0].Proyecto.ID != 1 || proyectos[0].Rol != domain.RolScrumMaster {
		t.Errorf("primero: se esperaba el proyecto 1 como scrum_master, se obtuvo %+v", proyectos[0])
	}
	if proyectos[1].Proyecto.ID != 3 || proyectos[1].Rol != domain.RolProductBuilder {
		t.Errorf("segundo: se esperaba el proyecto 3 como product_builder, se obtuvo %+v", proyectos[1])
	}
}

// Escenario 9: quien no integra ningún proyecto recibe un arreglo vacío, no un error.
func TestListarProyectos_SinProyectosDevuelveArregloNoNulo(t *testing.T) {
	proyectos, err := NewListarProyectos(&stubListadorProyectos{}).Ejecutar(context.Background(), 5)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if proyectos == nil || len(proyectos) != 0 {
		t.Fatalf("se esperaba un arreglo vacío (no nil), se obtuvo %#v", proyectos)
	}
}

func TestListarProyectos_PropagaErrorDePersistencia(t *testing.T) {
	errPersistencia := errors.New("fallo de persistencia")
	servicio := NewListarProyectos(&stubListadorProyectos{err: errPersistencia})

	_, err := servicio.Ejecutar(context.Background(), 1)

	if !errors.Is(err, errPersistencia) {
		t.Fatalf("se esperaba el error de persistencia, se obtuvo %v", err)
	}
}
