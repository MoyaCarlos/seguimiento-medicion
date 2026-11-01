package http

import (
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// registrarDefectoRequest no tiene estado ni Sprints: los asigna el sistema y,
// si llegan en el cuerpo, se ignoran.
type registrarDefectoRequest struct {
	HistoriaID  int64  `json:"historia_id"`
	Descripcion string `json:"descripcion"`
	Severidad   string `json:"severidad"`
}

type defectoResponse struct {
	ID                 int64  `json:"id"`
	ProyectoID         int64  `json:"proyecto_id"`
	HistoriaID         int64  `json:"historia_id"`
	Historia           string `json:"historia"`
	Descripcion        string `json:"descripcion"`
	Severidad          string `json:"severidad"`
	Estado             string `json:"estado"`
	SprintDeteccionID  int64  `json:"sprint_deteccion_id"`
	SprintResolucionID *int64 `json:"sprint_resolucion_id"`
	CreadoEn           string `json:"creado_en"`
}

// aRespuestaDefecto serializa sprint_resolucion_id como null mientras el
// defecto está Pendiente y creado_en en RFC 3339 UTC.
func aRespuestaDefecto(d domain.DefectoDetalle) defectoResponse {
	return defectoResponse{
		ID:                 d.ID,
		ProyectoID:         d.ProyectoID,
		HistoriaID:         d.HistoriaID,
		Historia:           d.TituloHistoria,
		Descripcion:        d.Descripcion,
		Severidad:          string(d.Severidad),
		Estado:             string(d.Estado),
		SprintDeteccionID:  d.SprintDeteccionID,
		SprintResolucionID: d.SprintResolucionID,
		CreadoEn:           d.CreadoEn.UTC().Format(time.RFC3339),
	}
}
