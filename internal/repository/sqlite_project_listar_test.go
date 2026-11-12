package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// HU-21, escenario 5: «Mis proyectos» trae solo los proyectos que integra el
// usuario, con el rol de su fila en integrantes, ordenados por id.
func TestSQLiteProjectRepository_ListarPorIntegrante(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	mustExec(t, db, `INSERT INTO usuarios (id, nombre, nombre_normalizado, creado_en) VALUES
		(1, 'Ana', 'ana', '2026-10-23T00:00:00Z'),
		(2, 'Bruno', 'bruno', '2026-10-23T00:00:00Z'),
		(3, 'Carla', 'carla', '2026-10-23T00:00:00Z')`)
	mustExec(t, db, `INSERT INTO proyectos (id, nombre, descripcion, fecha_inicio, fecha_fin, creado_en) VALUES
		(9, 'Nueve', '', NULL, NULL, '2026-10-23T00:00:00Z'),
		(2, 'Dos', 'Descripción', '2026-03-01', NULL, '2026-10-23T00:00:00Z'),
		(5, 'Cinco', '', NULL, NULL, '2026-10-23T00:00:00Z')`)
	mustExec(t, db, `INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES
		(9, 1, 'scrum_master', '2026-10-23T00:00:00Z'),
		(2, 1, 'product_builder', '2026-10-23T00:00:00Z'),
		(2, 2, 'scrum_master', '2026-10-23T00:00:00Z'),
		(5, 2, 'scrum_master', '2026-10-23T00:00:00Z')`)

	deAna, err := repo.ListarPorIntegrante(context.Background(), 1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(deAna) != 2 {
		t.Fatalf("se esperaban 2 proyectos de Ana, se obtuvieron %+v", deAna)
	}
	if deAna[0].Proyecto.ID != 2 || deAna[0].Rol != domain.RolProductBuilder {
		t.Errorf("primero: se esperaba el proyecto 2 como product_builder, se obtuvo %+v", deAna[0])
	}
	if deAna[1].Proyecto.ID != 9 || deAna[1].Rol != domain.RolScrumMaster {
		t.Errorf("segundo: se esperaba el proyecto 9 como scrum_master, se obtuvo %+v", deAna[1])
	}
	individual, err := repo.ObtenerPorID(context.Background(), 2)
	if err != nil {
		t.Fatalf("no se pudo consultar el proyecto 2: %v", err)
	}
	if !reflect.DeepEqual(deAna[0].Proyecto, individual) {
		t.Errorf("el proyecto del listado difiere de su consulta individual: %+v vs %+v", deAna[0].Proyecto, individual)
	}

	deCarla, err := repo.ListarPorIntegrante(context.Background(), 3)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if deCarla == nil || len(deCarla) != 0 {
		t.Fatalf("se esperaba un arreglo vacío (no nil) para quien no integra proyectos, se obtuvo %#v", deCarla)
	}
}
