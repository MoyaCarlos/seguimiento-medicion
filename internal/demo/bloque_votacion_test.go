package demo

import "testing"

func cargarVotaciones(t *testing.T) (*carga, func(consulta string, args ...any) int) {
	t.Helper()
	db := baseMigrada(t)
	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto, bloqueSprints, bloqueHistorias, bloqueVotacion)
	return c, func(consulta string, args ...any) int { return contar(t, db, consulta, args...) }
}

func TestBloqueVotacion_RondaReveladaConLosVotosDeLosCuatroIntegrantes(t *testing.T) {
	c, n := cargarVotaciones(t)

	aviso := c.historias["aviso"]
	if n(`SELECT COUNT(*) FROM rondas_votacion WHERE historia_id = ? AND numero = 1 AND revelada_en IS NOT NULL`, aviso) != 1 {
		t.Fatalf("falta la ronda revelada de la historia sin estimar")
	}
	votos := map[string]int{"carlos.moya": 3, "ana.iriarte": 5, "candela.vulcano": 5, "jimena.martinez": 8}
	for nombre, valor := range votos {
		if n(`SELECT COUNT(*) FROM votos v JOIN rondas_votacion r ON r.id = v.ronda_id WHERE r.historia_id = ? AND v.usuario_id = ? AND v.valor = ?`,
			aviso, c.usuarios[nombre], valor) != 1 {
			t.Errorf("falta el voto %d de %s", valor, nombre)
		}
	}
}

func TestBloqueVotacion_RondaAbiertaConDosVotosParciales(t *testing.T) {
	c, n := cargarVotaciones(t)

	digital := c.historias["digital"]
	if n(`SELECT COUNT(*) FROM rondas_votacion WHERE historia_id = ? AND numero = 1 AND revelada_en IS NULL`, digital) != 1 {
		t.Fatalf("falta la ronda abierta de la historia 'digital'")
	}
	votos := map[string]int{"ana.iriarte": 3, "jimena.martinez": 5}
	for nombre, valor := range votos {
		if n(`SELECT COUNT(*) FROM votos v JOIN rondas_votacion r ON r.id = v.ronda_id WHERE r.historia_id = ? AND v.usuario_id = ? AND v.valor = ?`,
			digital, c.usuarios[nombre], valor) != 1 {
			t.Errorf("falta el voto %d de %s", valor, nombre)
		}
	}
	if got := n(`SELECT COUNT(*) FROM votos v JOIN rondas_votacion r ON r.id = v.ronda_id WHERE r.historia_id = ?`, digital); got != 2 {
		t.Errorf("la ronda abierta tiene %d votos, se esperaban 2", got)
	}
}

func TestBloqueVotacion_ElResumenCuentaDosRondasYSeisVotos(t *testing.T) {
	c, _ := cargarVotaciones(t)

	if c.resumen.Rondas != 2 || c.resumen.Votos != 6 {
		t.Errorf("resumen: %+v, se esperaban 2 rondas y 6 votos", c.resumen)
	}
}

func TestBloqueVotacion_NoCambiaLaEstimacionDeLasHistoriasVotadas(t *testing.T) {
	c, n := cargarVotaciones(t)

	for _, clave := range []string{"aviso", "digital"} {
		if n(`SELECT COUNT(*) FROM historias WHERE id = ? AND estimacion_sp IS NULL`, c.historias[clave]) != 1 {
			t.Errorf("la historia %s debe seguir sin estimar", clave)
		}
	}
}
