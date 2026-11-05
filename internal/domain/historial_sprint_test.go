package domain

import (
	"errors"
	"slices"
	"testing"
)

func sprintDeHistorial(id int64, numero int, estado EstadoSprint) Sprint {
	return Sprint{ID: id, ProyectoID: 1, Numero: numero, Estado: estado}
}

func numerosDe(sprints []Sprint) []int {
	numeros := make([]int, 0, len(sprints))
	for _, s := range sprints {
		numeros = append(numeros, s.Numero)
	}
	return numeros
}

func TestSprintsFinalizados(t *testing.T) {
	casos := []struct {
		nombre   string
		sprints  []Sprint
		esperado []int
	}{
		{
			nombre: "solo los Finalizado, sin Activo ni Pendiente",
			sprints: []Sprint{
				sprintDeHistorial(1, 1, SprintFinalizado),
				sprintDeHistorial(2, 2, SprintActivo),
				sprintDeHistorial(3, 3, SprintPendiente),
			},
			esperado: []int{1},
		},
		{
			nombre: "ordenados por número aunque lleguen desordenados",
			sprints: []Sprint{
				sprintDeHistorial(3, 3, SprintFinalizado),
				sprintDeHistorial(1, 1, SprintFinalizado),
				sprintDeHistorial(2, 2, SprintFinalizado),
			},
			esperado: []int{1, 2, 3},
		},
		{
			nombre:   "sin finalizados",
			sprints:  []Sprint{sprintDeHistorial(1, 1, SprintActivo)},
			esperado: []int{},
		},
		{
			nombre:   "entrada vacía",
			sprints:  nil,
			esperado: []int{},
		},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			obtenidos := SprintsFinalizados(c.sprints)

			if obtenidos == nil {
				t.Fatal("se esperaba un slice vacío, nunca nil")
			}
			if numeros := numerosDe(obtenidos); !slices.Equal(numeros, c.esperado) {
				t.Errorf("se esperaban los Sprints %v, se obtuvieron %v", c.esperado, numeros)
			}
		})
	}
}

func TestValidarSprintFinalizado(t *testing.T) {
	casos := []struct {
		estado   EstadoSprint
		esperado error
	}{
		{SprintFinalizado, nil},
		{SprintActivo, ErrSprintNoFinalizado},
		{SprintPendiente, ErrSprintNoFinalizado},
	}
	for _, c := range casos {
		t.Run(string(c.estado), func(t *testing.T) {
			err := ValidarSprintFinalizado(sprintDeHistorial(1, 1, c.estado))

			if !errors.Is(err, c.esperado) {
				t.Errorf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
		})
	}
}
