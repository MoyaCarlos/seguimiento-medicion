package service

import (
	"context"
	"errors"
	"slices"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// AsignarIntegranteInput son los datos de entrada para vincular a un integrante.
type AsignarIntegranteInput struct {
	ProyectoID int64
	Nombre     string
	// Rol vacío es product_builder (HU-21).
	Rol domain.Role
	// ActorID es quien configura el equipo (la sesión, HU-21); debe ser el
	// Scrum Master. Nombre es la persona que se agrega.
	ActorID int64
}

// AsignarIntegrante suma al equipo a un usuario registrado (HU-21): no crea
// usuarios.
type AsignarIntegrante struct {
	proyectos    repository.ProjectRepository
	credenciales repository.BuscadorCredenciales
}

func NewAsignarIntegrante(proyectos repository.ProjectRepository, credenciales repository.BuscadorCredenciales) *AsignarIntegrante {
	return &AsignarIntegrante{proyectos: proyectos, credenciales: credenciales}
}

// Ejecutar verifica el proyecto y que actúe su Scrum Master, valida el rol,
// busca al usuario registrado y persiste la vinculación, rechazando duplicados.
func (s *AsignarIntegrante) Ejecutar(ctx context.Context, input AsignarIntegranteInput) (domain.Member, error) {
	proyecto, err := s.proyectos.ObtenerPorID(ctx, input.ProyectoID)
	if err != nil {
		return domain.Member{}, err
	}
	if err := exigirScrumMaster(ctx, s.proyectos, proyecto.ID, input.ActorID); err != nil {
		return domain.Member{}, err
	}

	rol := input.Rol
	if rol == "" {
		rol = domain.RolProductBuilder
	}
	if !rol.Valido() {
		return domain.Member{}, domain.ValidationError{Campo: "rol", Mensaje: "el rol debe ser scrum_master o product_builder"}
	}

	candidato, err := domain.NewUser(input.Nombre)
	if err != nil {
		return domain.Member{}, err
	}
	if rol == domain.RolScrumMaster {
		if err := s.verificarScrumMasterLibre(ctx, proyecto.ID, candidato); err != nil {
			return domain.Member{}, err
		}
	}
	usuario, err := s.usuarioRegistrado(ctx, candidato)
	if err != nil {
		return domain.Member{}, err
	}

	membresia, err := domain.NewMembership(proyecto.ID, usuario.ID, rol)
	if err != nil {
		return domain.Member{}, err
	}
	guardada, err := s.proyectos.AgregarIntegrante(ctx, membresia)
	if err != nil {
		if errors.Is(err, repository.ErrMiembroDuplicado) {
			return domain.Member{}, errIntegranteDuplicado()
		}
		return domain.Member{}, err
	}

	return domain.Member{
		UserID:    usuario.ID,
		ProjectID: proyecto.ID,
		Nombre:    usuario.Nombre,
		Role:      rol,
		CreadoEn:  guardada.CreadoEn,
	}, nil
}

// usuarioRegistrado busca al candidato entre los usuarios con cuenta: quien no
// existe y quien no tiene contraseña dan el mismo ErrUsuarioNoEncontrado.
func (s *AsignarIntegrante) usuarioRegistrado(ctx context.Context, candidato domain.User) (domain.User, error) {
	usuario, hash, err := s.credenciales.ObtenerPorNombreNormalizado(ctx, candidato.NombreNormalizado)
	if errors.Is(err, domain.ErrNoEncontrado) || (err == nil && hash == "") {
		return domain.User{}, domain.ErrUsuarioNoEncontrado
	}
	if err != nil {
		return domain.User{}, err
	}
	return usuario, nil
}

func errIntegranteDuplicado() error {
	return domain.ValidationError{Campo: "integrante", Mensaje: "el integrante ya pertenece al proyecto"}
}

// verificarScrumMasterLibre devuelve ErrScrumMasterYaExiste si el proyecto ya
// tiene Scrum Master. Si el candidato ya es integrante devuelve el duplicado,
// que tiene prioridad (SQLite informaría el índice en vez de la clave primaria).
func (s *AsignarIntegrante) verificarScrumMasterLibre(ctx context.Context, proyectoID int64, candidato domain.User) error {
	integrantes, err := s.proyectos.ListarIntegrantes(ctx, proyectoID)
	if err != nil {
		return err
	}
	if !slices.ContainsFunc(integrantes, func(m domain.Member) bool { return m.Role == domain.RolScrumMaster }) {
		return nil
	}
	existente, _, err := s.credenciales.ObtenerPorNombreNormalizado(ctx, candidato.NombreNormalizado)
	if errors.Is(err, domain.ErrNoEncontrado) {
		return domain.ErrScrumMasterYaExiste
	}
	if err != nil {
		return err
	}
	if slices.ContainsFunc(integrantes, func(m domain.Member) bool { return m.UserID == existente.ID }) {
		return errIntegranteDuplicado()
	}
	return domain.ErrScrumMasterYaExiste
}
