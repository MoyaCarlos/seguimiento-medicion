package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

type escenarioResolver struct {
	sprints  *fakeSprintsQueCuentan
	defectos *fakeDefectoRepository
}

func nuevoEscenarioResolver() escenarioResolver {
	return escenarioResolver{
		sprints:  &fakeSprintsQueCuentan{SprintRepositoryEnMemoria: repository.NewSprintRepositoryEnMemoria()},
		defectos: &fakeDefectoRepository{titulos: map[int64]string{3: "Carrito de compras"}},
	}
}

func (e escenarioResolver) sprint(t *testing.T, proyectoID int64, estado domain.EstadoSprint) domain.Sprint {
	t.Helper()
	s, err := e.sprints.Guardar(context.Background(), domain.Sprint{ProyectoID: proyectoID, SprintGoal: "g", Estado: estado})
	if err != nil {
		t.Fatalf("no se pudo sembrar el Sprint: %v", err)
	}
	return s
}

// pendiente siembra un defecto Pendiente del proyecto 1 detectado en el Sprint dado.
func (e escenarioResolver) pendiente(t *testing.T, deteccion domain.Sprint) domain.Defecto {
	t.Helper()
	d, err := e.defectos.Registrar(context.Background(), domain.Defecto{
		ProyectoID: 1, HistoriaID: 3, SprintDeteccionID: deteccion.ID,
		Descripcion: "El total no se actualiza", Severidad: domain.SeveridadAlta, Estado: domain.DefectoPendiente,
	})
	if err != nil {
		t.Fatalf("no se pudo sembrar el defecto: %v", err)
	}
	return d
}

func (e escenarioResolver) resolver() *ResolverDefecto {
	return NewResolverDefecto(e.defectos, e.sprints)
}

func TestResolverDefecto_ResuelveEnElSprintActivo(t *testing.T) {
	e := nuevoEscenarioResolver()
	activo := e.sprint(t, 1, domain.SprintActivo)
	d := e.pendiente(t, activo)

	detalle, err := e.resolver().Ejecutar(context.Background(), d.ID)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if e.defectos.resoluciones != 1 {
		t.Fatalf("se esperaba 1 llamada a Resolver, hubo %d", e.defectos.resoluciones)
	}
	guardado := e.defectos.defectos[0]
	if guardado.Estado != domain.DefectoResuelto || guardado.SprintResolucionID == nil || *guardado.SprintResolucionID != activo.ID {
		t.Errorf("defecto guardado inesperado: %+v", guardado)
	}
	if detalle.ID != d.ID || detalle.Estado != domain.DefectoResuelto || detalle.SprintResolucionID == nil ||
		*detalle.SprintResolucionID != activo.ID || detalle.SprintDeteccionID != activo.ID || detalle.TituloHistoria != "Carrito de compras" {
		t.Errorf("detalle inesperado: %+v", detalle)
	}
}

func TestResolverDefecto_DetectadoEnUnSprintFinalizado(t *testing.T) {
	e := nuevoEscenarioResolver()
	finalizado := e.sprint(t, 1, domain.SprintFinalizado)
	activo := e.sprint(t, 1, domain.SprintActivo)
	d := e.pendiente(t, finalizado)

	detalle, err := e.resolver().Ejecutar(context.Background(), d.ID)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if detalle.SprintDeteccionID != finalizado.ID || *detalle.SprintResolucionID != activo.ID {
		t.Errorf("se esperaba detectado en %d y resuelto en %d, se obtuvo %+v", finalizado.ID, activo.ID, detalle)
	}
}

func TestResolverDefecto_Inexistente(t *testing.T) {
	e := nuevoEscenarioResolver()
	e.sprint(t, 1, domain.SprintActivo)

	_, err := e.resolver().Ejecutar(context.Background(), 999)

	if !errors.Is(err, domain.ErrDefectoNoEncontrado) {
		t.Fatalf("se esperaba ErrDefectoNoEncontrado, se obtuvo %v", err)
	}
	if e.sprints.listados != 0 {
		t.Errorf("no se esperaba consultar los Sprints, se consultaron %d veces", e.sprints.listados)
	}
}

func TestResolverDefecto_YaResuelto(t *testing.T) {
	e := nuevoEscenarioResolver()
	activo := e.sprint(t, 1, domain.SprintActivo)
	d := e.pendiente(t, activo)
	e.defectos.defectos[0].Estado, e.defectos.defectos[0].SprintResolucionID = domain.DefectoResuelto, &activo.ID

	_, err := e.resolver().Ejecutar(context.Background(), d.ID)

	if !errors.Is(err, domain.ErrDefectoYaResuelto) {
		t.Fatalf("se esperaba ErrDefectoYaResuelto, se obtuvo %v", err)
	}
	if e.defectos.resoluciones != 0 {
		t.Errorf("no se esperaba llamar a Resolver, se llamó %d veces", e.defectos.resoluciones)
	}
}

func TestResolverDefecto_SinSprintActivo(t *testing.T) {
	casos := []struct {
		nombre        string
		otroSprintEn  int64
		estadoDelOtro domain.EstadoSprint
	}{
		{"solo Finalizado", 1, domain.SprintFinalizado},
		{"solo Pendiente", 1, domain.SprintPendiente},
		{"el Activo es de otro proyecto", 2, domain.SprintActivo},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			e := nuevoEscenarioResolver()
			deteccion := e.sprint(t, 1, domain.SprintFinalizado)
			e.sprint(t, c.otroSprintEn, c.estadoDelOtro)
			d := e.pendiente(t, deteccion)

			_, err := e.resolver().Ejecutar(context.Background(), d.ID)

			if !errors.Is(err, domain.ErrSinSprintActivoParaResolver) {
				t.Fatalf("se esperaba ErrSinSprintActivoParaResolver, se obtuvo %v", err)
			}
			if e.defectos.resoluciones != 0 {
				t.Errorf("no se esperaba llamar a Resolver, se llamó %d veces", e.defectos.resoluciones)
			}
		})
	}
}

func TestResolverDefecto_PropagaLasCarrerasDelRepositorio(t *testing.T) {
	for _, carrera := range []error{domain.ErrDefectoYaResuelto, domain.ErrSinSprintActivoParaResolver} {
		t.Run(carrera.Error(), func(t *testing.T) {
			e := nuevoEscenarioResolver()
			d := e.pendiente(t, e.sprint(t, 1, domain.SprintActivo))
			e.defectos.errResolver = carrera

			_, err := e.resolver().Ejecutar(context.Background(), d.ID)

			if !errors.Is(err, carrera) {
				t.Fatalf("se esperaba %v, se obtuvo %v", carrera, err)
			}
		})
	}
}
