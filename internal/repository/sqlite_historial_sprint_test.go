package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type historiaEsperada struct {
	id         int64
	titulo     string
	sp         any
	completada bool
}

func resumirHistorias(historias []domain.HistoriaComprometida) []historiaEsperada {
	resumen := make([]historiaEsperada, 0, len(historias))
	for _, h := range historias {
		var sp any
		if h.EstimacionSP != nil {
			sp = *h.EstimacionSP
		}
		resumen = append(resumen, historiaEsperada{h.ID, h.Titulo, sp, h.Completada})
	}
	return resumen
}

func TestSQLiteHistorialSprint_HistoriasComprometidas(t *testing.T) {
	db := nuevaBDMigrada(t)
	e := nuevoEscenarioMetricas(t, db, "Ana")
	primero := e.sprint("Finalizado")
	segundo := e.sprint("Finalizado")

	arrastradaYCompletada := e.enBacklog(3)
	e.arrastrada(primero, arrastradaYCompletada, 3)
	e.insertar(`UPDATE historias SET estado = 'Completada', estimacion_sp = 8, sprint_id = ?, fecha_finalizacion = '2026-10-25'
		WHERE id = ?`, segundo, arrastradaYCompletada)
	sinEstimar := e.completada(primero, nil)
	cinco := e.completada(primero, 5)
	arrastradaSinEstimar := e.enBacklog(nil)
	e.arrastrada(primero, arrastradaSinEstimar, nil)
	dos := e.completada(segundo, 2)

	repo := NewSQLiteHistorialSprintRepository(db)
	casos := []struct {
		nombre   string
		sprintID int64
		esperado []historiaEsperada
	}{
		{"completadas primero y arrastradas después, cada grupo por id, con la foto de las arrastradas", primero, []historiaEsperada{
			{sinEstimar, "Completada", nil, true},
			{cinco, "Completada", 5, true},
			{arrastradaYCompletada, "En el backlog", 3, false},
			{arrastradaSinEstimar, "En el backlog", nil, false},
		}},
		{"la arrastrada figura completada en el Sprint donde se completó", segundo, []historiaEsperada{
			{arrastradaYCompletada, "En el backlog", 8, true},
			{dos, "Completada", 2, true},
		}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			historias, err := repo.HistoriasComprometidas(context.Background(), c.sprintID)

			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if obtenido := resumirHistorias(historias); !reflect.DeepEqual(obtenido, c.esperado) {
				t.Errorf("historias inesperadas:\n se obtuvo %+v\n se esperaba %+v", obtenido, c.esperado)
			}
		})
	}
}

func TestSQLiteHistorialSprint_MinutosPorHistoriaSoloDelSprint(t *testing.T) {
	db := nuevaBDMigrada(t)
	e := nuevoEscenarioMetricas(t, db, "Ana")
	primero := e.sprint("Finalizado")
	segundo := e.sprint("Finalizado")
	historia := e.enBacklog(3)
	e.arrastrada(primero, historia, 3)
	otra := e.completada(primero, 5)
	e.horas(historia, primero, 30)
	e.horas(historia, primero, 20)
	e.horas(otra, primero, 90)
	e.horas(historia, segundo, 600)

	minutos, err := NewSQLiteHistorialSprintRepository(db).MinutosPorHistoria(context.Background(), primero)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if esperado := map[int64]int{historia: 50, otra: 90}; !reflect.DeepEqual(minutos, esperado) {
		t.Errorf("se esperaba %v, se obtuvo %v", esperado, minutos)
	}
}

func TestSQLiteHistorialSprint_DefectosDetectadosOResueltosEnElSprint(t *testing.T) {
	db := nuevaBDMigrada(t)
	e := nuevoEscenarioMetricas(t, db, "Ana")
	anterior := e.sprint("Finalizado")
	consultado := e.sprint("Finalizado")
	posterior := e.sprint("Activo")
	historia := e.completada(consultado, 5)

	e.defecto(historia, consultado, nil)
	e.defecto(historia, anterior, consultado)
	e.defecto(historia, consultado, consultado)
	e.defecto(historia, consultado, posterior)
	e.defecto(historia, anterior, nil)
	e.defecto(historia, posterior, posterior)

	defectos, err := NewSQLiteHistorialSprintRepository(db).DefectosDeSprint(context.Background(), consultado)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	var ids []int64
	for _, d := range defectos {
		ids = append(ids, d.ID)
		if d.TituloHistoria != "Completada" {
			t.Errorf("defecto %d: se esperaba el título de su historia, se obtuvo %q", d.ID, d.TituloHistoria)
		}
	}
	if esperado := []int64{1, 2, 3, 4}; !reflect.DeepEqual(ids, esperado) {
		t.Errorf("se esperaban los defectos %v (cada uno una vez, por id), se obtuvieron %v", esperado, ids)
	}
}

func TestSQLiteHistorialSprint_SprintSinDatosDevuelveVaciosNoNil(t *testing.T) {
	db := nuevaBDMigrada(t)
	e := nuevoEscenarioMetricas(t, db, "Ana")
	vacio := e.sprint("Finalizado")
	repo := NewSQLiteHistorialSprintRepository(db)
	ctx := context.Background()

	historias, err := repo.HistoriasComprometidas(ctx, vacio)
	if err != nil || historias == nil || len(historias) != 0 {
		t.Errorf("historias: se esperaba un slice vacío no nil, se obtuvo %#v (err %v)", historias, err)
	}
	minutos, err := repo.MinutosPorHistoria(ctx, vacio)
	if err != nil || minutos == nil || len(minutos) != 0 {
		t.Errorf("minutos: se esperaba un mapa vacío no nil, se obtuvo %#v (err %v)", minutos, err)
	}
	defectos, err := repo.DefectosDeSprint(ctx, vacio)
	if err != nil || defectos == nil || len(defectos) != 0 {
		t.Errorf("defectos: se esperaba un slice vacío no nil, se obtuvo %#v (err %v)", defectos, err)
	}
}
