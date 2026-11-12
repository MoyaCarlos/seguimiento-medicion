package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

const (
	productBuilderDePrueba int64 = 2
	noIntegranteDePrueba   int64 = 3
)

// relojDeCierre es el «hoy» de los tests de cierre (spec 044, #300).
func relojDeCierre() time.Time { return time.Date(2026, 11, 11, 15, 30, 0, 0, time.Local) }

var hoyDeCierre = time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)

// proyectoParaCerrar crea un proyecto con el actor de prueba como Scrum Master
// y otra persona como Product Builder.
func proyectoParaCerrar(t *testing.T) (*fakeProjectRepository, *repository.SprintRepositoryEnMemoria, int64) {
	t.Helper()
	ctx := context.Background()
	proyectos := &fakeProjectRepository{}
	creado, err := NewCrearProyecto(proyectos).Ejecutar(ctx, CrearProyectoInput{Nombre: "Tienda", CreadorID: actorDePrueba})
	if err != nil {
		t.Fatalf("no se pudo crear el proyecto: %v", err)
	}
	if _, err := proyectos.AgregarIntegrante(ctx, domain.Membership{ProjectID: creado.ID, UserID: productBuilderDePrueba, Role: domain.RolProductBuilder}); err != nil {
		t.Fatalf("no se pudo agregar al Product Builder: %v", err)
	}
	return proyectos, repository.NewSprintRepositoryEnMemoria(), creado.ID
}

func sembrarSprintEnEstado(t *testing.T, sprints *repository.SprintRepositoryEnMemoria, proyectoID int64, estado domain.EstadoSprint) {
	t.Helper()
	if _, err := sprints.Guardar(context.Background(), domain.Sprint{ProyectoID: proyectoID, Estado: estado}); err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
}

func TestCerrarProyecto_ElScrumMasterCierraUnProyectoAbierto(t *testing.T) {
	proyectos, sprints, id := proyectoParaCerrar(t)
	sembrarSprintEnEstado(t, sprints, id, domain.SprintFinalizado)
	sembrarSprintEnEstado(t, sprints, id, domain.SprintPendiente)

	cerrado, err := NewCerrarProyecto(proyectos, sprints, relojDeCierre).Ejecutar(context.Background(), id, actorDePrueba)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if cerrado.FechaCierre == nil || !cerrado.FechaCierre.Equal(hoyDeCierre) {
		t.Errorf("se esperaba la fecha de cierre %v, se obtuvo %v", hoyDeCierre, cerrado.FechaCierre)
	}
	guardado, _ := proyectos.ObtenerPorID(context.Background(), id)
	if guardado.FechaCierre == nil || !guardado.FechaCierre.Equal(hoyDeCierre) {
		t.Errorf("el cierre no se persistió: %v", guardado.FechaCierre)
	}
}

func TestCerrarProyecto_ConUnSprintActivoSeRechaza(t *testing.T) {
	proyectos, sprints, id := proyectoParaCerrar(t)
	sembrarSprintEnEstado(t, sprints, id, domain.SprintActivo)

	_, err := NewCerrarProyecto(proyectos, sprints, relojDeCierre).Ejecutar(context.Background(), id, actorDePrueba)

	if !errors.Is(err, domain.ErrSprintActivoAlCerrarProyecto) {
		t.Fatalf("se esperaba ErrSprintActivoAlCerrarProyecto, se obtuvo %v", err)
	}
	guardado, _ := proyectos.ObtenerPorID(context.Background(), id)
	if guardado.FechaCierre != nil {
		t.Error("el proyecto debe seguir abierto")
	}
}

func TestCerrarProyecto_YaCerradoSeRechaza(t *testing.T) {
	proyectos, sprints, id := proyectoParaCerrar(t)
	servicio := NewCerrarProyecto(proyectos, sprints, relojDeCierre)
	if _, err := servicio.Ejecutar(context.Background(), id, actorDePrueba); err != nil {
		t.Fatalf("no se esperaba error en el primer cierre: %v", err)
	}

	_, err := servicio.Ejecutar(context.Background(), id, actorDePrueba)

	if !errors.Is(err, domain.ErrProyectoYaCerrado) {
		t.Fatalf("se esperaba ErrProyectoYaCerrado, se obtuvo %v", err)
	}
}

// US4: solo el Scrum Master cierra, con el orden de errores de HU-21 (404 → 403 → 409).
func TestCerrarProyecto_Permisos(t *testing.T) {
	casos := []struct {
		nombre    string
		usuarioID int64
		proyecto  func(id int64) int64
		yaCerrado bool
		esperado  error
	}{
		{"Product Builder", productBuilderDePrueba, func(id int64) int64 { return id }, false, domain.ErrNoEsScrumMaster},
		{"no integrante", noIntegranteDePrueba, func(id int64) int64 { return id }, false, domain.ErrNoEsIntegrante},
		{"proyecto inexistente", actorDePrueba, func(int64) int64 { return 99 }, false, domain.ErrProyectoNoEncontrado},
		{"inexistente gana sobre no integrante", noIntegranteDePrueba, func(int64) int64 { return 99 }, false, domain.ErrProyectoNoEncontrado},
		{"Product Builder sobre un proyecto ya cerrado: 403 antes que 409", productBuilderDePrueba, func(id int64) int64 { return id }, true, domain.ErrNoEsScrumMaster},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			proyectos, sprints, id := proyectoParaCerrar(t)
			servicio := NewCerrarProyecto(proyectos, sprints, relojDeCierre)
			if c.yaCerrado {
				if _, err := servicio.Ejecutar(context.Background(), id, actorDePrueba); err != nil {
					t.Fatalf("no se pudo cerrar de antemano: %v", err)
				}
			}

			_, err := servicio.Ejecutar(context.Background(), c.proyecto(id), c.usuarioID)

			if !errors.Is(err, c.esperado) {
				t.Fatalf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if !c.yaCerrado {
				guardado, _ := proyectos.ObtenerPorID(context.Background(), id)
				if guardado.FechaCierre != nil {
					t.Error("el proyecto no debe cerrarse")
				}
			}
		})
	}
}
