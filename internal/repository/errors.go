package repository

import "errors"

// ErrNoEncontrado indica que el registro solicitado no existe.
var ErrNoEncontrado = errors.New("registro no encontrado")

// ErrMiembroDuplicado indica que el integrante ya pertenece al proyecto.
var ErrMiembroDuplicado = errors.New("el integrante ya pertenece al proyecto")

// ErrUsuarioDuplicado indica que ya existe un integrante con el mismo nombre normalizado.
var ErrUsuarioDuplicado = errors.New("el integrante ya existe")
