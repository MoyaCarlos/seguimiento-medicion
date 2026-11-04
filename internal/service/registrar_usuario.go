package service

import (
	"context"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// ResultadoSesion es el usuario autenticado y la sesión recién emitida, con el
// token en claro (HU-20).
type ResultadoSesion struct {
	Usuario domain.User
	Token   string
	VenceEn time.Time
}

// EmisorDeSesiones crea una sesión para un usuario y devuelve su token en claro.
type EmisorDeSesiones interface {
	Emitir(ctx context.Context, usuarioID int64) (string, domain.Sesion, error)
}

// RegistrarUsuario da de alta un usuario con contraseña y le abre una sesión (HU-20).
type RegistrarUsuario struct {
	credenciales repository.RegistroCredenciales
	hasheador    HasheadorContrasenas
	emisor       EmisorDeSesiones
}

func NewRegistrarUsuario(credenciales repository.RegistroCredenciales, hasheador HasheadorContrasenas, emisor EmisorDeSesiones) *RegistrarUsuario {
	return &RegistrarUsuario{credenciales: credenciales, hasheador: hasheador, emisor: emisor}
}

func (r *RegistrarUsuario) Ejecutar(ctx context.Context, nombre, contrasena string) (ResultadoSesion, error) {
	return ResultadoSesion{}, nil
}
