package service

import (
	"context"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

type fakeProjectRepository struct {
	proyectos     map[int64]domain.Project
	miembros      []domain.Membership
	listaMiembros []domain.Member
	err           error
	creates       int
	updates       int
}

func (f *fakeProjectRepository) Guardar(_ context.Context, p domain.Project) (domain.Project, error) {
	if f.err != nil {
		return domain.Project{}, f.err
	}
	f.creates++
	if f.proyectos == nil {
		f.proyectos = map[int64]domain.Project{}
	}
	if p.ID == 0 {
		p.ID = int64(f.creates)
	}
	p.CreadoEn = time.Now().UTC()
	f.proyectos[p.ID] = p
	return p, nil
}

func (f *fakeProjectRepository) GuardarConScrumMaster(ctx context.Context, p domain.Project, creadorID int64) (domain.Project, error) {
	guardado, err := f.Guardar(ctx, p)
	if err != nil {
		return domain.Project{}, err
	}
	if _, err := f.AgregarIntegrante(ctx, domain.Membership{ProjectID: guardado.ID, UserID: creadorID, Role: domain.RolScrumMaster}); err != nil {
		return domain.Project{}, err
	}
	return guardado, nil
}

func (f *fakeProjectRepository) ObtenerPorID(_ context.Context, id int64) (domain.Project, error) {
	if f.err != nil {
		return domain.Project{}, f.err
	}
	p, ok := f.proyectos[id]
	if !ok {
		return domain.Project{}, domain.ErrProyectoNoEncontrado
	}
	return p, nil
}

func (f *fakeProjectRepository) Listar(_ context.Context) ([]domain.Project, error) {
	if f.err != nil {
		return nil, f.err
	}
	proyectos := make([]domain.Project, 0, len(f.proyectos))
	for _, p := range f.proyectos {
		proyectos = append(proyectos, p)
	}
	sort.Slice(proyectos, func(i, j int) bool { return proyectos[i].ID < proyectos[j].ID })
	return proyectos, nil
}

func (f *fakeProjectRepository) Actualizar(_ context.Context, p domain.Project) error {
	if f.err != nil {
		return f.err
	}
	if _, ok := f.proyectos[p.ID]; !ok {
		return domain.ErrProyectoNoEncontrado
	}
	f.updates++
	f.proyectos[p.ID] = p
	return nil
}

func (f *fakeProjectRepository) AgregarIntegrante(_ context.Context, m domain.Membership) (domain.Membership, error) {
	if f.err != nil {
		return domain.Membership{}, f.err
	}
	for _, existente := range f.miembros {
		if existente.ProjectID == m.ProjectID && existente.UserID == m.UserID {
			return domain.Membership{}, repository.ErrMiembroDuplicado
		}
	}
	if m.CreadoEn.IsZero() {
		m.CreadoEn = time.Now().UTC()
	}
	f.miembros = append(f.miembros, m)
	return m, nil
}

func (f *fakeProjectRepository) ListarIntegrantes(_ context.Context, _ int64) ([]domain.Member, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.listaMiembros, nil
}

type fakeUserRepository struct {
	usuarios map[string]domain.User
	creates  int
	err      error
}

func (f *fakeUserRepository) ObtenerPorNombreNormalizado(_ context.Context, normalized string) (*domain.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	u, ok := f.usuarios[normalized]
	if !ok {
		return nil, domain.ErrNoEncontrado
	}
	return &u, nil
}

func (f *fakeUserRepository) Guardar(_ context.Context, u domain.User) (domain.User, error) {
	if f.err != nil {
		return domain.User{}, f.err
	}
	if f.usuarios == nil {
		f.usuarios = map[string]domain.User{}
	}
	if _, ok := f.usuarios[u.NombreNormalizado]; ok {
		return domain.User{}, repository.ErrUsuarioDuplicado
	}
	f.creates++
	u.ID = int64(f.creates)
	u.CreadoEn = time.Now().UTC()
	f.usuarios[u.NombreNormalizado] = u
	return u, nil
}

type fakeProductBacklogRepository struct {
	historias []domain.BacklogItem
	err       error
	llamadas  int
}

func (f *fakeProductBacklogRepository) ListarSinSprint(_ context.Context, _ int64) ([]domain.BacklogItem, error) {
	f.llamadas++
	if f.err != nil {
		return nil, f.err
	}
	return f.historias, nil
}

type fakeHistoriaRepository struct {
	historias       map[int64]domain.BacklogItem
	errActualizar   error
	actualizaciones []domain.BacklogItem
}

func (f *fakeHistoriaRepository) ObtenerPorID(_ context.Context, id int64) (domain.BacklogItem, error) {
	h, ok := f.historias[id]
	if !ok {
		return domain.BacklogItem{}, domain.ErrHistoriaNoEncontrada
	}
	return h, nil
}

func (f *fakeHistoriaRepository) Actualizar(_ context.Context, item domain.BacklogItem) error {
	if f.errActualizar != nil {
		return f.errActualizar
	}
	f.actualizaciones = append(f.actualizaciones, item)
	f.historias[item.ID] = item
	return nil
}

type fakeRegistroHorasRepository struct {
	registros    []domain.RegistroHoras
	listados     int
	errRegistrar error
}

func (f *fakeRegistroHorasRepository) Registrar(_ context.Context, r domain.RegistroHoras) (domain.RegistroHoras, error) {
	if f.errRegistrar != nil {
		return domain.RegistroHoras{}, f.errRegistrar
	}
	r.ID = int64(len(f.registros) + 1)
	f.registros = append(f.registros, r)
	return r, nil
}

func (f *fakeRegistroHorasRepository) ListarPorHistoria(_ context.Context, historiaID int64) ([]domain.RegistroHorasDetalle, error) {
	f.listados++
	detalle := make([]domain.RegistroHorasDetalle, 0)
	for _, r := range f.registros {
		if r.HistoriaID == historiaID {
			detalle = append(detalle, domain.RegistroHorasDetalle{RegistroHoras: r})
		}
	}
	return detalle, nil
}

type fakeIntegrantesRepository struct {
	miembros []domain.Member
}

func (f *fakeIntegrantesRepository) ListarIntegrantes(_ context.Context, proyectoID int64) ([]domain.Member, error) {
	var resultado []domain.Member
	for _, m := range f.miembros {
		if m.ProjectID == proyectoID {
			resultado = append(resultado, m)
		}
	}
	return resultado, nil
}

// fakeSprintsQueCuentan envuelve el repositorio en memoria y cuenta las lecturas.
type fakeSprintsQueCuentan struct {
	*repository.SprintRepositoryEnMemoria
	lecturas int
	listados int
}

func (f *fakeSprintsQueCuentan) ListarPorProyecto(ctx context.Context, proyectoID int64) ([]domain.Sprint, error) {
	f.listados++
	return f.SprintRepositoryEnMemoria.ListarPorProyecto(ctx, proyectoID)
}

func (f *fakeSprintsQueCuentan) ObtenerPorID(ctx context.Context, id int64) (domain.Sprint, error) {
	f.lecturas++
	return f.SprintRepositoryEnMemoria.ObtenerPorID(ctx, id)
}

// fakeVotacionRepository guarda la ronda actual de cada historia y las anteriores,
// y simula las escrituras condicionales del adaptador SQLite (research.md R6).
// Devuelve copias, así un cambio en la ronda del service no llega acá sin pasar
// por GuardarVoto.
type fakeVotacionRepository struct {
	rondas     map[int64]domain.RondaVotacion
	anteriores []domain.RondaVotacion
	err        error
}

func (f *fakeVotacionRepository) CrearRonda(_ context.Context, r domain.RondaVotacion) (domain.RondaVotacion, error) {
	if f.err != nil {
		return domain.RondaVotacion{}, f.err
	}
	actual, ok := f.rondas[r.HistoriaID]
	if ok && r.Numero <= actual.Numero {
		return domain.RondaVotacion{}, domain.ErrVotacionYaIniciada
	}
	if f.rondas == nil {
		f.rondas = map[int64]domain.RondaVotacion{}
	}
	if ok {
		f.anteriores = append(f.anteriores, actual)
	}
	r.ID = f.siguienteID()
	f.rondas[r.HistoriaID] = copiaDeRonda(r)
	return copiaDeRonda(r), nil
}

func (f *fakeVotacionRepository) siguienteID() int64 {
	var mayor int64
	for _, r := range f.rondas {
		mayor = max(mayor, r.ID)
	}
	for _, r := range f.anteriores {
		mayor = max(mayor, r.ID)
	}
	return mayor + 1
}

func (f *fakeVotacionRepository) ObtenerRondaActual(_ context.Context, historiaID int64) (domain.RondaVotacion, error) {
	if f.err != nil {
		return domain.RondaVotacion{}, f.err
	}
	r, ok := f.rondas[historiaID]
	if !ok {
		return domain.RondaVotacion{}, domain.ErrVotacionNoIniciada
	}
	return copiaDeRonda(r), nil
}

func (f *fakeVotacionRepository) RondaActualRevelada(ctx context.Context, historiaID int64) (bool, error) {
	r, err := f.ObtenerRondaActual(ctx, historiaID)
	if err != nil {
		return false, err
	}
	return r.Revelada(), nil
}

func (f *fakeVotacionRepository) GuardarVoto(_ context.Context, rondaID int64, v domain.Voto) error {
	if f.err != nil {
		return f.err
	}
	r, ok := f.rondaAbierta(rondaID)
	if !ok {
		return domain.ErrRondaRevelada
	}
	if i := slices.IndexFunc(r.Votos, func(x domain.Voto) bool { return x.UsuarioID == v.UsuarioID }); i >= 0 {
		r.Votos[i] = v
	} else {
		r.Votos = append(r.Votos, v)
	}
	f.rondas[r.HistoriaID] = r
	return nil
}

func (f *fakeVotacionRepository) MarcarRevelada(_ context.Context, rondaID int64, en time.Time) error {
	if f.err != nil {
		return f.err
	}
	r, ok := f.rondaAbierta(rondaID)
	if !ok {
		return domain.ErrRondaRevelada
	}
	r.ReveladaEn = &en
	f.rondas[r.HistoriaID] = r
	return nil
}

func (f *fakeVotacionRepository) ListarRondasReveladas(_ context.Context, historiaID int64) ([]domain.RondaVotacion, error) {
	if f.err != nil {
		return nil, f.err
	}
	rondas := make([]domain.RondaVotacion, 0)
	for _, r := range f.anteriores {
		if r.HistoriaID == historiaID && r.Revelada() {
			rondas = append(rondas, copiaDeRonda(r))
		}
	}
	if actual, ok := f.rondas[historiaID]; ok && actual.Revelada() {
		rondas = append(rondas, copiaDeRonda(actual))
	}
	return rondas, nil
}

func (f *fakeVotacionRepository) rondaAbierta(rondaID int64) (domain.RondaVotacion, bool) {
	for _, r := range f.rondas {
		if r.ID == rondaID && r.ReveladaEn == nil {
			return copiaDeRonda(r), true
		}
	}
	return domain.RondaVotacion{}, false
}

func copiaDeRonda(r domain.RondaVotacion) domain.RondaVotacion {
	r.Votos = slices.Clone(r.Votos)
	return r
}

// fakeAsignacionRepository simula la asignación atómica del repositorio real:
// el mutex hace indivisible "si no tiene Sprint, cargarlo", como el UPDATE
// condicional de SQLite.
type fakeAsignacionRepository struct {
	mu           sync.Mutex
	historias    map[int64]domain.BacklogItem
	asignaciones int
}

func (f *fakeAsignacionRepository) ObtenerPorID(_ context.Context, id int64) (domain.BacklogItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	h, ok := f.historias[id]
	if !ok {
		return domain.BacklogItem{}, domain.ErrHistoriaNoEncontrada
	}
	return h, nil
}

func (f *fakeAsignacionRepository) AsignarASprint(_ context.Context, historiaID, sprintID int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.asignaciones++
	h := f.historias[historiaID]
	if h.SprintID != nil {
		return domain.ErrHistoriaYaAsignada
	}
	h.SprintID = &sprintID
	f.historias[historiaID] = h
	return nil
}

func (f *fakeAsignacionRepository) llamadasAAsignar() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.asignaciones
}

