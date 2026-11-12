package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// Escenario 4 de HU-21: quien no integra el proyecto recibe ErrNoEsIntegrante y
// no cambia nada. El recurso inexistente (404) se informa antes (FR-006).

// sinIntegrantes es un proyecto sin nadie: cualquier actor es ajeno.
func sinIntegrantes() *fakeIntegrantesRepository { return &fakeIntegrantesRepository{} }

func esperarNoEsIntegrante(t *testing.T, err error) {
	t.Helper()
	if !errors.Is(err, domain.ErrNoEsIntegrante) {
		t.Fatalf("se esperaba ErrNoEsIntegrante, se obtuvo %v", err)
	}
}

func TestPertenencia_ObtenerProyectoComoIntegrante(t *testing.T) {
	proyectos := proyectoConEquipo()
	proyectos.proyectos = map[int64]domain.Project{proyectoPermisos: {ID: proyectoPermisos, Nombre: "P"}}
	servicio := NewObtenerProyecto(proyectos)

	proyecto, rol, err := servicio.ComoIntegrante(context.Background(), proyectoPermisos, productBuilderID)
	if err != nil || proyecto.ID != proyectoPermisos || rol != domain.RolProductBuilder {
		t.Fatalf("integrante: se esperaba el proyecto con rol product_builder, se obtuvo %+v %q %v", proyecto, rol, err)
	}

	_, _, err = servicio.ComoIntegrante(context.Background(), proyectoPermisos, ajenoAlProyectoID)
	esperarNoEsIntegrante(t, err)

	_, _, err = servicio.ComoIntegrante(context.Background(), 999, ajenoAlProyectoID)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("proyecto inexistente: se esperaba ErrProyectoNoEncontrado antes que el 403, se obtuvo %v", err)
	}
}

func TestPertenencia_ConsultarVotacion(t *testing.T) {
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votacionAbierta(domain.Voto{UsuarioID: carlos, Valor: 3}), domain.CoincidenciaExacta{})

	for _, consultante := range []int64{0, ajeno} {
		_, err := servicio.Ejecutar(context.Background(), 3, consultante)
		esperarNoEsIntegrante(t, err)
	}
	_, err := servicio.Ejecutar(context.Background(), 999, ajeno)
	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Fatalf("historia inexistente: se esperaba ErrHistoriaNoEncontrada antes que el 403, se obtuvo %v", err)
	}
}

func TestPertenencia_CrearHistoria(t *testing.T) {
	repo := &fakeBacklogRepository{}
	servicio := NewCrearHistoriaBacklog(repo, proyectoValido())

	_, err := servicio.Ejecutar(context.Background(), CrearHistoriaInput{
		UsuarioID: ajenoAlProyectoID, ProyectoID: 1, Titulo: "   ", Descripcion: "d", Prioridad: domain.PrioridadMust,
	})

	esperarNoEsIntegrante(t, err)
	if len(repo.guardados) != 0 {
		t.Error("no se debía guardar la historia")
	}
}

func TestPertenencia_EditarHistoria(t *testing.T) {
	repo := &fakeHistoriaRepository{historias: map[int64]domain.BacklogItem{3: historiaDelBacklog()}}
	servicio := NewEditarHistoria(repo, repository.NewSprintRepositoryEnMemoria(), sinIntegrantes())

	entrada := edicionValida(3)
	entrada.UsuarioID = ajenoAlProyectoID
	entrada.Titulo = ""
	_, err := servicio.Ejecutar(context.Background(), entrada)

	esperarNoEsIntegrante(t, err)
	if len(repo.actualizaciones) != 0 {
		t.Error("no se debía actualizar la historia")
	}
	entrada.ID = 999
	if _, err := servicio.Ejecutar(context.Background(), entrada); !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Fatalf("historia inexistente: se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
}

func TestPertenencia_CambiarEstadoHistoria(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	activo := sembrarSprintActivo(t, sprints, 1)
	historias := fakeCambioConHistorias(historiaEnSprint(1, activo.ID, domain.EstadoNueva))
	servicio := NewCambiarEstadoHistoria(historias, sprints, relojFijo, sinIntegrantes())

	_, err := servicio.Ejecutar(context.Background(), 1, ajenoAlProyectoID, domain.EstadoEnProgreso)

	esperarNoEsIntegrante(t, err)
	if len(historias.llamadas) != 0 || historias.historias[1].Estado != domain.EstadoNueva {
		t.Error("no se debía cambiar el estado")
	}
}

func TestPertenencia_AsignarHistoriaASprint(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	pendiente := sembrarSprintPendiente(t, sprints, 1)
	historias := fakeConHistorias(historiaSinSprint(1, 1, domain.EstadoNueva))
	servicio := NewAsignarHistoriaASprint(historias, sprints, sinIntegrantes())

	_, err := servicio.Ejecutar(context.Background(), 1, pendiente.ID, ajenoAlProyectoID)
	esperarNoEsIntegrante(t, err)

	// El Sprint del cuerpo se valida después del permiso: un ajeno no sabe si existe.
	_, err = servicio.Ejecutar(context.Background(), 1, 999, ajenoAlProyectoID)
	esperarNoEsIntegrante(t, err)
	if historias.asignaciones != 0 {
		t.Error("no se debía asignar la historia")
	}
}

func TestPertenencia_RegistrarDefecto(t *testing.T) {
	e := nuevoEscenarioDefecto()
	e.sprint(t, 1, domain.SprintActivo)
	servicio := NewRegistrarDefecto(e.sprints, e.historias, e.defectos, sinIntegrantes())

	entrada := defectoValido()
	entrada.UsuarioID = ajenoAlProyectoID
	entrada.Descripcion = ""
	_, err := servicio.Ejecutar(context.Background(), entrada)

	esperarNoEsIntegrante(t, err)
	if e.defectos.registros != 0 {
		t.Error("no se debía registrar el defecto")
	}
}

func TestPertenencia_ResolverDefecto(t *testing.T) {
	e := nuevoEscenarioResolver()
	d := e.pendiente(t, e.sprint(t, 1, domain.SprintActivo))
	servicio := NewResolverDefecto(e.defectos, e.sprints, sinIntegrantes())

	_, err := servicio.Ejecutar(context.Background(), d.ID, ajenoAlProyectoID)
	esperarNoEsIntegrante(t, err)
	if e.defectos.resoluciones != 0 {
		t.Error("no se debía resolver el defecto")
	}

	if _, err := servicio.Ejecutar(context.Background(), 999, ajenoAlProyectoID); !errors.Is(err, domain.ErrDefectoNoEncontrado) {
		t.Fatalf("defecto inexistente: se esperaba ErrDefectoNoEncontrado antes que el 403, se obtuvo %v", err)
	}
}
