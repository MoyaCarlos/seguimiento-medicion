package http

type sprintHistorialResponse struct {
	ID          int64            `json:"id"`
	ProyectoID  int64            `json:"proyecto_id"`
	Numero      int              `json:"numero"`
	SprintGoal  string           `json:"sprint_goal"`
	FechaInicio string           `json:"fecha_inicio"`
	FechaFin    string           `json:"fecha_fin"`
	Metricas    metricasResponse `json:"metricas"`
}
