package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type escenarioDefectos struct {
	db         *sql.DB
	repo       *SQLiteDefectoRepository
	proyectoID int64
	activo     int64
	finalizado int64
	historiaID int64
}

// nuevoEscenarioDefectos siembra un proyecto con un Sprint Activo, uno
// Finalizado y la historia "Carrito de compras" en el Product Backlog.
func nuevoEscenarioDefectos(t *testing.T) escenarioDefectos {
	t.Helper()
	db := nuevaBDMigrada(t)
	proyectoID := guardarProyecto(t, db)
	return escenarioDefectos{
		db:         db,
		repo:       NewSQLiteDefectoRepository(db),
		proyectoID: proyectoID,
		finalizado: sprintConEstado(t, db, proyectoID, domain.SprintFinalizado),
		activo:     sprintConEstado(t, db, proyectoID, domain.SprintActivo),
		historiaID: historiaDelBacklog(t, db, proyectoID, "Carrito de compras"),
	}
}

func sprintConEstado(t *testing.T, db *sql.DB, proyectoID int64, estado domain.EstadoSprint) int64 {
	t.Helper()
	var id int64
	err := db.QueryRow(
		"INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, 'g', '2026-10-05', '2026-10-19', ?) RETURNING id",
		proyectoID, string(estado),
	).Scan(&id)
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint %s: %v", estado, err)
	}
	return id
}

func historiaDelBacklog(t *testing.T, db *sql.DB, proyectoID int64, titulo string) int64 {
	t.Helper()
	var id int64
	err := db.QueryRow(
		"INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado) VALUES (?, ?, 'd', 'M', 'Nueva') RETURNING id",
		proyectoID, titulo,
	).Scan(&id)
	if err != nil {
		t.Fatalf("no se pudo sembrar la historia: %v", err)
	}
	return id
}

func cerrarSprintSQL(t *testing.T, db *sql.DB, sprintID int64) {
	t.Helper()
	if _, err := db.Exec("UPDATE sprints SET estado = 'Finalizado' WHERE id = ?", sprintID); err != nil {
		t.Fatalf("no se pudo cerrar el Sprint: %v", err)
	}
}

func (e escenarioDefectos) defecto(descripcion string, sprintID int64) domain.Defecto {
	return domain.Defecto{
		ProyectoID:        e.proyectoID,
		HistoriaID:        e.historiaID,
		SprintDeteccionID: sprintID,
		Descripcion:       descripcion,
		Severidad:         domain.SeveridadAlta,
		Estado:            domain.DefectoPendiente,
	}
}

func (e escenarioDefectos) registrar(t *testing.T, descripcion string) domain.Defecto {
	t.Helper()
	d, err := e.repo.Registrar(context.Background(), e.defecto(descripcion, e.activo))
	if err != nil {
		t.Fatalf("no se esperaba error al registrar: %v", err)
	}
	return d
}

func (e escenarioDefectos) resuelto(d domain.Defecto, sprintID int64) domain.Defecto {
	d.Estado, d.SprintResolucionID = domain.DefectoResuelto, &sprintID
	return d
}

func cantidadDeDefectos(t *testing.T, db *sql.DB) int {
	t.Helper()
	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM defectos").Scan(&n); err != nil {
		t.Fatalf("no se pudieron contar los defectos: %v", err)
	}
	return n
}

// estadoGuardado lee el estado y el Sprint de resolución persistidos.
func estadoGuardado(t *testing.T, db *sql.DB, id int64) (string, sql.NullInt64) {
	t.Helper()
	var (
		estado     string
		resolucion sql.NullInt64
	)
	if err := db.QueryRow("SELECT estado, sprint_resolucion_id FROM defectos WHERE id = ?", id).Scan(&estado, &resolucion); err != nil {
		t.Fatalf("no se pudo leer el defecto %d: %v", id, err)
	}
	return estado, resolucion
}

