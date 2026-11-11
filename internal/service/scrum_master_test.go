package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// Escenarios 1 y 2 de HU-21: las acciones de gestión solo las hace el Scrum
// Master. Un Product Builder recibe ErrNoEsScrumMaster y quien no integra el
// proyecto ErrNoEsIntegrante, aunque los datos sean inválidos, y no cambia nada.

// proyectoConEquipoGuardado es el proyecto de proyectoConEquipo, también guardado.
func proyectoConEquipoGuardado() *fakeProjectRepository {
	p := proyectoConEquipo()
	p.proyectos = map[int64]domain.Project{proyectoPermisos: {ID: proyectoPermisos, Nombre: "Original"}}
	return p
}

var rechazosDeGestion = []struct {
	nombre    string
	usuarioID int64
	esperado  error
}{
	{"Product Builder", productBuilderID, domain.ErrNoEsScrumMaster},
	{"quien no integra el proyecto", ajenoAlProyectoID, domain.ErrNoEsIntegrante},
}

func TestGestion_CrearSprint(t *testing.T) {
	for _, c := range rechazosDeGestion {
		t.Run(c.nombre, func(t *testing.T) {
			sprints := repository.NewSprintRepositoryEnMemoria()
			_, err := NewCrearSprint(proyectoConEquipoGuardado(), sprints).Ejecutar(context.Background(), proyectoPermisos, c.usuarioID)
			if !errors.Is(err, c.esperado) {
				t.Fatalf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if lista, _ := sprints.ListarPorProyecto(context.Background(), proyectoPermisos); len(lista) != 0 {
				t.Error("no se debía crear el Sprint")
			}
		})
	}
	sprints := repository.NewSprintRepositoryEnMemoria()
	if _, err := NewCrearSprint(proyectoConEquipoGuardado(), sprints).Ejecutar(context.Background(), proyectoPermisos, scrumMasterID); err != nil {
		t.Fatalf("el Scrum Master debía poder crear el Sprint: %v", err)
	}
}

func TestGestion_IniciarSprint(t *testing.T) {
	for _, c := range rechazosDeGestion {
		t.Run(c.nombre, func(t *testing.T) {
			sprints := repository.NewSprintRepositoryEnMemoria()
			pendiente := sembrarSprintPendiente(t, sprints, proyectoPermisos)
			servicio := NewIniciarSprint(sprints, proyectoConEquipo())
			entrada := inputValido(pendiente.ID)
			entrada.UsuarioID = c.usuarioID
			entrada.FechaFin = entrada.FechaInicio.AddDate(0, 0, -1) // fechas inválidas: igual gana el 403

			if _, err := servicio.Ejecutar(context.Background(), entrada); !errors.Is(err, c.esperado) {
				t.Fatalf("Ejecutar: se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if err := servicio.Autorizar(context.Background(), pendiente.ID, c.usuarioID); !errors.Is(err, c.esperado) {
				t.Fatalf("Autorizar: se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if estadoPersistido(t, sprints, pendiente.ID) != domain.SprintPendiente {
				t.Error("el Sprint no debía iniciarse")
			}
		})
	}
	sprints := repository.NewSprintRepositoryEnMemoria()
	pendiente := sembrarSprintPendiente(t, sprints, proyectoPermisos)
	servicio := NewIniciarSprint(sprints, proyectoConEquipo())
	if err := servicio.Autorizar(context.Background(), pendiente.ID, scrumMasterID); err != nil {
		t.Fatalf("Autorizar: el Scrum Master debía poder iniciar el Sprint: %v", err)
	}
	if err := servicio.Autorizar(context.Background(), 999, productBuilderID); !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("Autorizar: el Sprint inexistente (404) va antes que el 403, se obtuvo %v", err)
	}
}

func TestGestion_CerrarSprint(t *testing.T) {
	for _, c := range rechazosDeGestion {
		t.Run(c.nombre, func(t *testing.T) {
			sprints := repository.NewSprintRepositoryEnMemoria()
			activo := sembrarSprintActivo(t, sprints, proyectoPermisos)
			_, err := NewCerrarSprint(sprints, &historiasEnMemoria{}, proyectoConEquipo()).Ejecutar(context.Background(), activo.ID, c.usuarioID)
			if !errors.Is(err, c.esperado) {
				t.Fatalf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if estadoPersistido(t, sprints, activo.ID) != domain.SprintActivo {
				t.Error("el Sprint no debía cerrarse")
			}
		})
	}
}

func TestGestion_EditarProyecto(t *testing.T) {
	for _, c := range rechazosDeGestion {
		t.Run(c.nombre, func(t *testing.T) {
			proyectos := proyectoConEquipoGuardado()
			servicio := NewEditarProyecto(proyectos)

			_, err := servicio.Ejecutar(context.Background(), EditarProyectoInput{ID: proyectoPermisos, Nombre: "   ", UsuarioID: c.usuarioID})
			if !errors.Is(err, c.esperado) {
				t.Fatalf("Ejecutar con nombre vacío: se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if err := servicio.Autorizar(context.Background(), proyectoPermisos, c.usuarioID); !errors.Is(err, c.esperado) {
				t.Fatalf("Autorizar: se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if proyectos.updates != 0 {
				t.Error("no se debía actualizar el proyecto")
			}
		})
	}
	servicio := NewEditarProyecto(proyectoConEquipoGuardado())
	if err := servicio.Autorizar(context.Background(), 999, productBuilderID); !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("Autorizar: el proyecto inexistente (404) va antes que el 403, se obtuvo %v", err)
	}
}

func TestGestion_ConfigurarEquipo(t *testing.T) {
	for _, c := range rechazosDeGestion {
		t.Run(c.nombre, func(t *testing.T) {
			proyectos := proyectoConEquipoGuardado()
			usuarios := &fakeUserRepository{}

			_, err := NewAsignarIntegrante(proyectos, usuarios).Ejecutar(context.Background(), AsignarIntegranteInput{
				ProyectoID: proyectoPermisos, Nombre: "Nueva", Rol: "rol-invalido", ActorID: c.usuarioID,
			})
			if !errors.Is(err, c.esperado) {
				t.Fatalf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if len(proyectos.miembros) != 2 || usuarios.creates != 0 {
				t.Error("no se debía agregar a nadie ni crear usuarios")
			}
		})
	}
}
