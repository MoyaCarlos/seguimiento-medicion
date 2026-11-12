package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

const (
	proyectoPermisos  int64 = 1
	scrumMasterID     int64 = 10
	productBuilderID  int64 = 20
	ajenoAlProyectoID int64 = 30
)

func proyectoConEquipo() *fakeProjectRepository {
	return &fakeProjectRepository{miembros: []domain.Membership{
		{ProjectID: proyectoPermisos, UserID: scrumMasterID, Role: domain.RolScrumMaster},
		{ProjectID: proyectoPermisos, UserID: productBuilderID, Role: domain.RolProductBuilder},
	}}
}

func TestExigirIntegrante(t *testing.T) {
	casos := []struct {
		nombre    string
		usuarioID int64
		esperado  error
	}{
		{"el Scrum Master integra el proyecto", scrumMasterID, nil},
		{"un Product Builder integra el proyecto", productBuilderID, nil},
		{"quien no integra el proyecto", ajenoAlProyectoID, domain.ErrNoEsIntegrante},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			err := exigirIntegrante(context.Background(), proyectoConEquipo(), proyectoPermisos, c.usuarioID)
			if !errors.Is(err, c.esperado) {
				t.Errorf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
		})
	}
}

func TestExigirScrumMaster(t *testing.T) {
	casos := []struct {
		nombre    string
		usuarioID int64
		esperado  error
	}{
		{"el Scrum Master puede", scrumMasterID, nil},
		{"un Product Builder no puede", productBuilderID, domain.ErrNoEsScrumMaster},
		{"quien no integra el proyecto recibe primero el error de pertenencia", ajenoAlProyectoID, domain.ErrNoEsIntegrante},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			err := exigirScrumMaster(context.Background(), proyectoConEquipo(), proyectoPermisos, c.usuarioID)
			if !errors.Is(err, c.esperado) {
				t.Errorf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
		})
	}
}

func TestExigirPermisosPropagaElErrorDelRepositorio(t *testing.T) {
	fallo := errors.New("fallo de la base")
	repo := &fakeProjectRepository{err: fallo}
	if err := exigirIntegrante(context.Background(), repo, proyectoPermisos, scrumMasterID); !errors.Is(err, fallo) {
		t.Errorf("exigirIntegrante: se esperaba el error del repositorio, se obtuvo %v", err)
	}
	if err := exigirScrumMaster(context.Background(), repo, proyectoPermisos, scrumMasterID); !errors.Is(err, fallo) {
		t.Errorf("exigirScrumMaster: se esperaba el error del repositorio, se obtuvo %v", err)
	}
}

func TestMensajesDePermisos(t *testing.T) {
	if got, want := domain.ErrNoEsIntegrante.Error(), "No sos integrante de este proyecto"; got != want {
		t.Errorf("ErrNoEsIntegrante: se esperaba %q, se obtuvo %q", want, got)
	}
	if got, want := domain.ErrNoEsScrumMaster.Error(), "Solo el Scrum Master puede realizar esta acción"; got != want {
		t.Errorf("ErrNoEsScrumMaster: se esperaba %q, se obtuvo %q", want, got)
	}
}
