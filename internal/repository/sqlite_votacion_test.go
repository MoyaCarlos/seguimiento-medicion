package repository

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// escenarioVotacion es una base con un proyecto, una historia y dos usuarios.
type escenarioVotacion struct {
	repo            *SQLiteVotacionRepository
	db              *sql.DB
	historiaID      int64
	carlos, candela int64
}

func nuevoEscenarioVotacion(t *testing.T) escenarioVotacion {
	t.Helper()
	db := nuevaBDMigrada(t)
	proyectoID := guardarProyecto(t, db)
	historia := guardarHistoria(t, NewSQLiteBacklogRepository(db), proyectoID, "Historia", domain.PrioridadMust)
	return escenarioVotacion{
		repo:       NewSQLiteVotacionRepository(db),
		db:         db,
		historiaID: historia.ID,
		carlos:     guardarUsuario(t, db, "Carlos"),
		candela:    guardarUsuario(t, db, "Candela"),
	}
}

func guardarUsuario(t *testing.T, db *sql.DB, nombre string) int64 {
	t.Helper()
	u, err := domain.NewUser(nombre)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := NewSQLiteUserRepository(db).Guardar(context.Background(), u)
	if err != nil {
		t.Fatalf("no se esperaba error al guardar el usuario: %v", err)
	}
	return guardado.ID
}

func crearRonda(t *testing.T, e escenarioVotacion) domain.RondaVotacion {
	t.Helper()
	r, err := domain.NuevaRonda(e.historiaID)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	creada, err := e.repo.CrearRonda(context.Background(), r)
	if err != nil {
		t.Fatalf("no se esperaba error al crear la ronda: %v", err)
	}
	return creada
}

func guardarVoto(t *testing.T, e escenarioVotacion, rondaID int64, v domain.Voto) {
	t.Helper()
	if err := e.repo.GuardarVoto(context.Background(), rondaID, v); err != nil {
		t.Fatalf("no se esperaba error al guardar el voto: %v", err)
	}
}

type filaVoto struct {
	usuarioID int64
	valor     int
}

func votosDe(t *testing.T, db *sql.DB, rondaID int64) []filaVoto {
	t.Helper()
	filas, err := db.Query("SELECT usuario_id, valor FROM votos WHERE ronda_id = ? ORDER BY usuario_id", rondaID)
	if err != nil {
		t.Fatalf("no se pudieron leer los votos: %v", err)
	}
	defer filas.Close()
	var votos []filaVoto
	for filas.Next() {
		var v filaVoto
		if err := filas.Scan(&v.usuarioID, &v.valor); err != nil {
			t.Fatalf("no se pudo leer el voto: %v", err)
		}
		votos = append(votos, v)
	}
	return votos
}

func reveladaEnDe(t *testing.T, db *sql.DB, rondaID int64) sql.NullString {
	t.Helper()
	var revelada sql.NullString
	if err := db.QueryRow("SELECT revelada_en FROM rondas_votacion WHERE id = ?", rondaID).Scan(&revelada); err != nil {
		t.Fatalf("no se pudo leer la ronda %d: %v", rondaID, err)
	}
	return revelada
}

func TestSQLiteVotacionRepository_CrearRonda(t *testing.T) {
	e := nuevoEscenarioVotacion(t)

	r := crearRonda(t, e)

	if r.ID <= 0 {
		t.Fatalf("se esperaba un ID asignado, se obtuvo %d", r.ID)
	}
	var (
		historiaID int64
		numero     int
		revelada   sql.NullString
	)
	err := e.db.QueryRow("SELECT historia_id, numero, revelada_en FROM rondas_votacion WHERE id = ?", r.ID).
		Scan(&historiaID, &numero, &revelada)
	if err != nil {
		t.Fatalf("no se pudo releer la ronda: %v", err)
	}
	if historiaID != e.historiaID || numero != 1 || revelada.Valid {
		t.Errorf("se esperaba la ronda 1 sin revelar de la historia %d, se obtuvo historia %d, número %d, revelada %v",
			e.historiaID, historiaID, numero, revelada)
	}
}

func TestSQLiteVotacionRepository_CrearRonda_YaIniciada(t *testing.T) {
	e := nuevoEscenarioVotacion(t)
	crearRonda(t, e)
	r, _ := domain.NuevaRonda(e.historiaID)

	_, err := e.repo.CrearRonda(context.Background(), r)

	if !errors.Is(err, domain.ErrVotacionYaIniciada) {
		t.Errorf("se esperaba ErrVotacionYaIniciada, se obtuvo %v", err)
	}
}

func TestSQLiteVotacionRepository_CrearRonda_HistoriaInexistente(t *testing.T) {
	e := nuevoEscenarioVotacion(t)
	r, _ := domain.NuevaRonda(999)

	_, err := e.repo.CrearRonda(context.Background(), r)

	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Errorf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
}

func TestSQLiteVotacionRepository_ObtenerRondaActual_SinRondas(t *testing.T) {
	e := nuevoEscenarioVotacion(t)

	_, err := e.repo.ObtenerRondaActual(context.Background(), e.historiaID)

	if !errors.Is(err, domain.ErrVotacionNoIniciada) {
		t.Errorf("se esperaba ErrVotacionNoIniciada, se obtuvo %v", err)
	}
}

