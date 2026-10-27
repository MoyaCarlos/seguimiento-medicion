package repository

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"
	"testing"

	sqlite "modernc.org/sqlite"
)

func nuevaBDMigrada(t *testing.T) *sql.DB {
	t.Helper()
	db, err := AbrirSQLite(":memory:")
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := AplicarMigraciones(context.Background(), db); err != nil {
		t.Fatalf("no se pudieron aplicar las migraciones: %v", err)
	}
	return db
}

type columnaEsperada struct {
	nombre  string
	tipo    string
	notNull bool
	defecto string
	pk      int
}

type indiceEsperado struct {
	nombre   string
	columnas []string
	where    string
}

type tablaEsperada struct {
	columnas []columnaEsperada
	fks      []string
	uniques  []string
	indices  []indiceEsperado
}

// Estructura copiada de MODELO-DE-DATOS.md; cada tabla cita su sección.
// El test no lee el documento: si la migración se aparta de esto, falla.
var esquemaEsperado = map[string]tablaEsperada{
	// §4.1 proyectos
	"proyectos": {
		columnas: []columnaEsperada{
			{nombre: "id", tipo: "INTEGER", pk: 1},
			{nombre: "nombre", tipo: "TEXT", notNull: true},
			{nombre: "descripcion", tipo: "TEXT", notNull: true, defecto: "''"},
			{nombre: "fecha_inicio", tipo: "TEXT"},
			{nombre: "fecha_fin", tipo: "TEXT"},
			{nombre: "horas_por_story_point", tipo: "REAL"},
			{nombre: "creado_en", tipo: "TEXT", notNull: true},
		},
	},
	// §4.2 usuarios
	"usuarios": {
		columnas: []columnaEsperada{
			{nombre: "id", tipo: "INTEGER", pk: 1},
			{nombre: "nombre", tipo: "TEXT", notNull: true},
			{nombre: "nombre_normalizado", tipo: "TEXT", notNull: true},
			{nombre: "creado_en", tipo: "TEXT", notNull: true},
		},
		uniques: []string{"nombre_normalizado"},
	},
	// §4.3 integrantes
	"integrantes": {
		columnas: []columnaEsperada{
			{nombre: "proyecto_id", tipo: "INTEGER", notNull: true, pk: 1},
			{nombre: "usuario_id", tipo: "INTEGER", notNull: true, pk: 2},
			{nombre: "rol", tipo: "TEXT", notNull: true},
			{nombre: "creado_en", tipo: "TEXT", notNull: true},
		},
		fks: []string{
			"proyecto_id -> proyectos(id)",
			"usuario_id -> usuarios(id)",
		},
	},
	// §4.4 historias
	"historias": {
		columnas: []columnaEsperada{
			{nombre: "id", tipo: "INTEGER", pk: 1},
			{nombre: "proyecto_id", tipo: "INTEGER", notNull: true},
			{nombre: "titulo", tipo: "TEXT", notNull: true},
			{nombre: "descripcion", tipo: "TEXT", notNull: true},
			{nombre: "criterios_aceptacion", tipo: "TEXT"},
			{nombre: "prioridad", tipo: "TEXT", notNull: true},
			{nombre: "estado", tipo: "TEXT", notNull: true},
			{nombre: "valor_negocio", tipo: "INTEGER"},
			{nombre: "estimacion_sp", tipo: "INTEGER"},
			{nombre: "sprint_id", tipo: "INTEGER"},
			{nombre: "fecha_finalizacion", tipo: "TEXT"},
		},
		fks: []string{
			"proyecto_id -> proyectos(id)",
			"sprint_id,proyecto_id -> sprints(id,proyecto_id)",
		},
		uniques: []string{"id,proyecto_id"},
	},
	// §4.5 sprints
	"sprints": {
		columnas: []columnaEsperada{
			{nombre: "id", tipo: "INTEGER", pk: 1},
			{nombre: "proyecto_id", tipo: "INTEGER", notNull: true},
			{nombre: "sprint_goal", tipo: "TEXT"},
			{nombre: "fecha_inicio", tipo: "TEXT"},
			{nombre: "fecha_fin", tipo: "TEXT"},
			{nombre: "estado", tipo: "TEXT", notNull: true},
		},
		fks:     []string{"proyecto_id -> proyectos(id)"},
		uniques: []string{"id,proyecto_id"},
		indices: []indiceEsperado{
			{nombre: "sprints_un_activo_por_proyecto", columnas: []string{"proyecto_id"}, where: "WHERE estado = 'Activo'"},
			{nombre: "sprints_un_pendiente_por_proyecto", columnas: []string{"proyecto_id"}, where: "WHERE estado = 'Pendiente'"},
		},
	},
	// §4.7 rondas_votacion
	"rondas_votacion": {
		columnas: []columnaEsperada{
			{nombre: "id", tipo: "INTEGER", pk: 1},
			{nombre: "historia_id", tipo: "INTEGER", notNull: true},
			{nombre: "numero", tipo: "INTEGER", notNull: true},
			{nombre: "revelada_en", tipo: "TEXT"},
		},
		fks:     []string{"historia_id -> historias(id)"},
		uniques: []string{"historia_id,numero"},
	},
	// §4.7 votos
	"votos": {
		columnas: []columnaEsperada{
			{nombre: "ronda_id", tipo: "INTEGER", notNull: true, pk: 1},
			{nombre: "usuario_id", tipo: "INTEGER", notNull: true, pk: 2},
			{nombre: "valor", tipo: "INTEGER", notNull: true},
		},
		fks: []string{
			"ronda_id -> rondas_votacion(id)",
			"usuario_id -> usuarios(id)",
		},
	},
	// §9 punto 1 migraciones_aplicadas
	"migraciones_aplicadas": {
		columnas: []columnaEsperada{
			{nombre: "version", tipo: "INTEGER", pk: 1},
			{nombre: "aplicada_en", tipo: "TEXT", notNull: true},
		},
	},
}

