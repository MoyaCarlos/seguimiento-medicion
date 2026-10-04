package http

import (
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type crearSprintRequest struct {
	ProyectoID int64 `json:"proyecto_id"`
}

type iniciarSprintRequest struct {
	SprintGoal  string `json:"sprint_goal"`
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
}

type sprintResponse struct {
	ID          int64   `json:"id"`
	ProyectoID  int64   `json:"proyecto_id"`
	SprintGoal  *string `json:"sprint_goal"`
	FechaInicio *string `json:"fecha_inicio"`
	FechaFin    *string `json:"fecha_fin"`
	Estado      string  `json:"estado"`
}

func aRespuestaSprint(s domain.Sprint) sprintResponse {
	return sprintResponse{
		ID:          s.ID,
		ProyectoID:  s.ProyectoID,
		SprintGoal:  textoOpcional(s.SprintGoal),
		FechaInicio: fechaOpcional(s.FechaInicio),
		FechaFin:    fechaOpcional(s.FechaFin),
		Estado:      string(s.Estado),
	}
}

// Un Sprint Pendiente todavía no tiene Goal ni fechas: se devuelven como null.
func textoOpcional(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func fechaOpcional(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	f := t.Format(time.DateOnly)
	return &f
}
