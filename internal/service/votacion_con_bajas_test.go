package service

import (
	"context"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// darDeBaja quita a un Product Builder del proyecto 1 del fake (HU-16).
func darDeBaja(t *testing.T, equipo *fakeProjectRepository, usuarioID int64) {
	t.Helper()
	if err := equipo.DarDeBajaIntegrante(context.Background(), 1, usuarioID, time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("no se pudo dar de baja a %d: %v", usuarioID, err)
	}
}

// HU-16, FR-007: revelar una ronda exige solo los votos de los activos, así
// que no se espera a quien fue dado de baja sin votar.
func TestRevelarVotacion_NoEsperaAUnIntegranteDadoDeBaja(t *testing.T) {
	equipo := equipoDelProyecto()
	darDeBaja(t, equipo, jimena)
	votaciones := votacionAbierta(domain.Voto{UsuarioID: carlos, Valor: 3}, domain.Voto{UsuarioID: candela, Valor: 5})
	servicio := NewRevelarVotacion(historias(historiaDelBacklog()), equipo, votaciones)

	if _, err := servicio.Ejecutar(context.Background(), 3, carlos, momentoDeRevelar); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if votaciones.rondas[3].ReveladaEn == nil {
		t.Error("la ronda debía quedar revelada")
	}
}

// HU-16: en una ronda abierta, quien fue dado de baja sin votar no figura
// entre los que faltan votar.
func TestConsultarVotacion_Abierta_NoListaAUnDadoDeBajaQueNoVoto(t *testing.T) {
	equipo := equipoDelProyecto()
	darDeBaja(t, equipo, jimena)
	votaciones := votacionAbierta(domain.Voto{UsuarioID: carlos, Valor: 3})
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipo, votaciones, domain.CoincidenciaExacta{})

	vista := consultar(t, servicio, carlos)

	esperarVotosVistos(t, vista, []votoVisto{{carlos, true, 3}, {candela, false, 0}})
}

// HU-16, FR-002: el historial de rondas sigue mostrando con su nombre y rol el
// voto revelado de quien después fue dado de baja.
func TestConsultarHistorialRondas_ConservaElNombreDeUnDadoDeBaja(t *testing.T) {
	equipo := equipoDelProyecto()
	darDeBaja(t, equipo, candela)
	servicio := NewConsultarHistorialRondas(historias(historiaDelBacklog()), equipo, votacionConRondaUnoReveladaYDosAbierta())

	historial, err := servicio.Ejecutar(context.Background(), 3)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(historial.Rondas) != 1 {
		t.Fatalf("se esperaba la ronda 1, se obtuvo %+v", historial)
	}
	voto := historial.Rondas[0].Votos[1]
	if voto.UsuarioID != candela || voto.Nombre != "Candela" || voto.Rol != domain.RolProductBuilder || voto.Valor != 8 {
		t.Errorf("se esperaba el voto de Candela (product_builder, 8), se obtuvo %+v", voto)
	}
}

// HU-16, FR-002: en una ronda ya revelada, el voto de quien después fue dado de
// baja sigue en la vista, y el promedio coincide con los votos mostrados.
func TestConsultarVotacion_Revelada_IncluyeAUnDadoDeBajaQueVoto(t *testing.T) {
	equipo := equipoDelProyecto()
	votaciones := votacionRevelada(
		domain.Voto{UsuarioID: carlos, Valor: 3}, domain.Voto{UsuarioID: candela, Valor: 5}, domain.Voto{UsuarioID: jimena, Valor: 8})
	darDeBaja(t, equipo, jimena)
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipo, votaciones, domain.CoincidenciaExacta{})

	vista := consultar(t, servicio, carlos)

	esperarVotosVistos(t, vista, []votoVisto{{carlos, true, 3}, {candela, true, 5}, {jimena, true, 8}})
	if n := vista.Integrantes[2]; n.Nombre != "Jimena" || n.Rol != domain.RolProductBuilder {
		t.Errorf("se esperaba a Jimena (product_builder), se obtuvo %+v", n)
	}
	if vista.Promedio == nil || *vista.Promedio != domain.CalcularPromedio([]int{3, 5, 8}) {
		t.Errorf("se esperaba el promedio de 3, 5 y 8, se obtuvo %+v", vista.Promedio)
	}
}
