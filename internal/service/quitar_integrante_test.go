package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// equipoParaQuitar arma el proyecto 1 con Carlos Scrum Master y Candela y
// Jimena Product Builders, guardados como membresías (HU-16).
func equipoParaQuitar() *fakeProjectRepository {
	return &fakeProjectRepository{
		proyectos: map[int64]domain.Project{1: {ID: 1, Nombre: "Proyecto"}},
		miembros: []domain.Membership{
			{ProjectID: 1, UserID: carlos, Role: domain.RolScrumMaster},
			{ProjectID: 1, UserID: candela, Role: domain.RolProductBuilder},
			{ProjectID: 1, UserID: jimena, Role: domain.RolProductBuilder},
		},
	}
}

func idsActivos(t *testing.T, proyectos *fakeProjectRepository) []int64 {
	t.Helper()
	activos, err := proyectos.ListarIntegrantes(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListarIntegrantes: %v", err)
	}
	ids := make([]int64, len(activos))
	for i, m := range activos {
		ids[i] = m.UserID
	}
	return ids
}

func TestQuitarIntegrante_ScrumMasterQuitaAUnProductBuilder(t *testing.T) {
	proyectos := equipoParaQuitar()
	hoy := time.Date(2026, 11, 12, 15, 30, 0, 0, time.UTC)
	servicio := NewQuitarIntegrante(proyectos, func() time.Time { return hoy })

	err := servicio.Ejecutar(context.Background(), QuitarIntegranteInput{ProyectoID: 1, UsuarioID: candela, ActorID: carlos})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if got := idsActivos(t, proyectos); len(got) != 2 || got[0] != carlos || got[1] != jimena {
		t.Errorf("se esperaba el equipo [%d %d], se obtuvo %v", carlos, jimena, got)
	}
	if fecha := proyectos.bajas[[2]int64{1, candela}]; !fecha.Equal(hoy) {
		t.Errorf("se esperaba la fecha del reloj %v, se obtuvo %v", hoy, fecha)
	}
}

func TestQuitarIntegrante_Errores(t *testing.T) {
	casos := []struct {
		nombre    string
		input     QuitarIntegranteInput
		preparar  func(*fakeProjectRepository)
		esperado  error
		escenario string
	}{
		{"proyecto inexistente", QuitarIntegranteInput{ProyectoID: 7, UsuarioID: candela, ActorID: carlos}, nil, domain.ErrProyectoNoEncontrado, "US4"},
		{"actor no integrante", QuitarIntegranteInput{ProyectoID: 1, UsuarioID: candela, ActorID: ajeno}, nil, domain.ErrNoEsIntegrante, "HU-21"},
		{"actor Product Builder", QuitarIntegranteInput{ProyectoID: 1, UsuarioID: jimena, ActorID: candela}, nil, domain.ErrNoEsScrumMaster, "HU-21"},
		{"usuario que no integra", QuitarIntegranteInput{ProyectoID: 1, UsuarioID: ajeno, ActorID: carlos}, nil, domain.ErrIntegranteNoEncontrado, "US4"},
		{"usuario ya dado de baja", QuitarIntegranteInput{ProyectoID: 1, UsuarioID: candela, ActorID: carlos},
			func(f *fakeProjectRepository) { f.bajas = map[[2]int64]time.Time{{1, candela}: time.Now()} },
			domain.ErrIntegranteNoEncontrado, "US4"},
		{"el Scrum Master se quita a sí mismo", QuitarIntegranteInput{ProyectoID: 1, UsuarioID: carlos, ActorID: carlos}, nil, domain.ErrScrumMasterNoSeQuita, "US2"},
		{"orden: proyecto inexistente y actor no integrante → 404", QuitarIntegranteInput{ProyectoID: 7, UsuarioID: candela, ActorID: ajeno}, nil, domain.ErrProyectoNoEncontrado, "FR-013"},
		{"orden: actor Product Builder y objetivo inexistente → 403", QuitarIntegranteInput{ProyectoID: 1, UsuarioID: ajeno, ActorID: candela}, nil, domain.ErrNoEsScrumMaster, "FR-013"},
		{"proyecto cerrado", QuitarIntegranteInput{ProyectoID: 1, UsuarioID: candela, ActorID: carlos},
			func(f *fakeProjectRepository) { f.cerrados = map[int64]bool{1: true} },
			domain.ErrProyectoCerrado, "#300"},
		{"orden: proyecto cerrado antes del 404 de integrante", QuitarIntegranteInput{ProyectoID: 1, UsuarioID: ajeno, ActorID: carlos},
			func(f *fakeProjectRepository) { f.cerrados = map[int64]bool{1: true} },
			domain.ErrProyectoCerrado, "FR-013"},
		{"orden: no Scrum Master antes del proyecto cerrado", QuitarIntegranteInput{ProyectoID: 1, UsuarioID: jimena, ActorID: candela},
			func(f *fakeProjectRepository) { f.cerrados = map[int64]bool{1: true} },
			domain.ErrNoEsScrumMaster, "FR-013"},
	}
	for _, c := range casos {
		t.Run(c.escenario+" "+c.nombre, func(t *testing.T) {
			proyectos := equipoParaQuitar()
			if c.preparar != nil {
				c.preparar(proyectos)
			}
			antes := idsActivos(t, proyectos)
			servicio := NewQuitarIntegrante(proyectos, time.Now)

			err := servicio.Ejecutar(context.Background(), c.input)
			if !errors.Is(err, c.esperado) {
				t.Fatalf("se esperaba %v, se obtuvo %v", c.esperado, err)
			}
			if despues := idsActivos(t, proyectos); len(despues) != len(antes) {
				t.Errorf("el equipo no debe cambiar: antes %v, después %v", antes, despues)
			}
		})
	}
}
