package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func armarIngresar(t *testing.T) (*Ingresar, *fakeContrasenas, *fakeSesiones) {
	t.Helper()
	credenciales := &fakeCredenciales{}
	hasheador := &fakeContrasenas{}
	sesiones := &fakeSesiones{}
	hash, _ := hasheador.Hashear("clave-segura-1")
	if _, err := credenciales.Registrar(context.Background(), domain.User{Nombre: "Jimena", NombreNormalizado: "jimena"}, hash); err != nil {
		t.Fatalf("preparar usuario: %v", err)
	}
	// Usuario previo sin contraseña: hash vacío.
	if _, err := credenciales.Registrar(context.Background(), domain.User{Nombre: "Previo", NombreNormalizado: "previo"}, ""); err != nil {
		t.Fatalf("preparar usuario previo: %v", err)
	}
	caso := NewIngresar(credenciales, hasheador, NewEmisorSesiones(sesiones, relojDeEmision))
	return caso, hasheador, sesiones
}

func TestIngresar_ExitoConOtrasMayusculasDevuelveUsuarioYSesion(t *testing.T) {
	caso, _, sesiones := armarIngresar(t)

	resultado, err := caso.Ejecutar(context.Background(), "JIMENA", "clave-segura-1")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if resultado.Usuario.NombreNormalizado != "jimena" {
		t.Errorf("NombreNormalizado = %q, se esperaba %q", resultado.Usuario.NombreNormalizado, "jimena")
	}
	if resultado.Token == "" {
		t.Errorf("se esperaba un token")
	}
	if sesiones.creadas != 1 {
		t.Errorf("sesiones creadas = %d, se esperaba 1", sesiones.creadas)
	}
}

func TestIngresar_FallasDevuelvenElMismoErrorYVerificanUnaVez(t *testing.T) {
	casos := []struct {
		nombre     string
		usuario    string
		contrasena string
	}{
		{"contraseña incorrecta", "jimena", "otra-clave-22"},
		{"usuario inexistente", "fantasma", "clave-segura-1"},
		{"usuario previo sin contraseña", "previo", "clave-segura-1"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			caso, hasheador, sesiones := armarIngresar(t)
			antes := hasheador.llamadasAVerificar()

			_, err := caso.Ejecutar(context.Background(), c.usuario, c.contrasena)

			if !errors.Is(err, domain.ErrCredencialesInvalidas) {
				t.Fatalf("se esperaba ErrCredencialesInvalidas, se obtuvo %v", err)
			}
			if n := hasheador.llamadasAVerificar() - antes; n != 1 {
				t.Errorf("Verificar llamado %d veces, se esperaba 1", n)
			}
			if sesiones.creadas != 0 {
				t.Errorf("no se esperaba emitir sesión")
			}
		})
	}
}

func TestIngresar_DosIngresosGeneranDosSesionesDistintas(t *testing.T) {
	caso, _, sesiones := armarIngresar(t)

	primero, err := caso.Ejecutar(context.Background(), "jimena", "clave-segura-1")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	segundo, err := caso.Ejecutar(context.Background(), "jimena", "clave-segura-1")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if primero.Token == "" || primero.Token == segundo.Token {
		t.Errorf("se esperaban tokens distintos y no vacíos")
	}
	if sesiones.creadas != 2 {
		t.Errorf("sesiones creadas = %d, se esperaba 2", sesiones.creadas)
	}
}
