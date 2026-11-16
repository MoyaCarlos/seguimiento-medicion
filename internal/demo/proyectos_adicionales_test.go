package demo

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

const (
	pTurnos         = "Plataforma de turnos médicos (datos de prueba)"
	pBiblioteca     = "Sistema de biblioteca (datos de prueba)"
	pTienda         = "Tienda de barrio (datos de prueba)"
	pReservas       = "App de reservas (datos de prueba)"
	pPortal         = "Portal de alumnos (datos de prueba)"
	pEnBlanco       = "Proyecto en blanco (datos de prueba)"
	pDonaciones     = "Campaña de donaciones (datos de prueba)"
	proyectosNuevos = 7
)

func cargarTodo(t *testing.T) (*sql.DB, Resumen) {
	t.Helper()
	db := baseMigrada(t)
	resumen, err := cargador(db).Cargar(context.Background())
	if err != nil {
		t.Fatalf("cargar: %v", err)
	}
	return db, resumen
}

func idDelProyecto(t *testing.T, db *sql.DB, nombre string) int64 {
	t.Helper()
	var id int64
	if err := db.QueryRow(`SELECT id FROM proyectos WHERE nombre = ?`, nombre).Scan(&id); err != nil {
		t.Fatalf("proyecto %q: %v", nombre, err)
	}
	return id
}

func idDelUsuario(t *testing.T, db *sql.DB, nombre string) int64 {
	t.Helper()
	var id int64
	if err := db.QueryRow(`SELECT id FROM usuarios WHERE nombre_normalizado = ?`, nombre).Scan(&id); err != nil {
		t.Fatalf("usuario %q: %v", nombre, err)
	}
	return id
}

func TestProyectosAdicionales_ConviveConElProyecto25(t *testing.T) {
	db, _ := cargarTodo(t)

	if n := contar(t, db, `SELECT COUNT(*) FROM proyectos`); n != 1+proyectosNuevos {
		t.Fatalf("%d proyectos, se esperaban %d", n, 1+proyectosNuevos)
	}
	if n := contar(t, db, `SELECT COUNT(*) FROM historias WHERE proyecto_id = 25`); n != 12 {
		t.Errorf("el proyecto 25 quedó con %d historias, se esperaban 12", n)
	}
	if n := contar(t, db, `SELECT COUNT(*) FROM integrantes i JOIN usuarios u ON u.id = i.usuario_id
		WHERE i.proyecto_id = 25 AND i.rol = 'scrum_master' AND u.nombre_normalizado = 'carlos.moya'`); n != 1 {
		t.Errorf("Carlos debía seguir como Scrum Master del 25")
	}
	ana := idDelUsuario(t, db, "ana.iriarte")
	for _, nombre := range []string{pTurnos, pBiblioteca, pTienda, pReservas, pPortal, pEnBlanco, pDonaciones} {
		id := idDelProyecto(t, db, nombre)
		if id == 25 {
			t.Errorf("%s usa el id fijo 25", nombre)
		}
		if n := contar(t, db, `SELECT COUNT(*) FROM integrantes WHERE proyecto_id = ? AND rol = 'scrum_master' AND usuario_id = ?`, id, ana); n != 1 {
			t.Errorf("%s: Ana debía ser el único Scrum Master", nombre)
		}
		if n := contar(t, db, `SELECT COUNT(*) FROM integrantes WHERE proyecto_id = ? AND rol = 'scrum_master'`, id); n != 1 {
			t.Errorf("%s: %d Scrum Masters, se esperaba 1", nombre, n)
		}
	}
}