func TestEsquema_Estructura(t *testing.T) {
	db := nuevaBDMigrada(t)

	tablas := consultarTextos(t, db,
		"SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name")
	var esperadas []string
	for nombre := range esquemaEsperado {
		esperadas = append(esperadas, nombre)
	}
	slices.Sort(esperadas)
	if !slices.Equal(tablas, esperadas) {
		t.Fatalf("tablas: se esperaba %v, se obtuvo %v", esperadas, tablas)
	}

	for nombre, esperada := range esquemaEsperado {
		t.Run(nombre, func(t *testing.T) {
			if got := columnasDe(t, db, nombre); !slices.Equal(got, esperada.columnas) {
				t.Errorf("columnas:\n se esperaba %+v\n se obtuvo   %+v", esperada.columnas, got)
			}
			if got := fksDe(t, db, nombre); !slices.Equal(got, esperada.fks) {
				t.Errorf("FK: se esperaba %v, se obtuvo %v", esperada.fks, got)
			}
			uniques, indices := indicesDe(t, db, nombre)
			if !slices.Equal(uniques, esperada.uniques) {
				t.Errorf("UNIQUE: se esperaba %v, se obtuvo %v", esperada.uniques, uniques)
			}
			if !slices.EqualFunc(indices, esperada.indices, indicesIguales) {
				t.Errorf("índices: se esperaba %+v, se obtuvo %+v", esperada.indices, indices)
			}
		})
	}

	if got := versionesAplicadas(t, db); !slices.Contains(got, 1) || !slices.Contains(got, 3) {
		t.Errorf("se esperaban las versiones 1 y 3 registradas, se obtuvo %v", got)
	}
}

func indicesIguales(a, b indiceEsperado) bool {
	return a.nombre == b.nombre && slices.Equal(a.columnas, b.columnas) && a.where == b.where
}

