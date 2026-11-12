package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
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

func TestSQLiteProjectRepository_CrearYLeer(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(nuevaBDMigrada(t))
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

func TestSQLiteProjectRepository_FechasComoDia(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(nuevaBDMigrada(t))
	inicio := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	p, err := domain.NewProject("Proyecto", "", &inicio, &fin)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := repo.Guardar(ctx, p)
	if err != nil {
		t.Fatalf("no se esperaba error al crear: %v", err)
	}

	var fechaInicio, fechaFin string
	if err := repo.db.QueryRow(
		"SELECT fecha_inicio, fecha_fin FROM proyectos WHERE id = ?", guardado.ID,
	).Scan(&fechaInicio, &fechaFin); err != nil {
		t.Fatalf("no se pudieron leer las fechas: %v", err)
	}
	if fechaInicio != "2026-10-01" || fechaFin != "2026-12-01" {
		t.Errorf("se esperaban 2026-10-01 y 2026-12-01, se obtuvo %q y %q", fechaInicio, fechaFin)
	}
}

func TestSQLiteProjectRepository_IDsAutoincrementales(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(nuevaBDMigrada(t))

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
	repo := NewSQLiteProjectRepository(nuevaBDMigrada(t))

	p, _ := domain.NewProject("Proyecto", "", nil, nil)
	// creadorID inexistente: la FK de integrantes falla en la segunda escritura.
	if _, err := repo.GuardarConScrumMaster(ctx, p, 999999); err == nil {
		t.Fatal("se esperaba error por FK del integrante inexistente")
	}

	var total int
	if err := repo.db.QueryRow("SELECT COUNT(*) FROM proyectos").Scan(&total); err != nil {
		t.Fatal(err)
	}
	if total != 0 {
		t.Fatalf("se esperaban 0 proyectos (rollback de la transacción), se encontraron %d", total)
	}
}

func TestSQLiteProjectRepository_GuardarConScrumMaster_FKDelCreador(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(nuevaBDMigrada(t))

	p, _ := domain.NewProject("Proyecto", "", nil, nil)
	_, err := repo.GuardarConScrumMaster(ctx, p, 999999)
	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Fatalf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
}

func TestSQLiteProjectRepository_AgregarIntegrante_FKVioladaNoEsNoEncontrado(t *testing.T) {
	ctx := context.Background()
	proyectos := NewSQLiteProjectRepository(nuevaBDMigrada(t))

	p, _ := domain.NewProject("Proyecto", "", nil, nil)
	guardado, err := proyectos.Guardar(ctx, p)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	m := domain.Membership{ProjectID: guardado.ID, UserID: 999, Role: domain.RolScrumMaster}
	_, err = proyectos.AgregarIntegrante(ctx, m)
	if err == nil {
		t.Fatal("se esperaba error por FK")
	}
	if errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("una falla de FK no debe mapearse a 'no encontrado': %v", err)
	}
}

func TestSQLiteProjectRepository_NombresRepetidos(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(nuevaBDMigrada(t))

	for i := 0; i < 2; i++ {
		p, _ := domain.NewProject("Mismo nombre", "", nil, nil)
		if _, err := repo.Guardar(ctx, p); err != nil {
			t.Fatalf("no se esperaba error con nombres repetidos: %v", err)
		}
	}
	var total int
	if err := repo.db.QueryRow("SELECT COUNT(*) FROM proyectos WHERE nombre = ?", "Mismo nombre").Scan(&total); err != nil {
		t.Fatalf("no se pudo contar: %v", err)
	}
	if total != 2 {
		t.Errorf("se esperaban 2 proyectos con el mismo nombre, se obtuvieron %d", total)
	}
}

func TestSQLiteProjectRepository_Actualizar(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteProjectRepository(nuevaBDMigrada(t))

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
	repo := NewSQLiteProjectRepository(nuevaBDMigrada(t))
	p, _ := domain.NewProject("Fantasma", "", nil, nil)
	p.ID = 999
	if err := repo.Actualizar(ctx, p); !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
}

