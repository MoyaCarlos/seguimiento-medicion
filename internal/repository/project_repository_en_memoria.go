package repository

import (
	"context"
	"sync"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// ProjectRepositoryEnMemoria implementa ProjectRepository en memoria para que
// HU-05 y HU-13 usen un fake compartido en lugar de fakes propios. Valida la
// existencia de proyecto al vincular, rechaza duplicados y devuelve el nombre
// del integrante al listar.
type ProjectRepositoryEnMemoria struct {
	mu               sync.Mutex
	proyectos        map[int64]domain.Project
	usuarios         map[int64]domain.User
	miembros         []domain.Membership
	siguiente        int64
	siguienteUsuario int64
}

func NewProjectRepositoryEnMemoria() *ProjectRepositoryEnMemoria {
	return &ProjectRepositoryEnMemoria{
		proyectos: map[int64]domain.Project{},
		usuarios:  map[int64]domain.User{},
	}
}

// GuardarUsuario registra un integrante con su ID asignado, para que las
// vinculaciones y el listado puedan resolver el nombre.
func (r *ProjectRepositoryEnMemoria) GuardarUsuario(_ context.Context, u domain.User) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u.CreadoEn.IsZero() {
		u.CreadoEn = time.Now().UTC()
	}
	r.siguienteUsuario++
	u.ID = r.siguienteUsuario
	r.usuarios[u.ID] = u
	return u, nil
}

func (r *ProjectRepositoryEnMemoria) Guardar(_ context.Context, p domain.Project) (domain.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.guardar(p), nil
}

func (r *ProjectRepositoryEnMemoria) guardar(p domain.Project) domain.Project {
	if p.CreadoEn.IsZero() {
		p.CreadoEn = time.Now().UTC()
	}
	r.siguiente++
	p.ID = r.siguiente
	r.proyectos[p.ID] = p
	return p
}

func (r *ProjectRepositoryEnMemoria) GuardarConScrumMaster(_ context.Context, p domain.Project, creadorID int64) (domain.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	guardado := r.guardar(p)
	if _, err := r.agregarIntegrante(domain.Membership{ProjectID: guardado.ID, UserID: creadorID, Role: domain.RolScrumMaster}); err != nil {
		delete(r.proyectos, guardado.ID)
		return domain.Project{}, err
	}
	return guardado, nil
}

func (r *ProjectRepositoryEnMemoria) ObtenerPorID(_ context.Context, id int64) (domain.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.proyectos[id]
	if !ok {
		return domain.Project{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (r *ProjectRepositoryEnMemoria) Actualizar(_ context.Context, p domain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.proyectos[p.ID]; !ok {
		return domain.ErrProyectoNoEncontrado
	}
	r.proyectos[p.ID] = p
	return nil
}

func (r *ProjectRepositoryEnMemoria) AgregarIntegrante(_ context.Context, m domain.Membership) (domain.Membership, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.agregarIntegrante(m)
}

func (r *ProjectRepositoryEnMemoria) agregarIntegrante(m domain.Membership) (domain.Membership, error) {
	if _, ok := r.proyectos[m.ProjectID]; !ok {
		return domain.Membership{}, domain.ErrProyectoNoEncontrado
	}
	for _, existente := range r.miembros {
		if existente.ProjectID == m.ProjectID && existente.UserID == m.UserID {
			return domain.Membership{}, ErrMiembroDuplicado
		}
	}
	if m.CreadoEn.IsZero() {
		m.CreadoEn = time.Now().UTC()
	}
	r.miembros = append(r.miembros, m)
	return m, nil
}

func (r *ProjectRepositoryEnMemoria) ListarIntegrantes(_ context.Context, proyectoID int64) ([]domain.Member, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []domain.Member
	for _, m := range r.miembros {
		if m.ProjectID == proyectoID {
			u := r.usuarios[m.UserID]
			out = append(out, domain.Member{
				UserID:    m.UserID,
				ProjectID: m.ProjectID,
				Nombre:    u.Nombre,
				Role:      m.Role,
				CreadoEn:  m.CreadoEn,
			})
		}
	}
	return out, nil
}