func consultarTextos(t *testing.T, db *sql.DB, consulta string, args ...any) []string {
	t.Helper()
	filas, err := db.Query(consulta, args...)
	if err != nil {
		t.Fatalf("consulta %q: %v", consulta, err)
	}
	defer filas.Close()
	var valores []string
	for filas.Next() {
		var v string
		if err := filas.Scan(&v); err != nil {
			t.Fatalf("consulta %q: %v", consulta, err)
		}
		valores = append(valores, v)
	}
	return valores
}

func columnasDe(t *testing.T, db *sql.DB, tabla string) []columnaEsperada {
	t.Helper()
	filas, err := db.Query(`SELECT name, type, "notnull", dflt_value, pk FROM pragma_table_info(?) ORDER BY cid`, tabla)
	if err != nil {
		t.Fatalf("pragma_table_info(%s): %v", tabla, err)
	}
	defer filas.Close()
	var columnas []columnaEsperada
	for filas.Next() {
		var (
			c       columnaEsperada
			notNull int
			defecto sql.NullString
		)
		if err := filas.Scan(&c.nombre, &c.tipo, &notNull, &defecto, &c.pk); err != nil {
			t.Fatalf("pragma_table_info(%s): %v", tabla, err)
		}
		c.notNull = notNull == 1
		c.defecto = defecto.String
		columnas = append(columnas, c)
	}
	return columnas
}

func fksDe(t *testing.T, db *sql.DB, tabla string) []string {
	t.Helper()
	filas, err := db.Query(
		`SELECT id, "table", "from", "to", on_update, on_delete FROM pragma_foreign_key_list(?) ORDER BY id, seq`, tabla)
	if err != nil {
		t.Fatalf("pragma_foreign_key_list(%s): %v", tabla, err)
	}
	defer filas.Close()
	type fk struct {
		padre    string
		desde, a []string
	}
	porID := map[int]*fk{}
	var orden []int
	for filas.Next() {
		var (
			id                 int
			padre, desde       string
			a                  sql.NullString
			onUpdate, onDelete string
		)
		if err := filas.Scan(&id, &padre, &desde, &a, &onUpdate, &onDelete); err != nil {
			t.Fatalf("pragma_foreign_key_list(%s): %v", tabla, err)
		}
		if onUpdate != "NO ACTION" || onDelete != "NO ACTION" {
			t.Errorf("FK %s.%s: se esperaba ON UPDATE/ON DELETE NO ACTION (D9), se obtuvo %s/%s", tabla, desde, onUpdate, onDelete)
		}
		if porID[id] == nil {
			porID[id] = &fk{padre: padre}
			orden = append(orden, id)
		}
		porID[id].desde = append(porID[id].desde, desde)
		porID[id].a = append(porID[id].a, a.String)
	}
	var fks []string
	for _, id := range orden {
		f := porID[id]
		fks = append(fks, strings.Join(f.desde, ",")+" -> "+f.padre+"("+strings.Join(f.a, ",")+")")
	}
	slices.Sort(fks)
	return fks
}

