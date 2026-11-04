package http

import (
	"encoding/json"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// actorRequest identifica a quien inicia, vota o revela.
type actorRequest struct {
	UsuarioID int64 `json:"usuario_id"`
}

func (a actorRequest) usuario() int64 { return a.UsuarioID }

type votoRequest struct {
	actorRequest
	Valor json.RawMessage `json:"valor"`
}

// valor devuelve 0 si el campo no es un entero: queda fuera de la baraja y el
// service lo rechaza después de validar el permiso.
func (v votoRequest) valor() int {
	var valor int
	if err := json.Unmarshal(v.Valor, &valor); err != nil {
		return 0
	}
	return valor
}

type estimacionRequest struct {
	actorRequest
	EstimacionSP json.RawMessage `json:"estimacion_sp"`
}

// valor devuelve 0 si estimacion_sp no es un entero: queda fuera de la baraja y
// el service lo rechaza después de validar el permiso.
func (e estimacionRequest) valor() int {
	var valor int
	if err := json.Unmarshal(e.EstimacionSP, &valor); err != nil {
		return 0
	}
	return valor
}

type votacionResponse struct {
	HistoriaID  int64                      `json:"historia_id"`
	Ronda       int                        `json:"ronda"`
	ReveladaEn  *string                    `json:"revelada_en"`
	Integrantes []votoDeIntegranteResponse `json:"integrantes"`
	Promedio    *promedioResponse          `json:"promedio"`
	Diferencias *diferenciasResponse       `json:"diferencias"`
}

type votoDeIntegranteResponse struct {
	UsuarioID int64  `json:"usuario_id"`
	Nombre    string `json:"nombre"`
	Rol       string `json:"rol"`
	YaVoto    bool   `json:"ya_voto"`
	Valor     *int   `json:"valor"`
}

type promedioResponse struct {
	Exacto        float64 `json:"exacto"`
	CartaSugerida int     `json:"carta_sugerida"`
}

// diferenciasResponse informa los extremos solo si hay diferencias.
type diferenciasResponse struct {
	HayDiferencias bool `json:"hay_diferencias"`
	Minimo         *int `json:"minimo"`
	Maximo         *int `json:"maximo"`
}

func aRespuestaDiferencias(d domain.Diferencias) *diferenciasResponse {
	if !d.Hay {
		return &diferenciasResponse{}
	}
	return &diferenciasResponse{HayDiferencias: true, Minimo: &d.Minimo, Maximo: &d.Maximo}
}

type historialRondasResponse struct {
	HistoriaID int64                   `json:"historia_id"`
	Rondas     []rondaReveladaResponse `json:"rondas"`
}

type rondaReveladaResponse struct {
	Ronda      int                    `json:"ronda"`
	ReveladaEn string                 `json:"revelada_en"`
	Votos      []votoReveladoResponse `json:"votos"`
}

type votoReveladoResponse struct {
	UsuarioID int64  `json:"usuario_id"`
	Nombre    string `json:"nombre"`
	Rol       string `json:"rol"`
	Valor     int    `json:"valor"`
}

func aRespuestaHistorial(h service.HistorialRondas) historialRondasResponse {
	respuesta := historialRondasResponse{HistoriaID: h.HistoriaID, Rondas: make([]rondaReveladaResponse, 0, len(h.Rondas))}
	for _, r := range h.Rondas {
		ronda := rondaReveladaResponse{Ronda: r.Numero, ReveladaEn: r.ReveladaEn.UTC().Format(time.RFC3339), Votos: make([]votoReveladoResponse, 0, len(r.Votos))}
		for _, v := range r.Votos {
			ronda.Votos = append(ronda.Votos, votoReveladoResponse{UsuarioID: v.UsuarioID, Nombre: v.Nombre, Rol: string(v.Rol), Valor: v.Valor})
		}
		respuesta.Rondas = append(respuesta.Rondas, ronda)
	}
	return respuesta
}

func aRespuestaVotacion(v service.VistaVotacion) votacionResponse {
	respuesta := votacionResponse{
		HistoriaID:  v.HistoriaID,
		Ronda:       v.Numero,
		Integrantes: make([]votoDeIntegranteResponse, 0, len(v.Integrantes)),
	}
	if v.ReveladaEn != nil {
		revelada := v.ReveladaEn.UTC().Format(time.RFC3339)
		respuesta.ReveladaEn = &revelada
	}
	for _, i := range v.Integrantes {
		respuesta.Integrantes = append(respuesta.Integrantes, votoDeIntegranteResponse{
			UsuarioID: i.UsuarioID, Nombre: i.Nombre, Rol: string(i.Rol), YaVoto: i.YaVoto, Valor: i.Valor,
		})
	}
	if v.Promedio != nil {
		respuesta.Promedio = &promedioResponse{Exacto: v.Promedio.Exacto, CartaSugerida: v.Promedio.CartaSugerida}
	}
	if v.Diferencias != nil {
		respuesta.Diferencias = aRespuestaDiferencias(*v.Diferencias)
	}
	return respuesta
}
