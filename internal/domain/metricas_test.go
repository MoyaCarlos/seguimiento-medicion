package domain

import "testing"

func horas(valor float64) *float64 { return &valor }

func esperarValor(t *testing.T, nombre string, obtenido, esperado *float64) {
	t.Helper()
	switch {
	case esperado == nil && obtenido != nil:
		t.Errorf("%s: se esperaba sin datos, se obtuvo %v", nombre, *obtenido)
	case esperado != nil && obtenido == nil:
		t.Errorf("%s: se esperaba %v, se obtuvo sin datos", nombre, *esperado)
	case esperado != nil && *obtenido != *esperado:
		t.Errorf("%s: se esperaba %v, se obtuvo %v", nombre, *esperado, *obtenido)
	}
}

func esperarMetricas(t *testing.T, obtenidas, esperadas Metricas) {
	t.Helper()
	if obtenidas.SPPlanificados != esperadas.SPPlanificados {
		t.Errorf("SP planificados: se esperaba %d, se obtuvo %d", esperadas.SPPlanificados, obtenidas.SPPlanificados)
	}
	if obtenidas.SPCompletados != esperadas.SPCompletados {
		t.Errorf("SP completados: se esperaba %d, se obtuvo %d", esperadas.SPCompletados, obtenidas.SPCompletados)
	}
	if obtenidas.Velocidad != esperadas.Velocidad {
		t.Errorf("velocidad: se esperaba %v, se obtuvo %v", esperadas.Velocidad, obtenidas.Velocidad)
	}
	if obtenidas.PorcentajeCompletadas != esperadas.PorcentajeCompletadas {
		t.Errorf("porcentaje: se esperaba %v, se obtuvo %v", esperadas.PorcentajeCompletadas, obtenidas.PorcentajeCompletadas)
	}
	if obtenidas.HorasReales != esperadas.HorasReales {
		t.Errorf("horas reales: se esperaba %v, se obtuvo %v", esperadas.HorasReales, obtenidas.HorasReales)
	}
	if obtenidas.DefectosDetectados != esperadas.DefectosDetectados || obtenidas.DefectosResueltos != esperadas.DefectosResueltos {
		t.Errorf("defectos: se esperaba %d y %d, se obtuvo %d y %d", esperadas.DefectosDetectados, esperadas.DefectosResueltos, obtenidas.DefectosDetectados, obtenidas.DefectosResueltos)
	}
	esperarValor(t, "horas estimadas", obtenidas.HorasEstimadas, esperadas.HorasEstimadas)
	esperarValor(t, "desviación en horas", obtenidas.DesviacionHoras, esperadas.DesviacionHoras)
	esperarValor(t, "desviación porcentual", obtenidas.DesviacionPorcentaje, esperadas.DesviacionPorcentaje)
}

func TestCalcularMetricas_DeUnSprint(t *testing.T) {
	casos := []struct {
		nombre     string
		datos      DatosMetricasSprint
		horasPorSP *float64
		esperadas  Metricas
	}{
		{
			nombre: "completadas, arrastradas, horas y defectos",
			datos: DatosMetricasSprint{
				HistoriasCompletadas: 1, HistoriasArrastradas: 1, SPCompletados: 5, SPArrastrados: 3,
				MinutosReales: 30 * 60, DefectosDetectados: 2, DefectosResueltos: 1,
			},
			horasPorSP: horas(4),
			esperadas: Metricas{
				SPPlanificados: 8, SPCompletados: 5, Velocidad: 5, PorcentajeCompletadas: 50,
				HorasEstimadas: horas(32), HorasReales: 30, DesviacionHoras: horas(-2), DesviacionPorcentaje: horas(-6.25),
				DefectosDetectados: 2, DefectosResueltos: 1,
			},
		},
		{
			nombre:     "las historias sin Story Points cuentan en el porcentaje y suman cero",
			datos:      DatosMetricasSprint{HistoriasCompletadas: 1, HistoriasArrastradas: 3, SPCompletados: 5},
			horasPorSP: horas(4),
			esperadas: Metricas{
				SPPlanificados: 5, SPCompletados: 5, Velocidad: 5, PorcentajeCompletadas: 25,
				HorasEstimadas: horas(20), DesviacionHoras: horas(-20), DesviacionPorcentaje: horas(-100),
			},
		},
		{
			nombre:     "un Sprint sin historias comprometidas",
			horasPorSP: horas(4),
			esperadas:  Metricas{HorasEstimadas: horas(0), DesviacionHoras: horas(0)},
		},
		{
			nombre: "sin horas por Story Point no hay estimadas ni desviación",
			datos:  DatosMetricasSprint{HistoriasCompletadas: 1, SPCompletados: 5, MinutosReales: 90},
			esperadas: Metricas{
				SPPlanificados: 5, SPCompletados: 5, Velocidad: 5, PorcentajeCompletadas: 100, HorasReales: 1.5,
			},
		},
		{
			nombre:     "con cero Story Points planificados la desviación son las horas reales",
			datos:      DatosMetricasSprint{HistoriasArrastradas: 1, MinutosReales: 120},
			horasPorSP: horas(4),
			esperadas:  Metricas{HorasEstimadas: horas(0), HorasReales: 2, DesviacionHoras: horas(2)},
		},
		{
			nombre:     "sin horas imputadas las reales son cero",
			datos:      DatosMetricasSprint{HistoriasCompletadas: 1, SPCompletados: 2},
			horasPorSP: horas(2.5),
			esperadas: Metricas{
				SPPlanificados: 2, SPCompletados: 2, Velocidad: 2, PorcentajeCompletadas: 100,
				HorasEstimadas: horas(5), DesviacionHoras: horas(-5), DesviacionPorcentaje: horas(-100),
			},
		},
		{
			nombre:     "los decimales se redondean a dos cifras",
			datos:      DatosMetricasSprint{HistoriasCompletadas: 1, HistoriasArrastradas: 2, SPCompletados: 1, MinutosReales: 100},
			horasPorSP: horas(1),
			esperadas: Metricas{
				SPPlanificados: 1, SPCompletados: 1, Velocidad: 1, PorcentajeCompletadas: 33.33,
				HorasEstimadas: horas(1), HorasReales: 1.67, DesviacionHoras: horas(0.67), DesviacionPorcentaje: horas(66.67),
			},
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			c.datos.SprintID, c.datos.Numero = 7, 2

			resultado := CalcularMetricas([]DatosMetricasSprint{c.datos}, c.horasPorSP)

			if len(resultado.Sprints) != 1 {
				t.Fatalf("se esperaba un Sprint, se obtuvieron %d", len(resultado.Sprints))
			}
			sprint := resultado.Sprints[0]
			if sprint.SprintID != 7 || sprint.Numero != 2 {
				t.Errorf("se esperaba el Sprint 7 con número 2, se obtuvo %d con número %d", sprint.SprintID, sprint.Numero)
			}
			esperarMetricas(t, sprint.Metricas, c.esperadas)
		})
	}
}
