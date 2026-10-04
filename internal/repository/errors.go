package repository

import "errors"

// ErrMiembroDuplicado indica que el integrante ya pertenece al proyecto.
var ErrMiembroDuplicado = errors.New("el integrante ya pertenece al proyecto")

// ErrUsuarioDuplicado indica que ya existe un integrante con el mismo nombre normalizado.
var ErrUsuarioDuplicado = errors.New("el integrante ya existe")

// ErrIntegridadReferencial indica una violación de llave foránea en la persistencia.
var ErrIntegridadReferencial = errors.New("violación de integridad referencial")

// ErrEsquemaViejo indica que la base local tiene un esquema previo a la
// migración a IDs int64 y debe borrarse antes de reiniciar la aplicación.
var ErrEsquemaViejo = errors.New("seguimiento.db tiene un esquema viejo (IDs TEXT); borrá el archivo y reiniciá")
