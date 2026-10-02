package domain

import (
	"strings"
	"time"
)

// Project es el contenedor de todo lo demás: Sprints, historias, esfuerzo.
// Los campos de fechas/estado (HU-13) se agregan sobre esta misma struct,
// no se cambia el nombre ni se remueven los que ya hay.
type Project struct {
	ID          string
	Name        string
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
}

const (
	projectNameMaxLen        = 100
	projectDescriptionMaxLen = 2000
)

// NewProject valida y construye un Project nuevo (sin ID ni CreatedAt, que los
// asigna la persistencia).
func NewProject(name, description string, start, end *time.Time) (Project, error) {
	nombre, descripcion, err := validarProyecto(name, description, start, end)
	if err != nil {
		return Project{}, err
	}
	return Project{Name: nombre, Description: descripcion, StartDate: start, EndDate: end}, nil
}

// ConDatosEditados reutiliza las mismas validaciones del alta y devuelve una
// copia del proyecto con los datos editables actualizados, conservando ID y
// CreatedAt.
func (p Project) ConDatosEditados(name, description string, start, end *time.Time) (Project, error) {
	nombre, descripcion, err := validarProyecto(name, description, start, end)
	if err != nil {
		return Project{}, err
	}
	p.Name = nombre
	p.Description = descripcion
	p.StartDate = start
	p.EndDate = end
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
