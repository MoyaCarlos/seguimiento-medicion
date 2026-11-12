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
	// FechaCierre es nil mientras el proyecto está abierto (#300).
	FechaCierre *time.Time
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

// Cerrar devuelve una copia del proyecto cerrado con el día civil de ahora
// como fecha de cierre (#300). Un proyecto ya cerrado, o con un Sprint Activo,
// no se puede cerrar.
func (p Project) Cerrar(ahora time.Time, sprints []Sprint) (Project, error) {
	if p.FechaCierre != nil {
		return Project{}, ErrProyectoYaCerrado
	}
	for _, s := range sprints {
		if s.Estado == SprintActivo {
			return Project{}, ErrSprintActivoAlCerrarProyecto
		}
	}
	hoy := inicioDelDiaCivil(ahora)
	p.FechaCierre = &hoy
	return p, nil
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
