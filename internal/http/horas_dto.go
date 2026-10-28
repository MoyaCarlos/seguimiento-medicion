package http

import "encoding/json"

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
