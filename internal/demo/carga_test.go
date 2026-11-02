package demo

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestCargar_BaseVaciaDevuelveElResumenCompleto(t *testing.T) {
	db := baseMigrada(t)

	resumen, err := NewCargador(db, func() time.Time { return ahoraFijo }).Cargar(context.Background())

	if err != nil {
		t.Fatalf("cargar: %v", err)
	}
	esperado := Resumen{ProyectoID: 25, Usuarios: 4, Integrantes: 3, Sprints: 3, Historias: 10, RegistrosHoras: 15, Rondas: 1, Votos: 3}
	if resumen != esperado {
		t.Errorf("resumen %+v, se esperaba %+v", resumen, esperado)
	}
}

func TestCargar_ElResumenSerializaConLosNombresDelContrato(t *testing.T) {
	cuerpo, err := json.Marshal(Resumen{ProyectoID: 25, Usuarios: 4, Integrantes: 3, Sprints: 3, Historias: 10, RegistrosHoras: 15, Rondas: 1, Votos: 3})
	if err != nil {
		t.Fatal(err)
	}
	esperado := `{"proyecto_id":25,"usuarios":4,"integrantes":3,"sprints":3,"historias":10,"registros_horas":15,"rondas_votacion":1,"votos":3}`
	if string(cuerpo) != esperado {
		t.Errorf("JSON %s, se esperaba %s", cuerpo, esperado)
	}
}

func TestCargar_NoEscribeEnTablasAjenasALaVersion1(t *testing.T) {
	db := baseMigrada(t)
	if _, err := NewCargador(db, func() time.Time { return ahoraFijo }).Cargar(context.Background()); err != nil {
		t.Fatal(err)
	}

	propias := map[string]bool{"usuarios": true, "proyectos": true, "integrantes": true, "sprints": true, "historias": true,
		"historias_arrastradas": true, "registros_horas": true, "rondas_votacion": true, "votos": true,
		"migraciones_aplicadas": true, "sqlite_sequence": true}
	filas, err := db.Query(`SELECT name FROM sqlite_master WHERE type = 'table'`)
	if err != nil {
		t.Fatal(err)
	}
	var tablas []string
	for filas.Next() {
		var nombre string
		_ = filas.Scan(&nombre)
		tablas = append(tablas, nombre)
	}
	_ = filas.Close()
	for _, tabla := range tablas {
		if !propias[tabla] && contar(t, db, `SELECT COUNT(*) FROM `+tabla) != 0 {
			t.Errorf("la tabla %s no es de la versión 1 y recibió filas", tabla)
		}
	}
}

func TestCargar_LasFechasSonRelativasAlDiaDeLaCarga(t *testing.T) {
	inicioDelProyecto := func(ahora time.Time) string {
		db := baseMigrada(t)
		if _, err := NewCargador(db, func() time.Time { return ahora }).Cargar(context.Background()); err != nil {
			t.Fatal(err)
		}
		var inicio string
		if err := db.QueryRow(`SELECT fecha_inicio FROM proyectos WHERE id = 25`).Scan(&inicio); err != nil {
			t.Fatal(err)
		}
		return inicio
	}

	if a, b := inicioDelProyecto(ahoraFijo), inicioDelProyecto(ahoraFijo.AddDate(0, 0, 10)); a != "2026-10-15" || b != "2026-10-25" {
		t.Errorf("inicios %s y %s, se esperaba 2026-10-15 y 2026-10-25", a, b)
	}
}