func TestSQLiteProjectRepository_ListarIntegrantes_OrdenDeAlta(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	proyectos := NewSQLiteProjectRepository(db)
	usuarios := NewSQLiteUserRepository(db)

	p, _ := domain.NewProject("Proyecto", "", nil, nil)
	guardado, err := proyectos.Guardar(ctx, p)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	// Tres altas con el mismo creado_en (resolución de segundos) en orden no
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
		if _, err := proyectos.AgregarIntegrante(ctx, m); err != nil {
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
	db := nuevaBDMigrada(t)
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
	if _, err := proyectos.AgregarIntegrante(ctx, m); err != nil {
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
	db := nuevaBDMigrada(t)
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
	if _, err := proyectos.AgregarIntegrante(ctx, m); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if _, err := proyectos.AgregarIntegrante(ctx, m); !errors.Is(err, ErrMiembroDuplicado) {
		t.Fatalf("se esperaba ErrMiembroDuplicado, se obtuvo %v", err)
	}
}

func sembrarHorasPorStoryPoint(t *testing.T, db *sql.DB, proyectoID int64, horas any) {
	t.Helper()
	if _, err := db.Exec("UPDATE proyectos SET horas_por_story_point = ? WHERE id = ?", horas, proyectoID); err != nil {
		t.Fatalf("no se pudo sembrar el parámetro: %v", err)
	}
}

func TestSQLiteProjectRepository_LeeLasHorasPorStoryPoint(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	conValor := guardarProyecto(t, db)
	sinValor := guardarProyecto(t, db)
	sembrarHorasPorStoryPoint(t, db, conValor, 2.5)

	leido, err := repo.ObtenerPorID(ctx, conValor)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if leido.HorasPorStoryPoint == nil || *leido.HorasPorStoryPoint != 2.5 {
		t.Errorf("se esperaban 2.5 horas, se obtuvo %v", leido.HorasPorStoryPoint)
	}
	sinDefinir, err := repo.ObtenerPorID(ctx, sinValor)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if sinDefinir.HorasPorStoryPoint != nil {
		t.Errorf("se esperaba sin definir, se obtuvo %v", *sinDefinir.HorasPorStoryPoint)
	}

	mustExec(t, db, `INSERT INTO usuarios (id, nombre, nombre_normalizado, creado_en) VALUES (1, 'Ana', 'ana', '2026-10-23T00:00:00Z')`)
	mustExec(t, db, fmt.Sprintf(`INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES
		(%d, 1, 'scrum_master', '2026-10-23T00:00:00Z'), (%d, 1, 'scrum_master', '2026-10-23T00:00:00Z')`, conValor, sinValor))
	listados, err := repo.ListarPorIntegrante(ctx, 1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(listados) != 2 || listados[0].Proyecto.HorasPorStoryPoint == nil || *listados[0].Proyecto.HorasPorStoryPoint != 2.5 || listados[1].Proyecto.HorasPorStoryPoint != nil {
		t.Errorf("el listado no trae el parámetro: %+v", listados)
	}
}

func TestSQLiteProjectRepository_ActualizarLasHorasPorStoryPoint(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	id := guardarProyecto(t, db)
	proyecto, err := repo.ObtenerPorID(ctx, id)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	horas := 4.0
	proyecto.HorasPorStoryPoint = &horas
	if err := repo.Actualizar(ctx, proyecto); err != nil {
		t.Fatalf("no se esperaba error al guardar: %v", err)
	}
	guardado, _ := repo.ObtenerPorID(ctx, id)
	if guardado.HorasPorStoryPoint == nil || *guardado.HorasPorStoryPoint != 4 {
		t.Fatalf("se esperaban 4 horas, se obtuvo %v", guardado.HorasPorStoryPoint)
	}

	guardado.Nombre = "Otro nombre"
	if err := repo.Actualizar(ctx, guardado); err != nil {
		t.Fatalf("no se esperaba error al editar otro dato: %v", err)
	}
	editado, _ := repo.ObtenerPorID(ctx, id)
	if editado.HorasPorStoryPoint == nil || *editado.HorasPorStoryPoint != 4 {
		t.Fatalf("editar otro dato debía conservar el valor, se obtuvo %v", editado.HorasPorStoryPoint)
	}

	editado.HorasPorStoryPoint = nil
	if err := repo.Actualizar(ctx, editado); err != nil {
		t.Fatalf("no se esperaba error al borrar: %v", err)
	}
	borrado, _ := repo.ObtenerPorID(ctx, id)
	if borrado.HorasPorStoryPoint != nil {
		t.Fatalf("se esperaba sin definir, se obtuvo %v", *borrado.HorasPorStoryPoint)
	}
}

func sembrarFechaCierre(t *testing.T, db *sql.DB, proyectoID int64, fecha string) {
	t.Helper()
	if _, err := db.Exec("UPDATE proyectos SET fecha_cierre = ? WHERE id = ?", fecha, proyectoID); err != nil {
		t.Fatalf("no se pudo sembrar la fecha de cierre: %v", err)
	}
}

// Spec 044 (#300): el repositorio lee proyectos.fecha_cierre como día (00:00 UTC).
func TestSQLiteProjectRepository_LeeLaFechaDeCierre(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	cerrado := guardarProyecto(t, db)
	abierto := guardarProyecto(t, db)
	sembrarFechaCierre(t, db, cerrado, "2026-11-11")
	dia := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)

	leido, err := repo.ObtenerPorID(ctx, cerrado)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if leido.FechaCierre == nil || !leido.FechaCierre.Equal(dia) {
		t.Errorf("se esperaba la fecha de cierre %v, se obtuvo %v", dia, leido.FechaCierre)
	}
	sinCierre, err := repo.ObtenerPorID(ctx, abierto)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if sinCierre.FechaCierre != nil {
		t.Errorf("un proyecto abierto no debe tener fecha de cierre, se obtuvo %v", *sinCierre.FechaCierre)
	}

	mustExec(t, db, `INSERT INTO usuarios (id, nombre, nombre_normalizado, creado_en) VALUES (1, 'Ana', 'ana', '2026-10-23T00:00:00Z')`)
	mustExec(t, db, fmt.Sprintf(`INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES
		(%d, 1, 'scrum_master', '2026-10-23T00:00:00Z'), (%d, 1, 'scrum_master', '2026-10-23T00:00:00Z')`, cerrado, abierto))
	listados, err := repo.ListarPorIntegrante(ctx, 1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(listados) != 2 || listados[0].Proyecto.FechaCierre == nil || !listados[0].Proyecto.FechaCierre.Equal(dia) || listados[1].Proyecto.FechaCierre != nil {
		t.Errorf("el listado no trae la fecha de cierre: %+v", listados)
	}
}

// Editar el proyecto no toca fecha_cierre: solo la escriben cerrar y reabrir.
func TestSQLiteProjectRepository_ActualizarNoPisaLaFechaDeCierre(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	id := guardarProyecto(t, db)
	sembrarFechaCierre(t, db, id, "2026-11-11")
	proyecto, err := repo.ObtenerPorID(ctx, id)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	proyecto.Nombre = "Otro nombre"
	proyecto.FechaCierre = nil
	if err := repo.Actualizar(ctx, proyecto); err != nil {
		t.Fatalf("no se esperaba error al editar: %v", err)
	}

	editado, _ := repo.ObtenerPorID(ctx, id)
	if editado.Nombre != "Otro nombre" {
		t.Errorf("se esperaba el nombre editado, se obtuvo %q", editado.Nombre)
	}
	if editado.FechaCierre == nil {
		t.Error("editar el proyecto no debe borrar la fecha de cierre")
	}
}

// proyectoConUsuarios crea un proyecto sin integrantes y un usuario por nombre.
func proyectoConUsuarios(t *testing.T, db *sql.DB, nombres ...string) (int64, []int64) {
	t.Helper()
	ctx := context.Background()
	p, _ := domain.NewProject("Proyecto", "", nil, nil)
	guardado, err := NewSQLiteProjectRepository(db).Guardar(ctx, p)
	if err != nil {
		t.Fatalf("no se pudo guardar el proyecto: %v", err)
	}
	usuarios := NewSQLiteUserRepository(db)
	var ids []int64
	for _, nombre := range nombres {
		u, _ := domain.NewUser(nombre)
		guardadoU, err := usuarios.Guardar(ctx, u)
		if err != nil {
			t.Fatalf("no se pudo guardar el usuario %s: %v", nombre, err)
		}
		ids = append(ids, guardadoU.ID)
	}
	return guardado.ID, ids
}

func TestSQLiteProjectRepository_AgregarIntegrante_SegundoScrumMasterEsErrorDeDominio(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	proyectoID, ids := proyectoConUsuarios(t, db, "Ana", "Beto")
	ctx := context.Background()
	if _, err := repo.AgregarIntegrante(ctx, domain.Membership{ProjectID: proyectoID, UserID: ids[0], Role: domain.RolScrumMaster}); err != nil {
		t.Fatalf("no se esperaba error al agregar el primer Scrum Master: %v", err)
	}

	_, err := repo.AgregarIntegrante(ctx, domain.Membership{ProjectID: proyectoID, UserID: ids[1], Role: domain.RolScrumMaster})

	if !errors.Is(err, domain.ErrScrumMasterYaExiste) {
		t.Fatalf("se esperaba domain.ErrScrumMasterYaExiste, se obtuvo %v", err)
	}
	if errors.Is(err, ErrMiembroDuplicado) {
		t.Errorf("la violación del índice no debe confundirse con un integrante duplicado")
	}
}

func TestSQLiteProjectRepository_AgregarIntegrante_MismoUsuarioSigueSiendoDuplicado(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	proyectoID, ids := proyectoConUsuarios(t, db, "Ana")
	ctx := context.Background()
	m := domain.Membership{ProjectID: proyectoID, UserID: ids[0], Role: domain.RolProductBuilder}
	if _, err := repo.AgregarIntegrante(ctx, m); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	_, err := repo.AgregarIntegrante(ctx, m)

	if !errors.Is(err, ErrMiembroDuplicado) {
		t.Fatalf("se esperaba ErrMiembroDuplicado, se obtuvo %v", err)
	}
}

func TestSQLiteProjectRepository_AgregarIntegrante_ProductBuildersYScrumMasterEnOtroProyecto(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	proyectoID, ids := proyectoConUsuarios(t, db, "Ana", "Beto", "Cami")
	otroID, _ := proyectoConUsuarios(t, db, "Dani")
	ctx := context.Background()
	altas := []domain.Membership{
		{ProjectID: proyectoID, UserID: ids[0], Role: domain.RolScrumMaster},
		{ProjectID: proyectoID, UserID: ids[1], Role: domain.RolProductBuilder},
		{ProjectID: proyectoID, UserID: ids[2], Role: domain.RolProductBuilder},
		{ProjectID: otroID, UserID: ids[0], Role: domain.RolScrumMaster},
	}
	for _, m := range altas {
		if _, err := repo.AgregarIntegrante(ctx, m); err != nil {
			t.Fatalf("no se esperaba error para %+v: %v", m, err)
		}
	}
}

func TestSQLiteProjectRepository_AgregarIntegrante_DosAltasSimultaneasDeScrumMaster(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	proyectoID, ids := proyectoConUsuarios(t, db, "Ana", "Beto")
	ctx := context.Background()

	resultados := make([]error, len(ids))
	var espera sync.WaitGroup
	for i, id := range ids {
		espera.Add(1)
		go func() {
			defer espera.Done()
			_, resultados[i] = repo.AgregarIntegrante(ctx, domain.Membership{ProjectID: proyectoID, UserID: id, Role: domain.RolScrumMaster})
		}()
	}
	espera.Wait()

	var aceptadas, conflictos int
	for _, err := range resultados {
		switch {
		case err == nil:
			aceptadas++
		case errors.Is(err, domain.ErrScrumMasterYaExiste):
			conflictos++
		default:
			t.Errorf("error inesperado: %v", err)
		}
	}
	if aceptadas != 1 || conflictos != 1 {
		t.Fatalf("se esperaba 1 alta aceptada y 1 conflicto, se obtuvo %d y %d", aceptadas, conflictos)
	}
}

// Spec 044 (#300): Cerrar es un UPDATE condicional (sin cierre previo y sin Sprint Activo).
func TestSQLiteProjectRepository_Cerrar(t *testing.T) {
	ctx := context.Background()
	hoy := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)

	t.Run("abierto sin Sprint activo guarda la fecha", func(t *testing.T) {
		db := nuevaBDMigrada(t)
		repo := NewSQLiteProjectRepository(db)
		id := guardarProyecto(t, db)
		mustExec(t, db, `INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, 'g', '2026-10-01', '2026-10-15', 'Finalizado'), (?, NULL, NULL, NULL, 'Pendiente')`, id, id)

		if err := repo.Cerrar(ctx, id, hoy); err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}

		var guardada string
		if err := db.QueryRow("SELECT fecha_cierre FROM proyectos WHERE id = ?", id).Scan(&guardada); err != nil {
			t.Fatalf("no se pudo leer la fecha de cierre: %v", err)
		}
		if guardada != "2026-11-11" {
			t.Errorf("se esperaba '2026-11-11', se obtuvo %q", guardada)
		}
	})

	t.Run("ya cerrado se rechaza y conserva la fecha", func(t *testing.T) {
		db := nuevaBDMigrada(t)
		repo := NewSQLiteProjectRepository(db)
		id := guardarProyecto(t, db)
		sembrarFechaCierre(t, db, id, "2026-10-01")

		if err := repo.Cerrar(ctx, id, hoy); !errors.Is(err, domain.ErrProyectoYaCerrado) {
			t.Fatalf("se esperaba ErrProyectoYaCerrado, se obtuvo %v", err)
		}
		leido, _ := repo.ObtenerPorID(ctx, id)
		if leido.FechaCierre == nil || leido.FechaCierre.Format(time.DateOnly) != "2026-10-01" {
			t.Errorf("la fecha de cierre previa no debe cambiar, se obtuvo %v", leido.FechaCierre)
		}
	})

	t.Run("con un Sprint activo se rechaza sin cerrar", func(t *testing.T) {
		db := nuevaBDMigrada(t)
		repo := NewSQLiteProjectRepository(db)
		id := guardarProyecto(t, db)
		mustExec(t, db, `INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, 'g', '2026-11-01', '2026-11-15', 'Activo')`, id)

		if err := repo.Cerrar(ctx, id, hoy); !errors.Is(err, domain.ErrSprintActivoAlCerrarProyecto) {
			t.Fatalf("se esperaba ErrSprintActivoAlCerrarProyecto, se obtuvo %v", err)
		}
		leido, _ := repo.ObtenerPorID(ctx, id)
		if leido.FechaCierre != nil {
			t.Errorf("el proyecto debe seguir abierto, se obtuvo %v", *leido.FechaCierre)
		}
	})

	t.Run("el Sprint activo de otro proyecto no impide cerrar", func(t *testing.T) {
		db := nuevaBDMigrada(t)
		repo := NewSQLiteProjectRepository(db)
		id := guardarProyecto(t, db)
		otro := guardarProyecto(t, db)
		mustExec(t, db, `INSERT INTO sprints (proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (?, 'g', '2026-11-01', '2026-11-15', 'Activo')`, otro)

		if err := repo.Cerrar(ctx, id, hoy); err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
	})

	t.Run("inexistente", func(t *testing.T) {
		repo := NewSQLiteProjectRepository(nuevaBDMigrada(t))

		if err := repo.Cerrar(ctx, 99, hoy); !errors.Is(err, domain.ErrProyectoNoEncontrado) {
			t.Fatalf("se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
		}
	})

	t.Run("dos cierres seguidos: el segundo informa que ya está cerrado", func(t *testing.T) {
		db := nuevaBDMigrada(t)
		repo := NewSQLiteProjectRepository(db)
		id := guardarProyecto(t, db)

		if err := repo.Cerrar(ctx, id, hoy); err != nil {
			t.Fatalf("no se esperaba error en el primer cierre: %v", err)
		}
		if err := repo.Cerrar(ctx, id, hoy.AddDate(0, 0, 1)); !errors.Is(err, domain.ErrProyectoYaCerrado) {
			t.Fatalf("se esperaba ErrProyectoYaCerrado, se obtuvo %v", err)
		}
	})
}
