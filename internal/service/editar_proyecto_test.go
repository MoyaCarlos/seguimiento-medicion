package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestEditarProyecto_Exitosa(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	crear := NewCrearProyecto(proyectos, usuarios)
	creado, err := crear.Ejecutar(context.Background(), CrearProyectoInput{Nombre: "Original", Descripcion: "desc", Creador: "Ana"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	fin := time.Date(2026, 12, 15, 0, 0, 0, 0, time.UTC)
	servicio := NewEditarProyecto(proyectos)
	editado, err := servicio.Ejecutar(context.Background(), EditarProyectoInput{
		ID:          creado.ID,
		Nombre:      "Corregido",
		Descripcion: "nueva desc",
		FechaFin:    &fin,
	})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if editado.Nombre != "Corregido" || editado.Descripcion != "nueva desc" {
		t.Errorf("edición no aplicada: %+v", editado)
	}
	if editado.ID != creado.ID {
		t.Errorf("no se conservó el ID")
	}
	if proyectos.updates != 1 {
		t.Errorf("se esperaba 1 actualización, se obtuvieron %d", proyectos.updates)
	}
}

func TestEditarProyecto_Inexistente(t *testing.T) {
	servicio := NewEditarProyecto(&fakeProjectRepository{})
	_, err := servicio.Ejecutar(context.Background(), EditarProyectoInput{ID: 0, Nombre: "X"})
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
}

func TestEditarProyecto_NoPersisteSiInvalido(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	creado, _ := NewCrearProyecto(proyectos, usuarios).Ejecutar(context.Background(), CrearProyectoInput{Nombre: "Original", Creador: "Ana"})

	servicio := NewEditarProyecto(proyectos)
	_, err := servicio.Ejecutar(context.Background(), EditarProyectoInput{ID: creado.ID, Nombre: "   "})
	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "nombre" {
		t.Fatalf("se esperaba ValidationError de nombre, se obtuvo %v", err)
	}
	if proyectos.updates != 0 {
		t.Error("no se debía persistir una edición inválida")
	}
}
