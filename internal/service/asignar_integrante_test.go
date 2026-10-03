package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func proyectoDePrueba(t *testing.T, proyectos *fakeProjectRepository, usuarios *fakeUserRepository) int64 {
	t.Helper()
	crear := NewCrearProyecto(proyectos, usuarios)
	proyecto, err := crear.Ejecutar(context.Background(), CrearProyectoInput{Nombre: "Proyecto", Creador: "Ana"})
	if err != nil {
		t.Fatalf("no se pudo preparar el proyecto: %v", err)
	}
	return proyecto.ID
}

func TestAsignarIntegrante_Exitosa(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios)

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	miembro, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: proyectoID,
		Nombre:     "Jimena",
		Rol:        domain.RolProductBuilder,
	})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if miembro.Name != "Jimena" || miembro.Role != domain.RolProductBuilder {
		t.Errorf("integrante inesperado: %+v", miembro)
	}
	if len(proyectos.miembros) != 2 { // creador + Jimena
		t.Errorf("se esperaban 2 vinculaciones, se obtuvieron %d", len(proyectos.miembros))
	}
}

func TestAsignarIntegrante_ReutilizaUsuario(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios)

	otro, err := NewCrearProyecto(proyectos, usuarios).Ejecutar(context.Background(), CrearProyectoInput{Nombre: "Otro", Creador: "Ana"})
	if err != nil {
		t.Fatalf("no se pudo preparar el segundo proyecto: %v", err)
	}
	createsInicial := usuarios.creates

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	if _, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: proyectoID, Nombre: "jimena", Rol: domain.RolProductBuilder,
	}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	// En otro proyecto, "  Jimena " debe reutilizar el mismo usuario.
	if _, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: otro.ID, Nombre: "  Jimena ", Rol: domain.RolScrumMaster,
	}); err != nil {
		t.Fatalf("no se esperaba error al reutilizar el integrante: %v", err)
	}
	if usuarios.creates != createsInicial+1 {
		t.Errorf("se esperaba reutilizar el usuario; se crearon %d usuarios nuevos", usuarios.creates-createsInicial)
	}
}

func TestAsignarIntegrante_Duplicado(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios)

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	entrada := AsignarIntegranteInput{ProyectoID: proyectoID, Nombre: "Jimena", Rol: domain.RolProductBuilder}
	if _, err := servicio.Ejecutar(context.Background(), entrada); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	_, err := servicio.Ejecutar(context.Background(), entrada)
	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "integrante" {
		t.Fatalf("se esperaba ValidationError de integrante, se obtuvo %v", err)
	}
}

func TestAsignarIntegrante_RolInvalido(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios)

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	_, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: proyectoID, Nombre: "Jimena", Rol: domain.Role("Product Owner"),
	})
	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "rol" {
		t.Fatalf("se esperaba ValidationError de rol, se obtuvo %v", err)
	}
}

func TestAsignarIntegrante_ProyectoInexistente(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	_, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: 0, Nombre: "Jimena", Rol: domain.RolProductBuilder,
	})
	if !errors.Is(err, domain.ErrNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrNoEncontrado, se obtuvo %v", err)
	}
	if len(proyectos.miembros) != 0 {
		t.Error("no se debía persistir ninguna vinculación")
	}
}

func TestAsignarIntegrante_VariosScrumMaster(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios)

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	for _, nombre := range []string{"Candela", "Carlos"} {
		if _, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
			ProyectoID: proyectoID, Nombre: nombre, Rol: domain.RolScrumMaster,
		}); err != nil {
			t.Fatalf("se permiten varios Scrum Master: %v", err)
		}
	}
}
