package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

type escenarioDefecto struct {
	sprints   *repository.SprintRepositoryEnMemoria
	historias *fakeHistoriaRepository
	defectos  *fakeDefectoRepository
}

// nuevoEscenarioDefecto siembra la historia 3 del proyecto 1, en el Product
// Backlog, y la historia 4 del proyecto 2. Los Sprints los agrega cada test.
func nuevoEscenarioDefecto() escenarioDefecto {
	return escenarioDefecto{
		sprints: repository.NewSprintRepositoryEnMemoria(),
		historias: &fakeHistoriaRepository{historias: map[int64]domain.BacklogItem{
			3: {ID: 3, ProyectoID: 1, Titulo: "Carrito de compras", Estado: domain.EstadoNueva},
			4: {ID: 4, ProyectoID: 2, Titulo: "Login", Estado: domain.EstadoNueva},
		}},
		defectos: &fakeDefectoRepository{titulos: map[int64]string{3: "Carrito de compras"}},
	}
}

func (e escenarioDefecto) sprint(t *testing.T, proyectoID int64, estado domain.EstadoSprint) domain.Sprint {
	t.Helper()
	s, err := e.sprints.Guardar(context.Background(), domain.Sprint{ProyectoID: proyectoID, SprintGoal: "g", Estado: estado})
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
	return s
}

func (e escenarioDefecto) registrar() *RegistrarDefecto {
	return NewRegistrarDefecto(e.sprints, e.historias, e.defectos)
}

func defectoValido() RegistrarDefectoInput {
	return RegistrarDefectoInput{ProyectoID: 1, HistoriaID: 3, Descripcion: "El total no se actualiza", Severidad: "Alta"}
}

func TestRegistrarDefecto_RegistraEnElSprintActivo(t *testing.T) {
	e := nuevoEscenarioDefecto()
	e.sprint(t, 1, domain.SprintFinalizado)
	activo := e.sprint(t, 1, domain.SprintActivo)

	detalle, err := e.registrar().Ejecutar(context.Background(), defectoValido())

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(e.defectos.defectos) != 1 {
		t.Fatalf("se esperaba 1 defecto registrado, hay %d", len(e.defectos.defectos))
	}
	registrado := e.defectos.defectos[0]
	if registrado.ProyectoID != 1 || registrado.HistoriaID != 3 || registrado.SprintDeteccionID != activo.ID ||
		registrado.Estado != domain.DefectoPendiente || registrado.SprintResolucionID != nil ||
		registrado.Descripcion != "El total no se actualiza" || registrado.Severidad != domain.SeveridadAlta {
		t.Errorf("defecto registrado inesperado: %+v", registrado)
	}
	if detalle.ID == 0 || detalle.ID != registrado.ID || detalle.TituloHistoria != "Carrito de compras" {
		t.Errorf("detalle inesperado: %+v", detalle)
	}
}

func TestRegistrarDefecto_HistoriaDeUnSprintFinalizado(t *testing.T) {
	e := nuevoEscenarioDefecto()
	finalizado := e.sprint(t, 1, domain.SprintFinalizado)
	activo := e.sprint(t, 1, domain.SprintActivo)
	historia := e.historias.historias[3]
	historia.SprintID, historia.Estado = &finalizado.ID, domain.EstadoCompletada
	e.historias.historias[3] = historia

	detalle, err := e.registrar().Ejecutar(context.Background(), defectoValido())

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if detalle.SprintDeteccionID != activo.ID {
		t.Errorf("se esperaba el Sprint de detección %d, se obtuvo %d", activo.ID, detalle.SprintDeteccionID)
	}
}

