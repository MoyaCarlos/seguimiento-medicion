package repository

import (
	"context"
	"sync"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// ProjectRepositoryEnMemoria implementa ProjectRepository en memoria para que
// HU-05 y HU-13 usen un fake compartido en lugar de fakes propios.
type ProjectRepositoryEnMemoria struct {
	mu        sync.Mutex
	proyectos map[int64]domain.Project
	miembros  []domain.Membership
	siguiente int64
}

func NewProjectRepositoryEnMemoria() *ProjectRepositoryEnMemoria {
	return &ProjectRepositoryEnMemoria{proyectos: map[int64]domain.Project{}}
}

func (r *ProjectRepositoryEnMemoria) Guardar(_ context.Context, p domain.Project) (domain.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.CreadoEn.IsZero() {
		p.CreadoEn = time.Now().UTC()
	}
	r.siguiente++
	p.ID = r.siguiente
	r.proyectos[p.ID] = p
	return p, nil
}

func (r *ProjectRepositoryEnMemoria) GuardarConScrumMaster(ctx context.Context, p domain.Project, creadorID int64) (domain.Project, error) {
	guardado, err := r.Guardar(ctx, p)
	if err != nil {
		return domain.Project{}, err
	}
	if err := r.AgregarIntegrante(ctx, domain.Membership{ProjectID: guardado.ID, UserID: creadorID, Role: domain.RolScrumMaster}); err != nil {
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

func (r *ProjectRepositoryEnMemoria) AgregarIntegrante(_ context.Context, m domain.Membership) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existente := range r.miembros {
		if existente.ProjectID == m.ProjectID && existente.UserID == m.UserID {
			return ErrMiembroDuplicado
		}
	}
	if m.CreadoEn.IsZero() {
		m.CreadoEn = time.Now().UTC()
	}
	r.miembros = append(r.miembros, m)
	return nil
}

func (r *ProjectRepositoryEnMemoria) ListarIntegrantes(_ context.Context, proyectoID int64) ([]domain.Member, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []domain.Member
	for _, m := range r.miembros {
		if m.ProjectID == proyectoID {
			out = append(out, domain.Member{
				UserID:    m.UserID,
				ProjectID: m.ProjectID,
				Role:      m.Role,
				CreadoEn:  m.CreadoEn,
			})
		}
	}
	return out, nil
}
