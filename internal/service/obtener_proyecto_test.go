package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestObtenerProyecto_Ejecutar_Existente(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	crear := NewCrearProyecto(proyectos, &fakeUserRepository{})
	creado, err := crear.Ejecutar(context.Background(), CrearProyectoInput{Nombre: "Proyecto", Creador: "Ana"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	servicio := NewObtenerProyecto(proyectos)
	obtenido, err := servicio.Ejecutar(context.Background(), creado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if obtenido.ID != creado.ID || obtenido.Name != "Proyecto" {
		t.Errorf("proyecto inesperado: %+v", obtenido)
	}
}

func TestObtenerProyecto_Ejecutar_Inexistente(t *testing.T) {
	servicio := NewObtenerProyecto(&fakeProjectRepository{})

	_, err := servicio.Ejecutar(context.Background(), "no-existe")
	if !errors.Is(err, domain.ErrNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrNoEncontrado, se obtuvo %v", err)
	}
}
