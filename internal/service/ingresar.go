package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// Ingresar autentica a un usuario con nombre y contraseña y le abre una sesión (HU-20).
type Ingresar struct {
	credenciales repository.BuscadorCredenciales
	hasheador    HasheadorContrasenas
	emisor       EmisorDeSesiones
	hashSenuelo  string
}

// NewIngresar calcula el hash señuelo con los mismos parámetros del hasheador,
// para que verificar contra él cueste lo mismo que contra un hash real.
func NewIngresar(credenciales repository.BuscadorCredenciales, hasheador HasheadorContrasenas, emisor EmisorDeSesiones) *Ingresar {
	senuelo, _ := hasheador.Hashear("contrasena-senuelo")
	return &Ingresar{credenciales: credenciales, hasheador: hasheador, emisor: emisor, hashSenuelo: senuelo}
}

// Ejecutar devuelve el mismo ErrCredencialesInvalidas ante contraseña incorrecta,
// usuario inexistente o usuario sin contraseña, y verifica exactamente una vez
// en los tres casos.
func (i *Ingresar) Ejecutar(ctx context.Context, nombre, contrasena string) (ResultadoSesion, error) {
	usuario, hash, err := i.credenciales.ObtenerPorNombreNormalizado(ctx, domain.NormalizarNombre(nombre))
	if err != nil && !errors.Is(err, domain.ErrNoEncontrado) {
		return ResultadoSesion{}, fmt.Errorf("buscar credenciales: %w", err)
	}
	if err != nil || hash == "" {
		hash = i.hashSenuelo
		err = domain.ErrCredencialesInvalidas
	}
	coincide := i.hasheador.Verificar(contrasena, hash)
	if err != nil || !coincide {
		return ResultadoSesion{}, domain.ErrCredencialesInvalidas
	}
	token, sesion, err := i.emisor.Emitir(ctx, usuario.ID)
	if err != nil {
		return ResultadoSesion{}, err
	}
	return ResultadoSesion{Usuario: usuario, Token: token, VenceEn: sesion.VenceEn}, nil
}
