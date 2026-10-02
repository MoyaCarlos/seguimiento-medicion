package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestCrearProyecto_Ejecutar_Exitosa(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	servicio := NewCrearProyecto(proyectos, usuarios)

	proyecto, err := servicio.Ejecutar(context.Background(), CrearProyectoInput{
		Nombre:      "Software Metrics & Estimation",
		Descripcion: "Sistema de estimación",
		Creador:     "Ana Valentina",
	})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if proyecto.ID == "" {
		t.Error("se esperaba ID asignado")
	}
	if proyectos.creates != 1 {
		t.Errorf("se esperaba 1 proyecto persistido, se crearon %d", proyectos.creates)
	}
	if len(proyectos.miembros) != 1 {
		t.Fatalf("se esperaba 1 vinculación, se obtuvieron %d", len(proyectos.miembros))
	}
	if proyectos.miembros[0].Role != domain.RolScrumMaster {
		t.Errorf("se esperaba rol scrum_master, se obtuvo %q", proyectos.miembros[0].Role)
	}
	if proyectos.miembros[0].ProjectID != proyecto.ID {
		t.Errorf("la vinculación apunta a otro proyecto: %+v", proyectos.miembros[0])
	}
}

func TestCrearProyecto_Ejecutar_NoPersisteSiInvalido(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	servicio := NewCrearProyecto(proyectos, usuarios)

	_, err := servicio.Ejecutar(context.Background(), CrearProyectoInput{
		Nombre:  "   ",
		Creador: "Ana",
	})
	if err == nil {
		t.Fatal("se esperaba error de validación")
	}
	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "nombre" {
		t.Fatalf("se esperaba ValidationError de nombre, se obtuvo %v", err)
	}
	if proyectos.creates != 0 || len(proyectos.miembros) != 0 {
		t.Error("no se debía persistir nada con datos inválidos")
	}
}

func TestCrearProyecto_Ejecutar_CreadorObligatorio(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	servicio := NewCrearProyecto(proyectos, usuarios)

	_, err := servicio.Ejecutar(context.Background(), CrearProyectoInput{
		Nombre:  "Proyecto",
		Creador: "",
	})
	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "nombre" {
		t.Fatalf("se esperaba ValidationError de nombre del creador, se obtuvo %v", err)
	}
	if proyectos.creates != 0 {
		t.Error("no se debía persistir el proyecto sin creador válido")
	}
}

func TestCrearProyecto_Ejecutar_ReutilizaCreadorExistente(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	servicio := NewCrearProyecto(proyectos, usuarios)

	for i := 0; i < 2; i++ {
		if _, err := servicio.Ejecutar(context.Background(), CrearProyectoInput{
			Nombre:  "Proyecto",
			Creador: "ana",
		}); err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
	}
	if usuarios.creates != 1 {
		t.Errorf("se esperaba un único usuario reutilizado, se crearon %d", usuarios.creates)
	}
}
