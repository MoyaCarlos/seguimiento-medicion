package domain

import "testing"

func historiaDeAvance(id int64, titulo string, estado Estado, sp *int) BacklogItem {
	return BacklogItem{ID: id, ProyectoID: 1, Titulo: titulo, Estado: estado, EstimacionSP: sp}
}

func TestAgruparAvance_SiempreDevuelveLosTresGruposEnOrden(t *testing.T) {
	for nombre, historias := range map[string][]BacklogItem{
		"sin historias":                   nil,
		"con historias en un solo estado": {historiaDeAvance(1, "H1", EstadoNueva, intPtr(3))},
	} {
		t.Run(nombre, func(t *testing.T) {
			avance, err := AgruparAvance(historias, []Sprint{})
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}

			esperados := []Estado{EstadoNueva, EstadoEnProgreso, EstadoCompletada}
			if len(avance.Grupos) != len(esperados) {
				t.Fatalf("se esperaban %d grupos, se obtuvieron %d", len(esperados), len(avance.Grupos))
			}
			for i, estado := range esperados {
				if avance.Grupos[i].Estado != estado {
					t.Errorf("grupo %d: se esperaba %q, se obtuvo %q", i, estado, avance.Grupos[i].Estado)
				}
				if avance.Grupos[i].Historias == nil {
					t.Errorf("grupo %q: Historias no debe ser nil", estado)
				}
			}
		})
	}
}

func TestAgruparAvance_ProyectoSinHistoriasTieneTodoEnCero(t *testing.T) {
	avance, err := AgruparAvance(nil, []Sprint{})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(avance.Grupos) != 3 {
		t.Fatalf("se esperaban 3 grupos, se obtuvieron %d", len(avance.Grupos))
	}
	for _, g := range avance.Grupos {
		if g.Cantidad != 0 || g.StoryPoints != 0 || len(g.Historias) != 0 {
			t.Errorf("grupo %q: se esperaba vacío, se obtuvo %+v", g.Estado, g)
		}
	}
	if avance.Total != (TotalesAvance{}) {
		t.Errorf("se esperaban totales en cero, se obtuvo %+v", avance.Total)
	}
}

func TestAgruparAvance_CuentaYSumaStoryPointsPorGrupo(t *testing.T) {
	historias := []BacklogItem{
		historiaDeAvance(1, "H1", EstadoNueva, intPtr(3)),
		historiaDeAvance(2, "H2", EstadoEnProgreso, intPtr(5)),
		historiaDeAvance(3, "H3", EstadoCompletada, intPtr(8)),
		historiaDeAvance(4, "H4", EstadoCompletada, intPtr(2)),
	}

	avance, err := AgruparAvance(historias, []Sprint{})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	esperados := []struct {
		estado      Estado
		cantidad    int
		storyPoints int
	}{
		{EstadoNueva, 1, 3},
		{EstadoEnProgreso, 1, 5},
		{EstadoCompletada, 2, 10},
	}
	if len(avance.Grupos) != len(esperados) {
		t.Fatalf("se esperaban %d grupos, se obtuvieron %d", len(esperados), len(avance.Grupos))
	}
	for i, e := range esperados {
		g := avance.Grupos[i]
		if g.Estado != e.estado || g.Cantidad != e.cantidad || g.StoryPoints != e.storyPoints || len(g.Historias) != e.cantidad {
			t.Errorf("grupo %d: se esperaba %+v, se obtuvo estado %q cantidad %d SP %d historias %d",
				i, e, g.Estado, g.Cantidad, g.StoryPoints, len(g.Historias))
		}
	}
}

func TestAgruparAvance_TotalEsLaSumaDeLosTresGrupos(t *testing.T) {
	historias := []BacklogItem{
		historiaDeAvance(1, "H1", EstadoNueva, intPtr(3)),
		historiaDeAvance(2, "H2", EstadoEnProgreso, intPtr(5)),
		historiaDeAvance(3, "H3", EstadoCompletada, intPtr(10)),
		historiaDeAvance(4, "H4", EstadoNueva, nil),
	}

	avance, err := AgruparAvance(historias, []Sprint{})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if avance.Total != (TotalesAvance{Historias: 4, StoryPoints: 18}) {
		t.Errorf("se esperaban 4 historias y 18 Story Points, se obtuvo %+v", avance.Total)
	}
}

