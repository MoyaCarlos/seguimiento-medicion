package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

type escenarioEstimacion struct {
	sp         *int
	estado     domain.Estado
	sprint     domain.EstadoSprint
	votaciones *fakeVotacionRepository
}

func votacionReveladaConDiferencias() *fakeVotacionRepository {
	return votacionRevelada(domain.Voto{UsuarioID: carlos, Valor: 3}, domain.Voto{UsuarioID: candela, Valor: 8}, domain.Voto{UsuarioID: jimena, Valor: 8})
}

func nuevoRegistrarEstimacion(e escenarioEstimacion) (*RegistrarEstimacion, *fakeEstimacionRepository) {
	historia := historiaDelBacklog()
	historia.EstimacionSP = e.sp
	if e.estado != "" {
		historia.Estado = e.estado
	}
	sprints := repository.NewSprintRepositoryEnMemoria()
	if e.sprint != "" {
		sprint, _ := sprints.Guardar(context.Background(), domain.Sprint{ProyectoID: 1, Estado: e.sprint})
		historia.SprintID = &sprint.ID
	}
	repo := &fakeEstimacionRepository{historias: map[int64]domain.BacklogItem{historia.ID: historia}}
	return NewRegistrarEstimacion(repo, equipoDelProyecto(), e.votaciones, sprints), repo
}

func TestRegistrarEstimacion_Acepta(t *testing.T) {
	casos := []struct {
		nombre    string
		escenario escenarioEstimacion
		valor     int
	}{
		{"historia sin estimación", escenarioEstimacion{}, 5},
		{"reemplaza una carga directa sin Sprint", escenarioEstimacion{sp: intPtr(3)}, 8},
		{"reemplaza con un Sprint Pendiente", escenarioEstimacion{sp: intPtr(3), sprint: domain.SprintPendiente}, 8},
		{"primera carga con un Sprint Activo", escenarioEstimacion{sprint: domain.SprintActivo}, 5},
		{"carta que nadie votó", escenarioEstimacion{}, 13},
		{"historia Completada sin estimación", escenarioEstimacion{estado: domain.EstadoCompletada, sprint: domain.SprintActivo}, 5},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			c.escenario.votaciones = votacionReveladaConDiferencias()
			servicio, repo := nuevoRegistrarEstimacion(c.escenario)

			historia, err := servicio.Ejecutar(context.Background(), 3, carlos, c.valor)

			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if historia.EstimacionSP == nil || *historia.EstimacionSP != c.valor {
				t.Errorf("se esperaban %d Story Points en la respuesta, se obtuvo %v", c.valor, historia.EstimacionSP)
			}
			if sp := repo.historias[3].EstimacionSP; sp == nil || *sp != c.valor {
				t.Errorf("se esperaban %d Story Points guardados, hay %v", c.valor, sp)
			}
		})
	}
}

func TestRegistrarEstimacion_MismoValorConSprintActivo_NoEscribe(t *testing.T) {
	servicio, repo := nuevoRegistrarEstimacion(escenarioEstimacion{sp: intPtr(3), sprint: domain.SprintActivo, votaciones: votacionReveladaConDiferencias()})

	historia, err := servicio.Ejecutar(context.Background(), 3, carlos, 3)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if historia.EstimacionSP == nil || *historia.EstimacionSP != 3 {
		t.Errorf("se esperaban 3 Story Points, se obtuvo %v", historia.EstimacionSP)
	}
	if repo.escrituras != 0 {
		t.Errorf("el mismo valor no debía escribir, hubo %d escrituras", repo.escrituras)
	}
}

func TestRegistrarEstimacion_PropagaElConflictoDelRepositorio(t *testing.T) {
	servicio, repo := nuevoRegistrarEstimacion(escenarioEstimacion{sp: intPtr(3), votaciones: votacionReveladaConDiferencias()})
	repo.errGuardar = domain.ErrRegistroEstimacionDesactualizado

	_, err := servicio.Ejecutar(context.Background(), 3, carlos, 8)

	if !errors.Is(err, domain.ErrRegistroEstimacionDesactualizado) {
		t.Fatalf("se esperaba domain.ErrRegistroEstimacionDesactualizado, se obtuvo %v", err)
	}
	if repo.escrituras != 0 {
		t.Errorf("el conflicto no debía escribir, hubo %d escrituras", repo.escrituras)
	}
}