// fakeDefectoRepository guarda los defectos en memoria, asigna IDs y cuenta
// las escrituras.
type fakeDefectoRepository struct {
	defectos     []domain.Defecto
	titulos      map[int64]string
	registros    int
	resoluciones int
	listados     int
	errRegistrar error
	errResolver  error
}

func (f *fakeDefectoRepository) Registrar(_ context.Context, d domain.Defecto) (domain.Defecto, error) {
	f.registros++
	if f.errRegistrar != nil {
		return domain.Defecto{}, f.errRegistrar
	}
	d.ID = int64(len(f.defectos) + 1)
	d.CreadoEn = time.Now().UTC()
	f.defectos = append(f.defectos, d)
	return d, nil
}

func (f *fakeDefectoRepository) ObtenerPorID(_ context.Context, id int64) (domain.DefectoDetalle, error) {
	for _, d := range f.defectos {
		if d.ID == id {
			return domain.DefectoDetalle{Defecto: d, TituloHistoria: f.titulos[d.HistoriaID]}, nil
		}
	}
	return domain.DefectoDetalle{}, domain.ErrDefectoNoEncontrado
}

func (f *fakeDefectoRepository) Resolver(_ context.Context, d domain.Defecto) error {
	f.resoluciones++
	if f.errResolver != nil {
		return f.errResolver
	}
	for i := range f.defectos {
		if f.defectos[i].ID == d.ID {
			f.defectos[i] = d
			return nil
		}
	}
	return domain.ErrDefectoNoEncontrado
}

