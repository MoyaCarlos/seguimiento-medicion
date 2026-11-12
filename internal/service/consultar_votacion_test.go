package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type votoVisto struct {
	usuarioID int64
	yaVoto    bool
	valor     int // 0 = oculto
}

func votosVistos(vista VistaVotacion) []votoVisto {
	vistos := make([]votoVisto, len(vista.Integrantes))
	for i, v := range vista.Integrantes {
		vistos[i] = votoVisto{usuarioID: v.UsuarioID, yaVoto: v.YaVoto}
		if v.Valor != nil {
			vistos[i].valor = *v.Valor
		}
	}
	return vistos
}

func esperarVotosVistos(t *testing.T, vista VistaVotacion, esperados []votoVisto) {
	t.Helper()
	got := votosVistos(vista)
	if len(got) != len(esperados) {
		t.Fatalf("se esperaban %+v, se obtuvo %+v", esperados, got)
	}
	for i := range esperados {
		if got[i] != esperados[i] {
			t.Errorf("se esperaban %+v, se obtuvo %+v", esperados, got)
			return
		}
	}
}

func consultar(t *testing.T, servicio *ConsultarVotacion, consultanteID int64) VistaVotacion {
	t.Helper()
	vista, err := servicio.Ejecutar(context.Background(), 3, consultanteID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	return vista
}

func TestConsultarVotacion_HistoriaInexistente(t *testing.T) {
	servicio := NewConsultarVotacion(historias(), equipoDelProyecto(), votacionAbierta(), domain.CoincidenciaExacta{})

	_, err := servicio.Ejecutar(context.Background(), 999, carlos)

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Errorf("se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
}

func TestConsultarVotacion_SinVotacion(t *testing.T) {
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), &fakeVotacionRepository{}, domain.CoincidenciaExacta{})

	_, err := servicio.Ejecutar(context.Background(), 3, carlos)

	if !errors.Is(err, domain.ErrVotacionNoIniciada) {
		t.Errorf("se esperaba ErrVotacionNoIniciada, se obtuvo %v", err)
	}
}

func TestConsultarVotacion_Abierta_SoloElVotoPropioEsVisible(t *testing.T) {
	votaciones := votacionAbierta(domain.Voto{UsuarioID: carlos, Valor: 3}, domain.Voto{UsuarioID: candela, Valor: 8})
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones, domain.CoincidenciaExacta{})

	vista := consultar(t, servicio, candela)

	if vista.HistoriaID != 3 || vista.Numero != 1 || vista.ReveladaEn != nil || vista.Promedio != nil {
		t.Errorf("se esperaba la ronda 1 de la historia 3, oculta y sin promedio, se obtuvo %+v", vista)
	}
	esperarVotosVistos(t, vista, []votoVisto{{carlos, true, 0}, {candela, true, 8}, {jimena, false, 0}})
	candelaVista := vista.Integrantes[1]
	if candelaVista.Nombre != "Candela" || candelaVista.Rol != domain.RolProductBuilder {
		t.Errorf("se esperaba Candela, product_builder, se obtuvo %+v", candelaVista)
	}
}

func TestConsultarVotacion_Abierta_IntegranteQueNoVotoNoVeNingunValor(t *testing.T) {
	votaciones := votacionAbierta(domain.Voto{UsuarioID: carlos, Valor: 3}, domain.Voto{UsuarioID: candela, Valor: 8})
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones, domain.CoincidenciaExacta{})

	vista := consultar(t, servicio, jimena)

	esperarVotosVistos(t, vista, []votoVisto{{carlos, true, 0}, {candela, true, 0}, {jimena, false, 0}})
}

func TestConsultarVotacion_IntegranteSumadoConLaRondaAbierta(t *testing.T) {
	equipo := equipoDelProyecto()
	equipo.listaMiembros = append([]domain.Member{{UserID: 4, ProjectID: 1, Nombre: "Ana", Role: domain.RolProductBuilder}}, equipo.listaMiembros...)
	votaciones := votacionAbierta(domain.Voto{UsuarioID: carlos, Valor: 3})
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipo, votaciones, domain.CoincidenciaExacta{})

	vista := consultar(t, servicio, carlos)

	esperarVotosVistos(t, vista, []votoVisto{{carlos, true, 3}, {candela, false, 0}, {jimena, false, 0}, {4, false, 0}})
}

