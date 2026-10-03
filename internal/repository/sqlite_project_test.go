package repository

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestAbrirSQLite_ForeignKeysEnTodasLasConexiones(t *testing.T) {
	archivo := filepath.Join(t.TempDir(), "prueba.db")
	db, err := AbrirSQLite(archivo)
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(2)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("no se pudo iniciar transacción: %v", err)
	}
	defer tx.Rollback()

	var fkTx, fkOtra int
	if err := tx.QueryRow("PRAGMA foreign_keys").Scan(&fkTx); err != nil {
		t.Fatalf("no se pudo consultar foreign_keys (tx): %v", err)
	}
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&fkOtra); err != nil {
		t.Fatalf("no se pudo consultar foreign_keys (otra conexión): %v", err)
	}
	if fkTx != 1 || fkOtra != 1 {
		t.Fatalf("se esperaba foreign_keys = 1 en ambas conexiones, se obtuvo %d y %d", fkTx, fkOtra)
	}
}

func abrirBDDePrueba(t *testing.T) *sql.DB {
	t.Helper()
	db, err := AbrirSQLite(":memory:")
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	if err := Migrar(context.Background(), db); err != nil {
		t.Fatalf("no se pudo migrar: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestSQLiteProjectRepository_CrearYLeer(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(abrirBDDePrueba(t))
	inicio := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	p, err := domain.NewProject("Software Metrics & Estimation", "Descripción", &inicio, nil)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := repo.Guardar(ctx, p)
	if err != nil {
		t.Fatalf("no se esperaba error al crear: %v", err)
	}
	if guardado.ID == 0 {
		t.Fatal("se esperaba un ID asignado")
	}

	leido, err := repo.ObtenerPorID(ctx, guardado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al leer: %v", err)
	}
	if leido.Nombre != guardado.Nombre || leido.Descripcion != "Descripción" {
		t.Errorf("proyecto inesperado: %+v", leido)
	}
	if leido.FechaInicio == nil || !leido.FechaInicio.Equal(inicio) {
		t.Errorf("se esperaba fecha de inicio %v, se obtuvo %v", inicio, leido.FechaInicio)
	}
	if leido.FechaFin != nil {
		t.Errorf("se esperaba fecha de fin nula, se obtuvo %v", leido.FechaFin)
	}
}

func TestSQLiteProjectRepository_IDsAutoincrementales(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(abrirBDDePrueba(t))

	var ids []int64
	for i := 0; i < 3; i++ {
		p, _ := domain.NewProject("Proyecto", "", nil, nil)
		guardado, err := repo.Guardar(ctx, p)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		ids = append(ids, guardado.ID)
	}
	if ids[0] <= 0 {
		t.Fatalf("se esperaba un ID positivo, se obtuvo %d", ids[0])
	}
	for i := 1; i < len(ids); i++ {
		if ids[i] <= ids[i-1] {
			t.Fatalf("se esperaban IDs crecientes, se obtuvieron %v", ids)
		}
	}
}

func TestSQLiteProjectRepository_GuardarConScrumMaster_Rollback(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(abrirBDDePrueba(t))

	p, _ := domain.NewProject("Proyecto", "", nil, nil)
	// creadorID inexistente: la FK de project_members falla en la segunda escritura.
	if _, err := repo.GuardarConScrumMaster(ctx, p, 999999); err == nil {
		t.Fatal("se esperaba error por FK del integrante inexistente")
	}

	var total int
	if err := repo.db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&total); err != nil {
		t.Fatal(err)
	}
	if total != 0 {
		t.Fatalf("se esperaban 0 proyectos (rollback de la transacción), se encontraron %d", total)
	}
}

func TestSQLiteProjectRepository_AgregarIntegrante_FKVioladaNoEsNoEncontrado(t *testing.T) {
	ctx := context.Background()
	proyectos := NewSQLiteProjectRepository(abrirBDDePrueba(t))

	p, _ := domain.NewProject("Proyecto", "", nil, nil)
	guardado, err := proyectos.Guardar(ctx, p)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	m := domain.Membership{ProjectID: guardado.ID, UserID: 999, Role: domain.RolScrumMaster}
	err = proyectos.AgregarIntegrante(ctx, m)
	if err == nil {
		t.Fatal("se esperaba error por FK")
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("una falla de FK no debe mapearse a 'no encontrado': %v", err)
	}
}

func TestSQLiteProjectRepository_NombresRepetidos(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(abrirBDDePrueba(t))

	for i := 0; i < 2; i++ {
		p, _ := domain.NewProject("Mismo nombre", "", nil, nil)
		if _, err := repo.Guardar(ctx, p); err != nil {
			t.Fatalf("no se esperaba error con nombres repetidos: %v", err)
		}
	}
	var total int
	if err := repo.db.QueryRow("SELECT COUNT(*) FROM projects WHERE name = ?", "Mismo nombre").Scan(&total); err != nil {
		t.Fatalf("no se pudo contar: %v", err)
	}
	if total != 2 {
		t.Errorf("se esperaban 2 proyectos con el mismo nombre, se obtuvieron %d", total)
	}
}

func TestSQLiteProjectRepository_Actualizar(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(abrirBDDePrueba(t))

	p, _ := domain.NewProject("Original", "desc", nil, nil)
	guardado, err := repo.Guardar(ctx, p)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	editado, err := guardado.ConDatosEditados("Corregido", "nueva desc", nil, nil)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	if err := repo.Actualizar(ctx, editado); err != nil {
		t.Fatalf("no se esperaba error al actualizar: %v", err)
	}

	leido, _ := repo.ObtenerPorID(ctx, guardado.ID)
	if leido.Nombre != "Corregido" || leido.Descripcion != "nueva desc" {
		t.Errorf("no se persistió la edición: %+v", leido)
	}
}

func TestSQLiteProjectRepository_ActualizarInexistente(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(abrirBDDePrueba(t))
	p, _ := domain.NewProject("Fantasma", "", nil, nil)
	p.ID = 999
	if err := repo.Actualizar(ctx, p); !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
}

func TestSQLiteProjectRepository_ListarIntegrantes_OrdenDeAlta(t *testing.T) {
	ctx := context.Background()
	db := abrirBDDePrueba(t)
	proyectos := NewSQLiteProjectRepository(db)
	usuarios := NewSQLiteUserRepository(db)

	p, _ := domain.NewProject("Proyecto", "", nil, nil)
	guardado, err := proyectos.Guardar(ctx, p)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	// Tres altas con el mismo created_at (resolución de segundos) en orden no
	// alfabético: el orden real de alta debe prevalecer.
	nombres := []string{"Jimena", "Ana", "Candela"}
	fecha := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, n := range nombres {
		u, _ := domain.NewUser(n)
		guardadoU, err := usuarios.Guardar(ctx, u)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		m, _ := domain.NewMembership(guardado.ID, guardadoU.ID, domain.RolProductBuilder)
		m.CreadoEn = fecha
		if err := proyectos.AgregarIntegrante(ctx, m); err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
	}

	miembros, err := proyectos.ListarIntegrantes(ctx, guardado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	for i, n := range nombres {
		if miembros[i].Nombre != n {
			t.Fatalf("se esperaba orden de alta %v, se obtuvo %q en la posición %d", nombres, miembros[i].Nombre, i)
		}
	}
}

func TestSQLiteProjectRepository_Miembros(t *testing.T) {
	ctx := context.Background()
	db := abrirBDDePrueba(t)
	proyectos := NewSQLiteProjectRepository(db)
	usuarios := NewSQLiteUserRepository(db)

	p, _ := domain.NewProject("Proyecto", "", nil, nil)
	guardado, err := proyectos.Guardar(ctx, p)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	u, _ := domain.NewUser("Candela")
	guardadoU, err := usuarios.Guardar(ctx, u)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	m, _ := domain.NewMembership(guardado.ID, guardadoU.ID, domain.RolScrumMaster)
	if err := proyectos.AgregarIntegrante(ctx, m); err != nil {
		t.Fatalf("no se esperaba error al agregar: %v", err)
	}

	miembros, err := proyectos.ListarIntegrantes(ctx, guardado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al listar: %v", err)
	}
	if len(miembros) != 1 {
		t.Fatalf("se esperaba 1 integrante, se obtuvieron %d", len(miembros))
	}
	if miembros[0].Nombre != "Candela" || miembros[0].Role != domain.RolScrumMaster {
		t.Errorf("integrante inesperado: %+v", miembros[0])
	}
}

func TestSQLiteProjectRepository_MiembroDuplicado(t *testing.T) {
	ctx := context.Background()
	db := abrirBDDePrueba(t)
	proyectos := NewSQLiteProjectRepository(db)
	usuarios := NewSQLiteUserRepository(db)

	p, _ := domain.NewProject("Proyecto", "", nil, nil)
	guardado, err := proyectos.Guardar(ctx, p)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	u, _ := domain.NewUser("Candela")
	guardadoU, err := usuarios.Guardar(ctx, u)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	m, _ := domain.NewMembership(guardado.ID, guardadoU.ID, domain.RolProductBuilder)
	if err := proyectos.AgregarIntegrante(ctx, m); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if err := proyectos.AgregarIntegrante(ctx, m); !errors.Is(err, ErrMiembroDuplicado) {
		t.Fatalf("se esperaba ErrMiembroDuplicado, se obtuvo %v", err)
	}
}
