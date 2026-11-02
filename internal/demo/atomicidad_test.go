package demo

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

func TestCargar_UnBloqueQueFallaAlFinalNoDejaNingunDato(t *testing.T) {
	db := baseMigrada(t)
	errInyectado := errors.New("fallo inyectado")
	fallaAlFinal := func(context.Context, *sql.Tx, *carga) error { return errInyectado }
	cargador := &Cargador{
		db:      db,
		ahora:   func() time.Time { return ahoraFijo },
		bloques: append(append([]bloque{}, bloquesV1...), fallaAlFinal),
	}

	_, err := cargador.Cargar(context.Background())

	if !errors.Is(err, errInyectado) {
		t.Fatalf("error %v, se esperaba el fallo inyectado", err)
	}
	for tabla, filas := range fotoDeLaBase(t, db) {
		if filas != 0 {
			t.Errorf("%s quedó con %d filas tras una carga fallida", tabla, filas)
		}
	}
}
