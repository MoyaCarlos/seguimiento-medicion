package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

func strPtr(v string) *string { return &v }

func historiaDelBacklog() domain.BacklogItem {
	return domain.BacklogItem{
		ID: 3, ProyectoID: 1, Titulo: "Original", Descripcion: "Descripción original",
		Prioridad: domain.PrioridadMust, Estado: domain.EstadoNueva,
	}
}

func edicionValida(id int64) EditarHistoriaInput {
	return EditarHistoriaInput{
		ID: id, Titulo: "Nuevo", Descripcion: "Descripción nueva",
		CriteriosAceptacion: strPtr("Dado, cuando, entonces"), EstimacionSP: intPtr(5),
	}
}

func nuevoEditarHistoria(historias ...domain.BacklogItem) (*EditarHistoria, *fakeHistoriaRepository, *repository.SprintRepositoryEnMemoria) {
	repo := &fakeHistoriaRepository{historias: map[int64]domain.BacklogItem{}}
	for _, h := range historias {
		repo.historias[h.ID] = h
	}
	sprints := repository.NewSprintRepositoryEnMemoria()
	return NewEditarHistoria(repo, sprints), repo, sprints
}

func TestEditarHistoria_Inexistente_DominaSobreLaValidacion(t *testing.T) {
	servicio, repo, _ := nuevoEditarHistoria()

	_, err := servicio.Ejecutar(context.Background(), EditarHistoriaInput{ID: 999, Titulo: ""})

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Fatalf("se esperaba domain.ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
	if len(repo.actualizaciones) != 0 {
		t.Errorf("no se esperaba actualizar")
	}
}

func TestEditarHistoria_Valida_SinSprint(t *testing.T) {
	servicio, repo, _ := nuevoEditarHistoria(historiaDelBacklog())

	editada, err := servicio.Ejecutar(context.Background(), edicionValida(3))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if editada.ID != 3 || editada.Titulo != "Nuevo" || editada.EstimacionSP == nil || *editada.EstimacionSP != 5 {
		t.Errorf("historia editada inesperada: %+v", editada)
	}
	if editada.Estado != domain.EstadoNueva || editada.Prioridad != domain.PrioridadMust {
		t.Errorf("se modificaron el estado o la prioridad: %+v", editada)
	}
	if len(repo.actualizaciones) != 1 || repo.actualizaciones[0].Titulo != "Nuevo" {
		t.Errorf("se esperaba persistir la historia editada, se obtuvo %+v", repo.actualizaciones)
	}
}

func TestEditarHistoria_Invalida_NoActualiza(t *testing.T) {
	servicio, repo, _ := nuevoEditarHistoria(historiaDelBacklog())
	input := edicionValida(3)
	input.Titulo = "   "

	_, err := servicio.Ejecutar(context.Background(), input)

	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "titulo" {
		t.Fatalf("se esperaba ValidationError sobre titulo, se obtuvo %v", err)
	}
	if len(repo.actualizaciones) != 0 {
		t.Errorf("no se esperaba actualizar una edición inválida")
	}
}

func TestEditarHistoria_PropagaErrorDeActualizar(t *testing.T) {
	servicio, repo, _ := nuevoEditarHistoria(historiaDelBacklog())
	errPersistencia := errors.New("falla de persistencia")
	repo.errActualizar = errPersistencia

	_, err := servicio.Ejecutar(context.Background(), edicionValida(3))

	if !errors.Is(err, errPersistencia) {
		t.Fatalf("se esperaba el error de persistencia, se obtuvo %v", err)
	}
}

func TestEditarHistoria_StoryPointsSegunEstadoDelSprint(t *testing.T) {
	casos := []struct {
		nombre         string
		estadoSprint   domain.EstadoSprint
		spActual       *int
		spNuevo        *int
		debeRechazarse bool
	}{
		{"pendiente permite cambiarlos", domain.SprintPendiente, intPtr(3), intPtr(5), false},
		{"activo rechaza cambiarlos", domain.SprintActivo, intPtr(3), intPtr(5), true},
		{"finalizado rechaza cambiarlos", domain.SprintFinalizado, intPtr(3), intPtr(5), true},
		{"activo rechaza borrarlos", domain.SprintActivo, intPtr(3), nil, true},
		{"finalizado rechaza borrarlos", domain.SprintFinalizado, intPtr(3), nil, true},
		{"activo permite la carga inicial", domain.SprintActivo, nil, intPtr(5), false},
		{"finalizado permite la carga inicial", domain.SprintFinalizado, nil, intPtr(5), false},
		{"activo permite reenviar el mismo valor", domain.SprintActivo, intPtr(3), intPtr(3), false},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			historia := historiaDelBacklog()
			historia.EstimacionSP = tc.spActual
			servicio, repo, sprints := nuevoEditarHistoria()
			sprint, _ := sprints.Guardar(context.Background(), domain.Sprint{ProyectoID: 1, Estado: tc.estadoSprint})
			historia.SprintID = &sprint.ID
			repo.historias[historia.ID] = historia
			input := edicionValida(historia.ID)
			input.EstimacionSP = tc.spNuevo

			editada, err := servicio.Ejecutar(context.Background(), input)

			if tc.debeRechazarse {
				if !errors.Is(err, domain.ErrEstimacionCongelada) {
					t.Fatalf("se esperaba domain.ErrEstimacionCongelada, se obtuvo %v", err)
				}
				if len(repo.actualizaciones) != 0 {
					t.Errorf("no se esperaba actualizar la historia")
				}
				return
			}
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if editada.Titulo != "Nuevo" || len(repo.actualizaciones) != 1 {
				t.Errorf("se esperaba guardar la edición completa: %+v", editada)
			}
		})
	}
}