func TestAgruparAvance_HistoriaSinEstimarCuentaYSumaCero(t *testing.T) {
	historias := []BacklogItem{
		historiaDeAvance(1, "H1", EstadoNueva, intPtr(3)),
		historiaDeAvance(2, "H2", EstadoNueva, nil),
	}

	avance, err := AgruparAvance(historias, []Sprint{})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(avance.Grupos) != 3 || len(avance.Grupos[0].Historias) != 2 {
		t.Fatalf("se esperaban 3 grupos y 2 historias en Nueva, se obtuvo %+v", avance.Grupos)
	}
	nueva := avance.Grupos[0]
	if nueva.Cantidad != 2 || nueva.StoryPoints != 3 {
		t.Fatalf("se esperaban 2 historias y 3 Story Points, se obtuvo %d y %d", nueva.Cantidad, nueva.StoryPoints)
	}
	if nueva.Historias[1].EstimacionSP != nil {
		t.Errorf("la historia sin estimar debe conservar EstimacionSP nil, se obtuvo %d", *nueva.Historias[1].EstimacionSP)
	}
}

func TestAgruparAvance_ConservaElOrdenRecibidoDentroDeCadaGrupo(t *testing.T) {
	historias := []BacklogItem{
		historiaDeAvance(7, "H7", EstadoCompletada, intPtr(1)),
		historiaDeAvance(2, "H2", EstadoNueva, intPtr(1)),
		historiaDeAvance(5, "H5", EstadoCompletada, intPtr(1)),
		historiaDeAvance(3, "H3", EstadoNueva, intPtr(1)),
	}

	avance, err := AgruparAvance(historias, []Sprint{})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	assertIDs := func(grupo GrupoAvance, esperados ...int64) {
		t.Helper()
		if len(grupo.Historias) != len(esperados) {
			t.Fatalf("grupo %q: se esperaban %d historias, se obtuvieron %d", grupo.Estado, len(esperados), len(grupo.Historias))
		}
		for i, id := range esperados {
			if grupo.Historias[i].ID != id {
				t.Errorf("grupo %q posición %d: se esperaba la historia %d, se obtuvo %d", grupo.Estado, i, id, grupo.Historias[i].ID)
			}
		}
	}
	if len(avance.Grupos) != 3 {
		t.Fatalf("se esperaban 3 grupos, se obtuvieron %d", len(avance.Grupos))
	}
	assertIDs(avance.Grupos[0], 2, 3)
	assertIDs(avance.Grupos[2], 7, 5)
}

func TestAgruparAvance_CopiaLosDatosDeCadaHistoria(t *testing.T) {
	historias := []BacklogItem{historiaDeAvance(9, "Iniciar sesión", EstadoEnProgreso, intPtr(5))}

	avance, err := AgruparAvance(historias, []Sprint{})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(avance.Grupos) != 3 || len(avance.Grupos[1].Historias) != 1 {
		t.Fatalf("se esperaba la historia en el grupo En progreso, se obtuvo %+v", avance.Grupos)
	}
	h := avance.Grupos[1].Historias[0]
	if h.ID != 9 || h.Titulo != "Iniciar sesión" || h.Estado != EstadoEnProgreso || h.EstimacionSP == nil || *h.EstimacionSP != 5 {
		t.Errorf("datos de la historia incorrectos: %+v", h)
	}
}

func TestAgruparAvance_EstadoDesconocidoEsError(t *testing.T) {
	historias := []BacklogItem{historiaDeAvance(1, "H1", Estado("Archivada"), intPtr(3))}

	if _, err := AgruparAvance(historias, []Sprint{}); err == nil {
		t.Fatal("se esperaba un error por el estado desconocido")
	}
}
