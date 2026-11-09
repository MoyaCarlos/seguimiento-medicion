package service

import (
	"encoding/base64"
	"strconv"
	"strings"
	"testing"
)

const iteracionesDePrueba = 10

func TestServicioContrasenas_HashearUsaElFormatoPBKDF2(t *testing.T) {
	s := NewServicioContrasenas(iteracionesDePrueba)

	hash, err := s.Hashear("secreta-123")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	partes := strings.Split(hash, "$")
	if len(partes) != 4 || partes[0] != "pbkdf2-sha256" {
		t.Fatalf("formato inesperado: %q", hash)
	}
	if partes[1] != strconv.Itoa(iteracionesDePrueba) {
		t.Errorf("se esperaban %d iteraciones, se obtuvo %q", iteracionesDePrueba, partes[1])
	}
	if sal, err := base64.RawStdEncoding.DecodeString(partes[2]); err != nil || len(sal) != 16 {
		t.Errorf("se esperaba una sal de 16 bytes: %q (%v)", partes[2], err)
	}
	if clave, err := base64.RawStdEncoding.DecodeString(partes[3]); err != nil || len(clave) != 32 {
		t.Errorf("se esperaba una clave de 32 bytes: %q (%v)", partes[3], err)
	}
}

func TestServicioContrasenas_DosHashesDeLaMismaContrasenaDifieren(t *testing.T) {
	s := NewServicioContrasenas(iteracionesDePrueba)

	h1, _ := s.Hashear("secreta-123")
	h2, _ := s.Hashear("secreta-123")

	if h1 == "" || h1 == h2 {
		t.Errorf("la sal debería ser única por hash: %q y %q", h1, h2)
	}
}

func TestServicioContrasenas_VerificarAciertaConLaContrasenaCorrecta(t *testing.T) {
	s := NewServicioContrasenas(iteracionesDePrueba)
	hash, _ := s.Hashear("secreta-123")

	if !s.Verificar("secreta-123", hash) {
		t.Error("se esperaba que la contraseña correcta verifique")
	}
}

func TestServicioContrasenas_VerificarFallaConOtraContrasena(t *testing.T) {
	s := NewServicioContrasenas(iteracionesDePrueba)
	hash, _ := s.Hashear("secreta-123")

	if s.Verificar("otra-cosa", hash) {
		t.Error("una contraseña distinta no debería verificar")
	}
}

func TestServicioContrasenas_VerificarFallaSinErrorConFormatoCorruptoOVacio(t *testing.T) {
	s := NewServicioContrasenas(iteracionesDePrueba)
	corruptos := []string{
		"",
		"basura",
		"pbkdf2-sha256$10$sal",
		"pbkdf2-sha256$abc$c2Fs$Y2xhdmU",
		"pbkdf2-sha256$10$***$Y2xhdmU",
		"pbkdf2-sha256$10$c2Fs$***",
		"otro-algoritmo$10$c2Fs$Y2xhdmU",
	}

	for _, hash := range corruptos {
		if s.Verificar("secreta-123", hash) {
			t.Errorf("el hash %q no debería verificar", hash)
		}
	}
}

func TestServicioContrasenas_VerificaConLasIteracionesGuardadasEnElHash(t *testing.T) {
	hash, _ := NewServicioContrasenas(5).Hashear("secreta-123")

	if !NewServicioContrasenas(iteracionesDePrueba).Verificar("secreta-123", hash) {
		t.Error("debería verificar con las iteraciones del hash, no con las del servicio")
	}
}

func TestServicioContrasenasProduccion_UsaLosParametrosDeProduccion(t *testing.T) {
	s := NewServicioContrasenasProduccion()

	if s.iteraciones != 600_000 || s.tamanoSal != 16 || s.tamanoClave != 32 {
		t.Errorf("parámetros inesperados: %d iteraciones, sal de %d, clave de %d", s.iteraciones, s.tamanoSal, s.tamanoClave)
	}
}
