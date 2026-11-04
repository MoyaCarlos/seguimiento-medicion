package domain

import "unicode/utf8"

const (
	largoMinNombreUsuario = 3
	largoMaxNombreUsuario = 30
	largoMinContrasena    = 8
	largoMaxContrasena    = 128
)

// ValidarNombreUsuario aplica las reglas del nombre de usuario (HU-20): de 3 a
// 30 caracteres, solo letras ASCII, números, '.', '-' y '_'. No se recorta.
func ValidarNombreUsuario(nombre string) error {
	largo := utf8.RuneCountInString(nombre)
	if largo < largoMinNombreUsuario || largo > largoMaxNombreUsuario {
		return ValidationError{Campo: "nombre", Mensaje: "el nombre de usuario debe tener entre 3 y 30 caracteres"}
	}
	for _, r := range nombre {
		if !caracterDeNombreUsuario(r) {
			return ValidationError{Campo: "nombre", Mensaje: "el nombre de usuario solo admite letras sin tildes, números, punto, guion y guion bajo"}
		}
	}
	return nil
}

func caracterDeNombreUsuario(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
		r == '.' || r == '-' || r == '_'
}

// ValidarContrasena aplica las reglas de la contraseña (HU-20): de 8 a 128
// caracteres (runas), cualquiera.
func ValidarContrasena(contrasena string) error {
	largo := utf8.RuneCountInString(contrasena)
	if largo < largoMinContrasena || largo > largoMaxContrasena {
		return ValidationError{Campo: "contrasena", Mensaje: "la contraseña debe tener entre 8 y 128 caracteres"}
	}
	return nil
}
