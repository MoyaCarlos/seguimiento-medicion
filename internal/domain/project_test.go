package domain

import (
	"strings"
	"testing"
	"time"
)

func TestNewProject_Valido(t *testing.T) {
	inicio := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)

	conFechas, err := NewProject("  Software Metrics  ", "  Descripción  ", &inicio, &fin)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if conFechas.Nombre != "Software Metrics" || conFechas.Descripcion != "Descripción" {
		t.Errorf("no se recortaron los campos: %+v", conFechas)
	}
	if conFechas.FechaInicio == nil || !conFechas.FechaInicio.Equal(inicio) {
		t.Errorf("fecha de inicio inesperada: %v", conFechas.FechaInicio)
	}
	if conFechas.FechaFin == nil || !conFechas.FechaFin.Equal(fin) {
		t.Errorf("fecha de fin inesperada: %v", conFechas.FechaFin)
	}

	sinFechas, err := NewProject("Proyecto corto", "", nil, nil)
	if err != nil {
		t.Fatalf("no se esperaba error sin fechas: %v", err)
	}
	if sinFechas.FechaInicio != nil || sinFechas.FechaFin != nil {
		t.Errorf("se esperaban fechas nulas: %+v", sinFechas)
	}
}

func TestNewProject_NombreObligatorio(t *testing.T) {
	for _, nombre := range []string{"", "   "} {
		_, err := NewProject(nombre, "", nil, nil)
		if err == nil {
			t.Fatalf("se esperaba error para el nombre %q", nombre)
		}
		verr, ok := err.(ValidationError)
		if !ok || verr.Campo != "nombre" {
			t.Fatalf("se esperaba ValidationError de nombre, se obtuvo %v", err)
		}
	}
}

func TestNewProject_NombreDemasiadoLargo(t *testing.T) {
	_, err := NewProject(strings.Repeat("a", projectNameMaxLen+1), "", nil, nil)
	verr, ok := err.(ValidationError)
	if !ok || verr.Campo != "nombre" {
		t.Fatalf("se esperaba ValidationError de nombre, se obtuvo %v", err)
	}
}

func TestNewProject_DescripcionDemasiadoLarga(t *testing.T) {
	_, err := NewProject("Proyecto", strings.Repeat("d", projectDescriptionMaxLen+1), nil, nil)
	verr, ok := err.(ValidationError)
	if !ok || verr.Campo != "descripcion" {
		t.Fatalf("se esperaba ValidationError de descripcion, se obtuvo %v", err)
	}
}

func TestNewProject_FechaFinAnteriorAlInicio(t *testing.T) {
	inicio := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	_, err := NewProject("Proyecto", "", &inicio, &fin)
	verr, ok := err.(ValidationError)
	if !ok || verr.Campo != "fecha_fin" {
		t.Fatalf("se esperaba ValidationError de fecha_fin, se obtuvo %v", err)
	}
}

func TestNewProject_FechaFinIgualAlInicio(t *testing.T) {
	fecha := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	_, err := NewProject("Proyecto", "", &fecha, &fecha)
	verr, ok := err.(ValidationError)
	if !ok || verr.Campo != "fecha_fin" {
		t.Fatalf("se esperaba ValidationError de fecha_fin, se obtuvo %v", err)
	}
}

func TestNewProject_NombresRepetidosPermitidos(t *testing.T) {
	for i := 0; i < 2; i++ {
		if _, err := NewProject("Mismo nombre", "", nil, nil); err != nil {
			t.Fatalf("se permiten nombres repetidos: %v", err)
		}
	}
}

func TestConDatosEditados_ConservaIdentidad(t *testing.T) {
	creado := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	base := Project{ID: 1, Nombre: "Original", Descripcion: "desc", CreadoEn: creado}
	fin := time.Date(2026, 12, 15, 0, 0, 0, 0, time.UTC)

	editado, err := base.ConDatosEditados("  Nuevo  ", "", nil, &fin)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if editado.Nombre != "Nuevo" || editado.Descripcion != "" {
		t.Errorf("campos no actualizados: %+v", editado)
	}
	if editado.FechaInicio != nil || editado.FechaFin == nil || !editado.FechaFin.Equal(fin) {
		t.Errorf("fechas inesperadas: inicio=%v fin=%v", editado.FechaInicio, editado.FechaFin)
	}
	if editado.ID != 1 || !editado.CreadoEn.Equal(creado) {
		t.Errorf("no se conservó la identidad: %+v", editado)
	}
}

