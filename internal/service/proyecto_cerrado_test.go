package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// Spec 044 (#300): un proyecto cerrado es de solo lectura. El chequeo vive en
// los guardias de escritura de HU-21 y va después de los 403.

func proyectoConEquipoCerrado() *fakeProjectRepository {
	repo := proyectoConEquipo()
	repo.cerrados = map[int64]bool{proyectoPermisos: true}
	return repo
}

func TestGuardiasDeEscritura_ProyectoCerrado(t *testing.T) {
	casos := []struct {
		nombre    string
		guardia   func(context.Context, *fakeProjectRepository, int64, int64) error
		usuarioID int64
		esperado  error
	}{
		{"exigirIntegrante: un integrante recibe 409", func(ctx context.Context, r *fakeProjectRepository, p, u int64) error {
			return exigirIntegrante(ctx, r, p, u)
		}, productBuilderID, domain.ErrProyectoCerrado},
		{"exigirIntegrante: el Scrum Master recibe 409", func(ctx context.Context, r *fakeProjectRepository, p, u int64) error {
			return exigirIntegrante(ctx, r, p, u)
		}, scrumMasterID, domain.ErrProyectoCerrado},
		{"exigirIntegrante: un no integrante recibe primero el 403", func(ctx context.Context, r *fakeProjectRepository, p, u int64) error {
			return exigirIntegrante(ctx, r, p, u)
		}, ajenoAlProyectoID, domain.ErrNoEsIntegrante},
		{"exigirScrumMaster: el Scrum Master recibe 409", func(ctx context.Context, r *fakeProjectRepository, p, u int64) error {
			return exigirScrumMaster(ctx, r, p, u)
		}, scrumMasterID, domain.ErrProyectoCerrado},
		{"exigirScrumMaster: un Product Builder recibe primero el 403", func(ctx context.Context, r *fakeProjectRepository, p, u int64) error {
			return exigirScrumMaster(ctx, r, p, u)
		}, productBuilderID, domain.ErrNoEsScrumMaster},
		{"exigirScrumMaster: un no integrante recibe primero el 403", func(ctx context.Context, r *fakeProjectRepository, p, u int64) error {
			return exigirScrumMaster(ctx, r, p, u)
		}, ajenoAlProyectoID, domain.ErrNoEsIntegrante},
		{"verificarIntegrante no mira el cierre", func(ctx context.Context, r *fakeProjectRepository, p, u int64) error {
			return verificarIntegrante(ctx, r, p, u)
		}, productBuilderID, nil},
		{"verificarScrumMaster no mira el cierre", func(ctx context.Context, r *fakeProjectRepository, p, u int64) error {
			return verificarScrumMaster(ctx, r, p, u)
		}, scrumMasterID, nil},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			err := c.guardia(context.Background(), proyectoConEquipoCerrado(), proyectoPermisos, c.usuarioID)
			if !errors.Is(err, c.esperado) {
				t.Errorf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
		})
	}
}

func TestMensajeDeProyectoCerrado(t *testing.T) {
	if got, want := domain.ErrProyectoCerrado.Error(), "El proyecto está cerrado"; got != want {
		t.Errorf("se esperaba %q, se obtuvo %q", want, got)
	}
}

func TestGuardiasDeEscritura_PropaganElErrorDelCierre(t *testing.T) {
	fallo := errors.New("fallo de la base")
	repo := &fakeIntegrantesRepositoryQueFallaAlConsultarElCierre{fakeIntegrantesRepository{miembros: []domain.Member{
		{UserID: scrumMasterID, ProjectID: proyectoPermisos, Role: domain.RolScrumMaster},
	}}, fallo}

	if err := exigirScrumMaster(context.Background(), repo, proyectoPermisos, scrumMasterID); !errors.Is(err, fallo) {
		t.Errorf("se esperaba el error del repositorio, se obtuvo %v", err)
	}
}

type fakeIntegrantesRepositoryQueFallaAlConsultarElCierre struct {
	fakeIntegrantesRepository
	err error
}

func (f *fakeIntegrantesRepositoryQueFallaAlConsultarElCierre) ProyectoCerrado(context.Context, int64) (bool, error) {
	return false, f.err
}

func TestVotar_ProyectoCerradoNoGuardaElVoto(t *testing.T) {
	votaciones := votacionAbierta()
	equipo := equipoDelProyecto()
	equipo.cerrados = map[int64]bool{1: true}
	servicio := NewVotar(historias(historiaDelBacklog()), equipo, votaciones)

	_, err := servicio.Ejecutar(context.Background(), VotarInput{HistoriaID: 3, UsuarioID: candela, Valor: 8})

	if !errors.Is(err, domain.ErrProyectoCerrado) {
		t.Errorf("se esperaba ErrProyectoCerrado, se obtuvo %v", err)
	}
	if len(votosGuardados(votaciones)) != 0 {
		t.Errorf("no se esperaba guardar el voto, quedó %+v", votosGuardados(votaciones))
	}
}

func TestImputarHoras_ProyectoCerradoNoRegistraLaCarga(t *testing.T) {
	e := nuevoEscenarioImputar(t)
	e.integrantes.cerrados = map[int64]bool{1: true}

	_, err := e.service().Ejecutar(context.Background(), entradaValida(2, "4"), e.ahora)

	if !errors.Is(err, domain.ErrProyectoCerrado) {
		t.Errorf("se esperaba ErrProyectoCerrado, se obtuvo %v", err)
	}
	if len(e.horas.registros) != 0 {
		t.Errorf("no se esperaba registrar horas, hay %d", len(e.horas.registros))
	}
}

// Consultar la votación es una lectura: sigue funcionando con el proyecto cerrado.
func TestConsultarVotacion_ProyectoCerradoSePuedeConsultar(t *testing.T) {
	equipo := equipoDelProyecto()
	equipo.cerrados = map[int64]bool{1: true}
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipo, votacionAbierta(), domain.CoincidenciaExacta{})

	if _, err := servicio.Ejecutar(context.Background(), 3, candela); err != nil {
		t.Errorf("no se esperaba error al consultar, se obtuvo %v", err)
	}
}
