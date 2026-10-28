package http

import (
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// actorRequest identifica a quien inicia o revela la votación.
type actorRequest struct {
	UsuarioID int64 `json:"usuario_id"`
}

type votoRequest struct {
	UsuarioID int64 `json:"usuario_id"`
	Valor     int   `json:"valor"`
}

type votacionResponse struct {
	HistoriaID  int64                      `json:"historia_id"`
	Ronda       int                        `json:"ronda"`
	ReveladaEn  *string                    `json:"revelada_en"`
	Integrantes []votoDeIntegranteResponse `json:"integrantes"`
	Promedio    *promedioResponse          `json:"promedio"`
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
	return respuesta
}