func TestConDatosEditados_ReutilizaValidacionesDelAlta(t *testing.T) {
	inicio := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)
	fin := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	base, _ := NewProject("Original", "", nil, nil)

	casos := []struct {
		nombre      string
		name        string
		description string
		start       *time.Time
		end         *time.Time
		campo       string
	}{
		{"nombre vacío", "   ", "", nil, nil, "nombre"},
		{"nombre largo", strings.Repeat("a", projectNameMaxLen+1), "", nil, nil, "nombre"},
		{"descripción larga", "Proyecto", strings.Repeat("d", projectDescriptionMaxLen+1), nil, nil, "descripcion"},
		{"fin anterior al inicio", "Proyecto", "", &inicio, &fin, "fecha_fin"},
		{"fin igual al inicio", "Proyecto", "", &inicio, &inicio, "fecha_fin"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := base.ConDatosEditados(c.name, c.description, c.start, c.end)
			verr, ok := err.(ValidationError)
			if !ok || verr.Campo != c.campo {
				t.Fatalf("se esperaba ValidationError de %q, se obtuvo %v", c.campo, err)
			}
		})
	}
}

func TestProject_ConHorasPorStoryPoint_Validas(t *testing.T) {
	casos := map[string]float64{"4": 4, "2.5": 2.5, "1e2": 100}
	for texto, esperado := range casos {
		t.Run(texto, func(t *testing.T) {
			proyecto, err := Project{ID: 1}.ConHorasPorStoryPoint(&texto)
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if proyecto.HorasPorStoryPoint == nil || *proyecto.HorasPorStoryPoint != esperado {
				t.Errorf("se esperaba %v, se obtuvo %v", esperado, proyecto.HorasPorStoryPoint)
			}
		})
	}
}

func TestProject_ConHorasPorStoryPoint_SinValorLoDejaSinDefinir(t *testing.T) {
	cuatro := 4.0

	proyecto, err := Project{HorasPorStoryPoint: &cuatro}.ConHorasPorStoryPoint(nil)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if proyecto.HorasPorStoryPoint != nil {
		t.Errorf("se esperaba sin definir, se obtuvo %v", *proyecto.HorasPorStoryPoint)
	}
}

func TestProject_ConHorasPorStoryPoint_Invalidas(t *testing.T) {
	cuatro := 4.0
	base := Project{HorasPorStoryPoint: &cuatro}
	casos := []string{"0", "-1", `"abc"`, `""`, `"5"`, "true", "[]", "1e999", "Inf", "NaN"}

	for _, texto := range casos {
		t.Run(texto, func(t *testing.T) {
			_, err := base.ConHorasPorStoryPoint(&texto)

			if err != ErrHorasPorStoryPointInvalidas {
				t.Fatalf("se esperaba ErrHorasPorStoryPointInvalidas, se obtuvo %v", err)
			}
		})
	}
}

