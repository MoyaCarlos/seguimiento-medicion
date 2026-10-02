package http

import (
	"strings"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

const formatoFecha = "2006-01-02"

type crearProyectoRequest struct {
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Creador     string  `json:"creador"`
	FechaInicio *string `json:"fecha_inicio"`
	FechaFin    *string `json:"fecha_fin"`
}

type proyectoResponse struct {
	ID          string  `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	FechaInicio *string `json:"fecha_inicio"`
	FechaFin    *string `json:"fecha_fin"`
	CreadoEn    string  `json:"creado_en"`
}

type editarProyectoRequest struct {
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	FechaInicio *string `json:"fecha_inicio"`
	FechaFin    *string `json:"fecha_fin"`
}

type asignarIntegranteRequest struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type integranteResponse struct {
	ID          string `json:"id"`
	ProyectoID  string `json:"proyecto_id"`
	Nombre      string `json:"nombre"`
	Rol         string `json:"rol"`
	VinculadoEn string `json:"vinculado_en"`
}

func parsearFecha(valor *string, campo string) (*time.Time, error) {
	if valor == nil || strings.TrimSpace(*valor) == "" {
		return nil, nil
	}
	t, err := time.Parse(formatoFecha, strings.TrimSpace(*valor))
	if err != nil {
		return nil, domain.ValidationError{Campo: campo, Mensaje: "formato de fecha inválido (se espera AAAA-MM-DD)"}
	}
	return &t, nil
}

func formatearFecha(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(formatoFecha)
	return &s
}

func aProyectoResponse(p domain.Project) proyectoResponse {
	return proyectoResponse{
		ID:          p.ID,
		Nombre:      p.Name,
		Descripcion: p.Description,
		FechaInicio: formatearFecha(p.StartDate),
		FechaFin:    formatearFecha(p.EndDate),
		CreadoEn:    p.CreatedAt.Format(time.RFC3339),
	}
}

func aIntegranteResponse(m domain.Member) integranteResponse {
	return integranteResponse{
		ID:          m.UserID,
		ProyectoID:  m.ProjectID,
		Nombre:      m.Name,
		Rol:         string(m.Role),
		VinculadoEn: m.CreatedAt.Format(time.RFC3339),
	}
}
