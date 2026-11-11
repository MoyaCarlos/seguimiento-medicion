package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// equipoDePrueba registra a Ana (id 1 = actorDePrueba), que crea el proyecto y
// es su Scrum Master, y a Jimena y Carlos con cuenta. «Sin Cuenta» existe en
// usuarios pero sin contraseña.
func equipoDePrueba(t *testing.T) (*fakeProjectRepository, *fakeCredenciales, int64) {
	t.Helper()
	ctx := context.Background()
	proyectos := &fakeProjectRepository{}
	credenciales := &fakeCredenciales{}
	for _, nombre := range []string{"Ana", "Jimena", "Carlos"} {
		if _, err := credenciales.Registrar(ctx, domain.User{Nombre: nombre, NombreNormalizado: domain.NormalizarNombre(nombre)}, "hash"); err != nil {
			t.Fatalf("no se pudo registrar a %s: %v", nombre, err)
		}
	}
	if _, err := credenciales.Registrar(ctx, domain.User{Nombre: "Sin Cuenta", NombreNormalizado: "sin cuenta"}, ""); err != nil {
		t.Fatalf("no se pudo preparar al usuario sin cuenta: %v", err)
	}
	proyecto, err := NewCrearProyecto(proyectos).Ejecutar(ctx, CrearProyectoInput{Nombre: "Proyecto", CreadorID: actorDePrueba})
	if err != nil {
		t.Fatalf("no se pudo preparar el proyecto: %v", err)
	}
	return proyectos, credenciales, proyecto.ID
}

func agregar(proyectos *fakeProjectRepository, credenciales *fakeCredenciales, proyectoID int64, nombre string, rol domain.Role) (domain.Member, error) {
	return NewAsignarIntegrante(proyectos, credenciales).Ejecutar(context.Background(), AsignarIntegranteInput{
		ActorID: actorDePrueba, ProyectoID: proyectoID, Nombre: nombre, Rol: rol,
	})
}

// HU-21, escenario 8: un usuario registrado entra como Product Builder, con
// rol explícito o sin rol.
func TestAsignarIntegrante_UsuarioRegistradoEntraComoProductBuilder(t *testing.T) {
	for _, rol := range []domain.Role{domain.RolProductBuilder, ""} {
		proyectos, credenciales, proyectoID := equipoDePrueba(t)

		miembro, err := agregar(proyectos, credenciales, proyectoID, "  jimena ", rol)

		if err != nil {
			t.Fatalf("rol %q: no se esperaba error: %v", rol, err)
		}
		if miembro.Nombre != "Jimena" || miembro.Role != domain.RolProductBuilder || miembro.UserID != 2 {
			t.Errorf("rol %q: integrante inesperado: %+v", rol, miembro)
		}
		if miembro.CreadoEn.IsZero() {
			t.Errorf("rol %q: se esperaba vinculado_en distinto de la fecha cero", rol)
		}
		if len(proyectos.miembros) != 2 { // Ana + Jimena
			t.Errorf("rol %q: se esperaban 2 vinculaciones, se obtuvieron %d", rol, len(proyectos.miembros))
		}
	}
}

// Escenario 8: no se crean usuarios; quien no existe o no tiene cuenta da el
// mismo «usuario no encontrado» (404).
func TestAsignarIntegrante_SoloUsuariosRegistrados(t *testing.T) {
	for _, nombre := range []string{"Inexistente", "Sin Cuenta"} {
		proyectos, credenciales, proyectoID := equipoDePrueba(t)
		usuariosAntes := len(credenciales.usuarios)

		_, err := agregar(proyectos, credenciales, proyectoID, nombre, domain.RolProductBuilder)

		if !errors.Is(err, domain.ErrUsuarioNoEncontrado) {
			t.Fatalf("%s: se esperaba domain.ErrUsuarioNoEncontrado, se obtuvo %v", nombre, err)
		}
		if len(proyectos.miembros) != 1 || len(credenciales.usuarios) != usuariosAntes {
			t.Errorf("%s: no se debía agregar integrante ni crear usuario", nombre)
		}
	}
}

func TestAsignarIntegrante_Duplicado(t *testing.T) {
	proyectos, credenciales, proyectoID := equipoDePrueba(t)
	if _, err := agregar(proyectos, credenciales, proyectoID, "Jimena", domain.RolProductBuilder); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	_, err := agregar(proyectos, credenciales, proyectoID, "Jimena", domain.RolProductBuilder)

	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "integrante" {
		t.Fatalf("se esperaba ValidationError de integrante, se obtuvo %v", err)
	}
}

func TestAsignarIntegrante_RolInvalido(t *testing.T) {
	proyectos, credenciales, proyectoID := equipoDePrueba(t)

	_, err := agregar(proyectos, credenciales, proyectoID, "Jimena", domain.Role("Product Owner"))

	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "rol" {
		t.Fatalf("se esperaba ValidationError de rol, se obtuvo %v", err)
	}
}

func TestAsignarIntegrante_ProyectoInexistenteDominaRolInvalido(t *testing.T) {
	proyectos, credenciales, _ := equipoDePrueba(t)

	_, err := agregar(proyectos, credenciales, 999, "Jimena", domain.Role("Product Owner"))

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba proyecto no encontrado (404) antes que rol inválido (400), se obtuvo %v", err)
	}
}

// #236 y HU-21: todo proyecto ya tiene Scrum Master, así que pedir ese rol es
// 409, aunque la persona no tenga cuenta.
func TestAsignarIntegrante_PedirScrumMasterEs409(t *testing.T) {
	for _, nombre := range []string{"Carlos", "Inexistente"} {
		proyectos, credenciales, proyectoID := equipoDePrueba(t)

		_, err := agregar(proyectos, credenciales, proyectoID, nombre, domain.RolScrumMaster)

		if !errors.Is(err, domain.ErrScrumMasterYaExiste) {
			t.Fatalf("%s: se esperaba domain.ErrScrumMasterYaExiste, se obtuvo %v", nombre, err)
		}
		if len(proyectos.miembros) != 1 {
			t.Errorf("%s: no se debía agregar a nadie: hay %d integrantes", nombre, len(proyectos.miembros))
		}
	}
}

func TestAsignarIntegrante_VariosProductBuilder(t *testing.T) {
	proyectos, credenciales, proyectoID := equipoDePrueba(t)
	for _, nombre := range []string{"Jimena", "Carlos"} {
		if _, err := agregar(proyectos, credenciales, proyectoID, nombre, domain.RolProductBuilder); err != nil {
			t.Fatalf("se permiten varios Product Builder: %v", err)
		}
	}
}

func TestAsignarIntegrante_DuplicadoGanaAlScrumMasterRepetido(t *testing.T) {
	proyectos, credenciales, proyectoID := equipoDePrueba(t) // Ana es el Scrum Master
	if _, err := agregar(proyectos, credenciales, proyectoID, "Jimena", domain.RolProductBuilder); err != nil {
		t.Fatalf("no se pudo preparar a Jimena: %v", err)
	}

	for _, nombre := range []string{"Jimena", "Ana"} {
		_, err := agregar(proyectos, credenciales, proyectoID, nombre, domain.RolScrumMaster)
		var verr domain.ValidationError
		if !errors.As(err, &verr) || verr.Campo != "integrante" {
			t.Errorf("%s: se esperaba el duplicado (ValidationError de integrante), se obtuvo %v", nombre, err)
		}
	}
}
