package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

var relojDeEmision = func() time.Time { return time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC) }

func TestEmisorSesiones_EmitirDevuelveUnTokenBase64URLDe43Caracteres(t *testing.T) {
	emisor := NewEmisorSesiones(&fakeSesiones{}, relojDeEmision)

	token, _, err := emisor.Emitir(context.Background(), 7)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if !regexp.MustCompile(`^[A-Za-z0-9_-]{43}$`).MatchString(token) {
		t.Errorf("se esperaba base64url sin relleno de 43 caracteres, se obtuvo %q", token)
	}
}

func TestEmisorSesiones_EmitirGuardaLaHuellaYNuncaElToken(t *testing.T) {
	sesiones := &fakeSesiones{}
	emisor := NewEmisorSesiones(sesiones, relojDeEmision)

	token, _, err := emisor.Emitir(context.Background(), 7)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	suma := sha256.Sum256([]byte(token))
	huella := hex.EncodeToString(suma[:])
	if len(sesiones.sesiones) != 1 {
		t.Fatalf("se esperaba 1 sesión guardada, se obtuvo %d", len(sesiones.sesiones))
	}
	if _, ok := sesiones.sesiones[huella]; !ok {
		t.Errorf("se esperaba guardar la huella %q, se guardó %v", huella, sesiones.sesiones)
	}
	if _, ok := sesiones.sesiones[token]; ok {
		t.Error("el token en claro no debe guardarse")
	}
	if got := HuellaDeToken(token); got != huella {
		t.Errorf("HuellaDeToken = %q, se esperaba %q", got, huella)
	}
}

func TestEmisorSesiones_EmitirFijaCreadoEnYVenceEnSieteDiasDespues(t *testing.T) {
	emisor := NewEmisorSesiones(&fakeSesiones{}, relojDeEmision)

	_, sesion, err := emisor.Emitir(context.Background(), 7)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	ahora := relojDeEmision()
	if sesion.UsuarioID != 7 || sesion.ID == 0 {
		t.Errorf("sesión inesperada: %+v", sesion)
	}
	if !sesion.CreadoEn.Equal(ahora) {
		t.Errorf("CreadoEn = %v, se esperaba %v", sesion.CreadoEn, ahora)
	}
	if !sesion.VenceEn.Equal(ahora.Add(7*24*time.Hour)) || sesion.VenceEn.Sub(sesion.CreadoEn) != domain.DuracionSesion {
		t.Errorf("VenceEn = %v, se esperaba 7 días después de %v", sesion.VenceEn, ahora)
	}
}

func TestEmisorSesiones_DosEmisionesDanTokensDistintos(t *testing.T) {
	emisor := NewEmisorSesiones(&fakeSesiones{}, relojDeEmision)

	t1, _, _ := emisor.Emitir(context.Background(), 7)
	t2, _, _ := emisor.Emitir(context.Background(), 7)

	if t1 == "" || t1 == t2 {
		t.Errorf("los tokens deberían ser distintos: %q y %q", t1, t2)
	}
}
