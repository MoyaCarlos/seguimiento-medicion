package demo

import "testing"

func TestBloqueProyecto_CreaElProyecto25ConSusDatos(t *testing.T) {
	db := baseMigrada(t)

	ejecutar(t, db, bloqueUsuarios, bloqueProyecto)

	var nombre, descripcion string
	var horasPorSP float64
	err := db.QueryRow(`SELECT nombre, descripcion, horas_por_story_point FROM proyectos WHERE id = 25`).
		Scan(&nombre, &descripcion, &horasPorSP)
	if err != nil {
		t.Fatalf("leer el proyecto 25: %v", err)
	}
	if nombre != "Licencia de conducir (datos de prueba)" {
		t.Errorf("nombre: %q", nombre)
	}
	if descripcion == "" {
		t.Error("el proyecto debe tener descripción")
	}
	if horasPorSP != 4 {
		t.Errorf("horas por Story Point: se esperaba 4, se obtuvo %v", horasPorSP)
	}
}

func TestBloqueProyecto_CarlosEsElUnicoScrumMasterYLosOtrosTresProductBuilders(t *testing.T) {
	db := baseMigrada(t)

	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto)

	roles := map[string]string{"carlos.moya": "scrum_master", "ana.iriarte": "product_builder", "candela.vulcano": "product_builder", "jimena.martinez": "product_builder"}
	for nombre, rol := range roles {
		if n := contar(t, db, `SELECT COUNT(*) FROM integrantes WHERE proyecto_id = 25 AND usuario_id = ? AND rol = ?`, c.usuarios[nombre], rol); n != 1 {
			t.Errorf("%s debía ser %s en el proyecto 25", nombre, rol)
		}
	}
	if n := contar(t, db, `SELECT COUNT(*) FROM integrantes WHERE proyecto_id = 25 AND rol = 'scrum_master'`); n != 1 {
		t.Errorf("debía haber un solo Scrum Master, hay %d", n)
	}
	if c.resumen.ProyectoID != 25 || c.resumen.Integrantes != 4 {
		t.Errorf("resumen: %+v, se esperaban 4 integrantes", c.resumen)
	}
}

func TestBloqueProyecto_ComienzaTreintaDiasAntesDeLaCargaYTerminaCuarentaYCincoDespues(t *testing.T) {
	db := baseMigrada(t)

	ejecutar(t, db, bloqueUsuarios, bloqueProyecto)

	var inicio, fin string
	if err := db.QueryRow(`SELECT fecha_inicio, fecha_fin FROM proyectos WHERE id = 25`).Scan(&inicio, &fin); err != nil {
		t.Fatal(err)
	}
	if inicio != "2026-09-30" || fin != "2026-12-14" {
		t.Errorf("fechas: %s a %s, se esperaba 2026-09-30 a 2026-12-14 (hoy-30 a hoy+45)", inicio, fin)
	}
}
