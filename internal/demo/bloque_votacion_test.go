package demo

import "testing"

func TestBloqueVotacion_RondaReveladaConLosVotosDeLosTresIntegrantes(t *testing.T) {
	db := baseMigrada(t)

	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto, bloqueSprints, bloqueHistorias, bloqueVotacion)

	aviso := c.historias["aviso"]
	if n := contar(t, db, `SELECT COUNT(*) FROM rondas_votacion WHERE historia_id = ? AND numero = 1 AND revelada_en IS NOT NULL`, aviso); n != 1 {
		t.Fatalf("falta la ronda revelada de la historia sin estimar")
	}
	votos := map[string]int{"Beta Uno": 3, "Beta Dos": 5, "Beta Tres": 5}
	for nombre, valor := range votos {
		if n := contar(t, db, `SELECT COUNT(*) FROM votos v JOIN rondas_votacion r ON r.id = v.ronda_id WHERE r.historia_id = ? AND v.usuario_id = ? AND v.valor = ?`,
			aviso, c.usuarios[nombre], valor); n != 1 {
			t.Errorf("falta el voto %d de %s", valor, nombre)
		}
	}
	if c.resumen.Rondas != 1 || c.resumen.Votos != 3 {
		t.Errorf("resumen: %+v", c.resumen)
	}
}

func TestBloqueVotacion_NoCambiaLaEstimacionDeLaHistoria(t *testing.T) {
	db := baseMigrada(t)

	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto, bloqueSprints, bloqueHistorias, bloqueVotacion)

	if n := contar(t, db, `SELECT COUNT(*) FROM historias WHERE id = ? AND estimacion_sp IS NULL`, c.historias["aviso"]); n != 1 {
		t.Error("la historia votada debe seguir sin estimar")
	}
}
