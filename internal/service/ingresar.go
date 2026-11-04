package service

import (
	"context"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// Ingresar autentica a un usuario con nombre y contraseña y le abre una sesión (HU-20).
type Ingresar struct {
	credenciales repository.BuscadorCredenciales
	hasheador    HasheadorContrasenas
	emisor       EmisorDeSesiones
}

func NewIngresar(credenciales repository.BuscadorCredenciales, hasheador HasheadorContrasenas, emisor EmisorDeSesiones) *Ingresar {
	return &Ingresar{credenciales: credenciales, hasheador: hasheador, emisor: emisor}
}

func (i *Ingresar) Ejecutar(ctx context.Context, nombre, contrasena string) (ResultadoSesion, error) {
	return ResultadoSesion{}, nil
}
