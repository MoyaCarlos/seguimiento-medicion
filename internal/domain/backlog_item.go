package domain

import "strings"

const (
	longitudMaximaTitulo      = 200
	longitudMaximaDescripcion = 2000
)

// BacklogItem representa una historia de usuario del Product Backlog.
type BacklogItem struct {
	ID           int64
	ProyectoID   int64
	Titulo       string
	Descripcion  string
	Prioridad    Prioridad
	Estado       Estado
	ValorNegocio *int
	EstimacionSP *int
}

// NewBacklogItem construye una historia validando sus invariantes. El estado se
// fuerza a "Nueva" y la estimación queda sin asignar hasta Planning Poker (HU-02).
func NewBacklogItem(proyectoID int64, titulo, descripcion string, prioridad Prioridad, valorNegocio *int) (BacklogItem, error) {
	if proyectoID <= 0 {
		return BacklogItem{}, ValidationError{Campo: "proyecto_id", Mensaje: "el proyecto es obligatorio"}
	}

	titulo = strings.TrimSpace(titulo)
	if titulo == "" {
		return BacklogItem{}, ValidationError{Campo: "titulo", Mensaje: "el título es obligatorio"}
	}
	if len([]rune(titulo)) > longitudMaximaTitulo {
		return BacklogItem{}, ValidationError{Campo: "titulo", Mensaje: "el título no puede superar los 200 caracteres"}
	}

	descripcion = strings.TrimSpace(descripcion)
	if descripcion == "" {
		return BacklogItem{}, ValidationError{Campo: "descripcion", Mensaje: "la descripción es obligatoria"}
	}
	if len([]rune(descripcion)) > longitudMaximaDescripcion {
		return BacklogItem{}, ValidationError{Campo: "descripcion", Mensaje: "la descripción no puede superar los 2000 caracteres"}
	}

	if !prioridad.EsValida() {
		return BacklogItem{}, ValidationError{Campo: "prioridad", Mensaje: "la prioridad debe ser M, S, C o W"}
	}

	if valorNegocio != nil && !EsValorNegocioValido(*valorNegocio) {
		return BacklogItem{}, ValidationError{Campo: "valor_negocio", Mensaje: "el valor de negocio debe ser uno de 1, 2, 3, 5, 8, 13, 21"}
	}

	return BacklogItem{
		ProyectoID:   proyectoID,
		Titulo:       titulo,
		Descripcion:  descripcion,
		Prioridad:    prioridad,
		Estado:       EstadoNueva,
		ValorNegocio: valorNegocio,
		EstimacionSP: nil,
	}, nil
}