func (f *fakeDefectoRepository) ListarPorProyecto(_ context.Context, proyectoID int64) ([]domain.DefectoDetalle, error) {
	f.listados++
	detalle := make([]domain.DefectoDetalle, 0)
	for _, d := range f.defectos {
		if d.ProyectoID == proyectoID {
			detalle = append(detalle, domain.DefectoDetalle{Defecto: d, TituloHistoria: f.titulos[d.HistoriaID]})
		}
	}
	return detalle, nil
}

// llamadaCambioEstado registra los argumentos de una llamada a CambiarEstado.
type llamadaCambioEstado struct {
	desde, hacia domain.Estado
	fecha        *time.Time
}

// fakeCambioEstadoRepository simula el UPDATE condicional del repositorio real:
// el mutex hace indivisible "si sigue en desde, pasarla a hacia".
type fakeCambioEstadoRepository struct {
	mu        sync.Mutex
	historias map[int64]domain.BacklogItem
	llamadas  []llamadaCambioEstado
}

func (f *fakeCambioEstadoRepository) ObtenerPorID(_ context.Context, id int64) (domain.BacklogItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	h, ok := f.historias[id]
	if !ok {
		return domain.BacklogItem{}, domain.ErrHistoriaNoEncontrada
	}
	return h, nil
}

