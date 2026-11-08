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
	Rol        domain.Role
}

// AsignarIntegrante vincula a un integrante (existente o nuevo) con un rol.
type AsignarIntegrante struct {
	proyectos repository.ProjectRepository
	usuarios  repository.UserRepository
}

func NewAsignarIntegrante(proyectos repository.ProjectRepository, usuarios repository.UserRepository) *AsignarIntegrante {
	return &AsignarIntegrante{proyectos: proyectos, usuarios: usuarios}
}

// Ejecutar verifica el proyecto, resuelve el integrante y persiste la
// vinculación (el rol se valida en el dominio), rechazando duplicados.
func (s *AsignarIntegrante) Ejecutar(ctx context.Context, input AsignarIntegranteInput) (domain.Member, error) {
	proyecto, err := s.proyectos.ObtenerPorID(ctx, input.ProyectoID)
	if err != nil {
		return domain.Member{}, err
	}

	if !input.Rol.Valido() {
		return domain.Member{}, domain.ValidationError{Campo: "rol", Mensaje: "el rol debe ser scrum_master o product_builder"}
	}

	candidato, err := domain.NewUser(input.Nombre)
	if err != nil {
		return domain.Member{}, err
	}
	if input.Rol == domain.RolScrumMaster {
		if err := s.verificarScrumMasterLibre(ctx, proyecto.ID, candidato); err != nil {
			return domain.Member{}, err
		}
	}
	usuario, err := resolverUsuario(ctx, s.usuarios, candidato)
	if err != nil {
		return domain.Member{}, err
	}

	membresia, err := domain.NewMembership(proyecto.ID, usuario.ID, input.Rol)
	if err != nil {
		return domain.Member{}, err
	}
	guardada, err := s.proyectos.AgregarIntegrante(ctx, membresia)
	if err != nil {
		if errors.Is(err, repository.ErrMiembroDuplicado) {
			return domain.Member{}, domain.ValidationError{Campo: "integrante", Mensaje: "el integrante ya pertenece al proyecto"}
		}
		return domain.Member{}, err
	}

	return domain.Member{
		UserID:    usuario.ID,
		ProjectID: proyecto.ID,
		Nombre:    usuario.Nombre,
		Role:      input.Rol,
		CreadoEn:  guardada.CreadoEn,
	}, nil
}

// verificarScrumMasterLibre devuelve ErrScrumMasterYaExiste si el proyecto ya
// tiene Scrum Master. Si el candidato ya es integrante no se rechaza acá: el
// alta falla después como duplicado, que tiene prioridad. No crea usuarios.
func (s *AsignarIntegrante) verificarScrumMasterLibre(ctx context.Context, proyectoID int64, candidato domain.User) error {
	integrantes, err := s.proyectos.ListarIntegrantes(ctx, proyectoID)
	if err != nil {
		return err
	}
	if !slices.ContainsFunc(integrantes, func(m domain.Member) bool { return m.Role == domain.RolScrumMaster }) {
		return nil
	}
	existente, err := s.usuarios.ObtenerPorNombreNormalizado(ctx, candidato.NombreNormalizado)
	if errors.Is(err, domain.ErrNoEncontrado) {
		return domain.ErrScrumMasterYaExiste
	}
	if err != nil {
		return err
	}
	yaIntegrante := slices.ContainsFunc(integrantes, func(m domain.Member) bool { return m.UserID == existente.ID })
	if yaIntegrante {
		return nil
	}
	return domain.ErrScrumMasterYaExiste
}
