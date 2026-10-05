package domain

import (
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
	if start != nil && end != nil && end.Before(*start) {
		return "", "", ValidationError{Campo: "fecha_fin", Mensaje: "la fecha de fin no puede ser anterior a la fecha de inicio"}
	}
	return nombre, descripcion, nil
}
