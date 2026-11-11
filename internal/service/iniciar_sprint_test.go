package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

var (
	inicioSprint = time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	finSprint    = inicioSprint.AddDate(0, 0, 14)
)

func sembrarSprintPendiente(t *testing.T, repo repository.SprintRepository, proyectoID int64) domain.Sprint {
	t.Helper()
	s, err := domain.NewSprint(proyectoID)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := repo.Guardar(context.Background(), s)
	if err != nil {
		t.Fatalf("no se esperaba error al sembrar: %v", err)
	}
	return guardado
}

// proyectosSinFechas devuelve un proyecto sin fechas para cualquier ID, así
// que admite cualquier Sprint (#168). Solo implementa ObtenerPorID.
type proyectosSinFechas struct {
	repository.ProjectRepository
}

func (proyectosSinFechas) ObtenerPorID(_ context.Context, id int64) (domain.Project, error) {
	return domain.Project{ID: id, Nombre: "Software Metrics & Estimation"}, nil
}

func sembrarProyectoConPlazo(t *testing.T, repo repository.ProjectRepository, inicio, fin *time.Time) int64 {
	t.Helper()
	p, err := domain.NewProject("Software Metrics & Estimation", "", inicio, fin)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := repo.Guardar(context.Background(), p)
	if err != nil {
		t.Fatalf("no se esperaba error al sembrar el proyecto: %v", err)
	}
	return guardado.ID
}

func fecha(anio int, mes time.Month, dia int) time.Time {
	return time.Date(anio, mes, dia, 0, 0, 0, 0, time.UTC)
}

func inputValido(sprintID int64) IniciarSprintInput {
	return IniciarSprintInput{SprintID: sprintID, SprintGoal: "Entregar el MVP", FechaInicio: inicioSprint, FechaFin: finSprint}
}

func estadoPersistido(t *testing.T, repo repository.SprintRepository, id int64) domain.EstadoSprint {
	t.Helper()
	s, err := repo.ObtenerPorID(context.Background(), id)
	if err != nil {
		t.Fatalf("no se pudo releer el Sprint: %v", err)
	}
	return s.Estado
}

func esperarValidacion(t *testing.T, err error, campo string) {
	t.Helper()
	var verr domain.ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("se esperaba ValidationError en %q, se obtuvo %v", campo, err)
	}
	if verr.Campo != campo {
		t.Errorf("se esperaba campo %q, se obtuvo %q", campo, verr.Campo)
	}
}

func TestIniciarSprint_Exitoso(t *testing.T) {
	repo := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintPendiente(t, repo, 1)

	resultado, err := NewIniciarSprint(repo, proyectosSinFechas{}).Ejecutar(context.Background(), inputValido(sprint.ID))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if resultado.Estado != domain.SprintActivo || resultado.SprintGoal != "Entregar el MVP" {
		t.Errorf("resultado inesperado: %+v", resultado)
	}
	if estadoPersistido(t, repo, sprint.ID) != domain.SprintActivo {
		t.Errorf("el Sprint no quedó Activo en el repositorio")
	}
}

func TestIniciarSprint_RechazaSegundoActivoEnElProyecto(t *testing.T) {
	repo := repository.NewSprintRepositoryEnMemoria()
	servicio := NewIniciarSprint(repo, proyectosSinFechas{})
	primero := sembrarSprintPendiente(t, repo, 1)
	segundo := sembrarSprintPendiente(t, repo, 1)
	if _, err := servicio.Ejecutar(context.Background(), inputValido(primero.ID)); err != nil {
		t.Fatalf("no se esperaba error al iniciar el primero: %v", err)
	}

	_, err := servicio.Ejecutar(context.Background(), inputValido(segundo.ID))

	esperarValidacion(t, err, "estado")
	if estadoPersistido(t, repo, segundo.ID) != domain.SprintPendiente {
		t.Errorf("el segundo Sprint debía permanecer Pendiente")
	}
}

func TestIniciarSprint_UnActivoEnOtroProyectoNoBloquea(t *testing.T) {
	repo := repository.NewSprintRepositoryEnMemoria()
	servicio := NewIniciarSprint(repo, proyectosSinFechas{})
	otroProyecto := sembrarSprintPendiente(t, repo, 2)
	if _, err := servicio.Ejecutar(context.Background(), inputValido(otroProyecto.ID)); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	sprint := sembrarSprintPendiente(t, repo, 1)

	if _, err := servicio.Ejecutar(context.Background(), inputValido(sprint.ID)); err != nil {
		t.Fatalf("un Sprint activo de otro proyecto no debía bloquear: %v", err)
	}
}

func TestIniciarSprint_FechasInvalidasNoPersiste(t *testing.T) {
	repo := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintPendiente(t, repo, 1)
	input := inputValido(sprint.ID)
	input.FechaFin = input.FechaInicio

	_, err := NewIniciarSprint(repo, proyectosSinFechas{}).Ejecutar(context.Background(), input)

	esperarValidacion(t, err, "fecha_fin")
	if estadoPersistido(t, repo, sprint.ID) != domain.SprintPendiente {
		t.Errorf("el Sprint debía permanecer Pendiente")
	}
}

func TestIniciarSprint_NoEncontrado(t *testing.T) {
	repo := repository.NewSprintRepositoryEnMemoria()

	_, err := NewIniciarSprint(repo, proyectosSinFechas{}).Ejecutar(context.Background(), inputValido(999))

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("se esperaba ErrSprintNoEncontrado, se obtuvo %v", err)
	}
}