func TestSQLiteDefectoRepository_RegistrarPersisteLosCampos(t *testing.T) {
	e := nuevoEscenarioDefectos(t)

	d, err := e.repo.Registrar(context.Background(), e.defecto("El total no se actualiza", e.activo))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if d.ID <= 0 || d.CreadoEn.IsZero() {
		t.Fatalf("se esperaba el id y la fecha de registro asignados, se obtuvo %+v", d)
	}

	var (
		proyectoID, historiaID, deteccion int64
		descripcion, severidad, estado    string
		resolucion                        sql.NullInt64
		creadoEn                          string
	)
	err = e.db.QueryRow(
		`SELECT proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id, creado_en
		 FROM defectos WHERE id = ?`, d.ID,
	).Scan(&proyectoID, &historiaID, &descripcion, &severidad, &estado, &deteccion, &resolucion, &creadoEn)
	if err != nil {
		t.Fatalf("no se pudo leer el defecto: %v", err)
	}
	if proyectoID != e.proyectoID || historiaID != e.historiaID || descripcion != "El total no se actualiza" ||
		severidad != "Alta" || estado != "Pendiente" || deteccion != e.activo || resolucion.Valid {
		t.Errorf("fila inesperada: %d %d %q %q %q %d %v", proyectoID, historiaID, descripcion, severidad, estado, deteccion, resolucion)
	}
	instante, err := time.Parse(time.RFC3339, creadoEn)
	if err != nil || instante.Location() != time.UTC {
		t.Errorf("se esperaba creado_en en RFC 3339 UTC, se obtuvo %q", creadoEn)
	}
}

func TestSQLiteDefectoRepository_RegistrarConSprintNoActivo(t *testing.T) {
	e := nuevoEscenarioDefectos(t)

	_, err := e.repo.Registrar(context.Background(), e.defecto("El total no se actualiza", e.finalizado))

	if !errors.Is(err, domain.ErrSinSprintActivoParaDetectar) {
		t.Fatalf("se esperaba ErrSinSprintActivoParaDetectar, se obtuvo %v", err)
	}
	if n := cantidadDeDefectos(t, e.db); n != 0 {
		t.Errorf("no se esperaba ningún defecto, hay %d", n)
	}
}

func TestSQLiteDefectoRepository_RegistrarConHistoriaDeOtroProyecto(t *testing.T) {
	e := nuevoEscenarioDefectos(t)
	ajena := historiaDelBacklog(t, e.db, guardarProyecto(t, e.db), "Login")
	d := e.defecto("El total no se actualiza", e.activo)
	d.HistoriaID = ajena

	_, err := e.repo.Registrar(context.Background(), d)

	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Fatalf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
}

func TestSQLiteDefectoRepository_ObtenerPorID(t *testing.T) {
	e := nuevoEscenarioDefectos(t)
	registrado := e.registrar(t, "El total no se actualiza")

	got, err := e.repo.ObtenerPorID(context.Background(), registrado.ID)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if got.ID != registrado.ID || got.ProyectoID != e.proyectoID || got.HistoriaID != e.historiaID ||
		got.TituloHistoria != "Carrito de compras" || got.Descripcion != "El total no se actualiza" ||
		got.Severidad != domain.SeveridadAlta || got.Estado != domain.DefectoPendiente ||
		got.SprintDeteccionID != e.activo || got.SprintResolucionID != nil || !got.CreadoEn.Equal(registrado.CreadoEn) {
		t.Errorf("defecto inesperado: %+v", got)
	}
}

func TestSQLiteDefectoRepository_ObtenerPorIDInexistente(t *testing.T) {
	e := nuevoEscenarioDefectos(t)

	_, err := e.repo.ObtenerPorID(context.Background(), 999)

	if !errors.Is(err, domain.ErrDefectoNoEncontrado) {
		t.Fatalf("se esperaba ErrDefectoNoEncontrado, se obtuvo %v", err)
	}
}

func TestSQLiteDefectoRepository_ResolverPersisteElSprint(t *testing.T) {
	e := nuevoEscenarioDefectos(t)
	d := e.registrar(t, "El total no se actualiza")

	if err := e.repo.Resolver(context.Background(), e.resuelto(d, e.activo)); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	estado, resolucion := estadoGuardado(t, e.db, d.ID)
	if estado != "Resuelto" || !resolucion.Valid || resolucion.Int64 != e.activo {
		t.Errorf("se esperaba Resuelto en el Sprint %d, se obtuvo %s %v", e.activo, estado, resolucion)
	}
	got, err := e.repo.ObtenerPorID(context.Background(), d.ID)
	if err != nil || got.SprintResolucionID == nil || *got.SprintResolucionID != e.activo {
		t.Errorf("ObtenerPorID: se esperaba el Sprint de resolución %d, se obtuvo %+v (%v)", e.activo, got, err)
	}
}

