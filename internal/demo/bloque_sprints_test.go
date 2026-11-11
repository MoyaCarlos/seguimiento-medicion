package demo

import "testing"

func TestBloqueSprints_CreaDosFinalizadosUnActivoYUnPendienteSinFechas(t *testing.T) {
	db := baseMigrada(t)

	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto, bloqueSprints)

	if len(c.sprints) != 4 || c.resumen.Sprints != 4 {
		t.Fatalf("se esperaban 4 Sprints, hay %d (resumen %d)", len(c.sprints), c.resumen.Sprints)
	}
	casos := []struct {
		indice                    int
		estado, goal, inicio, fin string
	}{
		{0, "Finalizado", "Alta del trámite y validación médica", "2026-10-02", "2026-10-15"},
		{1, "Finalizado", "Turno y comprobante", "2026-10-16", "2026-10-22"},
		{2, "Activo", "Examen teórico y reprogramación de turnos", "2026-10-23", "2026-11-05"},
		{3, "Pendiente", "", "", ""},
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
		if caso.estado == "Pendiente" {
			if goal != nil || inicio != nil || fin != nil {
				t.Errorf("el Sprint pendiente no lleva goal ni fechas (%v %v %v)", goal, inicio, fin)
			}
			continue
		}
		if goal == nil || inicio == nil || fin == nil {
			t.Fatalf("Sprint %d: faltan goal o fechas (%v %v %v)", caso.indice+1, goal, inicio, fin)
		}
		if *goal != caso.goal || *inicio != caso.inicio || *fin != caso.fin {
			t.Errorf("Sprint %d: %q %s a %s, se esperaba %q %s a %s", caso.indice+1, *goal, *inicio, *fin, caso.goal, caso.inicio, caso.fin)
		}
	}
}
