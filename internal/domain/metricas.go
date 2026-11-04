package domain

import "math"

// DatosMetricasSprint son los datos ya sumados de un Sprint finalizado.
type DatosMetricasSprint struct {
	SprintID             int64
	Numero               int
	HistoriasCompletadas int
	HistoriasArrastradas int
	SPCompletados        int
	SPArrastrados        int
	MinutosReales        int
	DefectosDetectados   int
	DefectosResueltos    int
}

// Metricas son los valores calculados de un Sprint o del proyecto; no se
// guardan. Un puntero nil significa "sin datos".
type Metricas struct {
	SPPlanificados        int
	SPCompletados         int
	Velocidad             float64
	PorcentajeCompletadas float64
	HorasEstimadas        *float64
	HorasReales           float64
	DesviacionHoras       *float64
	DesviacionPorcentaje  *float64
	DefectosDetectados    int
	DefectosResueltos     int
}

// MetricasDeSprint son las métricas de un Sprint finalizado.
type MetricasDeSprint struct {
	SprintID int64
	Numero   int
	Metricas Metricas
}

// MetricasProyecto son las métricas del proyecto y el detalle por Sprint.
type MetricasProyecto struct {
	Proyecto Metricas
	Sprints  []MetricasDeSprint
}

// CalcularMetricas calcula las métricas de cada Sprint finalizado y del proyecto.
func CalcularMetricas(datos []DatosMetricasSprint, horasPorSP *float64) MetricasProyecto {
	sprints := make([]MetricasDeSprint, 0, len(datos))
	var total DatosMetricasSprint
	for _, d := range datos {
		sprints = append(sprints, MetricasDeSprint{SprintID: d.SprintID, Numero: d.Numero, Metricas: calcularDe(d, horasPorSP)})
		total = total.sumar(d)
	}

	proyecto := calcularDe(total, horasPorSP)
	proyecto.Velocidad = redondear(float64(total.SPCompletados) / float64(len(datos)))
	return MetricasProyecto{Proyecto: proyecto, Sprints: sprints}
}

func (d DatosMetricasSprint) sumar(otro DatosMetricasSprint) DatosMetricasSprint {
	d.HistoriasCompletadas += otro.HistoriasCompletadas
	d.HistoriasArrastradas += otro.HistoriasArrastradas
	d.SPCompletados += otro.SPCompletados
	d.SPArrastrados += otro.SPArrastrados
	d.MinutosReales += otro.MinutosReales
	d.DefectosDetectados += otro.DefectosDetectados
	d.DefectosResueltos += otro.DefectosResueltos
	return d
}

// calcularDe suma en enteros y redondea solo el valor final, para no acumular
// errores de redondeo (D8).
func calcularDe(d DatosMetricasSprint, horasPorSP *float64) Metricas {
	planificados := d.SPCompletados + d.SPArrastrados
	m := Metricas{
		SPPlanificados:        planificados,
		SPCompletados:         d.SPCompletados,
		Velocidad:             float64(d.SPCompletados),
		PorcentajeCompletadas: porcentajeCompletadas(d),
		HorasReales:           redondear(minutosAHoras(d.MinutosReales)),
		DefectosDetectados:    d.DefectosDetectados,
		DefectosResueltos:     d.DefectosResueltos,
	}
	if horasPorSP == nil {
		return m
	}

	estimadas := float64(planificados) * *horasPorSP
	desviacion := minutosAHoras(d.MinutosReales) - estimadas
	m.HorasEstimadas = valorRedondeado(estimadas)
	m.DesviacionHoras = valorRedondeado(desviacion)
	if estimadas > 0 {
		m.DesviacionPorcentaje = valorRedondeado(desviacion / estimadas * 100)
	}
	return m
}

func porcentajeCompletadas(d DatosMetricasSprint) float64 {
	comprometidas := d.HistoriasCompletadas + d.HistoriasArrastradas
	if comprometidas == 0 {
		return 0
	}
	return redondear(float64(d.HistoriasCompletadas) / float64(comprometidas) * 100)
}

func minutosAHoras(minutos int) float64 {
	return float64(minutos) / 60
}

func valorRedondeado(valor float64) *float64 {
	redondeado := redondear(valor)
	return &redondeado
}

// redondear devuelve 0 en lugar de -0, porque el JSON escribiría el signo.
func redondear(valor float64) float64 {
	redondeado := math.Round(valor*100) / 100
	if redondeado == 0 {
		return 0
	}
	return redondeado
}
