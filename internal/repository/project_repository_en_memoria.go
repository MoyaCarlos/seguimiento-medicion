package repository

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// ProjectRepositoryEnMemoria implementa ProjectRepository en memoria para que
// HU-05 y HU-13 usen un fake compartido en lugar de fakes propios. Valida la
// existencia de proyecto al vincular, rechaza duplicados y devuelve el nombre
// del integrante al listar.
type ProjectRepositoryEnMemoria struct {
	mu        sync.Mutex
	proyectos map[int64]domain.Project
	usuarios  map[int64]domain.User
	miembros  []domain.Membership
	// bajas guarda la fecha de baja de cada vínculo dado de baja (HU-16); sin
	// entrada, el integrante está activo. No se toca domain.Membership.
	bajas            map[claveIntegrante]time.Time
	siguiente        int64
	siguienteUsuario int64
}

// claveIntegrante identifica el vínculo de un usuario con un proyecto.
type claveIntegrante struct{ proyecto, usuario int64 }

func NewProjectRepositoryEnMemoria() *ProjectRepositoryEnMemoria {
	return &ProjectRepositoryEnMemoria{
		proyectos: map[int64]domain.Project{},
		usuarios:  map[int64]domain.User{},
		bajas:     map[claveIntegrante]time.Time{},
	}
}

// activo informa si el vínculo no tiene fecha de baja.
func (r *ProjectRepositoryEnMemoria) activo(m domain.Membership) bool {
	_, dadoDeBaja := r.bajas[claveIntegrante{proyecto: m.ProjectID, usuario: m.UserID}]
	return !dadoDeBaja
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

func (r *ProjectRepositoryEnMemoria) ListarPorIntegrante(_ context.Context, usuarioID int64) ([]domain.ProyectoConRol, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	proyectos := make([]domain.ProyectoConRol, 0)
	for _, m := range r.miembros {
		if m.UserID == usuarioID && r.activo(m) {
			proyectos = append(proyectos, domain.ProyectoConRol{Proyecto: r.proyectos[m.ProjectID], Rol: m.Role})
		}
	}
	sort.Slice(proyectos, func(i, j int) bool { return proyectos[i].Proyecto.ID < proyectos[j].Proyecto.ID })
	return proyectos, nil
}

func (r *ProjectRepositoryEnMemoria) Actualizar(_ context.Context, p domain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	actual, ok := r.proyectos[p.ID]
	if !ok {
		return domain.ErrProyectoNoEncontrado
	}
	p.FechaCierre = actual.FechaCierre
	r.proyectos[p.ID] = p
	return nil
}

// Reabrir deja vacía la fecha de cierre de un proyecto cerrado (#300).
func (r *ProjectRepositoryEnMemoria) Reabrir(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.proyectos[id]
	if !ok {
		return domain.ErrProyectoNoEncontrado
	}
	if p.FechaCierre == nil {
		return domain.ErrProyectoNoCerrado
	}
	p.FechaCierre = nil
	r.proyectos[id] = p
	return nil
}

// ProyectoCerrado informa si el proyecto tiene fecha de cierre (#300).
func (r *ProjectRepositoryEnMemoria) ProyectoCerrado(_ context.Context, id int64) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.proyectos[id].FechaCierre != nil, nil
}

// Cerrar guarda la fecha de cierre si el proyecto está abierto (#300). No
// conoce los Sprints: la regla del Sprint Activo la da Project.Cerrar.
func (r *ProjectRepositoryEnMemoria) Cerrar(_ context.Context, id int64, fecha time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.proyectos[id]
	if !ok {
		return domain.ErrProyectoNoEncontrado
	}
	if p.FechaCierre != nil {
		return domain.ErrProyectoYaCerrado
	}
	p.FechaCierre = &fecha
	r.proyectos[id] = p
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
	dadoDeBaja := -1
	for i, existente := range r.miembros {
		if existente.ProjectID == m.ProjectID && existente.UserID == m.UserID {
			if r.activo(existente) {
				return domain.Membership{}, ErrMiembroDuplicado
			}
			dadoDeBaja = i
			continue
		}
		if m.Role == domain.RolScrumMaster && existente.ProjectID == m.ProjectID && existente.Role == domain.RolScrumMaster && r.activo(existente) {
			return domain.Membership{}, domain.ErrScrumMasterYaExiste
		}
	}
	if dadoDeBaja >= 0 {
		// Reactivación (HU-16): mismo vínculo, rol pedido y CreadoEn original.
		delete(r.bajas, claveIntegrante{proyecto: m.ProjectID, usuario: m.UserID})
		r.miembros[dadoDeBaja].Role = m.Role
		return r.miembros[dadoDeBaja], nil
	}
	if m.CreadoEn.IsZero() {
		m.CreadoEn = time.Now().UTC()
	}
	r.miembros = append(r.miembros, m)
	return m, nil
}

// ListarIntegrantes devuelve los integrantes activos del proyecto (HU-16).
func (r *ProjectRepositoryEnMemoria) ListarIntegrantes(_ context.Context, proyectoID int64) ([]domain.Member, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.listarIntegrantes(proyectoID, true), nil
}

// ListarIntegrantesHistoricos devuelve los integrantes activos y los dados de
// baja del proyecto, en orden de alta (HU-16).
func (r *ProjectRepositoryEnMemoria) ListarIntegrantesHistoricos(_ context.Context, proyectoID int64) ([]domain.Member, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.listarIntegrantes(proyectoID, false), nil
}

func (r *ProjectRepositoryEnMemoria) listarIntegrantes(proyectoID int64, soloActivos bool) []domain.Member {
	var out []domain.Member
	for _, m := range r.miembros {
		if m.ProjectID == proyectoID && (!soloActivos || r.activo(m)) {
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
	return out
}

// DarDeBajaIntegrante guarda la fecha de baja de un integrante activo que no
// es el Scrum Master (HU-16). No descarta votos: viven en otro repositorio; esa
// garantía la cubre el repositorio SQLite.
func (r *ProjectRepositoryEnMemoria) DarDeBajaIntegrante(_ context.Context, proyectoID, usuarioID int64, fecha time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range r.miembros {
		if m.ProjectID != proyectoID || m.UserID != usuarioID || !r.activo(m) {
			continue
		}
		if m.Role == domain.RolScrumMaster {
			return domain.ErrScrumMasterNoSeQuita
		}
		r.bajas[claveIntegrante{proyecto: proyectoID, usuario: usuarioID}] = fecha
		return nil
	}
	return domain.ErrIntegranteNoEncontrado
}
