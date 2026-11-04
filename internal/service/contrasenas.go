package service

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

const (
	algoritmoContrasena   = "pbkdf2-sha256"
	iteracionesContrasena = 600_000
	tamanoSalContrasena   = 16
	tamanoClaveContrasena = 32
)

// HasheadorContrasenas deriva y verifica hashes de contraseñas (HU-20).
type HasheadorContrasenas interface {
	Hashear(contrasena string) (string, error)
	// Verificar devuelve false, sin error, si el hash está corrupto o vacío.
	Verificar(contrasena, hash string) bool
}

// ServicioContrasenas implementa HasheadorContrasenas con PBKDF2-SHA256.
type ServicioContrasenas struct {
	iteraciones int
	tamanoSal   int
	tamanoClave int
}

func NewServicioContrasenas(iteraciones int) *ServicioContrasenas {
	return &ServicioContrasenas{iteraciones: iteraciones, tamanoSal: tamanoSalContrasena, tamanoClave: tamanoClaveContrasena}
}

func NewServicioContrasenasProduccion() *ServicioContrasenas {
	return NewServicioContrasenas(iteracionesContrasena)
}

func (s *ServicioContrasenas) Hashear(contrasena string) (string, error) {
	sal := make([]byte, s.tamanoSal)
	if _, err := rand.Read(sal); err != nil {
		return "", fmt.Errorf("generar sal: %w", err)
	}
	clave, err := pbkdf2.Key(sha256.New, contrasena, sal, s.iteraciones, s.tamanoClave)
	if err != nil {
		return "", fmt.Errorf("derivar clave: %w", err)
	}
	return fmt.Sprintf("%s$%d$%s$%s", algoritmoContrasena, s.iteraciones,
		base64.RawStdEncoding.EncodeToString(sal), base64.RawStdEncoding.EncodeToString(clave)), nil
}

// Verificar usa las iteraciones, la sal y el largo de clave guardados en el
// hash, así los hashes viejos siguen verificando si cambian los parámetros.
func (s *ServicioContrasenas) Verificar(contrasena, hash string) bool {
	partes := strings.Split(hash, "$")
	if len(partes) != 4 || partes[0] != algoritmoContrasena {
		return false
	}
	iteraciones, err := strconv.Atoi(partes[1])
	if err != nil || iteraciones <= 0 {
		return false
	}
	sal, err := base64.RawStdEncoding.DecodeString(partes[2])
	if err != nil {
		return false
	}
	esperada, err := base64.RawStdEncoding.DecodeString(partes[3])
	if err != nil || len(esperada) == 0 {
		return false
	}
	obtenida, err := pbkdf2.Key(sha256.New, contrasena, sal, iteraciones, len(esperada))
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(obtenida, esperada) == 1
}