func TestSQLiteVotacionRepository_ObtenerRondaActual_EsLaDeMayorNumero(t *testing.T) {
	e := nuevoEscenarioVotacion(t)
	mustExec(t, e.db, "INSERT INTO rondas_votacion (id, historia_id, numero, revelada_en) VALUES (10, ?, 1, '2026-10-27T15:04:05Z')", e.historiaID)
	mustExec(t, e.db, "INSERT INTO rondas_votacion (id, historia_id, numero, revelada_en) VALUES (11, ?, 2, NULL)", e.historiaID)
	mustExec(t, e.db, "INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (10, ?, 3)", e.carlos)
	mustExec(t, e.db, "INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (11, ?, 8)", e.candela)
	mustExec(t, e.db, "INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (11, ?, 5)", e.carlos)

	r, err := e.repo.ObtenerRondaActual(context.Background(), e.historiaID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if r.ID != 11 || r.HistoriaID != e.historiaID || r.Numero != 2 || r.ReveladaEn != nil {
		t.Errorf("se esperaba la ronda 2 (id 11) sin revelar, se obtuvo %+v", r)
	}
	esperados := []domain.Voto{{UsuarioID: e.carlos, Valor: 5}, {UsuarioID: e.candela, Valor: 8}}
	if !slices.Equal(r.Votos, esperados) {
		t.Errorf("se esperaban los votos %+v, se obtuvo %+v", esperados, r.Votos)
	}
}

func TestSQLiteVotacionRepository_GuardarVoto_AgregaYReemplaza(t *testing.T) {
	e := nuevoEscenarioVotacion(t)
	r := crearRonda(t, e)

	guardarVoto(t, e, r.ID, domain.Voto{UsuarioID: e.candela, Valor: 8})
	if got := votosDe(t, e.db, r.ID); !slices.Equal(got, []filaVoto{{e.candela, 8}}) {
		t.Fatalf("se esperaba el voto 8 de Candela, se obtuvo %+v", got)
	}

	guardarVoto(t, e, r.ID, domain.Voto{UsuarioID: e.candela, Valor: 13})
	if got := votosDe(t, e.db, r.ID); !slices.Equal(got, []filaVoto{{e.candela, 13}}) {
		t.Errorf("se esperaba un solo voto de Candela, con 13, se obtuvo %+v", got)
	}
}

func TestSQLiteVotacionRepository_MarcarRevelada(t *testing.T) {
	e := nuevoEscenarioVotacion(t)
	r := crearRonda(t, e)
	enBuenosAires := time.Date(2026, 10, 27, 12, 4, 5, 0, time.FixedZone("ART", -3*60*60))

	if err := e.repo.MarcarRevelada(context.Background(), r.ID, enBuenosAires); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if got := reveladaEnDe(t, e.db, r.ID); got.String != "2026-10-27T15:04:05Z" {
		t.Errorf("se esperaba revelada_en en RFC 3339 UTC, se obtuvo %v", got)
	}
	actual, err := e.repo.ObtenerRondaActual(context.Background(), e.historiaID)
	if err != nil {
		t.Fatalf("no se esperaba error al releer: %v", err)
	}
	if actual.ReveladaEn == nil || !actual.ReveladaEn.Equal(enBuenosAires) {
		t.Errorf("se esperaba revelada en %v, se obtuvo %v", enBuenosAires, actual.ReveladaEn)
	}
}

func TestSQLiteVotacionRepository_RondaRevelada_NoCambia(t *testing.T) {
	e := nuevoEscenarioVotacion(t)
	r := crearRonda(t, e)
	guardarVoto(t, e, r.ID, domain.Voto{UsuarioID: e.carlos, Valor: 5})
	mustExec(t, e.db, "UPDATE rondas_votacion SET revelada_en = '2026-10-27T15:04:05Z' WHERE id = ?", r.ID)
	ctx := context.Background()

	if err := e.repo.GuardarVoto(ctx, r.ID, domain.Voto{UsuarioID: e.carlos, Valor: 13}); !errors.Is(err, domain.ErrRondaRevelada) {
		t.Errorf("cambiar un voto: se esperaba ErrRondaRevelada, se obtuvo %v", err)
	}
	if err := e.repo.GuardarVoto(ctx, r.ID, domain.Voto{UsuarioID: e.candela, Valor: 8}); !errors.Is(err, domain.ErrRondaRevelada) {
		t.Errorf("votar por primera vez: se esperaba ErrRondaRevelada, se obtuvo %v", err)
	}
	if err := e.repo.MarcarRevelada(ctx, r.ID, time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC)); !errors.Is(err, domain.ErrRondaRevelada) {
		t.Errorf("revelar de nuevo: se esperaba ErrRondaRevelada, se obtuvo %v", err)
	}

	if got := votosDe(t, e.db, r.ID); !slices.Equal(got, []filaVoto{{e.carlos, 5}}) {
		t.Errorf("los votos no debían cambiar, quedaron %+v", got)
	}
	if got := reveladaEnDe(t, e.db, r.ID); got.String != "2026-10-27T15:04:05Z" {
		t.Errorf("revelada_en no debía cambiar, quedó %v", got)
	}
}

