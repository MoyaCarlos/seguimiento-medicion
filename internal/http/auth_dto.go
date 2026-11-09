package http

type registroRequest struct {
	Nombre     string `json:"nombre"`
	Contrasena string `json:"contrasena"`
}

type usuarioSesionResponse struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
}

// sesionResponse entrega el token en claro una sola vez; nunca incluye la
// contraseña ni su hash.
type sesionResponse struct {
	Usuario usuarioSesionResponse `json:"usuario"`
	Token   string                `json:"token"`
	VenceEn string                `json:"vence_en"`
}