func TestSQLiteDefectoRepository_ResolverDosVecesConservaLaPrimera(t *testing.T) {
	e := nuevoEscenarioDefectos(t)
	d := e.registrar(t, "El total no se actualiza")
	if err := e.repo.Resolver(context.Background(), e.resuelto(d, e.activo)); err != nil {
		t.Fatalf("no se esperaba error en la primera resolución: %v", err)
	}
	cerrarSprintSQL(t, e.db, e.activo)
	otro := sprintConEstado(t, e.db, e.proyectoID, domain.SprintActivo)

	err := e.repo.Resolver(context.Background(), e.resuelto(d, otro))

	if !errors.Is(err, domain.ErrDefectoYaResuelto) {
		t.Fatalf("se esperaba ErrDefectoYaResuelto, se obtuvo %v", err)
	}
	if _, resolucion := estadoGuardado(t, e.db, d.ID); resolucion.Int64 != e.activo {
		t.Errorf("se esperaba conservar el Sprint de resolución %d, quedó %v", e.activo, resolucion)
	}
}

func TestSQLiteDefectoRepository_ResolverConSprintNoActivo(t *testing.T) {
	e := nuevoEscenarioDefectos(t)
	d := e.registrar(t, "El total no se actualiza")
	cerrarSprintSQL(t, e.db, e.activo)

	err := e.repo.Resolver(context.Background(), e.resuelto(d, e.activo))

	if !errors.Is(err, domain.ErrSinSprintActivoParaResolver) {
		t.Fatalf("se esperaba ErrSinSprintActivoParaResolver, se obtuvo %v", err)
	}
	if estado, resolucion := estadoGuardado(t, e.db, d.ID); estado != "Pendiente" || resolucion.Valid {
		t.Errorf("se esperaba que siga Pendiente, quedó %s %v", estado, resolucion)
	}
}

func TestSQLiteDefectoRepository_ListarPorProyecto(t *testing.T) {
	e := nuevoEscenarioDefectos(t)
	primero := e.registrar(t, "El total no se actualiza")
	segundo := e.registrar(t, "El botón no responde")
	if err := e.repo.Resolver(context.Background(), e.resuelto(primero, e.activo)); err != nil {
		t.Fatalf("no se esperaba error al resolver: %v", err)
	}
	otroProyecto := guardarProyecto(t, e.db)
	ajeno := domain.Defecto{
		ProyectoID:        otroProyecto,
		HistoriaID:        historiaDelBacklog(t, e.db, otroProyecto, "Login"),
		SprintDeteccionID: sprintConEstado(t, e.db, otroProyecto, domain.SprintActivo),
		Descripcion:       "Defecto ajeno",
		Severidad:         domain.SeveridadBaja,
		Estado:            domain.DefectoPendiente,
	}
	if _, err := e.repo.Registrar(context.Background(), ajeno); err != nil {
		t.Fatalf("no se pudo sembrar el defecto del otro proyecto: %v", err)
	}

	got, err := e.repo.ListarPorProyecto(context.Background(), e.proyectoID)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(got) != 2 || got[0].ID != primero.ID || got[1].ID != segundo.ID {
		t.Fatalf("se esperaban los defectos %d y %d en orden, se obtuvo %+v", primero.ID, segundo.ID, got)
	}
	if got[0].TituloHistoria != "Carrito de compras" || got[0].Estado != domain.DefectoResuelto ||
		got[0].SprintResolucionID == nil || *got[0].SprintResolucionID != e.activo {
		t.Errorf("primer defecto inesperado: %+v", got[0])
	}
	if got[1].Estado != domain.DefectoPendiente || got[1].SprintResolucionID != nil || got[1].CreadoEn.IsZero() {
		t.Errorf("segundo defecto inesperado: %+v", got[1])
	}
}

func TestSQLiteDefectoRepository_ListarProyectoSinDefectos(t *testing.T) {
	e := nuevoEscenarioDefectos(t)

	got, err := e.repo.ListarPorProyecto(context.Background(), e.proyectoID)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("se esperaba un slice vacío no nil, se obtuvo %#v", got)
	}
}