func (f *fakeCambioEstadoRepository) CambiarEstado(_ context.Context, id int64, desde, hacia domain.Estado, fecha *time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.llamadas = append(f.llamadas, llamadaCambioEstado{desde: desde, hacia: hacia, fecha: fecha})
	h := f.historias[id]
	if h.Estado != desde {
		return domain.ErrConflictoEstadoHistoria
	}
	h.Estado = hacia
	h.FechaFinalizacion = fecha
	f.historias[id] = h
	return nil
}

func (f *fakeCambioEstadoRepository) registroDeLlamadas() []llamadaCambioEstado {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.llamadas)
}

type fakeCredenciales struct {
	mu       sync.Mutex
	usuarios map[string]domain.User
	hashes   map[string]string
	err      error
}

func (f *fakeCredenciales) Registrar(_ context.Context, u domain.User, hash string) (domain.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return domain.User{}, f.err
	}
	if f.usuarios == nil {
		f.usuarios = map[string]domain.User{}
		f.hashes = map[string]string{}
	}
	if _, ok := f.usuarios[u.NombreNormalizado]; ok {
		return domain.User{}, repository.ErrUsuarioDuplicado
	}
	u.ID = int64(len(f.usuarios) + 1)
	u.CreadoEn = time.Now().UTC()
	f.usuarios[u.NombreNormalizado] = u
	f.hashes[u.NombreNormalizado] = hash
	return u, nil
}

func (f *fakeCredenciales) ObtenerPorNombreNormalizado(_ context.Context, normalizado string) (domain.User, string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return domain.User{}, "", f.err
	}
	u, ok := f.usuarios[normalizado]
	if !ok {
		return domain.User{}, "", domain.ErrNoEncontrado
	}
	return u, f.hashes[normalizado], nil
}

type fakeSesiones struct {
	sesiones map[string]domain.SesionConUsuario
	usuarios map[int64]domain.User
	creadas  int
	cerradas []int64
	err      error
}

func (f *fakeSesiones) Crear(_ context.Context, huella string, s domain.Sesion) (domain.Sesion, error) {
	if f.err != nil {
		return domain.Sesion{}, f.err
	}
	if f.sesiones == nil {
		f.sesiones = map[string]domain.SesionConUsuario{}
	}
	f.creadas++
	s.ID = int64(f.creadas)
	f.sesiones[huella] = domain.SesionConUsuario{Sesion: s, Usuario: f.usuarios[s.UsuarioID]}
	return s, nil
}

func (f *fakeSesiones) ObtenerPorHuella(_ context.Context, huella string) (domain.SesionConUsuario, error) {
	if f.err != nil {
		return domain.SesionConUsuario{}, f.err
	}
	s, ok := f.sesiones[huella]
	if !ok {
		return domain.SesionConUsuario{}, domain.ErrNoEncontrado
	}
	return s, nil
}

