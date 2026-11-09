package http

import (
	"bytes"
	"encoding/json"
	"log"
	nethttp "net/http"
	"strings"
	"testing"
)

// capturarLog redirige la salida del logger estándar a un buffer durante el test.
func capturarLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var salida bytes.Buffer
	anteriorSalida, anterioresFlags := log.Writer(), log.Flags()
	log.SetOutput(&salida)
	t.Cleanup(func() {
		log.SetOutput(anteriorSalida)
		log.SetFlags(anterioresFlags)
	})
	return &salida
}

// SC-006: ni los logs ni las respuestas de error contienen la contraseña, el
// hash ni el token, durante registro, ingreso, «quién soy» y cierre de sesión.
func TestAuth_ContrasenasHashesYTokensNoSeFiltranALogsNiAErrores(t *testing.T) {
	salida := capturarLog(t)
	h, credenciales := nuevoAuthHandlerDePrueba()

	const contrasena = "una-clave-larga"
	const contrasenaIncorrecta = "clave-equivocada-larga"

	var respuestas []string
	var errores []string

	reg := postRegistro(h, `{"nombre":"ana.perez","contrasena":"`+contrasena+`"}`)
	var sesion struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(reg.Body.Bytes(), &sesion); err != nil || sesion.Token == "" {
		t.Fatalf("preparar sesión: %v (cuerpo: %s)", err, reg.Body.String())
	}
	token := sesion.Token
	respuestas = append(respuestas, reg.Body.String())

	ing := postIngreso(h, `{"nombre":"ana.perez","contrasena":"`+contrasena+`"}`)
	var segunda struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(ing.Body.Bytes(), &segunda)
	respuestas = append(respuestas, ing.Body.String())

	respuestas = append(respuestas, getYo(h, "Bearer "+token).Body.String())

	// Caminos de error: ninguna respuesta puede devolver lo que se envió.
	errores = append(errores,
		postRegistro(h, `{"nombre":"ana.perez","contrasena":"`+contrasena+`"}`).Body.String(),
		postRegistro(h, `{"nombre":"otra.persona","contrasena":"corta"}`).Body.String(),
		postIngreso(h, `{"nombre":"ana.perez","contrasena":"`+contrasenaIncorrecta+`"}`).Body.String(),
		postIngreso(h, `{"nombre":"fantasma","contrasena":"`+contrasenaIncorrecta+`"}`).Body.String(),
		getYo(h, "Bearer "+strings.Repeat("A", 43)).Body.String(),
		getYo(h, "Bearer "+token+"x").Body.String(),
		postSalir(h, "Bearer no-es-un-token").Body.String(),
	)

	if salir := postSalir(h, "Bearer "+token); salir.Code != nethttp.StatusNoContent {
		t.Fatalf("salir: status = %d, se esperaba 204", salir.Code)
	}
	errores = append(errores, getYo(h, "Bearer "+token).Body.String())

	hash := credenciales.hashes["ana.perez"]
	if hash == "" {
		t.Fatal("no se guardó el hash")
	}
	secretos := map[string]string{
		"contraseña":                contrasena,
		"contraseña incorrecta":     contrasenaIncorrecta,
		"hash":                      hash,
		"token":                     token,
		"token del segundo ingreso": segunda.Token,
	}
	if segunda.Token == "" {
		t.Fatal("el ingreso no devolvió token")
	}

	registro := salida.String()
	for nombre, secreto := range secretos {
		if strings.Contains(registro, secreto) {
			t.Errorf("el log contiene la %s: %q", nombre, registro)
		}
	}
	for i, cuerpo := range errores {
		for nombre, secreto := range secretos {
			if strings.Contains(cuerpo, secreto) {
				t.Errorf("la respuesta de error %d contiene la %s: %s", i, nombre, cuerpo)
			}
		}
	}
	// Las respuestas exitosas solo pueden llevar su propio token (nunca la
	// contraseña ni el hash).
	for i, cuerpo := range respuestas {
		for _, nombre := range []string{"contraseña", "contraseña incorrecta", "hash"} {
			if strings.Contains(cuerpo, secretos[nombre]) {
				t.Errorf("la respuesta %d contiene la %s: %s", i, nombre, cuerpo)
			}
		}
	}
}