func TestRegistrarDefecto_SinSprintActivo(t *testing.T) {
	casos := []struct {
		nombre  string
		sembrar func(t *testing.T, e escenarioDefecto)
	}{
		{"sin Sprints", func(*testing.T, escenarioDefecto) {}},
		{"solo Pendiente", func(t *testing.T, e escenarioDefecto) { e.sprint(t, 1, domain.SprintPendiente) }},
		{"solo Finalizado", func(t *testing.T, e escenarioDefecto) { e.sprint(t, 1, domain.SprintFinalizado) }},
		{"Activo de otro proyecto", func(t *testing.T, e escenarioDefecto) { e.sprint(t, 2, domain.SprintActivo) }},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := nuevoEscenarioDefecto()
			c.sembrar(t, e)

			_, err := e.registrar().Ejecutar(context.Background(), defectoValido())

			if !errors.Is(err, domain.ErrSinSprintActivoParaDetectar) {
				t.Fatalf("se esperaba ErrSinSprintActivoParaDetectar, se obtuvo %v", err)
			}
			if e.defectos.registros != 0 {
				t.Errorf("no se esperaba llamar a Registrar, se llamó %d veces", e.defectos.registros)
			}
		})
	}
}

func TestRegistrarDefecto_Rechazos(t *testing.T) {
	casos := []struct {
		nombre  string
		cambiar func(*RegistrarDefectoInput)
		campo   string
		mensaje string
	}{
		{"sin historia", func(in *RegistrarDefectoInput) { in.HistoriaID = 0 }, "historia_id", domain.ErrHistoriaRelacionadaObligatoria.Mensaje},
		{"historia inexistente", func(in *RegistrarDefectoInput) { in.HistoriaID = 999 }, "historia_id", domain.ErrHistoriaRelacionadaAjena.Mensaje},
		{"historia de otro proyecto", func(in *RegistrarDefectoInput) { in.HistoriaID = 4 }, "historia_id", domain.ErrHistoriaRelacionadaAjena.Mensaje},
		{"descripción vacía", func(in *RegistrarDefectoInput) { in.Descripcion = "  " }, "descripcion", "la descripción es obligatoria"},
		{"severidad inválida", func(in *RegistrarDefectoInput) { in.Severidad = "alta" }, "severidad", "la severidad debe ser Alta, Media o Baja"},
		{"historia ajena antes que descripción", func(in *RegistrarDefectoInput) { in.HistoriaID, in.Descripcion = 4, "" }, "historia_id", domain.ErrHistoriaRelacionadaAjena.Mensaje},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := nuevoEscenarioDefecto()
			e.sprint(t, 1, domain.SprintActivo)
			in := defectoValido()
			c.cambiar(&in)

			_, err := e.registrar().Ejecutar(context.Background(), in)

			var verr domain.ValidationError
			if !errors.As(err, &verr) || verr.Campo != c.campo || verr.Mensaje != c.mensaje {
				t.Fatalf("se esperaba %s: %s, se obtuvo %v", c.campo, c.mensaje, err)
			}
			if e.defectos.registros != 0 {
				t.Errorf("no se esperaba llamar a Registrar, se llamó %d veces", e.defectos.registros)
			}
		})
	}
}

func TestRegistrarDefecto_SinSprintActivoVaPrimero(t *testing.T) {
	e := nuevoEscenarioDefecto()
	in := defectoValido()
	in.HistoriaID, in.Descripcion = 4, ""

	_, err := e.registrar().Ejecutar(context.Background(), in)

	if !errors.Is(err, domain.ErrSinSprintActivoParaDetectar) {
		t.Fatalf("se esperaba ErrSinSprintActivoParaDetectar, se obtuvo %v", err)
	}
}

func TestRegistrarDefecto_PropagaElCierreDelSprintAlRegistrar(t *testing.T) {
	e := nuevoEscenarioDefecto()
	e.sprint(t, 1, domain.SprintActivo)
	e.defectos.errRegistrar = domain.ErrSinSprintActivoParaDetectar

	_, err := e.registrar().Ejecutar(context.Background(), defectoValido())

	if !errors.Is(err, domain.ErrSinSprintActivoParaDetectar) {
		t.Fatalf("se esperaba ErrSinSprintActivoParaDetectar, se obtuvo %v", err)
	}
}
