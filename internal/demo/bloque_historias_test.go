package demo

import "testing"

func cargarHistorias(t *testing.T) (*carga, func(consulta string, args ...any) int) {
	t.Helper()
	db := baseMigrada(t)
	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto, bloqueSprints, bloqueHistorias)
	return c, func(consulta string, args ...any) int { return contar(t, db, consulta, args...) }
}

func TestBloqueHistorias_RepartoPorSprintYEstado(t *testing.T) {
	c, n := cargarHistorias(t)

	if c.resumen.Historias != 10 || n(`SELECT COUNT(*) FROM historias WHERE proyecto_id = 25`) != 10 {
		t.Fatalf("se esperaban 10 historias, resumen %d", c.resumen.Historias)
	}
	casos := []struct {
		descripcion string
		sprint      int
		estado      string
		esperadas   int
	}{
		{"Completadas en el Sprint 1", 0, "Completada", 3},
		{"En progreso en el Sprint 2", 1, "En progreso", 2},
		{"Nuevas en el Sprint 2", 1, "Nueva", 1},
		{"Nuevas en el Sprint 3", 2, "Nueva", 1},
	}
	for _, caso := range casos {
		if got := n(`SELECT COUNT(*) FROM historias WHERE sprint_id = ? AND estado = ?`, c.sprints[caso.sprint], caso.estado); got != caso.esperadas {
			t.Errorf("%s: %d, se esperaban %d", caso.descripcion, got, caso.esperadas)
		}
	}
	if got := n(`SELECT COUNT(*) FROM historias WHERE sprint_id IS NULL AND estado = 'Nueva'`); got != 3 {
		t.Errorf("en el Product Backlog: %d, se esperaban 3", got)
	}
}

func TestBloqueHistorias_VariedadDePrioridadesEstimacionesYCriterios(t *testing.T) {
	_, n := cargarHistorias(t)

	for _, prioridad := range []string{"M", "S", "C", "W"} {
		if n(`SELECT COUNT(*) FROM historias WHERE prioridad = ?`, prioridad) == 0 {
			t.Errorf("falta una historia de prioridad %s", prioridad)
		}
	}
	if n(`SELECT COUNT(*) FROM historias WHERE estimacion_sp IS NOT NULL AND estimacion_sp NOT IN (1,2,3,5,8,13,21)`) != 0 {
		t.Error("hay estimaciones fuera de la escala Fibonacci")
	}
	if n(`SELECT COUNT(*) FROM historias WHERE sprint_id IS NULL AND estimacion_sp IS NULL`) < 1 {
		t.Error("al menos una historia del backlog debe estar sin estimar")
	}
	if n(`SELECT COUNT(*) FROM historias WHERE estimacion_sp IS NOT NULL`) == 0 {
		t.Error("debe haber historias estimadas")
	}
	if n(`SELECT COUNT(*) FROM historias WHERE criterios_aceptacion IS NOT NULL`) == 0 ||
		n(`SELECT COUNT(*) FROM historias WHERE criterios_aceptacion IS NULL`) == 0 {
		t.Error("debe haber historias con y sin criterios de aceptación")
	}
	if n(`SELECT COUNT(*) FROM historias WHERE descripcion LIKE '%**%'`) == 0 {
		t.Error("alguna descripción debe usar Markdown")
	}
}

func TestBloqueHistorias_CompletadasTienenFechaDeFinalizacionDentroDelSprint1(t *testing.T) {
	_, n := cargarHistorias(t)

	if got := n(`SELECT COUNT(*) FROM historias WHERE estado = 'Completada' AND fecha_finalizacion BETWEEN '2026-10-16' AND '2026-10-22'`); got != 3 {
		t.Errorf("fechas de finalización dentro del Sprint 1: %d, se esperaban 3", got)
	}
	if got := n(`SELECT COUNT(*) FROM historias WHERE estado <> 'Completada' AND fecha_finalizacion IS NOT NULL`); got != 0 {
		t.Errorf("solo las Completadas llevan fecha de finalización, hay %d de más", got)
	}
}

func TestBloqueHistorias_LaArrastradaVolvioAlBacklogConSuRegistro(t *testing.T) {
	c, n := cargarHistorias(t)

	if got := n(`SELECT COUNT(*) FROM historias_arrastradas WHERE sprint_id = ? AND historia_id = ? AND proyecto_id = 25 AND estimacion_sp = 8`,
		c.sprints[0], c.historias["tasa"]); got != 1 {
		t.Fatalf("falta el registro de arrastrada de la historia 'tasa'")
	}
	if got := n(`SELECT COUNT(*) FROM historias WHERE id = ? AND sprint_id IS NULL AND estado = 'Nueva'`, c.historias["tasa"]); got != 1 {
		t.Error("la arrastrada debe estar en el backlog como Nueva")
	}
	if got := n(`SELECT COUNT(*) FROM historias_arrastradas`); got != 1 {
		t.Errorf("se esperaba una sola arrastrada, hay %d", got)
	}
}
