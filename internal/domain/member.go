package domain

import (
	"strings"
	"time"
)

// Role es la categoría de responsabilidad de un integrante dentro de un proyecto.
type Role string

const (
	RolScrumMaster    Role = "scrum_master"
	RolProductBuilder Role = "product_builder"
)

// Valido indica si el rol es uno de los aceptados por la historia.
func (r Role) Valido() bool {
	return r == RolScrumMaster || r == RolProductBuilder
}

// User es un integrante del equipo. Se identifica por su nombre normalizado,
// de modo que "Ana", "ana" y " Ana " resuelven al mismo usuario.
type User struct {
	ID                int64
	Nombre            string
	NombreNormalizado string
	CreadoEn          time.Time
}

const userNameMaxLen = 200

// NewUser valida y construye un integrante nuevo.
func NewUser(name string) (User, error) {
	nombre := strings.TrimSpace(name)
	if nombre == "" {
		return User{}, ValidationError{Campo: "nombre", Mensaje: "el nombre del integrante es obligatorio"}
	}
	if len([]rune(nombre)) > userNameMaxLen {
		return User{}, ValidationError{Campo: "nombre", Mensaje: "el nombre del integrante no puede superar los 200 caracteres"}
	}
	return User{Nombre: nombre, NombreNormalizado: NormalizarNombre(nombre)}, nil
}

// NormalizarNombre aplica la regla de identidad: minúsculas y sin espacios externos.
func NormalizarNombre(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// Membership es la relación Proyecto–Usuario–Rol: un integrante con exactamente
// un rol dentro de un proyecto.
type Membership struct {
	ProjectID int64
	UserID    int64
	Role      Role
	CreadoEn  time.Time
}

// NewMembership valida y construye una asignación integrante–proyecto–rol.
func NewMembership(projectID, userID int64, role Role) (Membership, error) {
	if projectID <= 0 {
		return Membership{}, ValidationError{Campo: "proyecto_id", Mensaje: "el proyecto es obligatorio"}
	}
	if userID <= 0 {
		return Membership{}, ValidationError{Campo: "integrante_id", Mensaje: "el integrante es obligatorio"}
	}
	if !role.Valido() {
		return Membership{}, ValidationError{Campo: "rol", Mensaje: "el rol debe ser scrum_master o product_builder"}
	}
	return Membership{ProjectID: projectID, UserID: userID, Role: role}, nil
}

// Member es la vista de lectura de un integrante dentro de un proyecto,
// combinando el nombre del usuario con su rol.
type Member struct {
	UserID    int64
	ProjectID int64
	Nombre    string
	Role      Role
	CreadoEn  time.Time
}
