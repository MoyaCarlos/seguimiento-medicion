package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// inicioSimultaneo simula que otro inicio ganó la carrera entre el chequeo
// previo y la escritura: no ve la ronda, pero CrearRonda choca con el UNIQUE.
type inicioSimultaneo struct {
	*fakeVotacionRepository
}

func (inicioSimultaneo) CrearRonda(context.Context, domain.RondaVotacion) (domain.RondaVotacion, error) {
	return domain.RondaVotacion{}, domain.ErrVotacionYaIniciada
}

func historiaEnEstado(estado domain.Estado) domain.BacklogItem {
	h := historiaDelBacklog()
	h.Estado = estado
	return h
}

func esperarSinRondas(t *testing.T, votaciones *fakeVotacionRepository) {
	t.Helper()
	if len(votaciones.rondas) != 0 {
		t.Errorf("no se esperaba crear rondas, quedaron %+v", votaciones.rondas)
	}
}

func TestIniciarVotacion_HistoriaInexistente(t *testing.T) {
	votaciones := &fakeVotacionRepository{}
	servicio := NewIniciarVotacion(historias(), equipoDelProyecto(), votaciones)

	_, err := servicio.Ejecutar(context.Background(), 999, carlos)

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Errorf("se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
	esperarSinRondas(t, votaciones)
}

func TestIniciarVotacion_NoEsScrumMaster(t *testing.T) {
	for _, usuarioID := range []int64{candela, ajeno} {
		votaciones := &fakeVotacionRepository{}
		servicio := NewIniciarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

		_, err := servicio.Ejecutar(context.Background(), 3, usuarioID)

		if !errors.Is(err, domain.ErrNoEsScrumMaster) {
			t.Errorf("usuario %d: se esperaba ErrNoEsScrumMaster, se obtuvo %v", usuarioID, err)
		}
		esperarSinRondas(t, votaciones)
	}
}

func TestIniciarVotacion_HistoriaCompletada(t *testing.T) {
	votaciones := &fakeVotacionRepository{}
	servicio := NewIniciarVotacion(historias(historiaEnEstado(domain.EstadoCompletada)), equipoDelProyecto(), votaciones)

	_, err := servicio.Ejecutar(context.Background(), 3, carlos)

	if !errors.Is(err, domain.ErrHistoriaCompletada) {
		t.Errorf("se esperaba ErrHistoriaCompletada, se obtuvo %v", err)
	}
	esperarSinRondas(t, votaciones)
}

func TestIniciarVotacion_AbreLaRonda1(t *testing.T) {
	enSprint := historiaEnEstado("En progreso")
	sprintID := int64(10)
	enSprint.SprintID = &sprintID

	for _, historia := range []domain.BacklogItem{historiaDelBacklog(), enSprint} {
		votaciones := &fakeVotacionRepository{}
		repoHistorias := historias(historia)
		servicio := NewIniciarVotacion(repoHistorias, equipoDelProyecto(), votaciones)

		ronda, err := servicio.Ejecutar(context.Background(), 3, carlos)
		if err != nil {
			t.Fatalf("historia %q: no se esperaba error: %v", historia.Estado, err)
		}

		guardada, ok := votaciones.rondas[3]
		if !ok || guardada.Numero != 1 || len(guardada.Votos) != 0 || guardada.ReveladaEn != nil {
			t.Errorf("historia %q: se esperaba guardar la ronda 1 sin votos ni revelar, quedó %+v", historia.Estado, guardada)
		}
		if ronda.ID != guardada.ID || ronda.HistoriaID != 3 || ronda.Numero != 1 {
			t.Errorf("historia %q: se esperaba devolver la ronda guardada, se obtuvo %+v", historia.Estado, ronda)
		}
		if len(repoHistorias.actualizaciones) != 0 {
			t.Errorf("historia %q: iniciar la votación no debe modificar la historia (FR-013)", historia.Estado)
		}
	}
}

func TestIniciarVotacion_YaIniciada(t *testing.T) {
	votaciones := votacionAbierta(domain.Voto{UsuarioID: candela, Valor: 8})
	servicio := NewIniciarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	_, err := servicio.Ejecutar(context.Background(), 3, carlos)

	if !errors.Is(err, domain.ErrVotacionYaIniciada) {
		t.Errorf("se esperaba ErrVotacionYaIniciada, se obtuvo %v", err)
	}
	if r := votaciones.rondas[3]; r.ID != 1 || len(r.Votos) != 1 {
		t.Errorf("la ronda existente no debía cambiar, quedó %+v", r)
	}
}

func TestIniciarVotacion_InicioSimultaneo(t *testing.T) {
	carrera := inicioSimultaneo{&fakeVotacionRepository{}}
	servicio := NewIniciarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), carrera)

	_, err := servicio.Ejecutar(context.Background(), 3, carlos)

	if !errors.Is(err, domain.ErrVotacionYaIniciada) {
		t.Errorf("se esperaba ErrVotacionYaIniciada, se obtuvo %v", err)
	}
}
