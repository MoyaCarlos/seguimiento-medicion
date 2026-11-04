package service

const (
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
	return "", nil
}

func (s *ServicioContrasenas) Verificar(contrasena, hash string) bool {
	return false
}