func estadoDe(t *testing.T, db *sql.DB, id int64, ahora time.Time) domain.EstadoProyecto {
	t.Helper()
	ctx := context.Background()
	proyecto, err := repository.NewSQLiteProjectRepository(db).ObtenerPorID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	sprints, err := repository.NewSQLiteSprintRepository(db).ListarPorProyecto(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	return domain.CalcularEstadoProyecto(proyecto, sprints, ahora)
}

func TestProyectosAdicionales_CadaProyectoTieneElEstadoEsperado(t *testing.T) {
	db, _ := cargarTodo(t)
	casos := map[string]domain.EstadoProyecto{
		pTurnos:     domain.ProyectoEnCurso,
		pBiblioteca: domain.ProyectoFinalizado,
		pTienda:     domain.ProyectoFinalizado,
		pReservas:   domain.ProyectoPlanificado,
		pPortal:     domain.ProyectoEnCurso,
		pEnBlanco:   domain.ProyectoPlanificado,
		pDonaciones: domain.ProyectoEnCurso,
	}
	for nombre, esperado := range casos {
		if visto := estadoDe(t, db, idDelProyecto(t, db, nombre), ahoraFijo); visto != esperado {
			t.Errorf("%s: estado %s, se esperaba %s", nombre, visto, esperado)
		}
	}
}

func TestProyectosAdicionales_UnaSegundaCargaNoDuplicaNada(t *testing.T) {
	db, _ := cargarTodo(t)
	antes := fotoDeLaBase(t, db)

	_, err := cargador(db).Cargar(context.Background())

	if !errors.Is(err, ErrYaCargado) {
		t.Fatalf("error %v, se esperaba ErrYaCargado", err)
	}
	for tabla, filas := range fotoDeLaBase(t, db) {
		if filas != antes[tabla] {
			t.Errorf("%s: %d filas tras la segunda carga, había %d", tabla, filas, antes[tabla])
		}
	}
}

func TestProyectosAdicionales_SiFallaElBloqueNuevoNoQuedaNada(t *testing.T) {
	db := baseMigrada(t)
	errInyectado := errors.New("fallo del bloque nuevo")
	fallaDespues := func(ctx context.Context, tx *sql.Tx, c *carga) error {
		if err := bloqueProyectosAdicionales(ctx, tx, c); err != nil {
			return err
		}
		return errInyectado
	}
	bloques := append(append([]bloque{}, bloquesDemo[:len(bloquesDemo)-1]...), fallaDespues)

	_, err := (&Cargador{db: db, ahora: func() time.Time { return ahoraFijo }, hasheador: hasheadorBarato, bloques: bloques}).Cargar(context.Background())

	if !errors.Is(err, errInyectado) {
		t.Fatalf("error %v, se esperaba el fallo inyectado", err)
	}
	for tabla, filas := range fotoDeLaBase(t, db) {
		if filas != 0 {
			t.Errorf("%s quedó con %d filas", tabla, filas)
		}
	}
}

func TestProyectosAdicionales_ElBloqueEstaAlFinalDeLaCarga(t *testing.T) {
	ultimo := bloquesDemo[len(bloquesDemo)-1]
	db := baseMigrada(t)
	// ponytail: se compara por efecto, no por identidad de función.
	c := ejecutar(t, db, bloquesDemo[:len(bloquesDemo)-1]...)
	if c.resumen.ProyectosAdicionales != 0 {
		t.Fatal("antes del último bloque no debía haber proyectos adicionales")
	}
	if err := func() error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer func() { _ = tx.Rollback() }()
		return ultimo(context.Background(), tx, c)
	}(); err != nil {
		t.Fatal(err)
	}
	if c.resumen.ProyectosAdicionales != proyectosNuevos {
		t.Errorf("el último bloque creó %d proyectos, se esperaban %d", c.resumen.ProyectosAdicionales, proyectosNuevos)
	}
}

