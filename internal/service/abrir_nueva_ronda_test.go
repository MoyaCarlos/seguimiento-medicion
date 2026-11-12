package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func votosDe(valores ...int) []domain.Voto {
	usuarios := []int64{carlos, candela, jimena}
	votos := make([]domain.Voto, len(valores))
	for i, v := range valores {
		votos[i] = domain.Voto{UsuarioID: usuarios[i], Valor: v}
	}
	return votos
}

func TestAbrirNuevaRonda_ConservaLaAnterior(t *testing.T) {
	casos := map[string][]domain.Voto{
		"con diferencias": votosDe(3, 5, 8),
		"sin diferencias": votosDe(5, 5, 5),
	}
	for nombre, votos := range casos {
		t.Run(nombre, func(t *testing.T) {
			votaciones := votacionRevelada(votos...)
			servicio := NewAbrirNuevaRonda(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

			ronda, err := servicio.Ejecutar(context.Background(), 3, carlos)

			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if ronda.Numero != 2 || len(ronda.Votos) != 0 || ronda.Revelada() {
				t.Errorf("se esperaba la ronda 2 sin votos ni revelar, se obtuvo %+v", ronda)
			}
			if actual := votaciones.rondas[3]; actual.Numero != 2 {
				t.Errorf("la ronda actual debía ser la 2, es la %d", actual.Numero)
			}
			if len(votaciones.anteriores) != 1 || len(votaciones.anteriores[0].Votos) != len(votos) || !votaciones.anteriores[0].Revelada() {
				t.Errorf("la ronda 1 debía conservarse revelada y con sus votos, hay %+v", votaciones.anteriores)
			}
		})
	}
}

func TestAbrirNuevaRonda_DespuesDeRegistrarLaEstimacion(t *testing.T) {
	historia := historiaDelBacklog()
	historia.EstimacionSP = intPtr(5)
	servicio := NewAbrirNuevaRonda(historias(historia), equipoDelProyecto(), votacionRevelada(votosDe(3, 5, 8)...))

	ronda, err := servicio.Ejecutar(context.Background(), 3, carlos)

	if err != nil {
		t.Fatalf("registrar la estimación no cierra la votación: %v", err)
	}
	if ronda.Numero != 2 {
		t.Errorf("se esperaba la ronda 2, se obtuvo %d", ronda.Numero)
	}
}

func TestAbrirNuevaRonda_LaRondaNuevaCalculaSoloConSusVotos(t *testing.T) {
	ctx := context.Background()
	repoHistorias, equipo := historias(historiaDelBacklog()), equipoDelProyecto()
	votaciones := votacionRevelada(votosDe(3, 5, 13)...)
	if _, err := NewAbrirNuevaRonda(repoHistorias, equipo, votaciones).Ejecutar(ctx, 3, carlos); err != nil {
		t.Fatalf("no se esperaba error al abrir la ronda: %v", err)
	}
	votar := NewVotar(repoHistorias, equipo, votaciones)
	for _, v := range votosDe(5, 5, 8) {
		if _, err := votar.Ejecutar(ctx, VotarInput{HistoriaID: 3, UsuarioID: v.UsuarioID, Valor: v.Valor}); err != nil {
			t.Fatalf("no se esperaba error al votar: %v", err)
		}
	}
	if _, err := NewRevelarVotacion(repoHistorias, equipo, votaciones).Ejecutar(ctx, 3, carlos, time.Now()); err != nil {
		t.Fatalf("no se esperaba error al revelar: %v", err)
	}

	vista, err := NewConsultarVotacion(repoHistorias, equipo, votaciones, domain.CoincidenciaExacta{}).Ejecutar(ctx, 3, carlos)

	if err != nil {
		t.Fatalf("no se esperaba error al consultar: %v", err)
	}
	esperadas := domain.Diferencias{Hay: true, Minimo: 5, Maximo: 8}
	if vista.Numero != 2 || vista.Diferencias == nil || *vista.Diferencias != esperadas {
		t.Errorf("se esperaban la ronda 2 y %+v, se obtuvo la %d y %+v", esperadas, vista.Numero, vista.Diferencias)
	}
	if vista.Promedio == nil || *vista.Promedio != domain.CalcularPromedio([]int{5, 5, 8}) {
		t.Errorf("el promedio debía salir de 5, 5 y 8, se obtuvo %+v", vista.Promedio)
	}
}

func TestAbrirNuevaRonda_Rechaza(t *testing.T) {
	completada := historiaDelBacklog()
	completada.Estado = domain.EstadoCompletada
	casos := []struct {
		nombre     string
		historia   domain.BacklogItem
		votaciones *fakeVotacionRepository
		historiaID int64
		usuarioID  int64
		esperado   error
	}{
		{"historia inexistente", historiaDelBacklog(), votacionRevelada(votosDe(3, 8)...), 999, carlos, domain.ErrHistoriaNoEncontrada},
		{"Product Builder", historiaDelBacklog(), votacionRevelada(votosDe(3, 8)...), 3, candela, domain.ErrNoEsScrumMaster},
		{"quien no integra el proyecto", historiaDelBacklog(), votacionRevelada(votosDe(3, 8)...), 3, ajeno, domain.ErrNoEsIntegrante},
		{"historia Completada", completada, votacionRevelada(votosDe(3, 8)...), 3, carlos, domain.ErrRondaEnHistoriaCompletada},
		{"ronda oculta", historiaDelBacklog(), votacionAbierta(votosDe(3, 8)...), 3, carlos, domain.ErrRondaNoRevelada},
		{"sin votación", historiaDelBacklog(), &fakeVotacionRepository{}, 3, carlos, domain.ErrVotacionNoIniciada},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			servicio := NewAbrirNuevaRonda(historias(c.historia), equipoDelProyecto(), c.votaciones)
			numeroAntes := c.votaciones.rondas[3].Numero

			_, err := servicio.Ejecutar(context.Background(), c.historiaID, c.usuarioID)

			if !errors.Is(err, c.esperado) {
				t.Errorf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if len(c.votaciones.anteriores) != 0 || c.votaciones.rondas[3].Numero != numeroAntes {
				t.Errorf("un rechazo no debía crear ninguna ronda: %+v", c.votaciones)
			}
		})
	}
}

// ganaOtraApertura lee la ronda revelada, pero rechaza la escritura: simula que
// otra apertura ganó entre la lectura y el INSERT.
type ganaOtraApertura struct {
	*fakeVotacionRepository
}

func (ganaOtraApertura) CrearRonda(context.Context, domain.RondaVotacion) (domain.RondaVotacion, error) {
	return domain.RondaVotacion{}, domain.ErrVotacionYaIniciada
}

func TestAbrirNuevaRonda_PierdeLaCarrera(t *testing.T) {
	votaciones := ganaOtraApertura{votacionRevelada(votosDe(3, 8)...)}
	servicio := NewAbrirNuevaRonda(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones)

	_, err := servicio.Ejecutar(context.Background(), 3, carlos)

	if !errors.Is(err, domain.ErrRondaYaAbierta) {
		t.Errorf("se esperaba ErrRondaYaAbierta, se obtuvo %v", err)
	}
}