func indicesDe(t *testing.T, db *sql.DB, tabla string) (uniques []string, indices []indiceEsperado) {
	t.Helper()
	filas, err := db.Query(`SELECT name, "unique", origin, partial FROM pragma_index_list(?)`, tabla)
	if err != nil {
		t.Fatalf("pragma_index_list(%s): %v", tabla, err)
	}
	type indice struct {
		nombre, origen string
		unico, parcial int
	}
	var lista []indice
	for filas.Next() {
		var i indice
		if err := filas.Scan(&i.nombre, &i.unico, &i.origen, &i.parcial); err != nil {
			t.Fatalf("pragma_index_list(%s): %v", tabla, err)
		}
		lista = append(lista, i)
	}
	filas.Close()

	for _, i := range lista {
		columnas := consultarTextos(t, db, "SELECT name FROM pragma_index_info(?) ORDER BY seqno", i.nombre)
		switch i.origen {
		case "u":
			uniques = append(uniques, strings.Join(columnas, ","))
		case "c":
			if i.unico != 1 || i.parcial != 1 {
				t.Errorf("índice %s: se esperaba único y parcial", i.nombre)
			}
			ddl := consultarTextos(t, db, "SELECT sql FROM sqlite_master WHERE type = 'index' AND name = ?", i.nombre)
			where := ""
			if len(ddl) == 1 {
				if pos := strings.Index(ddl[0], "WHERE"); pos >= 0 {
					where = strings.TrimSpace(ddl[0][pos:])
				}
			}
			indices = append(indices, indiceEsperado{nombre: i.nombre, columnas: columnas, where: where})
		}
	}
	slices.Sort(uniques)
	slices.SortFunc(indices, func(a, b indiceEsperado) int { return strings.Compare(a.nombre, b.nombre) })
	return uniques, indices
}

const (
	aceptada         = 0
	violaFK          = 787
	violaCheck       = 275
	violaNotNull     = 1299
	violaUnique      = 2067
	violaPK          = 1555
	ahora            = "'2026-10-23T00:00:00Z'"
	insertarSprint   = "INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES "
	insertarHistoria = "INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, sprint_id, fecha_finalizacion) VALUES "
)

// Datos base de cada caso: proyectos 1 y 2, usuario 1 (Scrum Master del 1),
// Sprints 10 (Activo) y 11 (Finalizado) del proyecto 1, 20 (Pendiente) y 21
// (Finalizado) del proyecto 2, historias 100 (Nueva) y 101 (En progreso en el 10).
const datosBase = `
INSERT INTO proyectos (id, nombre, creado_en) VALUES (1, 'P1', ` + ahora + `), (2, 'P2', ` + ahora + `);
INSERT INTO usuarios (id, nombre, nombre_normalizado, creado_en) VALUES (1, 'Ana', 'ana', ` + ahora + `);
INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (1, 1, 'scrum_master', ` + ahora + `);
INSERT INTO sprints (id, proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES
    (10, 1, 'g', '2026-10-01', '2026-10-15', 'Activo'),
    (11, 1, 'g', '2026-09-01', '2026-09-15', 'Finalizado'),
    (20, 2, NULL, NULL, NULL, 'Pendiente'),
    (21, 2, 'g', '2026-09-01', '2026-09-15', 'Finalizado');
INSERT INTO historias (id, proyecto_id, titulo, descripcion, prioridad, estado, sprint_id) VALUES
    (100, 1, 't', 'd', 'M', 'Nueva', NULL),
    (101, 1, 't', 'd', 'M', 'En progreso', 10);
`

