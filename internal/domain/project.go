package domain

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// Project es el contenedor de todo lo demás: Sprints, historias, esfuerzo.
// Las fechas de inicio/fin (HU-04) son campos persistidos; el estado del
// proyecto (HU-13) NO es un campo: se calcula on-demand sobre esta struct.
type Project struct {
	ID          int64
	Nombre      string
	Descripcion string
	FechaInicio *time.Time
	FechaFin    *time.Time
	CreadoEn    time.Time
	// HorasPorStoryPoint es nil mientras no esté definido (HU-03).
	HorasPorStoryPoint *float64
}

const (
	projectNameMaxLen        = 100
	projectDescriptionMaxLen = 2000
)

// NewProject valida y construye un Project nuevo (sin ID ni CreadoEn, que los
// asigna la persistencia).
func NewProject(name, description string, start, end *time.Time) (Project, error) {
	nombre, descripcion, err := validarProyecto(name, description, start, end)
	if err != nil {
		return Project{}, err
	}
	return Project{Nombre: nombre, Descripcion: descripcion, FechaInicio: start, FechaFin: end}, nil
}

// ConDatosEditados reutiliza las mismas validaciones del alta y devuelve una
// copia del proyecto con los datos editables actualizados, conservando ID y
// CreadoEn.
func (p Project) ConDatosEditados(name, description string, start, end *time.Time) (Project, error) {
	nombre, descripcion, err := validarProyecto(name, description, start, end)
	if err != nil {
		return Project{}, err
	}
	p.Nombre = nombre
	p.Descripcion = descripcion
	p.FechaInicio = start
	p.FechaFin = end
	return p, nil
}

func validarProyecto(name, description string, start, end *time.Time) (string, string, error) {
	nombre := strings.TrimSpace(name)
	if nombre == "" {
		return "", "", ValidationError{Campo: "nombre", Mensaje: "el nombre es obligatorio"}
	}
	if len([]rune(nombre)) > projectNameMaxLen {
		return "", "", ValidationError{Campo: "nombre", Mensaje: "el nombre no puede superar los 100 caracteres"}
	}
	descripcion := strings.TrimSpace(description)
	if len([]rune(descripcion)) > projectDescriptionMaxLen {
		return "", "", ValidationError{Campo: "descripcion", Mensaje: "la descripción no puede superar los 2000 caracteres"}
	}
	if start != nil && end != nil && !end.After(*start) {
		return "", "", ValidationError{Campo: "fecha_fin", Mensaje: "la fecha de fin debe ser posterior a la de inicio"}
	}
	return nombre, descripcion, nil
}

// AdmiteSprint verifica que un Sprint de inicio a fin caiga dentro del plazo
// del proyecto, con los bordes incluidos. Solo se compara contra las fechas
// que el proyecto tenga (#168).
func (p Project) AdmiteSprint(inicio, fin time.Time) error {
	if p.FechaInicio != nil && inicio.Before(*p.FechaInicio) {
		return ValidationError{Campo: "fecha_inicio", Mensaje: "el Sprint no puede empezar antes del inicio del proyecto (" + p.FechaInicio.Format(time.DateOnly) + ")"}
	}
	if p.FechaFin != nil && fin.After(*p.FechaFin) {
		return ValidationError{Campo: "fecha_fin", Mensaje: "el Sprint no puede terminar después del fin del proyecto (" + p.FechaFin.Format(time.DateOnly) + ")"}
	}
	return nil
}

// ConHorasPorStoryPoint devuelve una copia del proyecto con el parámetro
// indicado. valor es el número en texto; nil lo deja sin definir.
func (p Project) ConHorasPorStoryPoint(valor *string) (Project, error) {
	if valor == nil {
		p.HorasPorStoryPoint = nil
		return p, nil
	}
	horas, err := strconv.ParseFloat(*valor, 64)
	if err != nil || !(horas > 0) || math.IsInf(horas, 1) {
		return Project{}, ErrHorasPorStoryPointInvalidas
	}
	p.HorasPorStoryPoint = &horas
	return p, nil
}
