package demo

import "testing"

func TestBloqueProyecto_CreaElProyecto25ConSusDatos(t *testing.T) {
	db := baseMigrada(t)

	ejecutar(t, db, bloqueUsuarios, bloqueProyecto)

	var nombre, descripcion, inicio, fin string
	var horasPorSP float64
	err := db.QueryRow(`SELECT nombre, descripcion, fecha_inicio, fecha_fin, horas_por_story_point FROM proyectos WHERE id = 25`).
		Scan(&nombre, &descripcion, &inicio, &fin, &horasPorSP)
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
	if inicio != "2026-10-15" || fin != "2026-12-14" {
		t.Errorf("fechas: %s a %s, se esperaba 2026-10-15 a 2026-12-14 (hoy-15 a hoy+45)", inicio, fin)
	}
}

func TestBloqueProyecto_AsignaRolesYDejaABetaCuatroSinProyecto(t *testing.T) {
	db := baseMigrada(t)

	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto)

	roles := map[string]string{"Beta Uno": "scrum_master", "Beta Dos": "product_builder", "Beta Tres": "product_builder"}
	for nombre, rol := range roles {
		if n := contar(t, db, `SELECT COUNT(*) FROM integrantes WHERE proyecto_id = 25 AND usuario_id = ? AND rol = ?`, c.usuarios[nombre], rol); n != 1 {
			t.Errorf("%s debía ser %s en el proyecto 25", nombre, rol)
		}
	}
	if n := contar(t, db, `SELECT COUNT(*) FROM integrantes WHERE usuario_id = ?`, c.usuarios["Beta Cuatro"]); n != 0 {
		t.Errorf("Beta Cuatro no debe integrar ningún proyecto, tiene %d", n)
	}
	if c.resumen.ProyectoID != 25 || c.resumen.Integrantes != 3 {
		t.Errorf("resumen: %+v", c.resumen)
	}
}