func TestEsquema_Filas(t *testing.T) {
	casos := []struct {
		nombre    string
		sentencia string
		esperado  int
	}{
		// CA-6 proyectos
		{"proyecto fin igual al inicio", "INSERT INTO proyectos (nombre, fecha_inicio, fecha_fin, creado_en) VALUES ('x', '2026-10-01', '2026-10-01', " + ahora + ")", violaCheck},
		{"proyecto fin anterior al inicio", "INSERT INTO proyectos (nombre, fecha_inicio, fecha_fin, creado_en) VALUES ('x', '2026-10-02', '2026-10-01', " + ahora + ")", violaCheck},
		{"proyecto solo inicio", "INSERT INTO proyectos (nombre, fecha_inicio, creado_en) VALUES ('x', '2026-10-01', " + ahora + ")", aceptada},
		{"proyecto solo fin", "INSERT INTO proyectos (nombre, fecha_fin, creado_en) VALUES ('x', '2026-10-01', " + ahora + ")", aceptada},
		{"proyecto sin fechas", "INSERT INTO proyectos (nombre, fecha_inicio, fecha_fin, creado_en) VALUES ('x', NULL, NULL, " + ahora + ")", aceptada},
		{"horas por SP 0", "INSERT INTO proyectos (nombre, horas_por_story_point, creado_en) VALUES ('x', 0, " + ahora + ")", violaCheck},
		{"horas por SP negativas", "INSERT INTO proyectos (nombre, horas_por_story_point, creado_en) VALUES ('x', -1, " + ahora + ")", violaCheck},
		{"horas por SP texto", "INSERT INTO proyectos (nombre, horas_por_story_point, creado_en) VALUES ('x', 'abc', " + ahora + ")", violaCheck},
		{"horas por SP NULL", "INSERT INTO proyectos (nombre, horas_por_story_point, creado_en) VALUES ('x', NULL, " + ahora + ")", aceptada},
		{"horas por SP 2.5", "INSERT INTO proyectos (nombre, horas_por_story_point, creado_en) VALUES ('x', 2.5, " + ahora + ")", aceptada},

		// CA-7 sprints (los Activo van al proyecto 2, que no tiene uno)
		{"Activo sin fecha_fin", insertarSprint + "(2, 'g', '2026-10-01', NULL, 'Activo')", violaCheck},
		{"Activo sin fecha_inicio", insertarSprint + "(2, 'g', NULL, '2026-10-15', 'Activo')", violaCheck},
		{"Activo sin goal", insertarSprint + "(2, NULL, '2026-10-01', '2026-10-15', 'Activo')", violaCheck},
		{"Finalizado sin fecha_fin", insertarSprint + "(1, 'g', '2026-10-01', NULL, 'Finalizado')", violaCheck},
		{"Finalizado sin fecha_inicio", insertarSprint + "(1, 'g', NULL, '2026-10-15', 'Finalizado')", violaCheck},
		{"Finalizado sin goal", insertarSprint + "(1, NULL, '2026-10-01', '2026-10-15', 'Finalizado')", violaCheck},
		{"Pendiente con goal", insertarSprint + "(1, 'g', NULL, NULL, 'Pendiente')", violaCheck},
		{"Pendiente con fecha_inicio", insertarSprint + "(1, NULL, '2026-10-01', NULL, 'Pendiente')", violaCheck},
		{"Pendiente con fecha_fin", insertarSprint + "(1, NULL, NULL, '2026-10-15', 'Pendiente')", violaCheck},
		{"Sprint fin igual al inicio", insertarSprint + "(2, 'g', '2026-10-01', '2026-10-01', 'Activo')", violaCheck},
		{"Sprint fin anterior al inicio", insertarSprint + "(2, 'g', '2026-10-02', '2026-10-01', 'Activo')", violaCheck},

		// CA-8 historias
		{"En progreso sin Sprint", insertarHistoria + "(1, 't', 'd', 'M', 'En progreso', NULL, NULL)", violaCheck},
		{"Completada sin Sprint", insertarHistoria + "(1, 't', 'd', 'M', 'Completada', NULL, '2026-10-20')", violaCheck},
		{"Completada sin fecha_finalizacion", insertarHistoria + "(1, 't', 'd', 'M', 'Completada', 10, NULL)", violaCheck},
		{"Nueva con fecha_finalizacion", insertarHistoria + "(1, 't', 'd', 'M', 'Nueva', NULL, '2026-10-20')", violaCheck},
		{"En progreso con fecha_finalizacion", insertarHistoria + "(1, 't', 'd', 'M', 'En progreso', 10, '2026-10-20')", violaCheck},
		{"historia en Sprint de otro proyecto", insertarHistoria + "(1, 't', 'd', 'M', 'Nueva', 20, NULL)", violaFK},
		{"historia sin Sprint", insertarHistoria + "(1, 't', 'd', 'M', 'Nueva', NULL, NULL)", aceptada},

		// CA-9 índices parciales
		{"segundo Activo por INSERT", insertarSprint + "(1, 'g', '2026-11-01', '2026-11-15', 'Activo')", violaUnique},
		{"segundo Pendiente por INSERT", insertarSprint + "(2, NULL, NULL, NULL, 'Pendiente')", violaUnique},
		{"segundo Activo por UPDATE", "UPDATE sprints SET estado = 'Activo' WHERE id = 11", violaUnique},
		{"segundo Pendiente por UPDATE", "UPDATE sprints SET estado = 'Pendiente', sprint_goal = NULL, fecha_inicio = NULL, fecha_fin = NULL WHERE id = 21", violaUnique},
		{"varios Finalizado", insertarSprint + "(1, 'g', '2026-08-01', '2026-08-15', 'Finalizado')", aceptada},

		// CA-10 fechas (D12)
		{"proyecto fecha_inicio sin cero", "INSERT INTO proyectos (nombre, fecha_inicio, creado_en) VALUES ('x', '2026-10-9', " + ahora + ")", violaCheck},
		{"proyecto fecha_inicio basura", "INSERT INTO proyectos (nombre, fecha_inicio, creado_en) VALUES ('x', 'basura', " + ahora + ")", violaCheck},
		{"proyecto fecha_inicio inexistente", "INSERT INTO proyectos (nombre, fecha_inicio, creado_en) VALUES ('x', '2026-02-30', " + ahora + ")", violaCheck},
		{"proyecto fecha_inicio válida", "INSERT INTO proyectos (nombre, fecha_inicio, creado_en) VALUES ('x', '2026-10-09', " + ahora + ")", aceptada},
		{"proyecto fecha_fin sin cero", "INSERT INTO proyectos (nombre, fecha_fin, creado_en) VALUES ('x', '2026-10-9', " + ahora + ")", violaCheck},
		{"proyecto fecha_fin basura", "INSERT INTO proyectos (nombre, fecha_fin, creado_en) VALUES ('x', 'basura', " + ahora + ")", violaCheck},
		{"proyecto fecha_fin inexistente", "INSERT INTO proyectos (nombre, fecha_fin, creado_en) VALUES ('x', '2026-02-30', " + ahora + ")", violaCheck},
		{"proyecto fecha_fin válida", "INSERT INTO proyectos (nombre, fecha_fin, creado_en) VALUES ('x', '2026-10-09', " + ahora + ")", aceptada},
		{"Sprint fecha_inicio sin cero", insertarSprint + "(1, 'g', '2026-10-9', '2026-11-01', 'Finalizado')", violaCheck},
		{"Sprint fecha_inicio basura", insertarSprint + "(1, 'g', 'basura', '2026-11-01', 'Finalizado')", violaCheck},
		{"Sprint fecha_inicio inexistente", insertarSprint + "(1, 'g', '2026-02-30', '2026-03-01', 'Finalizado')", violaCheck},
		{"Sprint fecha_inicio válida", insertarSprint + "(1, 'g', '2026-10-09', '2026-11-01', 'Finalizado')", aceptada},
		{"Sprint fecha_fin sin cero", insertarSprint + "(1, 'g', '2026-01-01', '2026-10-9', 'Finalizado')", violaCheck},
		{"Sprint fecha_fin basura", insertarSprint + "(1, 'g', '2026-01-01', 'basura', 'Finalizado')", violaCheck},
		{"Sprint fecha_fin inexistente", insertarSprint + "(1, 'g', '2026-01-01', '2026-02-30', 'Finalizado')", violaCheck},
		{"Sprint fecha_fin válida", insertarSprint + "(1, 'g', '2026-01-01', '2026-10-09', 'Finalizado')", aceptada},
		{"fecha_finalizacion sin cero", insertarHistoria + "(1, 't', 'd', 'M', 'Completada', 10, '2026-10-9')", violaCheck},
		{"fecha_finalizacion basura", insertarHistoria + "(1, 't', 'd', 'M', 'Completada', 10, 'basura')", violaCheck},
		{"fecha_finalizacion inexistente", insertarHistoria + "(1, 't', 'd', 'M', 'Completada', 10, '2026-02-30')", violaCheck},
		{"fecha_finalizacion válida", insertarHistoria + "(1, 't', 'd', 'M', 'Completada', 10, '2026-10-09')", aceptada},

		// CA-11 enums y obligatoriedad
		{"rol fuera del enum", "INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (2, 1, 'admin', " + ahora + ")", violaCheck},
		{"prioridad fuera del enum", insertarHistoria + "(1, 't', 'd', 'X', 'Nueva', NULL, NULL)", violaCheck},
		{"estado de historia fuera del enum", insertarHistoria + "(1, 't', 'd', 'M', 'Cerrada', 10, NULL)", violaCheck},
		{"estado de Sprint fuera del enum", insertarSprint + "(1, 'g', '2026-11-01', '2026-11-15', 'Cancelado')", violaCheck},
		{"proyectos.nombre NULL", "INSERT INTO proyectos (nombre, creado_en) VALUES (NULL, " + ahora + ")", violaNotNull},
		{"proyectos.descripcion NULL", "INSERT INTO proyectos (nombre, descripcion, creado_en) VALUES ('x', NULL, " + ahora + ")", violaNotNull},
		{"proyectos.creado_en NULL", "INSERT INTO proyectos (nombre, creado_en) VALUES ('x', NULL)", violaNotNull},
		{"usuarios.nombre NULL", "INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES (NULL, 'b', " + ahora + ")", violaNotNull},
		{"usuarios.nombre_normalizado NULL", "INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES ('B', NULL, " + ahora + ")", violaNotNull},
		{"usuarios.creado_en NULL", "INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES ('B', 'b', NULL)", violaNotNull},
		{"integrantes.proyecto_id NULL", "INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (NULL, 1, 'scrum_master', " + ahora + ")", violaNotNull},
		{"integrantes.usuario_id NULL", "INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (2, NULL, 'scrum_master', " + ahora + ")", violaNotNull},
		{"integrantes.rol NULL", "INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (2, 1, NULL, " + ahora + ")", violaNotNull},
		{"integrantes.creado_en NULL", "INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (2, 1, 'scrum_master', NULL)", violaNotNull},
		{"historias.proyecto_id NULL", insertarHistoria + "(NULL, 't', 'd', 'M', 'Nueva', NULL, NULL)", violaNotNull},
		{"historias.titulo NULL", insertarHistoria + "(1, NULL, 'd', 'M', 'Nueva', NULL, NULL)", violaNotNull},
		{"historias.descripcion NULL", insertarHistoria + "(1, 't', NULL, 'M', 'Nueva', NULL, NULL)", violaNotNull},
		{"historias.prioridad NULL", insertarHistoria + "(1, 't', 'd', NULL, 'Nueva', NULL, NULL)", violaNotNull},
		{"historias.estado NULL", insertarHistoria + "(1, 't', 'd', 'M', NULL, NULL, NULL)", violaNotNull},
		{"sprints.proyecto_id NULL", insertarSprint + "(NULL, NULL, NULL, NULL, 'Pendiente')", violaNotNull},
		{"sprints.estado NULL", insertarSprint + "(1, NULL, NULL, NULL, NULL)", violaNotNull},

		// CA-12 cierre de Sprint
		{"cierre con UPDATE único", "UPDATE historias SET sprint_id = NULL, estado = 'Nueva' WHERE id = 101", aceptada},
		{"cierre quitando primero el Sprint", "UPDATE historias SET sprint_id = NULL WHERE id = 101", violaCheck},

		// D9: sin cascadas
		{"borrar proyecto con historias", "DELETE FROM proyectos WHERE id = 1", violaFK},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			db := nuevaBDMigrada(t)
			if _, err := db.Exec(datosBase); err != nil {
				t.Fatalf("no se pudieron cargar los datos base: %v", err)
			}

			_, err := db.Exec(c.sentencia)

			if got := codigoDeViolacion(t, err); got != c.esperado {
				t.Fatalf("se esperaba el código %d, se obtuvo %d (%v)", c.esperado, got, err)
			}
		})
	}
}

