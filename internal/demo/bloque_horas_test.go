package demo

import "testing"

func cargarHoras(t *testing.T) (*carga, func(consulta string, args ...any) int) {
	t.Helper()
	db := baseMigrada(t)
	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto, bloqueSprints, bloqueHistorias, bloqueHoras)
	if len(c.sprints) != 4 {
		t.Fatalf("se esperaban 4 Sprints, hay %d", len(c.sprints))
	}
	return c, func(consulta string, args ...any) int { return contar(t, db, consulta, args...) }
}

func TestBloqueHoras_SoloEnLosSprints1A3ConFechasDentroDelSprint(t *testing.T) {
	c, n := cargarHoras(t)

	if c.resumen.RegistrosHoras != 20 || n(`SELECT COUNT(*) FROM registros_horas`) != 20 {
		t.Fatalf("se esperaban 20 registros, resumen %d", c.resumen.RegistrosHoras)
	}
	if got := n(`SELECT COUNT(*) FROM registros_horas WHERE sprint_id NOT IN (?, ?, ?)`, c.sprints[0], c.sprints[1], c.sprints[2]); got != 0 {
		t.Errorf("%d registros fuera de los Sprints 1 a 3", got)
	}
	fueraDeRango := `SELECT COUNT(*) FROM registros_horas r JOIN sprints s ON s.id = r.sprint_id
		WHERE r.fecha < s.fecha_inicio OR r.fecha > s.fecha_fin OR r.fecha > '2026-10-30'`
	if got := n(fueraDeRango); got != 0 {
		t.Errorf("%d registros con fecha fuera de su Sprint o posterior a hoy", got)
	}
}

func TestBloqueHoras_RepartidasEntreLosCuatroIntegrantes(t *testing.T) {
	c, n := cargarHoras(t)

	for _, nombre := range nombresDeLaDemo {
		if n(`SELECT COUNT(*) FROM registros_horas WHERE usuario_id = ?`, c.usuarios[nombre]) == 0 {
			t.Errorf("%s no tiene horas imputadas", nombre)
		}
	}
}

func TestBloqueHoras_LoImputadoEnLasCompletadasEsRazonableFrenteALoEstimado(t *testing.T) {
	c, n := cargarHoras(t)

	for _, clave := range []string{"tipos", "datos", "apto", "turno", "comprobante"} {
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
