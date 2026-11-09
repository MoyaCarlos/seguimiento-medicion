package demo

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestCargar_BaseVaciaDevuelveElResumenCompleto(t *testing.T) {
	db := baseMigrada(t)

	resumen, err := NewCargador(db, func() time.Time { return ahoraFijo }, hasheadorBarato).Cargar(context.Background())

	if err != nil {
		t.Fatalf("cargar: %v", err)
	}
	esperado := Resumen{ProyectoID: 25, Usuarios: 4, Integrantes: 4, Sprints: 4, Historias: 12, RegistrosHoras: 20, Rondas: 2, Votos: 6}
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
	if _, err := NewCargador(db, func() time.Time { return ahoraFijo }, hasheadorBarato).Cargar(context.Background()); err != nil {
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

func TestCargar_ElHashDeLaContrasenaSeDerivaUnaSolaVezPorCarga(t *testing.T) {
	db := baseMigrada(t)
	contador := &hasheadorContador{}
	cargador := NewCargador(db, func() time.Time { return ahoraFijo }, contador)

	if _, err := cargador.Cargar(context.Background()); err != nil {
		t.Fatal(err)
	}
	if contador.llamadas != 1 {
		t.Errorf("una carga con cuatro usuarios hizo %d derivaciones, se esperaba 1", contador.llamadas)
	}

	if _, err := cargador.Cargar(context.Background()); err == nil {
		t.Fatal("la segunda carga debía responder ya cargado")
	}
	if contador.llamadas != 2 {
		t.Errorf("tras una segunda carga que termina en 409 hay %d derivaciones, se esperaban 2 (una por carga)", contador.llamadas)
	}
}

func TestCargar_ConLaContrasenaDeLaDemoSePuedeIngresarConLosCuatroUsuarios(t *testing.T) {
	db := baseMigrada(t)
	if _, err := NewCargador(db, func() time.Time { return ahoraFijo }, hasheadorBarato).Cargar(context.Background()); err != nil {
		t.Fatal(err)
	}

	for _, nombre := range nombresDeLaDemo {
		var hash string
		if err := db.QueryRow(`SELECT COALESCE(contrasena_hash, '') FROM usuarios WHERE nombre_normalizado = ?`, nombre).Scan(&hash); err != nil {
			t.Errorf("%s: %v", nombre, err)
			continue
		}
		if !hasheadorBarato.Verificar(contrasenaDemo, hash) {
			t.Errorf("%s: %s no verifica contra el hash guardado", nombre, contrasenaDemo)
		}
		if strings.Contains(hash, contrasenaDemo) {
			t.Errorf("%s: el hash guardado contiene la contraseña en claro", nombre)
		}
	}
}

func TestCargar_TodasLasFechasSonRelativasAlDiaDeLaCarga(t *testing.T) {
	for _, ahora := range []time.Time{ahoraFijo, ahoraOtroDia} {
		db := baseMigrada(t)
		if _, err := NewCargador(db, func() time.Time { return ahora }, hasheadorBarato).Cargar(context.Background()); err != nil {
			t.Fatal(err)
		}
		hoy := ahora.UTC().Truncate(24 * time.Hour)
		dia := func(n int) string { return hoy.AddDate(0, 0, n).Format(time.DateOnly) }

		var inicio string
		if err := db.QueryRow(`SELECT fecha_inicio FROM proyectos WHERE id = 25`).Scan(&inicio); err != nil {
			t.Fatal(err)
		}
		if inicio != dia(-30) {
			t.Errorf("%s: el proyecto empieza %s, se esperaba %s", hoy.Format(time.DateOnly), inicio, dia(-30))
		}
		esperados := [][2]string{{dia(-28), dia(-15)}, {dia(-14), dia(-8)}, {dia(-7), dia(6)}}
		filas, err := db.Query(`SELECT fecha_inicio, fecha_fin FROM sprints WHERE fecha_inicio IS NOT NULL ORDER BY id`)
		if err != nil {
			t.Fatal(err)
		}
		var vistos [][2]string
		for filas.Next() {
			var par [2]string
			_ = filas.Scan(&par[0], &par[1])
			vistos = append(vistos, par)
		}
		_ = filas.Close()
		if len(vistos) != len(esperados) {
			t.Fatalf("%s: %d Sprints con fechas, se esperaban 3", hoy.Format(time.DateOnly), len(vistos))
		}
		for i := range esperados {
			if vistos[i] != esperados[i] {
				t.Errorf("%s: Sprint %d con fechas %v, se esperaba %v", hoy.Format(time.DateOnly), i+1, vistos[i], esperados[i])
			}
		}
		if posterior := contar(t, db, `SELECT COUNT(*) FROM registros_horas WHERE fecha > ?`, dia(0)); posterior != 0 {
			t.Errorf("%s: %d registros de horas posteriores a hoy", hoy.Format(time.DateOnly), posterior)
		}
	}
}
