package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// proyectoCerradoParaReabrir crea y cierra el proyecto de proyectoParaCerrar.
func proyectoCerradoParaReabrir(t *testing.T) (*fakeProjectRepository, int64) {
	t.Helper()
	proyectos, sprints, id := proyectoParaCerrar(t)
	if _, err := NewCerrarProyecto(proyectos, sprints, relojDeCierre).Ejecutar(context.Background(), id, actorDePrueba); err != nil {
		t.Fatalf("no se pudo cerrar el proyecto: %v", err)
	}
	return proyectos, id
}

func TestReabrirProyecto_ElScrumMasterReabreUnProyectoCerrado(t *testing.T) {
	proyectos, id := proyectoCerradoParaReabrir(t)

	abierto, err := NewReabrirProyecto(proyectos).Ejecutar(context.Background(), id, actorDePrueba)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if abierto.FechaCierre != nil {
		t.Errorf("se esperaba sin fecha de cierre, se obtuvo %v", *abierto.FechaCierre)
	}
	guardado, _ := proyectos.ObtenerPorID(context.Background(), id)
	if guardado.FechaCierre != nil {
		t.Error("la reapertura no se persistió")
	}
}

func TestReabrirProyecto_AbiertoSeRechaza(t *testing.T) {
	proyectos, _, id := proyectoParaCerrar(t)

	_, err := NewReabrirProyecto(proyectos).Ejecutar(context.Background(), id, actorDePrueba)

	if !errors.Is(err, domain.ErrProyectoNoCerrado) {
		t.Fatalf("se esperaba ErrProyectoNoCerrado, se obtuvo %v", err)
	}
}

// US4: solo el Scrum Master reabre. Reabrir no recibe «El proyecto está
// cerrado»: es la única escritura permitida sobre un proyecto cerrado (FR-011).
func TestReabrirProyecto_Permisos(t *testing.T) {
	casos := []struct {
		nombre    string
		usuarioID int64
		proyecto  func(id int64) int64
		esperado  error
	}{
		{"Product Builder", productBuilderDePrueba, func(id int64) int64 { return id }, domain.ErrNoEsScrumMaster},
		{"no integrante", noIntegranteDePrueba, func(id int64) int64 { return id }, domain.ErrNoEsIntegrante},
		{"proyecto inexistente", actorDePrueba, func(int64) int64 { return 99 }, domain.ErrProyectoNoEncontrado},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			proyectos, id := proyectoCerradoParaReabrir(t)

			_, err := NewReabrirProyecto(proyectos).Ejecutar(context.Background(), c.proyecto(id), c.usuarioID)

			if !errors.Is(err, c.esperado) {
				t.Fatalf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if guardado, _ := proyectos.ObtenerPorID(context.Background(), id); guardado.FechaCierre == nil {
				t.Error("el proyecto debe seguir cerrado")
			}
		})
	}
}
