package report

import (
	"errors"
	"strconv"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

var (
	colorCabecera = &props.Color{Red: 232, Green: 232, Blue: 232}
	textoCabecera = props.Text{Size: 8.5, Style: fontstyle.Bold, Top: 1.4, Bottom: 1.4, Left: 1.2, Right: 1.2}
	textoCelda    = props.Text{Size: 8, Top: 1.2, Bottom: 1.2, Left: 1.2, Right: 1.2}
)

type columna struct {
	titulo  string
	ancho   int
	alinear align.Type
}

// Generador implementa la generación de los dos reportes como un tipo que
// puede inyectarse en el handler HTTP.
type Generador struct{}

func (Generador) GenerarSprint(reporte domain.ReporteSprint) ([]byte, error) {
	return GenerarSprint(reporte)
}

func (Generador) GenerarProyecto(reporte domain.ReporteProyecto) ([]byte, error) {
	return GenerarProyecto(reporte)
}

// GenerarSprint devuelve el PDF del reporte de un Sprint finalizado.
func GenerarSprint(reporte domain.ReporteSprint) ([]byte, error) {
	m := nuevaMaroto("Reporte de Sprint")
	agregarEncabezadoSprint(m, reporte)
	agregarSprint(m, reporte.Sprint, true)
	return generar(m)
}

// GenerarProyecto devuelve el PDF del reporte completo de un proyecto.
func GenerarProyecto(reporte domain.ReporteProyecto) ([]byte, error) {
	m := nuevaMaroto("Reporte de Proyecto")
	agregarSeccion(m, "Datos del proyecto")
	agregarDato(m, "Nombre", reporte.Proyecto.Nombre)
	if reporte.Proyecto.Descripcion != "" {
		agregarDatoAuto(m, "Descripción", reporte.Proyecto.Descripcion)
	}
	agregarDato(m, "Fecha de inicio", fechaDe(reporte.Proyecto.FechaInicio))
	agregarDato(m, "Fecha de fin", fechaDe(reporte.Proyecto.FechaFin))
	agregarDato(m, "Estado", string(reporte.Estado))

	agregarSeccion(m, "Equipo")
	agregarIntegrantes(m, reporte.Integrantes)

	agregarSeccion(m, "Product Backlog")
	agregarTablaHistorias(m, reporte.ProductBacklog, false)

	for _, sprint := range reporte.Sprints {
		agregarSalto(m)
		agregarSprint(m, sprint, false)
	}
	if len(reporte.Sprints) == 0 {
		agregarSeccion(m, "Sprints finalizados")
		agregarVacio(m, "Sin datos: el proyecto no tiene Sprints finalizados")
	}

	if reporte.SprintActivo == nil {
		agregarSeccion(m, "Sprint Activo")
		agregarVacio(m, "Sin Sprint activo")
	} else {
		agregarSeccion(m, "Sprint Activo")
		agregarDato(m, "Sprint", "Sprint "+strconv.Itoa(reporte.SprintActivo.Sprint.Numero))
		agregarDato(m, "Objetivo", reporte.SprintActivo.Sprint.SprintGoal)
		agregarDato(m, "Fechas", fechasDelSprint(reporte.SprintActivo.Sprint))
		agregarTablaHistorias(m, reporte.SprintActivo.Historias, false)
	}

	agregarSalto(m)
	agregarSeccion(m, "Totales del proyecto")
	if reporte.SinSprintsFinalizados {
		agregarVacio(m, "Sin datos: el proyecto no tiene Sprints finalizados")
	} else {
		agregarMetricas(m, reporte.Metricas)
	}
	return generar(m)
}

func agregarEncabezadoSprint(m core.Maroto, reporte domain.ReporteSprint) {
	m.AddRow(12, text.NewCol(12, "Reporte de Sprint", props.Text{
		Size: 17, Style: fontstyle.Bold, Align: align.Center,
	}))
	m.AddRow(7, text.NewCol(12, reporte.Proyecto.Nombre, props.Text{
		Size: 11, Align: align.Center,
	}))
}

func agregarSprint(m core.Maroto, reporte domain.SprintReporte, encabezado bool) {
	if !encabezado {
		agregarSeccion(m, "Sprint "+strconv.Itoa(reporte.Sprint.Numero)+" — "+reporte.Sprint.SprintGoal)
	} else {
		agregarSeccion(m, "Sprint "+strconv.Itoa(reporte.Sprint.Numero))
	}
	agregarDato(m, "Objetivo", reporte.Sprint.SprintGoal)
	agregarDato(m, "Fechas", fechasDelSprint(reporte.Sprint))
	agregarDato(m, "Estado", string(reporte.Sprint.Estado))

	agregarSubtitulo(m, "Historias planificadas")
	agregarTablaHistorias(m, reporte.Planificadas, true)

	agregarSubtitulo(m, "Métricas")
	agregarMetricas(m, reporte.Metricas)

	agregarSubtitulo(m, "Defectos detectados")
	agregarTablaDefectos(m, reporte.DefectosDetectados, "Detección", false)
	agregarSubtitulo(m, "Defectos resueltos")
	agregarTablaDefectos(m, reporte.DefectosResueltos, "Resolución", true)
}

func nuevaMaroto(titulo string) core.Maroto {
	cfg := config.NewBuilder().
		WithPageNumber(props.PageNumber{Pattern: "Página {current} de {total}", Place: props.RightBottom}).
		WithTitle(titulo, true).
		WithAuthor("Software Metrics & Estimation", true).
		Build()
	return maroto.New(cfg)
}

func generar(m core.Maroto) ([]byte, error) {
	documento, err := m.Generate()
	if err != nil {
		return nil, err
	}
	contenido := documento.GetBytes()
	if len(contenido) == 0 {
		return nil, errors.New("el PDF quedó vacío")
	}
	return contenido, nil
}

func agregarTitulo(m core.Maroto, texto string) {
	m.AddRow(12, text.NewCol(12, texto, props.Text{Size: 17, Style: fontstyle.Bold, Align: align.Center}))
}

func agregarSeccion(m core.Maroto, texto string) {
	m.AddRow(9, text.NewCol(12, texto, props.Text{Size: 12, Style: fontstyle.Bold, Top: 2}))
}

func agregarSubtitulo(m core.Maroto, texto string) {
	m.AddRow(7, text.NewCol(12, texto, props.Text{Size: 9.5, Style: fontstyle.Bold, Top: 1.5}))
}

func agregarSalto(m core.Maroto) {
	m.AddRow(4)
}

func agregarDato(m core.Maroto, etiqueta, valor string) {
	m.AddRow(5.5,
		text.NewCol(4, etiqueta, props.Text{Size: 8.5, Style: fontstyle.Bold}),
		text.NewCol(8, valor, props.Text{Size: 8.5}),
	)
}

func agregarDatoAuto(m core.Maroto, etiqueta, valor string) {
	m.AddAutoRow(
		text.NewCol(4, etiqueta, props.Text{Size: 8.5, Style: fontstyle.Bold, Top: 1.2}),
		text.NewCol(8, valor, props.Text{Size: 8.5, Top: 1.2, Bottom: 1.2}),
	)
}

func agregarVacio(m core.Maroto, texto string) {
	m.AddRow(6, text.NewCol(12, texto, props.Text{Size: 8.5, Style: fontstyle.Italic}))
}

func agregarIntegrantes(m core.Maroto, integrantes []domain.Member) {
	if len(integrantes) == 0 {
		agregarVacio(m, "Sin integrantes")
		return
	}
	filas := make([][]string, 0, len(integrantes))
	for _, integrante := range integrantes {
		filas = append(filas, []string{integrante.Nombre, nombreDelRol(integrante.Role)})
	}
	agregarTabla(m, []columna{{"Nombre", 8, align.Left}, {"Rol", 4, align.Left}}, filas)
}

func agregarTablaHistorias(m core.Maroto, historias []domain.HistoriaReporte, conCompletada bool) {
	if len(historias) == 0 {
		agregarVacio(m, "Sin historias")
		return
	}
	columnas := []columna{
		{"Título", 5, align.Left},
		{"Prioridad", 1, align.Center},
		{"SP", 1, align.Center},
		{"Estado", 2, align.Left},
	}
	if conCompletada {
		columnas[0].ancho = 4
		columnas = append(columnas, columna{"Completada", 2, align.Center})
	}
	filas := make([][]string, 0, len(historias))
	for _, historia := range historias {
		fila := []string{historia.Titulo, string(historia.Prioridad), estimacionDe(historia.EstimacionSP), string(historia.Estado)}
		if conCompletada {
			fila = append(fila, siNo(historia.Completada))
		}
		filas = append(filas, fila)
	}
	agregarTabla(m, columnas, filas)
}

func agregarTablaDefectos(m core.Maroto, defectos []domain.DefectoReporte, columnaSprint string, resolucion bool) {
	if len(defectos) == 0 {
		agregarVacio(m, "Sin defectos")
		return
	}
	columnas := []columna{
		{"Descripción", 5, align.Left},
		{"Severidad", 1, align.Center},
		{"Estado", 2, align.Left},
		{"Historia", 3, align.Left},
		{columnaSprint, 1, align.Center},
	}
	filas := make([][]string, 0, len(defectos))
	for _, defecto := range defectos {
		numero := strconv.Itoa(defecto.SprintDeteccion)
		if resolucion && defecto.SprintResolucion != nil {
			numero = strconv.Itoa(*defecto.SprintResolucion)
		}
		filas = append(filas, []string{
			defecto.Descripcion, string(defecto.Severidad), string(defecto.Estado), defecto.Historia, numero,
		})
	}
	agregarTabla(m, columnas, filas)
}

func agregarTabla(m core.Maroto, columnas []columna, filas [][]string) {
	cabecera := row.New(6.5)
	for _, columna := range columnas {
		texto := textoCabecera
		texto.Align = columna.alinear
		cabecera.Add(text.NewCol(columna.ancho, columna.titulo, texto))
	}
	m.AddRows(cabecera.WithStyle(&props.Cell{BackgroundColor: colorCabecera}))
	for _, fila := range filas {
		columnasDeLaFila := make([]core.Col, 0, len(columnas))
		for i, valor := range fila {
			texto := textoCelda
			texto.Align = columnas[i].alinear
			columnasDeLaFila = append(columnasDeLaFila, text.NewCol(columnas[i].ancho, valor, texto))
		}
		m.AddAutoRow(columnasDeLaFila...)
	}
}

func agregarMetricas(m core.Maroto, metricas domain.Metricas) {
	agregarDato(m, "Story Points planificados", strconv.Itoa(metricas.SPPlanificados))
	agregarDato(m, "Story Points completados", strconv.Itoa(metricas.SPCompletados))
	agregarDato(m, "Velocidad", decimalDe(metricas.Velocidad))
	agregarDato(m, "Historias completadas", decimalDe(metricas.PorcentajeCompletadas)+"%")
	agregarDato(m, "Horas estimadas", horasDe(metricas.HorasEstimadas))
	agregarDato(m, "Horas reales", decimalDe(metricas.HorasReales)+" h")
	agregarDato(m, "Desviación en horas", horasDe(metricas.DesviacionHoras))
	agregarDato(m, "Desviación porcentual", porcentajeDe(metricas.DesviacionPorcentaje))
	agregarDato(m, "Defectos detectados", strconv.Itoa(metricas.DefectosDetectados))
	agregarDato(m, "Defectos resueltos", strconv.Itoa(metricas.DefectosResueltos))
}

func fechasDelSprint(sprint domain.Sprint) string {
	if sprint.FechaInicio.IsZero() && sprint.FechaFin.IsZero() {
		return "sin definir"
	}
	return fechaDe(&sprint.FechaInicio) + " a " + fechaDe(&sprint.FechaFin)
}

func fechaDe(fecha *time.Time) string {
	if fecha == nil || fecha.IsZero() {
		return "sin definir"
	}
	return fecha.Format(time.DateOnly)
}

func estimacionDe(sp *int) string {
	if sp == nil {
		return "sin estimar"
	}
	return strconv.Itoa(*sp)
}

func horasDe(valor *float64) string {
	if valor == nil {
		return "sin datos"
	}
	return decimalDe(*valor) + " h"
}

func porcentajeDe(valor *float64) string {
	if valor == nil {
		return "sin datos"
	}
	return decimalDe(*valor) + "%"
}

func decimalDe(valor float64) string {
	return strconv.FormatFloat(valor, 'f', -1, 64)
}

func siNo(valor bool) string {
	if valor {
		return "Sí"
	}
	return "No"
}

func nombreDelRol(rol domain.Role) string {
	if rol == domain.RolScrumMaster {
		return "Scrum Master"
	}
	return "Product Builder"
}
