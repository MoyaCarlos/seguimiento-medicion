package domain

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
	return MetricasProyecto{}
}
