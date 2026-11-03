package http

import (
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
	Valor int `json:"valor"`
}

type estimacionRequest struct {
	actorRequest
	EstimacionSP int `json:"estimacion_sp"`
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
