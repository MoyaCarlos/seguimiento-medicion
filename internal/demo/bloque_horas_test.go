package demo

import "testing"

func cargarHoras(t *testing.T) (*carga, func(consulta string, args ...any) int) {
	t.Helper()
	db := baseMigrada(t)
	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto, bloqueSprints, bloqueHistorias, bloqueHoras)
	return c, func(consulta string, args ...any) int { return contar(t, db, consulta, args...) }
}

func TestBloqueHoras_SoloEnHistoriasDeLosSprints1Y2ConFechasDentroDelSprint(t *testing.T) {
	c, n := cargarHoras(t)

	if c.resumen.RegistrosHoras != 15 || n(`SELECT COUNT(*) FROM registros_horas`) != 15 {
		t.Fatalf("se esperaban 15 registros, resumen %d", c.resumen.RegistrosHoras)
	}
	if got := n(`SELECT COUNT(*) FROM registros_horas WHERE sprint_id NOT IN (?, ?)`, c.sprints[0], c.sprints[1]); got != 0 {
		t.Errorf("%d registros fuera de los Sprints 1 y 2", got)
	}
	fueraDeRango := `SELECT COUNT(*) FROM registros_horas r JOIN sprints s ON s.id = r.sprint_id
		WHERE r.fecha < s.fecha_inicio OR r.fecha > s.fecha_fin OR r.fecha > '2026-10-30'`
	if got := n(fueraDeRango); got != 0 {
		t.Errorf("%d registros con fecha fuera de su Sprint o posterior a hoy", got)
	}
}

func TestBloqueHoras_RepartidasEntreLosTresIntegrantes(t *testing.T) {
	c, n := cargarHoras(t)

	for _, nombre := range []string{"Beta Uno", "Beta Dos", "Beta Tres"} {
		if n(`SELECT COUNT(*) FROM registros_horas WHERE usuario_id = ?`, c.usuarios[nombre]) == 0 {
			t.Errorf("%s no tiene horas imputadas", nombre)
		}
	}
	if n(`SELECT COUNT(*) FROM registros_horas WHERE usuario_id = ?`, c.usuarios["Beta Cuatro"]) != 0 {
		t.Error("Beta Cuatro no integra el proyecto y no debe tener horas")
	}
}

func TestBloqueHoras_LoImputadoEnLasCompletadasEsRazonableFrenteALoEstimado(t *testing.T) {
	c, n := cargarHoras(t)

	for _, clave := range []string{"turno", "datos", "apto"} {
		var puntos int
		for _, h := range historiasDemo {
			if h.clave == clave {
				puntos = h.puntos
			}
		}
		estimado := puntos * horasPorSP * 60
		minutos := n(`SELECT SUM(minutos) FROM registros_horas WHERE historia_id = ?`, c.historias[clave])
		if minutos < estimado/2 || minutos > estimado*3/2 {
			t.Errorf("%s: %d min imputados contra %d estimados (debe estar entre 50 %% y 150 %%)", clave, minutos, estimado)
		}
	}
}
