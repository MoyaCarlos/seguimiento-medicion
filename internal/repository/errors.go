package repository

import "errors"

// ErrMiembroDuplicado indica que el integrante ya pertenece al proyecto.
var ErrMiembroDuplicado = errors.New("el integrante ya pertenece al proyecto")

// ErrUsuarioDuplicado indica que ya existe un integrante con el mismo nombre normalizado.
var ErrUsuarioDuplicado = errors.New("el integrante ya existe")

// ErrIntegridadReferencial indica una violación de llave foránea en la persistencia.
var ErrIntegridadReferencial = errors.New("violación de integridad referencial")
