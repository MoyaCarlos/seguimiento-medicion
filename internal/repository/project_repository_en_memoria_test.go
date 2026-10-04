package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestProjectRepositoryEnMemoria_AgregarIntegranteProyectoInexistente(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	_, err := repo.AgregarIntegrante(context.Background(), domain.Membership{
		ProjectID: 1, UserID: 1, Role: domain.RolScrumMaster,
	})
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
}

func TestProjectRepositoryEnMemoria_AgregarIntegranteUsuarioInexistente(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	guardado, err := repo.Guardar(context.Background(), domain.Project{Nombre: "Proyecto"})
	if err != nil {
		t.Fatalf("no se esperaba error al guardar: %v", err)
	}
	_, err = repo.AgregarIntegrante(context.Background(), domain.Membership{
		ProjectID: guardado.ID, UserID: 42, Role: domain.RolScrumMaster,
	})
	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Fatalf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
}