func TestSQLiteVotacionRepository_GuardarVoto_UsuarioInexistente(t *testing.T) {
	e := nuevoEscenarioVotacion(t)
	r := crearRonda(t, e)

	err := e.repo.GuardarVoto(context.Background(), r.ID, domain.Voto{UsuarioID: 999, Valor: 5})

	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Errorf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
}

func TestSQLiteVotacionRepository_SegundaRonda(t *testing.T) {
	ctx := context.Background()
	e := nuevoEscenarioVotacion(t)
	primera := crearRonda(t, e)
	if err := e.repo.GuardarVoto(ctx, primera.ID, domain.Voto{UsuarioID: e.carlos, Valor: 3}); err != nil {
		t.Fatalf("no se esperaba error al votar: %v", err)
	}
	if err := e.repo.MarcarRevelada(ctx, primera.ID, time.Now()); err != nil {
		t.Fatalf("no se esperaba error al revelar: %v", err)
	}

	segunda, err := e.repo.CrearRonda(ctx, domain.RondaVotacion{HistoriaID: e.historiaID, Numero: 2})
	if err != nil {
		t.Fatalf("no se esperaba error al crear la ronda 2: %v", err)
	}

	actual, err := e.repo.ObtenerRondaActual(ctx, e.historiaID)
	if err != nil {
		t.Fatalf("no se esperaba error al leer la ronda actual: %v", err)
	}
	if actual.ID != segunda.ID || actual.Numero != 2 || len(actual.Votos) != 0 || actual.Revelada() {
		t.Errorf("la ronda actual debía ser la 2, sin votos ni revelar, es %+v", actual)
	}
	if _, err := e.repo.CrearRonda(ctx, domain.RondaVotacion{HistoriaID: e.historiaID, Numero: 2}); !errors.Is(err, domain.ErrVotacionYaIniciada) {
		t.Errorf("repetir la ronda 2 debía dar ErrVotacionYaIniciada, se obtuvo %v", err)
	}
}

func TestSQLiteVotacionRepository_ListarRondasReveladas(t *testing.T) {
	ctx := context.Background()
	e := nuevoEscenarioVotacion(t)
	en := time.Date(2026, 11, 3, 15, 4, 5, 0, time.UTC)
	for i, votos := range [][]domain.Voto{
		{{UsuarioID: e.candela, Valor: 8}, {UsuarioID: e.carlos, Valor: 3}},
		nil,
	} {
		ronda, err := e.repo.CrearRonda(ctx, domain.RondaVotacion{HistoriaID: e.historiaID, Numero: i + 1})
		if err != nil {
			t.Fatalf("no se esperaba error al crear la ronda %d: %v", i+1, err)
		}
		for _, v := range votos {
			if err := e.repo.GuardarVoto(ctx, ronda.ID, v); err != nil {
				t.Fatalf("no se esperaba error al votar: %v", err)
			}
		}
		if err := e.repo.MarcarRevelada(ctx, ronda.ID, en); err != nil {
			t.Fatalf("no se esperaba error al revelar: %v", err)
		}
	}
	abierta, err := e.repo.CrearRonda(ctx, domain.RondaVotacion{HistoriaID: e.historiaID, Numero: 3})
	if err != nil {
		t.Fatalf("no se esperaba error al crear la ronda 3: %v", err)
	}
	if err := e.repo.GuardarVoto(ctx, abierta.ID, domain.Voto{UsuarioID: e.carlos, Valor: 13}); err != nil {
		t.Fatalf("no se esperaba error al votar: %v", err)
	}

	rondas, err := e.repo.ListarRondasReveladas(ctx, e.historiaID)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(rondas) != 2 || rondas[0].Numero != 1 || rondas[1].Numero != 2 {
		t.Fatalf("se esperaban las rondas 1 y 2, se obtuvo %+v", rondas)
	}
	if rondas[0].ReveladaEn == nil || !rondas[0].ReveladaEn.Equal(en) {
		t.Errorf("se esperaba el instante de revelación, se obtuvo %v", rondas[0].ReveladaEn)
	}
	esperados := []domain.Voto{{UsuarioID: e.carlos, Valor: 3}, {UsuarioID: e.candela, Valor: 8}}
	if !slices.Equal(rondas[0].Votos, esperados) {
		t.Errorf("se esperaban los votos %+v, se obtuvo %+v", esperados, rondas[0].Votos)
	}
	if len(rondas[1].Votos) != 0 {
		t.Errorf("la ronda 2 no tiene votos, se obtuvo %+v", rondas[1].Votos)
	}
}

func TestSQLiteVotacionRepository_ListarRondasReveladas_SinRondas(t *testing.T) {
	e := nuevoEscenarioVotacion(t)

	rondas, err := e.repo.ListarRondasReveladas(context.Background(), e.historiaID)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if rondas == nil || len(rondas) != 0 {
		t.Errorf("se esperaba una lista vacía, se obtuvo %#v", rondas)
	}
}