func TestConsultarVotacion_Revelada_TodosLosValoresYElPromedio(t *testing.T) {
	votaciones := votacionRevelada(
		domain.Voto{UsuarioID: carlos, Valor: 3}, domain.Voto{UsuarioID: candela, Valor: 5}, domain.Voto{UsuarioID: jimena, Valor: 8})
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones, domain.CoincidenciaExacta{})

	vista := consultar(t, servicio, jimena)

	esperarVotosVistos(t, vista, []votoVisto{{carlos, true, 3}, {candela, true, 5}, {jimena, true, 8}})
	if vista.ReveladaEn == nil || !vista.ReveladaEn.Equal(*votaciones.rondas[3].ReveladaEn) {
		t.Errorf("se esperaba el instante de revelación, se obtuvo %v", vista.ReveladaEn)
	}
	if vista.Promedio == nil || *vista.Promedio != domain.CalcularPromedio([]int{3, 5, 8}) {
		t.Errorf("se esperaba el promedio de 3, 5 y 8, se obtuvo %+v", vista.Promedio)
	}
}

func TestConsultarVotacion_NoEscribe(t *testing.T) {
	historiasRepo := historias(historiaDelBacklog())
	equipo := equipoDelProyecto()
	votaciones := votacionAbierta(domain.Voto{UsuarioID: candela, Valor: 8})
	servicio := NewConsultarVotacion(historiasRepo, equipo, votaciones, domain.CoincidenciaExacta{})

	consultar(t, servicio, candela)

	if len(historiasRepo.actualizaciones) != 0 || equipo.creates != 0 || equipo.updates != 0 || len(equipo.miembros) != 0 {
		t.Error("la consulta no debía escribir en historias ni en proyectos")
	}
	ronda := votaciones.rondas[3]
	if ronda.ReveladaEn != nil || len(ronda.Votos) != 1 || ronda.Votos[0] != (domain.Voto{UsuarioID: candela, Valor: 8}) {
		t.Errorf("la consulta no debía cambiar la ronda, quedó %+v", ronda)
	}
}

func TestConsultarVotacion_Oculta_NoInformaDiferencias(t *testing.T) {
	votaciones := votacionAbierta(domain.Voto{UsuarioID: carlos, Valor: 3}, domain.Voto{UsuarioID: candela, Valor: 8})
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones, domain.CoincidenciaExacta{})

	for _, consultanteID := range []int64{jimena, candela} {
		if vista := consultar(t, servicio, consultanteID); vista.Diferencias != nil {
			t.Errorf("una ronda oculta no debía informar diferencias, se obtuvo %+v", vista.Diferencias)
		}
	}
}

func TestConsultarVotacion_Revelada_InformaDiferencias(t *testing.T) {
	votaciones := votacionRevelada(
		domain.Voto{UsuarioID: carlos, Valor: 3}, domain.Voto{UsuarioID: candela, Valor: 5}, domain.Voto{UsuarioID: jimena, Valor: 8})
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones, domain.CoincidenciaExacta{})

	vista := consultar(t, servicio, jimena)

	esperadas := domain.Diferencias{Hay: true, Minimo: 3, Maximo: 8}
	if vista.Diferencias == nil || *vista.Diferencias != esperadas {
		t.Errorf("se esperaba %+v, se obtuvo %+v", esperadas, vista.Diferencias)
	}
}

func TestConsultarVotacion_Revelada_SinDiferencias(t *testing.T) {
	votaciones := votacionRevelada(
		domain.Voto{UsuarioID: carlos, Valor: 5}, domain.Voto{UsuarioID: candela, Valor: 5}, domain.Voto{UsuarioID: jimena, Valor: 5})
	servicio := NewConsultarVotacion(historias(historiaDelBacklog()), equipoDelProyecto(), votaciones, domain.CoincidenciaExacta{})

	vista := consultar(t, servicio, carlos)

	if vista.Diferencias == nil || vista.Diferencias.Hay {
		t.Errorf("se esperaba que no hubiera diferencias, se obtuvo %+v", vista.Diferencias)
	}
}