func TestProyectosAdicionales_IntegranteDadoDeBajaEnElPortal(t *testing.T) {
	db, _ := cargarTodo(t)
	ctx := context.Background()
	id := idDelProyecto(t, db, pPortal)
	repo := repository.NewSQLiteProjectRepository(db)
	jimena := idDelUsuario(t, db, "jimena.martinez")

	activos, err := repo.ListarIntegrantes(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	historicos, err := repo.ListarIntegrantesHistoricos(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if len(activos) != 3 || len(historicos) != 4 {
		t.Fatalf("activos %d e históricos %d, se esperaban 3 y 4", len(activos), len(historicos))
	}
	for _, m := range activos {
		if m.UserID == jimena {
			t.Error("Jimena no debía figurar entre los activos")
		}
	}
	if n := contar(t, db, `SELECT COUNT(*) FROM registros_horas WHERE proyecto_id = ? AND usuario_id = ?`, id, jimena); n == 0 {
		t.Error("las horas de la persona dada de baja deben seguir existiendo")
	}
	if n := contar(t, db, `SELECT COUNT(*) FROM votos v JOIN rondas_votacion r ON r.id = v.ronda_id JOIN historias h ON h.id = r.historia_id
		WHERE h.proyecto_id = ? AND v.usuario_id = ? AND r.revelada_en IS NULL`, id, jimena); n != 0 {
		t.Errorf("%d votos sin revelar de la persona dada de baja", n)
	}
	if n := contar(t, db, `SELECT COUNT(*) FROM votos v JOIN rondas_votacion r ON r.id = v.ronda_id JOIN historias h ON h.id = r.historia_id
		WHERE h.proyecto_id = ? AND v.usuario_id = ? AND r.revelada_en IS NOT NULL`, id, jimena); n == 0 {
		t.Error("debía quedar un voto ya revelado de la persona dada de baja")
	}
}

func TestProyectosAdicionales_ElProyectoCerradoRechazaEscrituras(t *testing.T) {
	db, _ := cargarTodo(t)
	id := idDelProyecto(t, db, pBiblioteca)
	ana := idDelUsuario(t, db, "ana.iriarte")
	repos := repository.NewSQLiteProjectRepository(db)
	crear := service.NewCrearSprint(repos, repository.NewSQLiteSprintRepository(db))

	_, err := crear.Ejecutar(context.Background(), id, ana)

	if !errors.Is(err, domain.ErrProyectoCerrado) {
		t.Fatalf("error %v, se esperaba ErrProyectoCerrado", err)
	}
	if err.Error() != "El proyecto está cerrado" {
		t.Errorf("mensaje %q", err.Error())
	}
}

func TestProyectosAdicionales_CampanaEnCursoYAlCerrarloFinalizado(t *testing.T) {
	db, _ := cargarTodo(t)
	id := idDelProyecto(t, db, pDonaciones)
	if visto := estadoDe(t, db, id, ahoraFijo); visto != domain.ProyectoEnCurso {
		t.Fatalf("estado %s, se esperaba En curso (regla de #315)", visto)
	}
	cerrar := service.NewCerrarProyecto(repository.NewSQLiteProjectRepository(db), repository.NewSQLiteSprintRepository(db),
		func() time.Time { return ahoraFijo })

	if _, err := cerrar.Ejecutar(context.Background(), id, idDelUsuario(t, db, "ana.iriarte")); err != nil {
		t.Fatalf("cerrar: %v", err)
	}

	if visto := estadoDe(t, db, id, ahoraFijo); visto != domain.ProyectoFinalizado {
		t.Errorf("estado %s tras cerrar, se esperaba Finalizado", visto)
	}
}

func TestProyectosAdicionales_JimenaNoIntegraLaTienda(t *testing.T) {
	db, _ := cargarTodo(t)
	id := idDelProyecto(t, db, pTienda)

	if n := contar(t, db, `SELECT COUNT(*) FROM integrantes WHERE proyecto_id = ?`, id); n != 3 {
		t.Errorf("%d integrantes, se esperaban 3", n)
	}
	if n := contar(t, db, `SELECT COUNT(*) FROM integrantes WHERE proyecto_id = ? AND usuario_id = ?`, id, idDelUsuario(t, db, "jimena.martinez")); n != 0 {
		t.Error("Jimena no debía integrar la tienda")
	}
}

func TestProyectosAdicionales_VotacionEnVivoEnLaPlataforma(t *testing.T) {
	db, _ := cargarTodo(t)
	id := idDelProyecto(t, db, pTurnos)
	carlos := idDelUsuario(t, db, "carlos.moya")

	abiertaSinCarlos := contar(t, db, `SELECT COUNT(*) FROM rondas_votacion r JOIN historias h ON h.id = r.historia_id
		WHERE h.proyecto_id = ? AND r.revelada_en IS NULL
		AND NOT EXISTS (SELECT 1 FROM votos v WHERE v.ronda_id = r.id AND v.usuario_id = ?)
		AND EXISTS (SELECT 1 FROM votos v WHERE v.ronda_id = r.id)`, id, carlos)
	if abiertaSinCarlos < 1 {
		t.Error("falta una ronda abierta con votos donde Carlos no votó")
	}
	segundaTrasDivergencia := contar(t, db, `SELECT COUNT(*) FROM rondas_votacion r2 JOIN rondas_votacion r1
		ON r1.historia_id = r2.historia_id AND r1.numero = 1 AND r2.numero = 2
		JOIN historias h ON h.id = r1.historia_id
		WHERE h.proyecto_id = ? AND r1.revelada_en IS NOT NULL
		AND (SELECT COUNT(DISTINCT valor) FROM votos WHERE ronda_id = r1.id) > 1`, id)
	if segundaTrasDivergencia < 1 {
		t.Error("falta una ronda con divergencia seguida de una segunda ronda")
	}
	consenso := contar(t, db, `SELECT COUNT(*) FROM rondas_votacion r JOIN historias h ON h.id = r.historia_id
		WHERE h.proyecto_id = ? AND r.revelada_en IS NOT NULL AND h.estimacion_sp IS NOT NULL
		AND (SELECT COUNT(DISTINCT valor) FROM votos WHERE ronda_id = r.id) = 1`, id)
	if consenso < 1 {
		t.Error("falta una ronda revelada con consenso y estimación registrada")
	}
	sinEstimar := contar(t, db, `SELECT COUNT(*) FROM historias h JOIN sprints s ON s.id = h.sprint_id
		WHERE h.proyecto_id = ? AND h.estimacion_sp IS NULL AND s.estado IN ('Activo', 'Pendiente')
		AND NOT EXISTS (SELECT 1 FROM rondas_votacion r WHERE r.historia_id = h.id)`, id)
	if sinEstimar < 1 {
		t.Error("faltan historias sin estimar ni rondas en un Sprint Activo o Pendiente")
	}
}

func TestProyectosAdicionales_LaPlataformaCubreLasFeaturesDeLaMatriz(t *testing.T) {
	db, _ := cargarTodo(t)
	id := idDelProyecto(t, db, pTurnos)
	consultas := map[string]string{
		"Sprints Finalizado x2":      `SELECT COUNT(*) = 2 FROM sprints WHERE proyecto_id = ? AND estado = 'Finalizado'`,
		"un Sprint Activo":           `SELECT COUNT(*) = 1 FROM sprints WHERE proyecto_id = ? AND estado = 'Activo'`,
		"un Sprint Pendiente":        `SELECT COUNT(*) = 1 FROM sprints WHERE proyecto_id = ? AND estado = 'Pendiente'`,
		"unas 25 historias":          `SELECT COUNT(*) BETWEEN 24 AND 26 FROM historias WHERE proyecto_id = ?`,
		"una historia arrastrada":    `SELECT COUNT(*) >= 1 FROM historias_arrastradas WHERE proyecto_id = ?`,
		"una historia sin criterios": `SELECT COUNT(*) >= 1 FROM historias WHERE proyecto_id = ? AND criterios_aceptacion IS NULL AND sprint_id IS NOT NULL`,
		"un título largo":            `SELECT COUNT(*) >= 1 FROM historias WHERE proyecto_id = ? AND length(titulo) > 80`,
		"Markdown variado":           `SELECT COUNT(*) >= 4 FROM historias WHERE proyecto_id = ? AND (descripcion LIKE '%**%' OR descripcion LIKE '%# %' OR descripcion LIKE '%](%' OR descripcion LIKE '%` + "`" + `%' OR descripcion LIKE '%> %')`,
		"sin imágenes externas":      `SELECT COUNT(*) = 0 FROM historias WHERE proyecto_id = ? AND (descripcion LIKE '%![%' OR criterios_aceptacion LIKE '%![%')`,
		"horas de los 4":             `SELECT COUNT(DISTINCT usuario_id) = 4 FROM registros_horas WHERE proyecto_id = ?`,
		"defectos 5 o 6":             `SELECT COUNT(*) BETWEEN 5 AND 6 FROM defectos WHERE proyecto_id = ?`,
		"defectos abiertos":          `SELECT COUNT(*) >= 1 FROM defectos WHERE proyecto_id = ? AND estado = 'Pendiente'`,
		"defectos resueltos":         `SELECT COUNT(*) >= 1 FROM defectos WHERE proyecto_id = ? AND estado = 'Resuelto'`,
		"horas por story point":      `SELECT COUNT(*) = 1 FROM proyectos WHERE id = ? AND horas_por_story_point IS NOT NULL`,
		"horas h:mm y decimales":     `SELECT COUNT(DISTINCT minutos) >= 4 FROM registros_horas WHERE proyecto_id = ? AND minutos % 60 <> 0`,
	}
	for nombre, consulta := range consultas {
		if contar(t, db, consulta, id) != 1 {
			t.Errorf("no se cumple: %s", nombre)
		}
	}
}

func TestProyectosAdicionales_ElResumenSumaLosConteosNuevos(t *testing.T) {
	db, resumen := cargarTodo(t)

	if resumen.ProyectoID != 25 || resumen.Historias != 12 || resumen.Defectos != 5 {
		t.Errorf("el resumen del proyecto 25 cambió: %+v", resumen)
	}
	esperados := map[string]int{
		"proyectos_adicionales":       resumen.ProyectosAdicionales,
		"historias_adicionales":       resumen.HistoriasAdicionales,
		"sprints_adicionales":         resumen.SprintsAdicionales,
		"registros_horas_adicionales": resumen.RegistrosHorasAdicionales,
		"rondas_adicionales":          resumen.RondasAdicionales,
		"votos_adicionales":           resumen.VotosAdicionales,
		"defectos_adicionales":        resumen.DefectosAdicionales,
	}
	reales := map[string]int{
		"proyectos_adicionales":       contar(t, db, `SELECT COUNT(*) FROM proyectos WHERE id <> 25`),
		"historias_adicionales":       contar(t, db, `SELECT COUNT(*) FROM historias WHERE proyecto_id <> 25`),
		"sprints_adicionales":         contar(t, db, `SELECT COUNT(*) FROM sprints WHERE proyecto_id <> 25`),
		"registros_horas_adicionales": contar(t, db, `SELECT COUNT(*) FROM registros_horas WHERE proyecto_id <> 25`),
		"rondas_adicionales":          contar(t, db, `SELECT COUNT(*) FROM rondas_votacion r JOIN historias h ON h.id = r.historia_id WHERE h.proyecto_id <> 25`),
		"votos_adicionales":           contar(t, db, `SELECT COUNT(*) FROM votos v JOIN rondas_votacion r ON r.id = v.ronda_id JOIN historias h ON h.id = r.historia_id WHERE h.proyecto_id <> 25`),
		"defectos_adicionales":        contar(t, db, `SELECT COUNT(*) FROM defectos WHERE proyecto_id <> 25`),
	}
	for campo, real := range reales {
		if esperados[campo] != real {
			t.Errorf("%s: el resumen dice %d y la base tiene %d", campo, esperados[campo], real)
		}
	}
	if resumen.ProyectosAdicionales != proyectosNuevos || resumen.RondasAdicionales != 4 || resumen.VotosAdicionales < 13 {
		t.Errorf("conteos fuera de lo esperado: %+v", resumen)
	}
}

func TestProyectosAdicionales_DosCargasSimultaneasDejanUnSoloJuego(t *testing.T) {
	ruta := baseEnArchivo(t)
	var dbs []*sql.DB
	for i := 0; i < 2; i++ {
		db, err := repository.AbrirSQLite(ruta)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		dbs = append(dbs, db)
	}

	exitosas, repetidas := cargarEnParalelo(t, dbs, 2)

	if exitosas != 1 || repetidas != 1 {
		t.Errorf("exitosas %d, repetidas %d; se esperaba 1 y 1", exitosas, repetidas)
	}
	if n := contar(t, dbs[0], `SELECT COUNT(*) FROM proyectos`); n != 1+proyectosNuevos {
		t.Errorf("%d proyectos, se esperaban %d", n, 1+proyectosNuevos)
	}
}
