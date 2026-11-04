package domain

import "time"

// DuracionSesion es lo que dura una sesión desde que se emite (HU-20).
const DuracionSesion = 7 * 24 * time.Hour

// Sesion es el ingreso de un usuario, identificado por la huella de su token.
type Sesion struct {
	ID        int64
	UsuarioID int64
	CreadoEn  time.Time
	VenceEn   time.Time
	CerradaEn *time.Time
}

// SesionConUsuario es una sesión junto con su usuario.
type SesionConUsuario struct {
	Sesion
	Usuario User
}

// Vigente indica si la sesión sirve para autenticar en el instante dado: no
// está cerrada y todavía no llegó su vencimiento.
func (s Sesion) Vigente(ahora time.Time) bool {
	return s.CerradaEn == nil && ahora.Before(s.VenceEn)
}