// Spec 044 (#300): reglas del cierre del proyecto.
func TestProject_Cerrar(t *testing.T) {
	ahora := time.Date(2026, 11, 11, 15, 30, 0, 0, time.Local)
	hoy := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)
	cierrePrevio := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	sprint := func(e EstadoSprint) Sprint { return Sprint{Estado: e} }

	t.Run("abierto sin Sprints se cierra con el día civil de hoy", func(t *testing.T) {
		abierto := Project{ID: 1, Nombre: "Tienda"}

		cerrado, err := abierto.Cerrar(ahora, nil)

		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if cerrado.FechaCierre == nil || !cerrado.FechaCierre.Equal(hoy) {
			t.Errorf("se esperaba la fecha de cierre %v, se obtuvo %v", hoy, cerrado.FechaCierre)
		}
		if cerrado.ID != 1 || cerrado.Nombre != "Tienda" {
			t.Errorf("cerrar no debe cambiar los demás datos: %+v", cerrado)
		}
		if abierto.FechaCierre != nil {
			t.Error("cerrar debe devolver una copia y no modificar el proyecto original")
		}
	})

	t.Run("abierto con Sprints pendientes y finalizados se cierra", func(t *testing.T) {
		cerrado, err := Project{}.Cerrar(ahora, []Sprint{sprint(SprintFinalizado), sprint(SprintPendiente)})

		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if cerrado.FechaCierre == nil {
			t.Error("se esperaba el proyecto cerrado")
		}
	})

	t.Run("con un Sprint activo se rechaza", func(t *testing.T) {
		_, err := Project{}.Cerrar(ahora, []Sprint{sprint(SprintFinalizado), sprint(SprintActivo)})

		if err != ErrSprintActivoAlCerrarProyecto {
			t.Fatalf("se esperaba ErrSprintActivoAlCerrarProyecto, se obtuvo %v", err)
		}
		if err.Error() != "Cerrá el Sprint activo antes de cerrar el proyecto" {
			t.Errorf("mensaje inesperado: %q", err.Error())
		}
	})

	t.Run("ya cerrado se rechaza y conserva la fecha original", func(t *testing.T) {
		yaCerrado := Project{FechaCierre: &cierrePrevio}

		_, err := yaCerrado.Cerrar(ahora, nil)

		if err != ErrProyectoYaCerrado {
			t.Fatalf("se esperaba ErrProyectoYaCerrado, se obtuvo %v", err)
		}
		if err.Error() != "El proyecto ya está cerrado" {
			t.Errorf("mensaje inesperado: %q", err.Error())
		}
		if !yaCerrado.FechaCierre.Equal(cierrePrevio) {
			t.Errorf("la fecha de cierre original no debe cambiar, se obtuvo %v", *yaCerrado.FechaCierre)
		}
	})

	t.Run("ya cerrado con un Sprint activo informa primero que ya está cerrado", func(t *testing.T) {
		_, err := Project{FechaCierre: &cierrePrevio}.Cerrar(ahora, []Sprint{sprint(SprintActivo)})

		if err != ErrProyectoYaCerrado {
			t.Fatalf("se esperaba ErrProyectoYaCerrado, se obtuvo %v", err)
		}
	})
}

// Spec 044 (#300): reabrir deja vacía la fecha de cierre.
func TestProject_Reabrir(t *testing.T) {
	cierre := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)

	t.Run("cerrado vuelve a estar abierto", func(t *testing.T) {
		cerrado := Project{ID: 1, Nombre: "Tienda", FechaCierre: &cierre}

		abierto, err := cerrado.Reabrir()

		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if abierto.FechaCierre != nil {
			t.Errorf("se esperaba sin fecha de cierre, se obtuvo %v", *abierto.FechaCierre)
		}
		if abierto.ID != 1 || abierto.Nombre != "Tienda" {
			t.Errorf("reabrir no debe cambiar los demás datos: %+v", abierto)
		}
		if cerrado.FechaCierre == nil {
			t.Error("reabrir debe devolver una copia y no modificar el proyecto original")
		}
	})

	t.Run("abierto se rechaza", func(t *testing.T) {
		_, err := Project{}.Reabrir()

		if err != ErrProyectoNoCerrado {
			t.Fatalf("se esperaba ErrProyectoNoCerrado, se obtuvo %v", err)
		}
		if err.Error() != "El proyecto no está cerrado" {
			t.Errorf("mensaje inesperado: %q", err.Error())
		}
	})

	t.Run("reabierto con la fecha de fin vencida vuelve a las reglas de HU-13", func(t *testing.T) {
		inicio := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
		fin := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
		cerrado := Project{FechaInicio: &inicio, FechaFin: &fin, FechaCierre: &cierre}

		abierto, err := cerrado.Reabrir()

		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if got := CalcularEstadoProyecto(abierto, nil, civilLocal(2026, 11, 11)); got != ProyectoFinalizado {
			t.Errorf("con la fecha de fin vencida se esperaba %q, se obtuvo %q", ProyectoFinalizado, got)
		}
		if got := CalcularEstadoProyecto(abierto, []Sprint{{Estado: SprintActivo}}, civilLocal(2026, 11, 11)); got != ProyectoEnCurso {
			t.Errorf("ya no está cerrado: con un Sprint activo se esperaba %q, se obtuvo %q", ProyectoEnCurso, got)
		}
	})
}
