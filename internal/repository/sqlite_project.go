package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// SQLiteProjectRepository implementa ProjectRepository sobre SQLite.
type SQLiteProjectRepository struct {
	db *sql.DB
}

func NewSQLiteProjectRepository(db *sql.DB) *SQLiteProjectRepository {
	return &SQLiteProjectRepository{db: db}
}

// Guardar persiste el proyecto y devuelve una copia con el ID asignado.
func (r *SQLiteProjectRepository) Guardar(ctx context.Context, p domain.Project) (domain.Project, error) {
	return insertarProyecto(ctx, r.db, p)
}

// GuardarConScrumMaster inserta el proyecto y su membresía de Scrum Master en
// una sola transacción, de modo que si la segunda escritura falla no queda un
// proyecto sin Scrum Master.
func (r *SQLiteProjectRepository) GuardarConScrumMaster(ctx context.Context, p domain.Project, creadorID int64) (domain.Project, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Project{}, fmt.Errorf("iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	p, err = insertarProyecto(ctx, tx, p)
	if err != nil {
		return domain.Project{}, err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, ?, ?)`,
		p.ID, creadorID, string(domain.RolScrumMaster), time.Now().UTC().Format(time.RFC3339),
	); err != nil {
		return domain.Project{}, fmt.Errorf("vincular scrum master: %w", errorDeEscritura(err))
	}

	if err := tx.Commit(); err != nil {
		return domain.Project{}, fmt.Errorf("confirmar transacción: %w", err)
	}
	return p, nil
}

// ejecutorSQL abstrae *sql.DB y *sql.Tx para compartir el INSERT de proyectos.
type ejecutorSQL interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// insertarProyecto inserta el proyecto en el ejecutor dado y asigna el ID.
func insertarProyecto(ctx context.Context, ex ejecutorSQL, p domain.Project) (domain.Project, error) {
	if p.CreadoEn.IsZero() {
		p.CreadoEn = time.Now().UTC()
	}
	res, err := ex.ExecContext(ctx,
		`INSERT INTO proyectos (nombre, descripcion, fecha_inicio, fecha_fin, creado_en)
		 VALUES (?, ?, ?, ?, ?)`,
		p.Nombre, p.Descripcion,
		tiempoOpcionalAValor(p.FechaInicio), tiempoOpcionalAValor(p.FechaFin),
		p.CreadoEn.Format(time.RFC3339),
	)
	if err != nil {
		return domain.Project{}, fmt.Errorf("crear proyecto: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Project{}, fmt.Errorf("obtener id del proyecto: %w", err)
	}
	p.ID = id
	return p, nil
}

// ObtenerPorID devuelve el proyecto con el identificador indicado.
func (r *SQLiteProjectRepository) ObtenerPorID(ctx context.Context, id int64) (domain.Project, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, nombre, descripcion, fecha_inicio, fecha_fin, creado_en, horas_por_story_point, fecha_cierre FROM proyectos WHERE id = ?`,
		id,
	)
	return escanearProyecto(row)
}

// ListarPorIntegrante devuelve los proyectos que integra usuarioID con el rol
// de su fila en integrantes, ordenados por id.
func (r *SQLiteProjectRepository) ListarPorIntegrante(ctx context.Context, usuarioID int64) ([]domain.ProyectoConRol, error) {
	filas, err := r.db.QueryContext(ctx,
		`SELECT p.id, p.nombre, p.descripcion, p.fecha_inicio, p.fecha_fin, p.creado_en, p.horas_por_story_point, p.fecha_cierre, i.rol
		 FROM proyectos p
		 JOIN integrantes i ON i.proyecto_id = p.id
		 WHERE i.usuario_id = ?
		 ORDER BY p.id`,
		usuarioID,
	)
	if err != nil {
		return nil, fmt.Errorf("listar proyectos del integrante: %w", err)
	}
	defer filas.Close()

	proyectos := make([]domain.ProyectoConRol, 0)
	for filas.Next() {
		var rol domain.Role
		p, err := escanearProyecto(conRol{filas, &rol})
		if err != nil {
			return nil, err
		}
		proyectos = append(proyectos, domain.ProyectoConRol{Proyecto: p, Rol: rol})
	}
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("listar proyectos del integrante: %w", err)
	}
	return proyectos, nil
}

// conRol lee, después de las columnas del proyecto, el rol del integrante.
type conRol struct {
	fila escaneable
	rol  *domain.Role
}

func (c conRol) Scan(destinos ...any) error {
	return c.fila.Scan(append(destinos, c.rol)...)
}

// Actualizar reemplaza los datos editables del proyecto.
func (r *SQLiteProjectRepository) Actualizar(ctx context.Context, p domain.Project) error {
	resultado, err := r.db.ExecContext(ctx,
		`UPDATE proyectos SET nombre = ?, descripcion = ?, fecha_inicio = ?, fecha_fin = ?, horas_por_story_point = ? WHERE id = ?`,
		p.Nombre, p.Descripcion,
		tiempoOpcionalAValor(p.FechaInicio), tiempoOpcionalAValor(p.FechaFin),
		p.HorasPorStoryPoint, p.ID,
	)
	if err != nil {
		return fmt.Errorf("actualizar proyecto: %w", err)
	}
	afectadas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("actualizar proyecto: %w", err)
	}
	if afectadas == 0 {
		return domain.ErrProyectoNoEncontrado
	}
	return nil
}