const (
	insertarRonda = "INSERT INTO rondas_votacion (historia_id, numero, revelada_en) VALUES "
	insertarVoto  = "INSERT INTO votos (ronda_id, usuario_id, valor) VALUES "
)

// Sobre datosBase: usuario 2 (Product Builder del proyecto 1) y la ronda 500
// (número 1, sin revelar) de la historia 100, con el voto del usuario 1.
const datosVotacion = `
INSERT INTO usuarios (id, nombre, nombre_normalizado, creado_en) VALUES (2, 'Candela', 'candela', ` + ahora + `);
INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (1, 2, 'product_builder', ` + ahora + `);
INSERT INTO rondas_votacion (id, historia_id, numero, revelada_en) VALUES (500, 100, 1, NULL);
INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (500, 1, 8);
`

func TestEsquema_FilasVotacion(t *testing.T) {
	casos := []struct {
		nombre    string
		sentencia string
		esperado  int
	}{
		// §4.7 rondas_votacion
		{"ronda sin revelar", insertarRonda + "(101, 1, NULL)", aceptada},
		{"ronda revelada", insertarRonda + "(101, 1, '2026-10-27T15:04:05Z')", aceptada},
		{"ronda siguiente de la misma historia", insertarRonda + "(100, 2, NULL)", aceptada},
		{"ronda con número repetido", insertarRonda + "(100, 1, NULL)", violaUnique},
		{"ronda de historia inexistente", insertarRonda + "(999, 1, NULL)", violaFK},
		{"rondas_votacion.historia_id NULL", insertarRonda + "(NULL, 1, NULL)", violaNotNull},
		{"rondas_votacion.numero NULL", insertarRonda + "(101, NULL, NULL)", violaNotNull},

		// §4.7 votos
		{"voto de otro integrante", insertarVoto + "(500, 2, 5)", aceptada},
		{"voto fuera de la escala (lo valida el dominio, D11)", insertarVoto + "(500, 2, 4)", aceptada},
		{"segundo voto del mismo integrante", insertarVoto + "(500, 1, 3)", violaPK},
		{"voto en ronda inexistente", insertarVoto + "(999, 2, 5)", violaFK},
		{"voto de usuario inexistente", insertarVoto + "(500, 999, 5)", violaFK},
		{"votos.valor NULL", insertarVoto + "(500, 2, NULL)", violaNotNull},
		// La PK compuesta de SQLite admite NULL sin el NOT NULL explícito.
		{"votos.ronda_id NULL", insertarVoto + "(NULL, 2, 5)", violaNotNull},
		{"votos.usuario_id NULL", insertarVoto + "(500, NULL, 5)", violaNotNull},

		// D9: sin cascadas
		{"borrar ronda con votos", "DELETE FROM rondas_votacion WHERE id = 500", violaFK},
		{"borrar historia con rondas", "DELETE FROM historias WHERE id = 100", violaFK},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			db := nuevaBDMigrada(t)
			if _, err := db.Exec(datosBase + datosVotacion); err != nil {
				t.Fatalf("no se pudieron cargar los datos base: %v", err)
			}

			_, err := db.Exec(c.sentencia)

			if got := codigoDeViolacion(t, err); got != c.esperado {
				t.Fatalf("se esperaba el código %d, se obtuvo %d (%v)", c.esperado, got, err)
			}
		})
	}
}

func codigoDeViolacion(t *testing.T, err error) int {
	t.Helper()
	if err == nil {
		return aceptada
	}
	var sqliteErr *sqlite.Error
	if !errors.As(err, &sqliteErr) {
		t.Fatalf("error que no es de SQLite: %v", err)
	}
	return sqliteErr.Code()
}
