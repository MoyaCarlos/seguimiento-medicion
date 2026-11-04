package service

import (
	"context"
	"errors"
	"fmt"
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
	if err := domain.ValidarNombreUsuario(nombre); err != nil {
		return ResultadoSesion{}, err
	}
	if err := domain.ValidarContrasena(contrasena); err != nil {
		return ResultadoSesion{}, err
	}
	hash, err := r.hasheador.Hashear(contrasena)
	if err != nil {
		return ResultadoSesion{}, fmt.Errorf("hashear contraseña: %w", err)
	}
	usuario, err := r.credenciales.Registrar(ctx, domain.User{Nombre: nombre, NombreNormalizado: domain.NormalizarNombre(nombre)}, hash)
	if errors.Is(err, repository.ErrUsuarioDuplicado) {
		return ResultadoSesion{}, domain.ErrNombreUsuarioEnUso
	}
	if err != nil {
		return ResultadoSesion{}, fmt.Errorf("registrar usuario: %w", err)
	}
	token, sesion, err := r.emisor.Emitir(ctx, usuario.ID)
	if err != nil {
		return ResultadoSesion{}, err
	}
	return ResultadoSesion{Usuario: usuario, Token: token, VenceEn: sesion.VenceEn}, nil
}
