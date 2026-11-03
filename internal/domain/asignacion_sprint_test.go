package domain

import (
	"errors"
	"testing"
)

func historiaEnBacklog() BacklogItem {
	sp := 5
	return BacklogItem{ID: 1, ProyectoID: 1, Titulo: "Iniciar sesión", Descripcion: "Con email",
		Prioridad: PrioridadMust, Estado: EstadoNueva, EstimacionSP: &sp}
}

func sprintEn(estado EstadoSprint) Sprint {
	return Sprint{ID: 10, ProyectoID: 1, Estado: estado}
}

func enSprint(h BacklogItem, sprintID int64, estado Estado) BacklogItem {
	h.SprintID, h.Estado = &sprintID, estado
	return h
}

func TestValidarAsignacion(t *testing.T) {
	sinEstimar := historiaEnBacklog()
	sinEstimar.EstimacionSP = nil
	deOtroProyecto := sprintEn(SprintFinalizado)
	deOtroProyecto.ProyectoID = 2

	casos := []struct {
		nombre   string
		historia BacklogItem
		sprint   Sprint
		esperado error
	}{
		{"historia Nueva a Sprint Pendiente", historiaEnBacklog(), sprintEn(SprintPendiente), nil},
		{"historia Nueva a Sprint Activo", historiaEnBacklog(), sprintEn(SprintActivo), nil},
		{"historia sin Story Points", sinEstimar, sprintEn(SprintPendiente), nil},
		{"Sprint de otro proyecto", historiaEnBacklog(), Sprint{ID: 10, ProyectoID: 2, Estado: SprintPendiente}, ErrSprintDeOtroProyecto},
		{"Sprint Finalizado", historiaEnBacklog(), sprintEn(SprintFinalizado), ErrSprintFinalizado},
		{"historia Completada", enSprint(historiaEnBacklog(), 7, EstadoCompletada), sprintEn(SprintPendiente), ErrHistoriaNoPendiente},
		{"historia Nueva que ya tiene Sprint", enSprint(historiaEnBacklog(), 7, EstadoNueva), sprintEn(SprintPendiente), ErrHistoriaYaAsignada},
		{"historia que ya está en el Sprint destino", enSprint(historiaEnBacklog(), 10, EstadoNueva), sprintEn(SprintActivo), ErrHistoriaYaAsignada},
		{"otro proyecto antes que Finalizado y Completada", enSprint(historiaEnBacklog(), 7, EstadoCompletada), deOtroProyecto, ErrSprintDeOtroProyecto},
		{"Finalizado antes que Completada", enSprint(historiaEnBacklog(), 7, EstadoCompletada), sprintEn(SprintFinalizado), ErrSprintFinalizado},
		{"Completada antes que ya asignada", enSprint(historiaEnBacklog(), 10, EstadoCompletada), sprintEn(SprintActivo), ErrHistoriaNoPendiente},
		{"historia En progreso con Sprint", enSprint(historiaEnBacklog(), 7, Estado("En progreso")), sprintEn(SprintActivo), ErrHistoriaYaAsignada},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			err := ValidarAsignacion(c.historia, c.sprint)
			if c.esperado == nil {
				if err != nil {
					t.Fatalf("no se esperaba error: %v", err)
				}
				return
			}
			if !errors.Is(err, c.esperado) {
				t.Errorf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
		})
	}
}
