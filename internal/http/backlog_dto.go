package http

import "github.com/MoyaCarlos/seguimiento-medicion/internal/domain"

type crearHistoriaRequest struct {
	ProyectoID   int64  `json:"proyecto_id"`
	Titulo       string `json:"titulo"`
	Descripcion  string `json:"descripcion"`
	Prioridad    string `json:"prioridad"`
	ValorNegocio *int   `json:"valor_negocio"`
}

type historiaResponse struct {
	ID           int64  `json:"id"`
	ProyectoID   int64  `json:"proyecto_id"`
	Titulo       string `json:"titulo"`
	Descripcion  string `json:"descripcion"`
	Prioridad    string `json:"prioridad"`
	Estado       string `json:"estado"`
	ValorNegocio *int   `json:"valor_negocio"`
	EstimacionSP *int   `json:"estimacion_sp"`
}

type errorResponse struct {
	Campo   string `json:"campo"`
	Mensaje string `json:"mensaje"`
}

func aRespuesta(item domain.BacklogItem) historiaResponse {
	return historiaResponse{
		ID:           item.ID,
		ProyectoID:   item.ProyectoID,
		Titulo:       item.Titulo,
		Descripcion:  item.Descripcion,
		Prioridad:    string(item.Prioridad),
		Estado:       string(item.Estado),
		ValorNegocio: item.ValorNegocio,
		EstimacionSP: item.EstimacionSP,
	}
}
