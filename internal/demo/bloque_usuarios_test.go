package demo

import "testing"

func TestBloqueUsuarios_CreaLosCuatroUsuariosBeta(t *testing.T) {
	db := baseMigrada(t)

	c := ejecutar(t, db, bloqueUsuarios)

	if n := contar(t, db, `SELECT COUNT(*) FROM usuarios WHERE nombre LIKE 'Beta %'`); n != 4 {
		t.Errorf("se esperaban 4 usuarios Beta, hay %d", n)
	}
	for _, nombre := range []string{"Beta Uno", "Beta Dos", "Beta Tres", "Beta Cuatro"} {
		if c.usuarios[nombre] == 0 {
			t.Errorf("falta el id de %s", nombre)
		}
	}
	if c.resumen.Usuarios != 4 {
		t.Errorf("el resumen cuenta %d usuarios creados, se esperaban 4", c.resumen.Usuarios)
	}
}

func TestBloqueUsuarios_UnUsuarioExistenteSeReutilizaSinDuplicar(t *testing.T) {
	db := baseMigrada(t)
	if _, err := db.Exec(`INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES ('Beta Dos', 'beta dos', '2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}

	c := ejecutar(t, db, bloqueUsuarios)

	if n := contar(t, db, `SELECT COUNT(*) FROM usuarios`); n != 4 {
		t.Errorf("se esperaban 4 usuarios en total, hay %d", n)
	}
	if c.usuarios["Beta Dos"] != 1 {
		t.Errorf("Beta Dos debía reutilizar el id 1, tiene %d", c.usuarios["Beta Dos"])
	}
	if c.resumen.Usuarios != 3 {
		t.Errorf("el resumen cuenta %d usuarios creados, se esperaban 3", c.resumen.Usuarios)
	}
}
