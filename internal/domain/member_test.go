package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestRole_Valido(t *testing.T) {
	validos := []Role{RolScrumMaster, RolProductBuilder}
	for _, r := range validos {
		if !r.Valido() {
			t.Errorf("se esperaba que el rol %q fuera válido", r)
		}
	}
	if Role("Product Owner").Valido() {
		t.Error("se esperaba que un rol fuera de Scrum Master / Product Builder fuera inválido")
	}
}

func TestNewUser_NormalizaNombre(t *testing.T) {
	u, err := NewUser("  Ana Valentina  ")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if u.Name != "Ana Valentina" {
		t.Errorf("se esperaba el nombre recortado, se obtuvo %q", u.Name)
	}
	if u.NormalizedName != "ana valentina" {
		t.Errorf("se esperaba nombre normalizado en minúsculas, se obtuvo %q", u.NormalizedName)
	}
}

func TestNewUser_NombreVacio(t *testing.T) {
	for _, nombre := range []string{"", "   "} {
		_, err := NewUser(nombre)
		if err == nil {
			t.Fatalf("se esperaba error para el nombre %q", nombre)
		}
		var verr ValidationError
		if !errors.As(err, &verr) {
			t.Fatalf("se esperaba ValidationError, se obtuvo %T", err)
		}
		if verr.Campo != "nombre" {
			t.Errorf("se esperaba campo nombre, se obtuvo %q", verr.Campo)
		}
	}
}

func TestNewUser_NombreDemasiadoLargo(t *testing.T) {
	_, err := NewUser(strings.Repeat("a", userNameMaxLen+1))
	if err == nil {
		t.Fatal("se esperaba error por longitud")
	}
	var verr ValidationError
	if !errors.As(err, &verr) || verr.Campo != "nombre" {
		t.Fatalf("se esperaba ValidationError de nombre, se obtuvo %v", err)
	}
}

func TestNewMembership_Valida(t *testing.T) {
	m, err := NewMembership("proy-1", "user-1", RolScrumMaster)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if m.ProjectID != "proy-1" || m.UserID != "user-1" || m.Role != RolScrumMaster {
		t.Errorf("asignación inesperada: %+v", m)
	}
}

func TestNewMembership_CamposObligatorios(t *testing.T) {
	casos := []struct {
		nombre     string
		proyectoID string
		userID     string
		rol        Role
		campo      string
	}{
		{"proyecto vacío", "", "user-1", RolScrumMaster, "proyecto_id"},
		{"integrante vacío", "proy-1", "", RolScrumMaster, "integrante_id"},
		{"rol inválido", "proy-1", "user-1", Role("Product Owner"), "rol"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := NewMembership(c.proyectoID, c.userID, c.rol)
			if err == nil {
				t.Fatal("se esperaba error de validación")
			}
			var verr ValidationError
			if !errors.As(err, &verr) {
				t.Fatalf("se esperaba ValidationError, se obtuvo %T", err)
			}
			if verr.Campo != c.campo {
				t.Errorf("se esperaba campo %q, se obtuvo %q", c.campo, verr.Campo)
			}
		})
	}
}
