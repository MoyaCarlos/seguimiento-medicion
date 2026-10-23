package repository

import (
	"errors"

	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// ErrMiembroDuplicado indica que el integrante ya pertenece al proyecto.
var ErrMiembroDuplicado = errors.New("el integrante ya pertenece al proyecto")

// ErrUsuarioDuplicado indica que ya existe un integrante con el mismo nombre normalizado.
var ErrUsuarioDuplicado = errors.New("el integrante ya existe")

// ErrIntegridadReferencial indica una violación de llave foránea en la persistencia.
var ErrIntegridadReferencial = errors.New("violación de integridad referencial")

// ErrEsquemaViejo indica que la base local tiene un esquema previo a la
// migración 001 y debe borrarse antes de reiniciar la aplicación.
var ErrEsquemaViejo = errors.New("seguimiento.db tiene un esquema viejo; borrá el archivo y reiniciá")

// errorDeEscritura traduce una violación de llave foránea a
// ErrIntegridadReferencial; cualquier otro error vuelve sin cambios.
func errorDeEscritura(err error) error {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY {
		return ErrIntegridadReferencial
	}
	return err
}
