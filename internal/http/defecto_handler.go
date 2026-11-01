package http

import (
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// DefectoHandler es el adaptador de entrada para los defectos de un proyecto (HU-08).
type DefectoHandler struct {
	obtenerProyecto *service.ObtenerProyecto
	registrar       *service.RegistrarDefecto
}

func NewDefectoHandler(obtenerProyecto *service.ObtenerProyecto, registrar *service.RegistrarDefecto) *DefectoHandler {
	return &DefectoHandler{obtenerProyecto: obtenerProyecto, registrar: registrar}
}

// Registrar atiende POST /projects/{id}/defectos.
func (h *DefectoHandler) Registrar(w nethttp.ResponseWriter, r *nethttp.Request) {
}
