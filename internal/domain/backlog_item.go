package domain

import (
	"strings"
	"time"
)

const (
	longitudMaximaTitulo      = 200
	longitudMaximaDescripcion = 2000
	longitudMaximaCriterios   = 2000
)

// BacklogItem representa una historia de usuario del Product Backlog.
type BacklogItem struct {
	ID                  int64
	ProyectoID          int64
	Titulo              string
	Descripcion         string
	CriteriosAceptacion *string
	Prioridad           Prioridad
	Estado              Estado
	ValorNegocio        *int
	EstimacionSP        *int
	// SprintID es nil mientras la historia está en el Product Backlog. Su
	// asignación (con sus validaciones) es alcance de HU-06.
	SprintID *int64
	// FechaFinalizacion es nil salvo que la historia esté Completada.
	FechaFinalizacion *time.Time
}

// NewBacklogItem construye una historia validando sus invariantes. El estado se
// fuerza a "Nueva" y la estimación queda sin asignar hasta Planning Poker (HU-02).
func NewBacklogItem(proyectoID int64, titulo, descripcion string, prioridad Prioridad, valorNegocio *int) (BacklogItem, error) {
	if proyectoID <= 0 {
		return BacklogItem{}, ValidationError{Campo: "proyecto_id", Mensaje: "el proyecto es obligatorio"}
	}

	titulo, err := validarTitulo(titulo)
	if err != nil {
		return BacklogItem{}, err
	}
	descripcion, err = validarDescripcion(descripcion)
	if err != nil {
		return BacklogItem{}, err
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

// ConDatosEditados devuelve una copia de la historia con los campos editables
// reemplazados, validados con las mismas reglas que el alta. El resto de los
// campos se conserva; ante un error la historia no cambia.
func (b BacklogItem) ConDatosEditados(titulo, descripcion string, criterios *string, estimacionSP *int) (BacklogItem, error) {
	titulo, err := validarTitulo(titulo)
	if err != nil {
		return BacklogItem{}, err
	}
	descripcion, err = validarDescripcion(descripcion)
	if err != nil {
		return BacklogItem{}, err
	}

	var criteriosNormalizados *string
	if criterios != nil {
		if texto := strings.TrimSpace(*criterios); texto != "" {
			if len([]rune(texto)) > longitudMaximaCriterios {
				return BacklogItem{}, ValidationError{Campo: "criterios_aceptacion", Mensaje: "los criterios de aceptación no pueden superar los 2000 caracteres"}
			}
			criteriosNormalizados = &texto
		}
	}

	if estimacionSP != nil && !EsEstimacionValida(*estimacionSP) {
		return BacklogItem{}, ErrEstimacionInvalida
	}

	b.Titulo, b.Descripcion = titulo, descripcion
	b.CriteriosAceptacion, b.EstimacionSP = criteriosNormalizados, estimacionSP
	return b, nil
}

func validarTitulo(titulo string) (string, error) {
	titulo = strings.TrimSpace(titulo)
	if titulo == "" {
		return "", ValidationError{Campo: "titulo", Mensaje: "el título es obligatorio"}
	}
	if len([]rune(titulo)) > longitudMaximaTitulo {
		return "", ValidationError{Campo: "titulo", Mensaje: "el título no puede superar los 200 caracteres"}
	}
	return titulo, nil
}

func validarDescripcion(descripcion string) (string, error) {
	descripcion = strings.TrimSpace(descripcion)
	if descripcion == "" {
		return "", ValidationError{Campo: "descripcion", Mensaje: "la descripción es obligatoria"}
	}
	if len([]rune(descripcion)) > longitudMaximaDescripcion {
		return "", ValidationError{Campo: "descripcion", Mensaje: "la descripción no puede superar los 2000 caracteres"}
	}
	return descripcion, nil
}
