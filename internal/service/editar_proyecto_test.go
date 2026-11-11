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
	crear := NewCrearProyecto(proyectos)
	creado, err := crear.Ejecutar(context.Background(), CrearProyectoInput{Nombre: "Original", Descripcion: "desc", CreadorID: 1})
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
	creado, _ := NewCrearProyecto(proyectos).Ejecutar(context.Background(), CrearProyectoInput{Nombre: "Original", CreadorID: 1})

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

func TestEditarProyecto_HorasPorStoryPoint(t *testing.T) {
	cuatro, cinco := "4", "5"
	guardadas := 2.5
	casos := []struct {
		nombre   string
		enviadas bool
		valor    *string
		esperado *float64
	}{
		{nombre: "sin enviar conserva el valor", esperado: &guardadas},
		{nombre: "enviadas sin valor las borra", enviadas: true},
		{nombre: "enviadas con un valor lo guarda", enviadas: true, valor: &cuatro, esperado: float64Ptr(4)},
		{nombre: "cambiar el valor lo reemplaza", enviadas: true, valor: &cinco, esperado: float64Ptr(5)},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1, Nombre: "Original", HorasPorStoryPoint: &guardadas}}}

			editado, err := NewEditarProyecto(proyectos).Ejecutar(context.Background(), EditarProyectoInput{
				ID: 1, Nombre: "Corregido", HorasPorSPEnviadas: c.enviadas, HorasPorSP: c.valor,
			})

			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			for nombre, obtenido := range map[string]*float64{"devuelto": editado.HorasPorStoryPoint, "guardado": proyectos.proyectos[1].HorasPorStoryPoint} {
				if (obtenido == nil) != (c.esperado == nil) || (obtenido != nil && *obtenido != *c.esperado) {
					t.Errorf("%s: se esperaba %v, se obtuvo %v", nombre, c.esperado, obtenido)
				}
			}
		})
	}
}

func TestEditarProyecto_HorasPorStoryPointInvalidasNoActualiza(t *testing.T) {
	cero := "0"
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1, Nombre: "Original"}}}

	_, err := NewEditarProyecto(proyectos).Ejecutar(context.Background(), EditarProyectoInput{
		ID: 1, Nombre: "Corregido", HorasPorSPEnviadas: true, HorasPorSP: &cero,
	})

	if !errors.Is(err, domain.ErrHorasPorStoryPointInvalidas) {
		t.Fatalf("se esperaba ErrHorasPorStoryPointInvalidas, se obtuvo %v", err)
	}
	if proyectos.updates != 0 {
		t.Errorf("no se debía actualizar, hubo %d actualizaciones", proyectos.updates)
	}
}

func TestEditarProyecto_LasHorasPorStoryPointSeValidanAlFinal(t *testing.T) {
	cero := "0"
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1, Nombre: "Original"}}}

	_, err := NewEditarProyecto(proyectos).Ejecutar(context.Background(), EditarProyectoInput{
		ID: 1, Nombre: "  ", HorasPorSPEnviadas: true, HorasPorSP: &cero,
	})

	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "nombre" {
		t.Fatalf("se esperaba el error del nombre, se obtuvo %v", err)
	}
}

func float64Ptr(valor float64) *float64 { return &valor }