// Cerrar guarda la fecha de cierre con un UPDATE condicional, así dos cierres
// simultáneos, o un Sprint que se inicia mientras tanto, no pasan los dos
// (#300, como el inicio de Sprint en #158). Si no cerró, explica por qué.
func (r *SQLiteProjectRepository) Cerrar(ctx context.Context, id int64, fecha time.Time) error {
	resultado, err := r.db.ExecContext(ctx,
		`UPDATE proyectos SET fecha_cierre = ?
		 WHERE id = ? AND fecha_cierre IS NULL
		   AND NOT EXISTS (SELECT 1 FROM sprints WHERE proyecto_id = proyectos.id AND estado = 'Activo')`,
		tiempoOpcionalAValor(&fecha), id,
	)
	if err != nil {
		return fmt.Errorf("cerrar proyecto: %w", err)
	}
	afectadas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("cerrar proyecto: %w", err)
	}
	if afectadas == 1 {
		return nil
	}
	return r.motivoDeNoCerrar(ctx, id)
}

// ProyectoCerrado informa si el proyecto tiene fecha de cierre (#300).
func (r *SQLiteProjectRepository) ProyectoCerrado(ctx context.Context, id int64) (bool, error) {
	var cerrado bool
	err := r.db.QueryRowContext(ctx, `SELECT fecha_cierre IS NOT NULL FROM proyectos WHERE id = ?`, id).Scan(&cerrado)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("consultar cierre del proyecto: %w", err)
	}
	return cerrado, nil
}

func (r *SQLiteProjectRepository) motivoDeNoCerrar(ctx context.Context, id int64) error {
	proyecto, err := r.ObtenerPorID(ctx, id)
	if err != nil {
		return err
	}
	if proyecto.FechaCierre != nil {
		return domain.ErrProyectoYaCerrado
	}
	return domain.ErrSprintActivoAlCerrarProyecto
}

// AgregarIntegrante vincula un integrante con un rol al proyecto y devuelve la
// membresía guardada (con su CreadoEn asignado).
func (r *SQLiteProjectRepository) AgregarIntegrante(ctx context.Context, m domain.Membership) (domain.Membership, error) {
	if m.CreadoEn.IsZero() {
		m.CreadoEn = time.Now().UTC()
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES (?, ?, ?, ?)`,
		m.ProjectID, m.UserID, string(m.Role), m.CreadoEn.Format(time.RFC3339),
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.Code() {
			case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
				return domain.Membership{}, ErrMiembroDuplicado
			case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
				return domain.Membership{}, domain.ErrScrumMasterYaExiste
			}
		}
		return domain.Membership{}, fmt.Errorf("agregar integrante: %w", errorDeEscritura(err))
	}
	return m, nil
}

// ListarIntegrantes devuelve los integrantes del proyecto con su nombre y rol.
func (r *SQLiteProjectRepository) ListarIntegrantes(ctx context.Context, proyectoID int64) ([]domain.Member, error) {
	filas, err := r.db.QueryContext(ctx,
		`SELECT i.usuario_id, i.proyecto_id, u.nombre, i.rol, i.creado_en
		 FROM integrantes i
		 JOIN usuarios u ON u.id = i.usuario_id
		 WHERE i.proyecto_id = ?
		 ORDER BY i.rowid`,
		proyectoID,
	)
	if err != nil {
		return nil, fmt.Errorf("listar integrantes: %w", err)
	}
	defer filas.Close()

	var miembros []domain.Member
	for filas.Next() {
		var (
			m         domain.Member
			createdAt string
		)
		if err := filas.Scan(&m.UserID, &m.ProjectID, &m.Nombre, &m.Role, &createdAt); err != nil {
			return nil, fmt.Errorf("leer integrante: %w", err)
		}
		m.CreadoEn = parsearTiempo(createdAt)
		miembros = append(miembros, m)
	}
	if err := filas.Err(); err != nil {
		return nil, fmt.Errorf("listar integrantes: %w", err)
	}
	return miembros, nil
}

func escanearProyecto(fila escaneable) (domain.Project, error) {
	var (
		p                              domain.Project
		inicio, fin, createdAt, cierre sql.NullString
		horas                          sql.NullFloat64
	)
	if err := fila.Scan(&p.ID, &p.Nombre, &p.Descripcion, &inicio, &fin, &createdAt, &horas, &cierre); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, domain.ErrProyectoNoEncontrado
		}
		return domain.Project{}, fmt.Errorf("leer proyecto: %w", err)
	}
	p.FechaInicio = valorATiempoOpcional(inicio)
	p.FechaFin = valorATiempoOpcional(fin)
	p.FechaCierre = valorATiempoOpcional(cierre)
	p.CreadoEn = parsearTiempo(createdAt.String)
	if horas.Valid {
		p.HorasPorStoryPoint = &horas.Float64
	}
	return p, nil
}

func tiempoOpcionalAValor(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.DateOnly)
}

func valorATiempoOpcional(v sql.NullString) *time.Time {
	if !v.Valid {
		return nil
	}
	t, err := time.Parse(time.DateOnly, v.String)
	if err != nil {
		return nil
	}
	return &t
}
