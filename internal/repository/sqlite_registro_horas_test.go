package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// integrante inserta un usuario y lo asigna al proyecto. Devuelve su id.
func integrante(t *testing.T, db *sql.DB, proyectoID int64, nombre string) int64 {
	t.Helper()
	usuarioID := usuarioSinProyecto(t, db, nombre)
	if _, err := db.Exec(
		"INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, 'product_builder', '2026-10-01T00:00:00Z')",
		proyectoID, usuarioID,
	); err != nil {
		t.Fatalf("no se pudo asignar el integrante: %v", err)
	}
	return usuarioID
}

func usuarioSinProyecto(t *testing.T, db *sql.DB, nombre string) int64 {
	t.Helper()
	var id int64
	err := db.QueryRow(
		"INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES (?, lower(?), '2026-10-01T00:00:00Z') RETURNING id",
		nombre, nombre,
	).Scan(&id)
	if err != nil {
		t.Fatalf("no se pudo sembrar el usuario: %v", err)
	}
	return id
}

type escenarioHoras struct {
	db        *sql.DB
	repo      *SQLiteRegistroHorasRepository
	sprints   *SQLiteSprintRepository
	activo    domain.Sprint
	historia  domain.BacklogItem
	usuarioID int64
}

// nuevoEscenarioHoras siembra un proyecto con la integrante Ana, un Sprint
// Activo del 2026-10-05 al 2026-10-19 y una historia Nueva asignada a él.
func nuevoEscenarioHoras(t *testing.T) escenarioHoras {
	t.Helper()
	db := nuevaBDMigrada(t)
	sprints := NewSQLiteSprintRepository(db)
	proyectoID := guardarProyecto(t, db)
	activo := iniciarSprint(t, sprints, guardarSprint(t, sprints, proyectoID))
	return escenarioHoras{
		db:        db,
		repo:      NewSQLiteRegistroHorasRepository(db),
		sprints:   sprints,
		activo:    activo,
		historia:  historiaEnSprint(t, db, proyectoID, activo.ID, domain.EstadoNueva, nil),
		usuarioID: integrante(t, db, proyectoID, "Ana"),
	}
}

func (e escenarioHoras) carga(fecha string, minutos int) domain.RegistroHoras {
	dia, _ := time.Parse(time.DateOnly, fecha)
	return domain.RegistroHoras{
		ProyectoID: e.historia.ProyectoID,
		HistoriaID: e.historia.ID,
		UsuarioID:  e.usuarioID,
		SprintID:   e.activo.ID,
		Fecha:      dia,
		Actividad:  "Implementar",
		Minutos:    minutos,
	}
}

func cantidadDeCargas(t *testing.T, db *sql.DB) int {
	t.Helper()
	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM registros_horas").Scan(&n); err != nil {
		t.Fatalf("no se pudieron contar las cargas: %v", err)
	}
	return n
}

func TestSQLiteRegistroHorasRepository_RegistrarPersisteLosCampos(t *testing.T) {
	e := nuevoEscenarioHoras(t)

	guardado, err := e.repo.Registrar(context.Background(), e.carga("2026-10-10", 240))

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if guardado.ID <= 0 {
		t.Fatalf("se esperaba un ID asignado, se obtuvo %d", guardado.ID)
	}
	var (
		proyectoID, historiaID, usuarioID, sprintID int64
		fecha, actividad, creadoEn                  string
		minutos                                     int
	)
	err = e.db.QueryRow(
		`SELECT proyecto_id, historia_id, usuario_id, sprint_id, fecha, actividad, minutos, creado_en
		 FROM registros_horas WHERE id = ?`, guardado.ID,
	).Scan(&proyectoID, &historiaID, &usuarioID, &sprintID, &fecha, &actividad, &minutos, &creadoEn)
	if err != nil {
		t.Fatalf("no se pudo leer la carga: %v", err)
	}
	if proyectoID != e.historia.ProyectoID || historiaID != e.historia.ID || usuarioID != e.usuarioID || sprintID != e.activo.ID {
		t.Errorf("ids: se obtuvo proyecto %d, historia %d, usuario %d, sprint %d", proyectoID, historiaID, usuarioID, sprintID)
	}
	if fecha != "2026-10-10" || actividad != "Implementar" || minutos != 240 {
		t.Errorf("datos: se obtuvo fecha %q, actividad %q, minutos %d", fecha, actividad, minutos)
	}
	if _, err := time.Parse(time.RFC3339, creadoEn); err != nil || creadoEn[len(creadoEn)-1] != 'Z' {
		t.Errorf("creado_en: se esperaba RFC 3339 en UTC, se obtuvo %q", creadoEn)
	}
}

