package http

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

const formatoFecha = "2006-01-02"

type crearProyectoRequest struct {
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	FechaInicio *string `json:"fecha_inicio"`
	FechaFin    *string `json:"fecha_fin"`
}

type proyectoResponse struct {
	ID          int64   `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	FechaInicio *string `json:"fecha_inicio"`
	FechaFin    *string `json:"fecha_fin"`
	CreadoEn    string  `json:"creado_en"`

	HorasPorStoryPoint *float64 `json:"horas_por_story_point"`

	// FechaCierre es null mientras el proyecto está abierto (#300).
	FechaCierre *string `json:"fecha_cierre"`

	// Rol es el de quien consulta en el proyecto (HU-21); vacío si no aplica.
	Rol string `json:"rol,omitempty"`
}

// proyectoConEstadoResponse es la respuesta de cerrar y reabrir (#300).
type proyectoConEstadoResponse struct {
	proyectoResponse
	Estado string `json:"estado"`
}

// horasPorStoryPointEnviadas distingue un campo ausente de uno en null. El valor
// pasa sin interpretar al dominio, para que un texto se rechace sobre ese campo
// y después de las demás validaciones, en lugar de romper el Unmarshal.
func horasPorStoryPointEnviadas(claves map[string]json.RawMessage) (bool, *string) {
	valor, enviadas := claves["horas_por_story_point"]
	if !enviadas || string(valor) == "null" {
		return enviadas, nil
	}
	texto := string(valor)
	return true, &texto
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

type estadoProyectoResponse struct {
	Estado string `json:"estado"`
}

type integranteResponse struct {
	ID          int64  `json:"id"`
	ProyectoID  int64  `json:"proyecto_id"`
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
		Nombre:      p.Nombre,
		Descripcion: p.Descripcion,
		FechaInicio: formatearFecha(p.FechaInicio),
		FechaFin:    formatearFecha(p.FechaFin),
		CreadoEn:    p.CreadoEn.Format(time.RFC3339),

		HorasPorStoryPoint: p.HorasPorStoryPoint,
		FechaCierre:        formatearFecha(p.FechaCierre),
	}
}

func aIntegranteResponse(m domain.Member) integranteResponse {
	return integranteResponse{
		ID:          m.UserID,
		ProyectoID:  m.ProjectID,
		Nombre:      m.Nombre,
		Rol:         string(m.Role),
		VinculadoEn: m.CreadoEn.Format(time.RFC3339),
	}
}
