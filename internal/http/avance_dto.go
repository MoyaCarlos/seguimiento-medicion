package http

import "github.com/MoyaCarlos/seguimiento-medicion/internal/domain"

type totalesAvanceResponse struct {
	Historias   int `json:"historias"`
	StoryPoints int `json:"story_points"`
}

type sprintRefResponse struct {
	ID     int64 `json:"id"`
	Numero int   `json:"numero"`
}

type historiaAvanceResponse struct {
	ID           int64              `json:"id"`
	Titulo       string             `json:"titulo"`
	EstimacionSP *int               `json:"estimacion_sp"`
	Estado       string             `json:"estado"`
	Sprint       *sprintRefResponse `json:"sprint"`
}

type grupoAvanceResponse struct {
	Estado      string                   `json:"estado"`
	Cantidad    int                      `json:"cantidad"`
	StoryPoints int                      `json:"story_points"`
	Historias   []historiaAvanceResponse `json:"historias"`
}

type avanceResponse struct {
	Total  totalesAvanceResponse `json:"total"`
	Grupos []grupoAvanceResponse `json:"grupos"`
}

// aAvanceResponse solo copia: la agrupación y las sumas ya vienen del dominio.
func aAvanceResponse(avance domain.AvanceHistorias) avanceResponse {
	grupos := make([]grupoAvanceResponse, 0, len(avance.Grupos))
	for _, g := range avance.Grupos {
		historias := make([]historiaAvanceResponse, 0, len(g.Historias))
		for _, h := range g.Historias {
			historias = append(historias, historiaAvanceResponse{
				ID: h.ID, Titulo: h.Titulo, EstimacionSP: h.EstimacionSP, Estado: string(h.Estado), Sprint: aSprintRefResponse(h.Sprint),
			})
		}
		grupos = append(grupos, grupoAvanceResponse{
			Estado: string(g.Estado), Cantidad: g.Cantidad, StoryPoints: g.StoryPoints, Historias: historias,
		})
	}
	return avanceResponse{
		Total:  totalesAvanceResponse{Historias: avance.Total.Historias, StoryPoints: avance.Total.StoryPoints},
		Grupos: grupos,
	}
}

func aSprintRefResponse(sprint *domain.SprintRef) *sprintRefResponse {
	if sprint == nil {
		return nil
	}
	return &sprintRefResponse{ID: sprint.ID, Numero: sprint.Numero}
}
