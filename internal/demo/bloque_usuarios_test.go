package demo

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

var nombresDeLaDemo = []string{"carlos.moya", "ana.iriarte", "candela.vulcano", "jimena.martinez"}

func TestBloqueUsuarios_CreaLosCuatroUsuariosDeLaDemoConElMismoHash(t *testing.T) {
	db := baseMigrada(t)
	hash := hashDePrueba(t)

	c, err := ejecutarBloques(context.Background(), db, ahoraFijo, hash, []bloque{bloqueUsuarios})

	if err != nil {
		t.Fatalf("ejecutar: %v", err)
	}
	for _, nombre := range nombresDeLaDemo {
		var guardado string
		if err := db.QueryRow(`SELECT COALESCE(contrasena_hash, '') FROM usuarios WHERE nombre_normalizado = ?`, nombre).Scan(&guardado); err != nil {
			t.Errorf("falta el usuario %s: %v", nombre, err)
			continue
		}
		if guardado != hash {
			t.Errorf("%s: el hash guardado %q no es el hash único de la carga", nombre, guardado)
		}
		if !hasheadorBarato.Verificar(contrasenaDemo, guardado) {
			t.Errorf("%s: la contraseña %s no verifica contra su hash", nombre, contrasenaDemo)
		}
		if c.usuarios[nombre] == 0 {
			t.Errorf("falta el id de %s", nombre)
		}
	}
	if c.resumen.Usuarios != 4 {
		t.Errorf("el resumen cuenta %d usuarios creados, se esperaban 4", c.resumen.Usuarios)
	}
}

func TestBloqueUsuarios_NoQuedanUsuariosBeta(t *testing.T) {
	db := baseMigrada(t)

	ejecutar(t, db, bloqueUsuarios)

	if n := contar(t, db, `SELECT COUNT(*) FROM usuarios WHERE nombre LIKE 'Beta %'`); n != 0 {
		t.Errorf("no debe haber usuarios Beta, hay %d", n)
	}
}

func TestBloqueUsuarios_LosNombresCumplenLaReglaDeHU20(t *testing.T) {
	db := baseMigrada(t)

	ejecutar(t, db, bloqueUsuarios)

	filas, err := db.Query(`SELECT nombre FROM usuarios`)
	if err != nil {
		t.Fatal(err)
	}
	defer filas.Close()
	for filas.Next() {
		var nombre string
		_ = filas.Scan(&nombre)
		if err := domain.ValidarNombreUsuario(nombre); err != nil {
			t.Errorf("%q no cumple la regla del nombre de usuario: %v", nombre, err)
		}
	}
}

func TestBloqueUsuarios_UnUsuarioExistenteAbortaLaTransaccion(t *testing.T) {
	db := baseMigrada(t)
	if _, err := db.Exec(`INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES ('ana.iriarte', 'ana.iriarte', '2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}

	_, err := ejecutarBloques(context.Background(), db, ahoraFijo, hashDePrueba(t), []bloque{bloqueUsuarios})

	if !esViolacionDeClave(err) {
		t.Fatalf("error %v, se esperaba una violación de clave única", err)
	}
	if n := contar(t, db, `SELECT COUNT(*) FROM usuarios`); n != 1 {
		t.Errorf("se esperaba solo el usuario preexistente, hay %d", n)
	}
	var existente ErrUsuarioExistente
	if errors.As(err, &existente) {
		t.Error("el bloque no clasifica: eso lo hace el Cargador tras el rollback")
	}
}
