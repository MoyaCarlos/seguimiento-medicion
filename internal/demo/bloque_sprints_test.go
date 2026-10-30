package demo

import "testing"

func TestBloqueSprints_CreaFinalizadoActivoYPendienteSinFechas(t *testing.T) {
	db := baseMigrada(t)

	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto, bloqueSprints)

	if len(c.sprints) != 3 || c.resumen.Sprints != 3 {
		t.Fatalf("se esperaban 3 Sprints, hay %d (resumen %d)", len(c.sprints), c.resumen.Sprints)
	}
	casos := []struct {
		indice              int
		estado, inicio, fin string
		conGoal, conFechas  bool
	}{
		{0, "Finalizado", "2026-10-16", "2026-10-22", true, true},
		{1, "Activo", "2026-10-23", "2026-11-05", true, true},
		{2, "Pendiente", "", "", false, false},
	}
	for _, caso := range casos {
		var estado string
		var goal, inicio, fin *string
		err := db.QueryRow(`SELECT estado, sprint_goal, fecha_inicio, fecha_fin FROM sprints WHERE id = ? AND proyecto_id = 25`, c.sprints[caso.indice]).
			Scan(&estado, &goal, &inicio, &fin)
		if err != nil {
			t.Fatalf("Sprint %d: %v", caso.indice+1, err)
		}
		if estado != caso.estado {
			t.Errorf("Sprint %d: estado %s, se esperaba %s", caso.indice+1, estado, caso.estado)
		}
		if (goal != nil) != caso.conGoal || (inicio != nil) != caso.conFechas || (fin != nil) != caso.conFechas {
			t.Errorf("Sprint %d: goal/fechas inesperados (%v %v %v)", caso.indice+1, goal, inicio, fin)
		}
		if caso.conFechas && (*inicio != caso.inicio || *fin != caso.fin) {
			t.Errorf("Sprint %d: fechas %s a %s, se esperaba %s a %s", caso.indice+1, *inicio, *fin, caso.inicio, caso.fin)
		}
	}
}
