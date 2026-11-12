package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// HU-16 (US1 esc. 4, FR-001, SC-002): el equipo del reporte lista solo a los
// activos, y las horas y métricas no cambian cuando se quita a alguien que
// imputó horas en un Sprint finalizado.
func TestGenerarReporte_EjecutarProyecto_ExcluyeAlDadoDeBajaSinCambiarSusHoras(t *testing.T) {
	ctx := context.Background()
	proyectos := &fakeProjectRepository{
		proyectos: map[int64]domain.Project{1: {ID: 1, Nombre: "Proyecto", HorasPorStoryPoint: punteroADecimal(4)}},
		listaMiembros: []domain.Member{
			{UserID: carlos, ProjectID: 1, Nombre: "Carlos", Role: domain.RolScrumMaster},
			{UserID: candela, ProjectID: 1, Nombre: "Candela", Role: domain.RolProductBuilder},
		},
	}
	sprints := repository.NewSprintRepositoryEnMemoria()
	finalizado := sprintDeReporte(t, sprints, 1, "Cerrar el carrito", domain.SprintFinalizado)
	historias := &fakeHistoriasDeSprintRepository{porSprint: map[int64][]domain.BacklogItem{
		finalizado.ID: {historiaDeReporte(10, 1, &finalizado.ID, "Carrito de compras", domain.PrioridadMust, domain.EstadoCompletada, 5)},
	}}
	// Los 360 minutos son las horas que imputó Candela en el Sprint.
	datos := &fakeDatosMetricasRepository{datos: []domain.DatosMetricasSprint{{
		SprintID: finalizado.ID, Numero: finalizado.Numero, HistoriasCompletadas: 1, SPCompletados: 5, MinutosReales: 360,
	}}}
	sut := nuevoGenerarReporteDePrueba(proyectos, sprints, &fakeProductBacklogRepository{}, historias,
		&fakeArrastradasRepository{}, datos, &fakeDefectoRepository{})

	antes, err := sut.EjecutarProyecto(ctx, 1, fechaDeReporte(20))
	if err != nil {
		t.Fatalf("reporte antes de la baja: %v", err)
	}
	darDeBaja(t, proyectos, candela)
	despues, err := sut.EjecutarProyecto(ctx, 1, fechaDeReporte(20))
	if err != nil {
		t.Fatalf("reporte después de la baja: %v", err)
	}

	if len(despues.Integrantes) != 1 || despues.Integrantes[0].UserID != carlos {
		t.Errorf("el equipo del reporte debe listar solo a Carlos, se obtuvo %+v", despues.Integrantes)
	}
	if despues.Metricas.HorasReales != 6 || !reflect.DeepEqual(antes.Metricas, despues.Metricas) {
		t.Errorf("las métricas del proyecto no deben cambiar (6 horas reales):\n antes   %+v\n después %+v", antes.Metricas, despues.Metricas)
	}
	if !reflect.DeepEqual(antes.Sprints, despues.Sprints) {
		t.Errorf("los Sprints del reporte no deben cambiar:\n antes   %+v\n después %+v", antes.Sprints, despues.Sprints)
	}
}
