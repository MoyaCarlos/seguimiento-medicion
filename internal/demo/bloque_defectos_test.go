package demo

import (
	"context"
	"testing"
	"time"
)

func cargarDefectos(t *testing.T) (*carga, func(consulta string, args ...any) int) {
	t.Helper()
	db := baseMigrada(t)
	c := ejecutar(t, db, bloqueUsuarios, bloqueProyecto, bloqueSprints, bloqueHistorias, bloqueVotacion, bloqueDefectos)
	return c, func(consulta string, args ...any) int { return contar(t, db, consulta, args...) }
}

func TestBloqueDefectos_CreaCincoConLasTresSeveridades(t *testing.T) {
	c, n := cargarDefectos(t)

	if c.resumen.Defectos != 5 || n(`SELECT COUNT(*) FROM defectos WHERE proyecto_id = 25`) != 5 {
		t.Fatalf("se esperaban 5 defectos, resumen %d", c.resumen.Defectos)
	}
	for _, severidad := range []string{"Alta", "Media", "Baja"} {
		if n(`SELECT COUNT(*) FROM defectos WHERE severidad = ?`, severidad) == 0 {
			t.Errorf("falta un defecto de severidad %s", severidad)
		}
	}
}

func TestBloqueDefectos_TresResueltosYDosPendientes(t *testing.T) {
	_, n := cargarDefectos(t)

	if got := n(`SELECT COUNT(*) FROM defectos WHERE estado = 'Resuelto'`); got != 3 {
		t.Errorf("Resueltos: %d, se esperaban 3", got)
	}
	if got := n(`SELECT COUNT(*) FROM defectos WHERE estado = 'Pendiente'`); got != 2 {
		t.Errorf("Pendientes: %d, se esperaban 2", got)
	}
}

func TestBloqueDefectos_CadaDefectoEstaLigadoASuHistoriaYSusSprints(t *testing.T) {
	c, n := cargarDefectos(t)

	casos := []struct {
		historia, severidad, estado string
		deteccion, resolucion       int // Sprint 1 a 4; 0 = sin resolución
	}{
		{"apto", "Alta", "Resuelto", 1, 2},
		{"turno", "Media", "Resuelto", 2, 2},
		{"comprobante", "Baja", "Pendiente", 2, 0},
		{"examen", "Alta", "Pendiente", 3, 0},
		{"reprogramar", "Baja", "Resuelto", 3, 3},
	}
	for _, caso := range casos {
		var resolucion any
		consulta := `SELECT COUNT(*) FROM defectos WHERE historia_id = ? AND severidad = ? AND estado = ?
			AND sprint_deteccion_id = ? AND sprint_resolucion_id IS NULL`
		args := []any{c.historias[caso.historia], caso.severidad, caso.estado, c.sprints[caso.deteccion-1]}
		if caso.resolucion > 0 {
			resolucion = c.sprints[caso.resolucion-1]
			consulta = `SELECT COUNT(*) FROM defectos WHERE historia_id = ? AND severidad = ? AND estado = ?
				AND sprint_deteccion_id = ? AND sprint_resolucion_id = ?`
			args = append(args, resolucion)
		}
		if n(consulta, args...) != 1 {
			t.Errorf("falta el defecto %s de la historia %s (detectado en el Sprint %d, resolución %d)",
				caso.severidad, caso.historia, caso.deteccion, caso.resolucion)
		}
	}
}

func TestBloqueDefectos_LaResolucionNuncaEsAnteriorALaDeteccion(t *testing.T) {
	_, n := cargarDefectos(t)
	if n(`SELECT COUNT(*) FROM defectos`) == 0 {
		t.Fatal("no hay defectos cargados")
	}

	if got := n(`SELECT COUNT(*) FROM defectos WHERE sprint_resolucion_id < sprint_deteccion_id`); got != 0 {
		t.Errorf("%d defectos resueltos antes de detectarse", got)
	}
	if got := n(`SELECT COUNT(*) FROM defectos WHERE (estado = 'Resuelto') <> (sprint_resolucion_id IS NOT NULL)`); got != 0 {
		t.Errorf("%d defectos con estado y Sprint de resolución incoherentes", got)
	}
}

func TestBloqueDefectos_LaFechaDeRegistroCaeDentroDeSuSprintDeDeteccionYNoEsFutura(t *testing.T) {
	_, n := cargarDefectos(t)
	if n(`SELECT COUNT(*) FROM defectos`) == 0 {
		t.Fatal("no hay defectos cargados")
	}

	fueraDeRango := `SELECT COUNT(*) FROM defectos d JOIN sprints s ON s.id = d.sprint_deteccion_id
		WHERE substr(d.creado_en, 1, 10) < s.fecha_inicio OR substr(d.creado_en, 1, 10) > s.fecha_fin OR d.creado_en > ?`
	if got := n(fueraDeRango, ahoraFijo.Format(time.RFC3339)); got != 0 {
		t.Errorf("%d defectos con fecha fuera de su Sprint de detección o posterior a hoy", got)
	}
}

func TestCargar_UnaFallaAMitadDelBloqueDeDefectosDeshaceTodaLaCarga(t *testing.T) {
	db := baseMigrada(t)
	if _, err := db.Exec(`CREATE TRIGGER falla_en_defecto_baja BEFORE INSERT ON defectos WHEN NEW.severidad = 'Baja'
		BEGIN SELECT RAISE(ABORT, 'fallo inyectado'); END`); err != nil {
		t.Fatal(err)
	}

	_, err := cargador(db).Cargar(context.Background())

	if err == nil {
		t.Fatal("la carga debía fallar al insertar un defecto de severidad Baja")
	}
	for tabla, filas := range fotoDeLaBase(t, db) {
		if filas != 0 {
			t.Errorf("%s quedó con %d filas tras la falla en los defectos", tabla, filas)
		}
	}
}
