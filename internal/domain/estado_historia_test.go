package domain

import (
	"errors"
	"testing"
)

func TestValidarCambioDeEstado(t *testing.T) {
	var sprintID int64 = 1
	sprint := func(e EstadoSprint) *Sprint { return &Sprint{ID: sprintID, Estado: e} }
	historia := func(e Estado, conSprint bool) BacklogItem {
		h := BacklogItem{ID: 1, Estado: e}
		if conSprint {
			h.SprintID = &sprintID
		}
		return h
	}

	casos := []struct {
		nombre string
		h      BacklogItem
		s      *Sprint
		nuevo  Estado
		want   error
	}{
		{"a Nueva a En progreso", historia(EstadoNueva, true), sprint(SprintActivo), EstadoEnProgreso, nil},
		{"b En progreso a Completada", historia(EstadoEnProgreso, true), sprint(SprintActivo), EstadoCompletada, nil},
		{"c estado destino inválido", historia(EstadoNueva, true), sprint(SprintActivo), Estado("Pendiente"), ErrEstadoInvalido},
		{"d sin Sprint a En progreso", historia(EstadoNueva, false), nil, EstadoEnProgreso, ErrHistoriaSinSprint},
		{"e sin Sprint a Completada", historia(EstadoNueva, false), nil, EstadoCompletada, ErrHistoriaSinSprint},
		{"f sin Sprint a Nueva", historia(EstadoNueva, false), nil, EstadoNueva, ErrEstadoRepetido},
		{"g Sprint Pendiente", historia(EstadoNueva, true), sprint(SprintPendiente), EstadoEnProgreso, ErrSprintNoActivoParaEstado},
		{"h Sprint Finalizado", historia(EstadoEnProgreso, true), sprint(SprintFinalizado), EstadoCompletada, ErrSprintNoActivoParaEstado},
		{"i En progreso repetido", historia(EstadoEnProgreso, true), sprint(SprintActivo), EstadoEnProgreso, ErrEstadoRepetido},
		{"j Completada repetido", historia(EstadoCompletada, true), sprint(SprintActivo), EstadoCompletada, ErrEstadoRepetido},
		{"k En progreso a Nueva", historia(EstadoEnProgreso, true), sprint(SprintActivo), EstadoNueva, ErrRetrocesoEstado},
		{"l Completada a En progreso", historia(EstadoCompletada, true), sprint(SprintActivo), EstadoEnProgreso, ErrRetrocesoEstado},
		{"m Nueva a Completada", historia(EstadoNueva, true), sprint(SprintActivo), EstadoCompletada, ErrSaltoDirecto},
		{"n Sprint Finalizado y estado repetido", historia(EstadoEnProgreso, true), sprint(SprintFinalizado), EstadoEnProgreso, ErrSprintNoActivoParaEstado},
		{"o sin Sprint y salto directo", historia(EstadoNueva, false), nil, EstadoCompletada, ErrHistoriaSinSprint},
		{"p estado inválido y sin Sprint", historia(EstadoNueva, false), nil, Estado("Pendiente"), ErrEstadoInvalido},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			err := ValidarCambioDeEstado(c.h, c.s, c.nuevo)
			if c.want == nil {
				if err != nil {
					t.Fatalf("se esperaba nil, llegó %v", err)
				}
				return
			}
			if !errors.Is(err, c.want) {
				t.Fatalf("se esperaba %v, llegó %v", c.want, err)
			}
			var ve ValidationError
			if !errors.As(err, &ve) || ve.Campo != "estado" {
				t.Fatalf("el campo del error debe ser estado, llegó %v", err)
			}
		})
	}
}