// Plazo del proyecto (#168): del 2026-10-31 al 2026-11-16.
func iniciarEnProyectoConPlazo(t *testing.T, inicio, fin *time.Time, input func(sprintID int64) IniciarSprintInput) (*repository.SprintRepositoryEnMemoria, domain.Sprint, error) {
	t.Helper()
	proyectos := repository.NewProjectRepositoryEnMemoria()
	sprints := repository.NewSprintRepositoryEnMemoria()
	proyectoID := sembrarProyectoConPlazo(t, proyectos, inicio, fin)
	pendiente := sembrarSprintPendiente(t, sprints, proyectoID)

	_, err := NewIniciarSprint(sprints, proyectos).Ejecutar(context.Background(), input(pendiente.ID))
	return sprints, pendiente, err
}

func conFechas(goal string, inicio, fin time.Time) func(int64) IniciarSprintInput {
	return func(sprintID int64) IniciarSprintInput {
		return IniciarSprintInput{SprintID: sprintID, SprintGoal: goal, FechaInicio: inicio, FechaFin: fin}
	}
}

func esperarSinCambios(t *testing.T, repo repository.SprintRepository, id int64) {
	t.Helper()
	s, err := repo.ObtenerPorID(context.Background(), id)
	if err != nil {
		t.Fatalf("no se pudo releer el Sprint: %v", err)
	}
	if s.Estado != domain.SprintPendiente || s.SprintGoal != "" || !s.FechaInicio.IsZero() || !s.FechaFin.IsZero() {
		t.Errorf("el Sprint debía seguir Pendiente, sin Goal ni fechas: %+v", s)
	}
}

func TestIniciarSprint_PlazoDelProyecto(t *testing.T) {
	inicioProyecto, finProyecto := fecha(2026, 10, 31), fecha(2026, 11, 16)
	casos := map[string]struct {
		inicio, fin time.Time
		campo       string
	}{
		"empieza antes que el proyecto (caso de Render)": {fecha(2026, 10, 30), fecha(2026, 11, 6), "fecha_inicio"},
		"empieza el mismo día que el proyecto":           {fecha(2026, 10, 31), fecha(2026, 11, 6), ""},
		"termina después que el proyecto":                {fecha(2026, 11, 10), fecha(2026, 11, 17), "fecha_fin"},
		"termina el mismo día que el proyecto":           {fecha(2026, 11, 10), fecha(2026, 11, 16), ""},
	}
	for nombre, caso := range casos {
		t.Run(nombre, func(t *testing.T) {
			sprints, pendiente, err := iniciarEnProyectoConPlazo(t, &inicioProyecto, &finProyecto, conFechas("Goal", caso.inicio, caso.fin))

			if caso.campo == "" {
				if err != nil {
					t.Fatalf("no se esperaba error: %v", err)
				}
				if estadoPersistido(t, sprints, pendiente.ID) != domain.SprintActivo {
					t.Errorf("el Sprint debía quedar Activo")
				}
				return
			}
			esperarValidacion(t, err, caso.campo)
			esperarSinCambios(t, sprints, pendiente.ID)
		})
	}
}

func TestIniciarSprint_ProyectoSinFechasNoRestringe(t *testing.T) {
	sprints, pendiente, err := iniciarEnProyectoConPlazo(t, nil, nil, conFechas("Goal", fecha(2026, 1, 1), fecha(2026, 1, 15)))

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if estadoPersistido(t, sprints, pendiente.ID) != domain.SprintActivo {
		t.Errorf("el Sprint debía quedar Activo")
	}
}

func TestIniciarSprint_ErroresPropiosAntesQueElPlazo(t *testing.T) {
	inicioProyecto, finProyecto := fecha(2026, 10, 31), fecha(2026, 11, 16)
	casos := map[string]struct {
		input func(int64) IniciarSprintInput
		campo string
	}{
		"goal vacío":                {conFechas("", fecha(2026, 10, 1), fecha(2026, 12, 1)), "sprint_goal"},
		"fin no posterior a inicio": {conFechas("Goal", fecha(2026, 10, 30), fecha(2026, 10, 29)), "fecha_fin"},
	}
	for nombre, caso := range casos {
		t.Run(nombre, func(t *testing.T) {
			sprints, pendiente, err := iniciarEnProyectoConPlazo(t, &inicioProyecto, &finProyecto, caso.input)

			esperarValidacion(t, err, caso.campo)
			if verr, _ := err.(domain.ValidationError); strings.Contains(verr.Mensaje, "proyecto") {
				t.Errorf("se esperaba el error propio del Sprint, no el del plazo: %v", err)
			}
			esperarSinCambios(t, sprints, pendiente.ID)
		})
	}
}

func TestIniciarSprint_OtroActivoAntesQueElPlazo(t *testing.T) {
	inicioProyecto, finProyecto := fecha(2026, 10, 31), fecha(2026, 11, 16)
	proyectos := repository.NewProjectRepositoryEnMemoria()
	sprints := repository.NewSprintRepositoryEnMemoria()
	proyectoID := sembrarProyectoConPlazo(t, proyectos, &inicioProyecto, &finProyecto)
	servicio := NewIniciarSprint(sprints, proyectos)
	activo := sembrarSprintPendiente(t, sprints, proyectoID)
	if _, err := servicio.Ejecutar(context.Background(), conFechas("Goal", fecha(2026, 11, 1), fecha(2026, 11, 8))(activo.ID)); err != nil {
		t.Fatalf("no se esperaba error al iniciar el primero: %v", err)
	}
	segundo := sembrarSprintPendiente(t, sprints, proyectoID)

	_, err := servicio.Ejecutar(context.Background(), conFechas("Goal", fecha(2026, 10, 1), fecha(2026, 12, 1))(segundo.ID))

	if !errors.Is(err, domain.ErrYaHaySprintActivo) {
		t.Fatalf("se esperaba ErrYaHaySprintActivo, se obtuvo %v", err)
	}
}
