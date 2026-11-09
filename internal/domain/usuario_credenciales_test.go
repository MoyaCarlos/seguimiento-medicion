package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestValidarNombreUsuario(t *testing.T) {
	casos := []struct {
		nombre string
		valor  string
		valido bool
	}{
		{"tres caracteres", "ana", true},
		{"signos permitidos", "a.b-c_9", true},
		{"30 caracteres", strings.Repeat("a", 30), true},
		{"2 caracteres", "an", false},
		{"31 caracteres", strings.Repeat("a", 31), false},
		{"vacío", "", false},
		{"espacio interno", "ana maria", false},
		{"espacio inicial", " ana", false},
		{"espacio final", "ana ", false},
		{"arroba", "ana@x", false},
		{"tilde", "josé", false},
		{"emoji", "ana😀", false},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			err := ValidarNombreUsuario(c.valor)
			if c.valido {
				if err != nil {
					t.Fatalf("no se esperaba error: %v", err)
				}
				return
			}
			var ve ValidationError
			if !errors.As(err, &ve) || ve.Campo != "nombre" {
				t.Fatalf("se esperaba ValidationError del campo nombre, se obtuvo %v", err)
			}
		})
	}
}

func TestValidarContrasena(t *testing.T) {
	casos := []struct {
		nombre string
		valor  string
		valido bool
	}{
		{"8 caracteres", "abcdefgh", true},
		{"8 runas con tilde y emoji", "ñandú😀ok", true},
		{"con espacios", "a b c d e", true},
		{"128 runas", strings.Repeat("é", 128), true},
		{"vacía", "", false},
		{"7 caracteres", "abcdefg", false},
		{"7 runas que superan 8 bytes", "ñññññññ", false},
		{"129 caracteres", strings.Repeat("a", 129), false},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			err := ValidarContrasena(c.valor)
			if c.valido {
				if err != nil {
					t.Fatalf("no se esperaba error: %v", err)
				}
				return
			}
			var ve ValidationError
			if !errors.As(err, &ve) || ve.Campo != "contrasena" {
				t.Fatalf("se esperaba ValidationError del campo contrasena, se obtuvo %v", err)
			}
		})
	}
}
