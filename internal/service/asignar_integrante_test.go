package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func proyectoDePrueba(t *testing.T, proyectos *fakeProjectRepository, usuarios *fakeUserRepository) int64 {
	t.Helper()
	ana, err := usuarios.Guardar(context.Background(), domain.User{Nombre: "Ana", NombreNormalizado: "ana"})
	if err != nil {
		t.Fatalf("no se pudo preparar a Ana: %v", err)
	}
	proyecto, err := NewCrearProyecto(proyectos).Ejecutar(context.Background(), CrearProyectoInput{Nombre: "Proyecto", CreadorID: ana.ID})
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
	if miembro.Nombre != "Jimena" || miembro.Role != domain.RolProductBuilder {
		t.Errorf("integrante inesperado: %+v", miembro)
	}
	if len(proyectos.miembros) != 2 { // creador + Jimena
		t.Errorf("se esperaban 2 vinculaciones, se obtuvieron %d", len(proyectos.miembros))
	}
}

func TestAsignarIntegrante_VinculadoEnNoCero(t *testing.T) {
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
	if miembro.CreadoEn.IsZero() {
		t.Error("se esperaba que vinculado_en no fuera la fecha cero")
	}
}

func TestAsignarIntegrante_ReutilizaUsuario(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios)

	ana, err := usuarios.ObtenerPorNombreNormalizado(context.Background(), "ana")
	if err != nil {
		t.Fatalf("no se encontró a Ana: %v", err)
	}
	otro, err := NewCrearProyecto(proyectos).Ejecutar(context.Background(), CrearProyectoInput{Nombre: "Otro", CreadorID: ana.ID})
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
		ProyectoID: otro.ID, Nombre: "  Jimena ", Rol: domain.RolProductBuilder,
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

func TestAsignarIntegrante_RolInvalidoNoCreaUsuario(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios)
	createsInicial := usuarios.creates

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	_, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: proyectoID, Nombre: "Jimena", Rol: domain.Role("Product Owner"),
	})
	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "rol" {
		t.Fatalf("se esperaba ValidationError de rol, se obtuvo %v", err)
	}
	if usuarios.creates != createsInicial {
		t.Fatalf("un rol inválido no debe crear usuario: se crearon %d usuarios nuevos", usuarios.creates-createsInicial)
	}
}

func TestAsignarIntegrante_ProyectoInexistente(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	_, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: 0, Nombre: "Jimena", Rol: domain.RolProductBuilder,
	})
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
	if len(proyectos.miembros) != 0 {
		t.Error("no se debía persistir ninguna vinculación")
	}
}

func TestAsignarIntegrante_ProyectoInexistenteDominaRolInvalido(t *testing.T) {
	servicio := NewAsignarIntegrante(&fakeProjectRepository{}, &fakeUserRepository{})

	_, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: 0, Nombre: "Jimena", Rol: domain.Role("Product Owner"),
	})
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba proyecto no encontrado (404) antes que rol inválido (400), se obtuvo %v", err)
	}
}

func TestAsignarIntegrante_RechazaUnSegundoScrumMaster(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios)

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	_, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: proyectoID, Nombre: "Candela", Rol: domain.RolScrumMaster,
	})

	if !errors.Is(err, domain.ErrScrumMasterYaExiste) {
		t.Fatalf("se esperaba domain.ErrScrumMasterYaExiste, se obtuvo %v", err)
	}
	if len(proyectos.miembros) != 1 { // solo el creador
		t.Errorf("no se debía agregar a nadie: hay %d integrantes", len(proyectos.miembros))
	}
}

func TestAsignarIntegrante_SegundoScrumMasterRechazadoNoCreaUsuario(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios)
	createsInicial := usuarios.creates

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	_, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: proyectoID, Nombre: "Candela", Rol: domain.RolScrumMaster,
	})

	if !errors.Is(err, domain.ErrScrumMasterYaExiste) {
		t.Fatalf("se esperaba domain.ErrScrumMasterYaExiste, se obtuvo %v", err)
	}
	if usuarios.creates != createsInicial {
		t.Errorf("un rechazo no debe crear usuario: se crearon %d usuarios nuevos", usuarios.creates-createsInicial)
	}
}

func TestAsignarIntegrante_ProductBuilderConScrumMasterPresente(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios)

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	for _, nombre := range []string{"Candela", "Carlos"} {
		if _, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
			ProyectoID: proyectoID, Nombre: nombre, Rol: domain.RolProductBuilder,
		}); err != nil {
			t.Fatalf("se permiten varios Product Builder: %v", err)
		}
	}
}

func TestAsignarIntegrante_ScrumMasterEnProyectoSinScrumMaster(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyecto, err := proyectos.Guardar(context.Background(), domain.Project{Nombre: "Sin SM"})
	if err != nil {
		t.Fatalf("no se pudo preparar el proyecto: %v", err)
	}

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	if _, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: proyecto.ID, Nombre: "Candela", Rol: domain.RolScrumMaster,
	}); err != nil {
		t.Fatalf("un proyecto sin Scrum Master admite uno: %v", err)
	}
}

func TestAsignarIntegrante_DuplicadoGanaAlScrumMasterRepetido(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	usuarios := &fakeUserRepository{}
	proyectoID := proyectoDePrueba(t, proyectos, usuarios) // Ana es el Scrum Master

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	if _, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: proyectoID, Nombre: "Jimena", Rol: domain.RolProductBuilder,
	}); err != nil {
		t.Fatalf("no se pudo preparar a Jimena: %v", err)
	}

	for _, nombre := range []string{"Jimena", "Ana"} {
		_, err := servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
			ProyectoID: proyectoID, Nombre: nombre, Rol: domain.RolScrumMaster,
		})
		var verr domain.ValidationError
		if !errors.As(err, &verr) || verr.Campo != "integrante" {
			t.Errorf("%s: se esperaba el duplicado (ValidationError de integrante), se obtuvo %v", nombre, err)
		}
		if errors.Is(err, domain.ErrScrumMasterYaExiste) {
			t.Errorf("%s: el duplicado debe ganar al Scrum Master repetido", nombre)
		}
	}
}

func TestAsignarIntegrante_CarreraDeScrumMasterDevuelveElErrorDeDominio(t *testing.T) {
	proyectos := &fakeProjectRepository{errAgregar: domain.ErrScrumMasterYaExiste}
	usuarios := &fakeUserRepository{}
	proyecto, err := proyectos.Guardar(context.Background(), domain.Project{Nombre: "Sin SM"})
	if err != nil {
		t.Fatalf("no se pudo preparar el proyecto: %v", err)
	}

	servicio := NewAsignarIntegrante(proyectos, usuarios)
	_, err = servicio.Ejecutar(context.Background(), AsignarIntegranteInput{
		ProyectoID: proyecto.ID, Nombre: "Candela", Rol: domain.RolScrumMaster,
	})

	if !errors.Is(err, domain.ErrScrumMasterYaExiste) {
		t.Fatalf("el error del repositorio debe llegar sin cambios, se obtuvo %v", err)
	}
}
