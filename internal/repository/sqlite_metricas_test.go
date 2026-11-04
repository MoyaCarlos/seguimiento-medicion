package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type escenarioMetricas struct {
	t         *testing.T
	db        *sql.DB
	proyecto  int64
	usuarioID int64
}

func nuevoEscenarioMetricas(t *testing.T, db *sql.DB, nombreUsuario string) escenarioMetricas {
	t.Helper()
	proyecto := guardarProyecto(t, db)
	return escenarioMetricas{t: t, db: db, proyecto: proyecto, usuarioID: integrante(t, db, proyecto, nombreUsuario)}
}

func (e escenarioMetricas) insertar(consulta string, args ...any) int64 {
	e.t.Helper()
	resultado, err := e.db.Exec(consulta, args...)
	if err != nil {
		e.t.Fatalf("no se pudo sembrar %q: %v", consulta, err)
	}
	id, err := resultado.LastInsertId()
	if err != nil {
		e.t.Fatalf("no se pudo leer el id: %v", err)
	}
	return id
}

func (e escenarioMetricas) sprint(estado string) int64 {
	e.t.Helper()
	if estado == "Pendiente" {
		return e.insertar("INSERT INTO sprints (proyecto_id, estado) VALUES (?, 'Pendiente')", e.proyecto)
	}
	return e.insertar(
		"INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, 'Objetivo', '2026-10-05', '2026-10-19', ?)",
		e.proyecto, estado)
}

func (e escenarioMetricas) completada(sprintID int64, puntos any) int64 {
	e.t.Helper()
	return e.insertar(
		`INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp, sprint_id, fecha_finalizacion)
		 VALUES (?, 'Completada', '', 'M', 'Completada', ?, ?, '2026-10-10')`,
		e.proyecto, puntos, sprintID)
}

func (e escenarioMetricas) enBacklog(puntos any) int64 {
	e.t.Helper()
	return e.insertar(
		"INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp) VALUES (?, 'En el backlog', '', 'M', 'Nueva', ?)",
		e.proyecto, puntos)
}

func (e escenarioMetricas) arrastrada(sprintID, historiaID int64, foto any) {
	e.t.Helper()
	e.insertar("INSERT INTO historias_arrastradas (sprint_id, historia_id, proyecto_id, estimacion_sp) VALUES (?, ?, ?, ?)",
		sprintID, historiaID, e.proyecto, foto)
}

func (e escenarioMetricas) horas(historiaID, sprintID int64, minutos int) {
	e.t.Helper()
	e.insertar(
		`INSERT INTO registros_horas (proyecto_id, historia_id, usuario_id, sprint_id, fecha, actividad, minutos, creado_en)
		 VALUES (?, ?, ?, ?, '2026-10-10', 'Trabajo', ?, '2026-10-10T00:00:00Z')`,
		e.proyecto, historiaID, e.usuarioID, sprintID, minutos)
}

func (e escenarioMetricas) defecto(historiaID, deteccion int64, resolucion any) {
	e.t.Helper()
	estado := "Pendiente"
	if resolucion != nil {
		estado = "Resuelto"
	}
	e.insertar(
		`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id, creado_en)
		 VALUES (?, ?, 'Falla', 'Alta', ?, ?, ?, '2026-10-10T00:00:00Z')`,
		e.proyecto, historiaID, estado, deteccion, resolucion)
}