type votacionQueNoDebeCargarVotos struct {
	*fakeVotacionRepository
	t *testing.T
}

func (v votacionQueNoDebeCargarVotos) ObtenerRondaActual(ctx context.Context, historiaID int64) (domain.RondaVotacion, error) {
	v.t.Error("registrar la estimación no debía cargar los votos de la ronda")
	return v.fakeVotacionRepository.ObtenerRondaActual(ctx, historiaID)
}

func TestRegistrarEstimacion_NoCargaLosVotos(t *testing.T) {
	historia := historiaDelBacklog()
	repo := &fakeEstimacionRepository{historias: map[int64]domain.BacklogItem{historia.ID: historia}}
	votaciones := votacionQueNoDebeCargarVotos{votacionReveladaConDiferencias(), t}
	servicio := NewRegistrarEstimacion(repo, equipoDelProyecto(), votaciones, repository.NewSprintRepositoryEnMemoria())

	if _, err := servicio.Ejecutar(context.Background(), 3, carlos, 5); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestRegistrarEstimacion_Rechaza(t *testing.T) {
	casos := []struct {
		nombre     string
		escenario  escenarioEstimacion
		historiaID int64
		usuarioID  int64
		valor      int
		esperado   error
	}{
		{"historia inexistente", escenarioEstimacion{votaciones: votacionReveladaConDiferencias()}, 999, carlos, 5, domain.ErrHistoriaNoEncontrada},
		{"Product Builder", escenarioEstimacion{sp: intPtr(3), votaciones: votacionReveladaConDiferencias()}, 3, candela, 5, domain.ErrNoEsScrumMaster},
		{"quien no integra el proyecto", escenarioEstimacion{sp: intPtr(3), votaciones: votacionReveladaConDiferencias()}, 3, ajeno, 5, domain.ErrNoEsIntegrante},
		{"carta fuera de la baraja", escenarioEstimacion{sp: intPtr(3), votaciones: votacionReveladaConDiferencias()}, 3, carlos, 4, domain.ErrEstimacionInvalida},
		{"ronda oculta", escenarioEstimacion{sp: intPtr(3), votaciones: votacionAbierta()}, 3, carlos, 5, domain.ErrRondaNoRevelada},
		{"sin votación", escenarioEstimacion{sp: intPtr(3), votaciones: &fakeVotacionRepository{}}, 3, carlos, 5, domain.ErrVotacionNoIniciada},
		{"ronda anterior revelada y actual abierta", escenarioEstimacion{sp: intPtr(3), votaciones: votacionConRondaUnoReveladaYDosAbierta()}, 3, carlos, 5, domain.ErrRondaNoRevelada},
		{"valor distinto con Sprint Activo", escenarioEstimacion{sp: intPtr(3), sprint: domain.SprintActivo, votaciones: votacionReveladaConDiferencias()}, 3, carlos, 5, domain.ErrEstimacionCongelada},
		{"valor distinto con Sprint Finalizado", escenarioEstimacion{sp: intPtr(3), sprint: domain.SprintFinalizado, votaciones: votacionReveladaConDiferencias()}, 3, carlos, 5, domain.ErrEstimacionCongelada},
		{"primera carga con un Sprint Finalizado", escenarioEstimacion{sprint: domain.SprintFinalizado, votaciones: votacionReveladaConDiferencias()}, 3, carlos, 5, domain.ErrCargaEnSprintFinalizado},
		{"primera carga en una historia Completada de un Sprint Finalizado", escenarioEstimacion{estado: domain.EstadoCompletada, sprint: domain.SprintFinalizado, votaciones: votacionReveladaConDiferencias()}, 3, carlos, 5, domain.ErrCargaEnSprintFinalizado},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			servicio, repo := nuevoRegistrarEstimacion(c.escenario)

			_, err := servicio.Ejecutar(context.Background(), c.historiaID, c.usuarioID, c.valor)

			if !errors.Is(err, c.esperado) {
				t.Errorf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if repo.escrituras != 0 {
				t.Errorf("un rechazo no debía escribir, hubo %d escrituras", repo.escrituras)
			}
		})
	}
}
