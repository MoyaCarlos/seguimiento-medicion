package http

import (
	"encoding/json"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// imputarHorasRequest lee horas como json.RawMessage para pasar su texto al
// dominio sin float64 y sin que un string rompa el Unmarshal (research.md R3).
type imputarHorasRequest struct {
	IntegranteID int64           `json:"integrante_id"`
	Fecha        string          `json:"fecha"`
	Actividad    string          `json:"actividad"`
	Horas        json.RawMessage `json:"horas"`
}

type registroHorasResponse struct {
	ID           int64   `json:"id"`
	HistoriaID   int64   `json:"historia_id"`
	SprintID     int64   `json:"sprint_id"`
	IntegranteID int64   `json:"integrante_id"`
	Integrante   string  `json:"integrante"`
	Fecha        string  `json:"fecha"`
	Actividad    string  `json:"actividad"`
	Minutos      int     `json:"minutos"`
	Horas        float64 `json:"horas"`
}

func aRespuestaHoras(d domain.RegistroHorasDetalle) registroHorasResponse {
	return registroHorasResponse{
		ID:           d.ID,
		HistoriaID:   d.HistoriaID,
		SprintID:     d.SprintID,
		IntegranteID: d.UsuarioID,
		Integrante:   d.NombreIntegrante,
		Fecha:        d.Fecha.Format(time.DateOnly),
		Actividad:    d.Actividad,
		Minutos:      d.Minutos,
		Horas:        float64(d.Minutos) / 60,
	}
}

type esfuerzoRealResponse struct {
	HistoriaID   int64                   `json:"historia_id"`
	TotalMinutos int                     `json:"total_minutos"`
	TotalHoras   float64                 `json:"total_horas"`
	Registros    []registroHorasResponse `json:"registros"`
}