func listarDatosMetricas(t *testing.T, db *sql.DB, proyectoID int64) []domain.DatosMetricasSprint {
	t.Helper()
	datos, err := NewSQLiteMetricasRepository(db).ListarDatosSprintsFinalizados(context.Background(), proyectoID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	return datos
}

func TestSQLiteMetricas_SumaLasFuentesDeCadaSprintFinalizado(t *testing.T) {
	db := nuevaBDMigrada(t)
	e := nuevoEscenarioMetricas(t, db, "Ana")
	primero := e.sprint("Finalizado")
	segundo := e.sprint("Finalizado")
	activo := e.sprint("Activo")
	e.sprint("Pendiente")

	cinco := e.completada(primero, 5)
	e.completada(primero, nil)
	arrastrada := e.enBacklog(8)
	e.arrastrada(primero, arrastrada, 3)
	e.arrastrada(primero, e.enBacklog(nil), nil)
	e.completada(segundo, 2)

	e.horas(cinco, primero, 90)
	e.horas(arrastrada, primero, 30)
	e.horas(arrastrada, activo, 500)

	e.defecto(cinco, primero, primero)
	e.defecto(cinco, primero, segundo)
	e.defecto(cinco, activo, nil)
	e.defecto(cinco, primero, activo)

	otro := nuevoEscenarioMetricas(t, db, "Beto")
	otroSprint := otro.sprint("Finalizado")
	otraHistoria := otro.completada(otroSprint, 13)
	otro.horas(otraHistoria, otroSprint, 60)
	otro.defecto(otraHistoria, otroSprint, nil)

	datos := listarDatosMetricas(t, db, e.proyecto)

	if len(datos) != 2 {
		t.Fatalf("se esperaban los 2 Sprints finalizados, se obtuvieron %d: %+v", len(datos), datos)
	}
	esperadoPrimero := domain.DatosMetricasSprint{
		SprintID: primero, Numero: 1, HistoriasCompletadas: 2, HistoriasArrastradas: 2,
		SPCompletados: 5, SPArrastrados: 3, MinutosReales: 120, DefectosDetectados: 3, DefectosResueltos: 1,
	}
	if datos[0] != esperadoPrimero {
		t.Errorf("primer Sprint: se esperaba %+v, se obtuvo %+v", esperadoPrimero, datos[0])
	}
	esperadoSegundo := domain.DatosMetricasSprint{
		SprintID: segundo, Numero: 2, HistoriasCompletadas: 1, SPCompletados: 2, DefectosResueltos: 1,
	}
	if datos[1] != esperadoSegundo {
		t.Errorf("segundo Sprint: se esperaba %+v, se obtuvo %+v", esperadoSegundo, datos[1])
	}
}

func TestSQLiteMetricas_NumeroEntreTodosLosSprintsDelProyecto(t *testing.T) {
	db := nuevaBDMigrada(t)
	e := nuevoEscenarioMetricas(t, db, "Ana")
	otro := nuevoEscenarioMetricas(t, db, "Beto")

	e.sprint("Activo")
	finalizadoDelOtro := otro.sprint("Finalizado")
	finalizado := e.sprint("Finalizado")

	datos := listarDatosMetricas(t, db, e.proyecto)
	if len(datos) != 1 || datos[0].SprintID != finalizado || datos[0].Numero != 2 {
		t.Errorf("se esperaba el Sprint %d con número 2, se obtuvo %+v", finalizado, datos)
	}
	delOtro := listarDatosMetricas(t, db, otro.proyecto)
	if len(delOtro) != 1 || delOtro[0].SprintID != finalizadoDelOtro || delOtro[0].Numero != 1 {
		t.Errorf("se esperaba el Sprint %d con número 1, se obtuvo %+v", finalizadoDelOtro, delOtro)
	}
}

func TestSQLiteMetricas_SprintSinDatosDaCeros(t *testing.T) {
	db := nuevaBDMigrada(t)
	e := nuevoEscenarioMetricas(t, db, "Ana")
	sprint := e.sprint("Finalizado")

	datos := listarDatosMetricas(t, db, e.proyecto)

	esperado := domain.DatosMetricasSprint{SprintID: sprint, Numero: 1}
	if len(datos) != 1 || datos[0] != esperado {
		t.Errorf("se esperaba %+v, se obtuvo %+v", esperado, datos)
	}
}

func TestSQLiteMetricas_SinSprintsFinalizadosDevuelveVacio(t *testing.T) {
	db := nuevaBDMigrada(t)
	e := nuevoEscenarioMetricas(t, db, "Ana")
	e.sprint("Activo")

	datos := listarDatosMetricas(t, db, e.proyecto)

	if datos == nil || len(datos) != 0 {
		t.Errorf("se esperaba un slice vacío y no nil, se obtuvo %#v", datos)
	}
}
