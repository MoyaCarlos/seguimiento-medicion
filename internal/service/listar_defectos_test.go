package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func nuevoListarDefectos() (*ListarDefectos, *fakeDefectoRepository) {
	proyectos := &fakeProjectRepository{proyectos: map[int64]domain.Project{1: {ID: 1}, 2: {ID: 2}}}
	defectos := &fakeDefectoRepository{titulos: map[int64]string{3: "Carrito de compras"}}
	return NewListarDefectos(proyectos, defectos), defectos
}

func TestListarDefectos_DevuelveLosDelProyecto(t *testing.T) {
	listar, defectos := nuevoListarDefectos()
	for _, proyectoID := range []int64{1, 2, 1} {
		if _, err := defectos.Registrar(context.Background(), domain.Defecto{ProyectoID: proyectoID, HistoriaID: 3}); err != nil {
			t.Fatalf("no se pudo sembrar el defecto: %v", err)
		}
	}

	got, err := listar.Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 3 || got[0].TituloHistoria != "Carrito de compras" {
		t.Errorf("se esperaban los defectos 1 y 3 del proyecto 1, se obtuvo %+v", got)
	}
}

func TestListarDefectos_ProyectoSinDefectos(t *testing.T) {
	listar, _ := nuevoListarDefectos()

	got, err := listar.Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("se esperaba un slice vacío no nil, se obtuvo %#v", got)
	}
}

func TestListarDefectos_ProyectoInexistente(t *testing.T) {
	listar, defectos := nuevoListarDefectos()

	_, err := listar.Ejecutar(context.Background(), 999)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
	if defectos.listados != 0 {
		t.Errorf("no se esperaba consultar los defectos, se consultaron %d veces", defectos.listados)
	}
}
