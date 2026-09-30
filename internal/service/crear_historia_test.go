package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type fakeBacklogRepository struct {
	guardados []domain.BacklogItem
	err       error
}

func (f *fakeBacklogRepository) Guardar(_ context.Context, item domain.BacklogItem) (domain.BacklogItem, error) {
	if f.err != nil {
		return domain.BacklogItem{}, f.err
	}
	item.ID = int64(len(f.guardados) + 1)
	f.guardados = append(f.guardados, item)
	return item, nil
}

func intPtr(v int) *int { return &v }

func TestCrearHistoriaBacklog_Ejecutar_Exitosa(t *testing.T) {
	repo := &fakeBacklogRepository{}
	servicio := NewCrearHistoriaBacklog(repo)

	resultado, err := servicio.Ejecutar(context.Background(), CrearHistoriaInput{
		ProyectoID:  1,
		Titulo:      "Como usuario quiero iniciar sesión",
		Descripcion: "Autenticación con email y contraseña",
		Prioridad:   domain.PrioridadMust,
	})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if resultado.ID != 1 {
		t.Errorf("se esperaba ID 1, se obtuvo %d", resultado.ID)
	}
	if len(repo.guardados) != 1 {
		t.Fatalf("se esperaba 1 historia guardada, se obtuvieron %d", len(repo.guardados))
	}
	if repo.guardados[0].Estado != domain.EstadoNueva {
		t.Errorf("se esperaba estado %q", domain.EstadoNueva)
	}
}

func TestCrearHistoriaBacklog_Ejecutar_NoPersisteSiEsInvalida(t *testing.T) {
	repo := &fakeBacklogRepository{}
	servicio := NewCrearHistoriaBacklog(repo)

	_, err := servicio.Ejecutar(context.Background(), CrearHistoriaInput{
		ProyectoID:  1,
		Titulo:      "   ",
		Descripcion: "Descripción",
		Prioridad:   domain.PrioridadMust,
	})
	if err == nil {
		t.Fatal("se esperaba un error de validación")
	}
	var verr domain.ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("se esperaba ValidationError, se obtuvo %T", err)
	}
	if verr.Campo != "titulo" {
		t.Errorf("se esperaba campo titulo, se obtuvo %q", verr.Campo)
	}
	if len(repo.guardados) != 0 {
		t.Errorf("no se debía persistir ninguna historia, se guardaron %d", len(repo.guardados))
	}
}

func TestCrearHistoriaBacklog_Ejecutar_ValorNegocio(t *testing.T) {
	repo := &fakeBacklogRepository{}
	servicio := NewCrearHistoriaBacklog(repo)

	if _, err := servicio.Ejecutar(context.Background(), CrearHistoriaInput{
		ProyectoID:   1,
		Titulo:       "Con valor",
		Descripcion:  "Descripción",
		Prioridad:    domain.PrioridadShould,
		ValorNegocio: intPtr(13),
	}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if repo.guardados[0].ValorNegocio == nil || *repo.guardados[0].ValorNegocio != 13 {
		t.Errorf("se esperaba valor de negocio 13")
	}

	repo2 := &fakeBacklogRepository{}
	servicio2 := NewCrearHistoriaBacklog(repo2)
	if _, err := servicio2.Ejecutar(context.Background(), CrearHistoriaInput{
		ProyectoID:  1,
		Titulo:      "Sin valor",
		Descripcion: "Descripción",
		Prioridad:   domain.PrioridadCould,
	}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if repo2.guardados[0].ValorNegocio != nil {
		t.Errorf("se esperaba valor de negocio nil")
	}
}