func TestSQLiteRegistroHorasRepository_RegistrarConHistoriaFueraDelSprint(t *testing.T) {
	e := nuevoEscenarioHoras(t)
	if _, err := e.db.Exec("UPDATE historias SET sprint_id = NULL WHERE id = ?", e.historia.ID); err != nil {
		t.Fatalf("no se pudo sacar la historia del Sprint: %v", err)
	}

	_, err := e.repo.Registrar(context.Background(), e.carga("2026-10-10", 60))

	if !errors.Is(err, domain.ErrHistoriaFueraDelSprintActivo) {
		t.Fatalf("se esperaba ErrHistoriaFueraDelSprintActivo, se obtuvo %v", err)
	}
	if n := cantidadDeCargas(t, e.db); n != 0 {
		t.Errorf("no se esperaba ninguna carga, hay %d", n)
	}
}

func TestSQLiteRegistroHorasRepository_RegistrarConSprintFinalizado(t *testing.T) {
	e := nuevoEscenarioHoras(t)
	if _, err := e.db.Exec("UPDATE sprints SET estado = 'Finalizado' WHERE id = ?", e.activo.ID); err != nil {
		t.Fatalf("no se pudo finalizar el Sprint: %v", err)
	}

	_, err := e.repo.Registrar(context.Background(), e.carga("2026-10-10", 60))

	if !errors.Is(err, domain.ErrHistoriaFueraDelSprintActivo) {
		t.Fatalf("se esperaba ErrHistoriaFueraDelSprintActivo, se obtuvo %v", err)
	}
	if n := cantidadDeCargas(t, e.db); n != 0 {
		t.Errorf("no se esperaba ninguna carga, hay %d", n)
	}
}

func TestSQLiteRegistroHorasRepository_RegistrarConUsuarioQueNoEsIntegrante(t *testing.T) {
	e := nuevoEscenarioHoras(t)
	carga := e.carga("2026-10-10", 60)
	carga.UsuarioID = usuarioSinProyecto(t, e.db, "Pedro")

	_, err := e.repo.Registrar(context.Background(), carga)

	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Fatalf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
}

func TestSQLiteRegistroHorasRepository_ListarPorHistoriaOrdenaPorFechaYCarga(t *testing.T) {
	e := nuevoEscenarioHoras(t)
	ctx := context.Background()
	otra := historiaEnSprint(t, e.db, e.historia.ProyectoID, e.activo.ID, domain.EstadoNueva, nil)
	deOtra := e.carga("2026-10-06", 15)
	deOtra.HistoriaID = otra.ID
	var ids []int64
	for _, c := range []domain.RegistroHoras{e.carga("2026-10-10", 30), deOtra, e.carga("2026-10-08", 90), e.carga("2026-10-10", 120)} {
		guardado, err := e.repo.Registrar(ctx, c)
		if err != nil {
			t.Fatalf("no se esperaba error al registrar: %v", err)
		}
		ids = append(ids, guardado.ID)
	}

	registros, err := e.repo.ListarPorHistoria(ctx, e.historia.ID)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	esperados := []struct {
		id      int64
		fecha   string
		minutos int
	}{{ids[2], "2026-10-08", 90}, {ids[0], "2026-10-10", 30}, {ids[3], "2026-10-10", 120}}
	if len(registros) != len(esperados) {
		t.Fatalf("se esperaban %d cargas, se obtuvieron %d: %+v", len(esperados), len(registros), registros)
	}
	for i, esp := range esperados {
		r := registros[i]
		if r.ID != esp.id || r.Fecha.Format(time.DateOnly) != esp.fecha || r.Minutos != esp.minutos {
			t.Errorf("carga %d: se esperaba %+v, se obtuvo id %d, fecha %s, minutos %d",
				i, esp, r.ID, r.Fecha.Format(time.DateOnly), r.Minutos)
		}
		if r.NombreIntegrante != "Ana" || r.UsuarioID != e.usuarioID || r.SprintID != e.activo.ID ||
			r.HistoriaID != e.historia.ID || r.ProyectoID != e.historia.ProyectoID || r.Actividad != "Implementar" {
			t.Errorf("carga %d: datos inesperados %+v", i, r)
		}
	}
}

func TestSQLiteRegistroHorasRepository_ListarPorHistoriaSinCargas(t *testing.T) {
	e := nuevoEscenarioHoras(t)

	registros, err := e.repo.ListarPorHistoria(context.Background(), e.historia.ID)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if registros == nil || len(registros) != 0 {
		t.Fatalf("se esperaba un slice vacío no nil, se obtuvo %#v", registros)
	}
}

func TestSQLiteRegistroHorasRepository_DosCargasSeguidasQuedanLasDos(t *testing.T) {
	e := nuevoEscenarioHoras(t)
	ctx := context.Background()

	primera, err1 := e.repo.Registrar(ctx, e.carga("2026-10-10", 60))
	segunda, err2 := e.repo.Registrar(ctx, e.carga("2026-10-10", 60))

	if err1 != nil || err2 != nil {
		t.Fatalf("no se esperaban errores: %v, %v", err1, err2)
	}
	if primera.ID == segunda.ID {
		t.Errorf("se esperaban IDs distintos, ambos son %d", primera.ID)
	}
	if n := cantidadDeCargas(t, e.db); n != 2 {
		t.Errorf("se esperaban 2 cargas, hay %d", n)
	}
}