func (f *fakeSesiones) Cerrar(_ context.Context, sesionID int64, ahora time.Time) error {
	if f.err != nil {
		return f.err
	}
	for huella, s := range f.sesiones {
		if s.ID == sesionID {
			s.CerradaEn = &ahora
			f.sesiones[huella] = s
		}
	}
	f.cerradas = append(f.cerradas, sesionID)
	return nil
}

// fakeContrasenas hashea con un prefijo fijo y cuenta las llamadas a Verificar.
type fakeContrasenas struct {
	mu          sync.Mutex
	verificadas int
}

func (f *fakeContrasenas) Hashear(contrasena string) (string, error) {
	return "hash:" + contrasena, nil
}

func (f *fakeContrasenas) Verificar(contrasena, hash string) bool {
	f.mu.Lock()
	f.verificadas++
	f.mu.Unlock()
	return hash == "hash:"+contrasena
}

func (f *fakeContrasenas) llamadasAVerificar() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.verificadas
}

// fakeEstimacionRepository guarda las historias y cuenta las escrituras de la
// estimación, para comprobar que un rechazo no escribe.
type fakeEstimacionRepository struct {
	historias  map[int64]domain.BacklogItem
	escrituras int
}

func (f *fakeEstimacionRepository) ObtenerPorID(_ context.Context, id int64) (domain.BacklogItem, error) {
	h, ok := f.historias[id]
	if !ok {
		return domain.BacklogItem{}, domain.ErrHistoriaNoEncontrada
	}
	return h, nil
}

func (f *fakeEstimacionRepository) GuardarEstimacion(_ context.Context, id int64, sp int) error {
	h, ok := f.historias[id]
	if !ok {
		return domain.ErrHistoriaNoEncontrada
	}
	h.EstimacionSP = &sp
	f.historias[id] = h
	f.escrituras++
	return nil
}

// fakeCalculadorMetricas simula el cálculo de métricas de HU-03 (HU-12).
type fakeCalculadorMetricas struct {
	resultado domain.MetricasProyecto
	err       error
}

func (f *fakeCalculadorMetricas) Ejecutar(_ context.Context, _ int64) (domain.MetricasProyecto, error) {
	return f.resultado, f.err
}

// fakeListadorSprints cuenta las llamadas para comprobar que no se lista
// cuando el proyecto no existe (HU-12).
type fakeListadorSprints struct {
	sprints  []domain.Sprint
	err      error
	llamadas int
}

func (f *fakeListadorSprints) ListarPorProyecto(_ context.Context, _ int64) ([]domain.Sprint, error) {
	f.llamadas++
	return f.sprints, f.err
}

// fakeObtenedorSprint devuelve los Sprints sembrados por id (HU-12).
type fakeObtenedorSprint struct {
	sprints map[int64]domain.Sprint
	err     error
}

func (f *fakeObtenedorSprint) ObtenerPorID(_ context.Context, id int64) (domain.Sprint, error) {
	if f.err != nil {
		return domain.Sprint{}, f.err
	}
	s, ok := f.sprints[id]
	if !ok {
		return domain.Sprint{}, domain.ErrSprintNoEncontrado
	}
	return s, nil
}

// fakeHistorialSprintRepository simula la lectura del detalle de un Sprint y
// cuenta las lecturas (HU-12).
type fakeHistorialSprintRepository struct {
	historias   []domain.HistoriaComprometida
	minutos     map[int64]int
	defectos    []domain.DefectoDetalle
	errHistoria error
	errMinutos  error
	errDefectos error
	lecturas    int
}

func (f *fakeHistorialSprintRepository) HistoriasComprometidas(_ context.Context, _ int64) ([]domain.HistoriaComprometida, error) {
	f.lecturas++
	return f.historias, f.errHistoria
}

func (f *fakeHistorialSprintRepository) MinutosPorHistoria(_ context.Context, _ int64) (map[int64]int, error) {
	f.lecturas++
	return f.minutos, f.errMinutos
}

func (f *fakeHistorialSprintRepository) DefectosDeSprint(_ context.Context, _ int64) ([]domain.DefectoDetalle, error) {
	f.lecturas++
	return f.defectos, f.errDefectos
}

// fakeAvanceHistoriasRepository devuelve las historias sembradas y cuenta las
// lecturas (HU-18).
type fakeAvanceHistoriasRepository struct {
	historias []domain.BacklogItem
	err       error
	llamadas  int
}

func (f *fakeAvanceHistoriasRepository) ListarPorProyecto(_ context.Context, _ int64) ([]domain.BacklogItem, error) {
	f.llamadas++
	return f.historias, f.err
}
