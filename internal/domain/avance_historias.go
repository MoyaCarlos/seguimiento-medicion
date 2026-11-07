package domain

// SprintRef identifica el Sprint de una historia en el avance (HU-18).
type SprintRef struct {
	ID     int64
	Numero int
}

// HistoriaAvance es una historia tal como se muestra en el avance del
// proyecto. EstimacionSP en nil significa "sin estimar"; Sprint en nil, que la
// historia está en el Product Backlog.
type HistoriaAvance struct {
	ID           int64
	Titulo       string
	EstimacionSP *int
	Estado       Estado
	Sprint       *SprintRef
}

// GrupoAvance reúne las historias de un mismo estado, con su cantidad y la
// suma de sus Story Points (las historias sin estimar suman cero).
type GrupoAvance struct {
	Estado      Estado
	Cantidad    int
	StoryPoints int
	Historias   []HistoriaAvance
}

// TotalesAvance son la cantidad de historias y los Story Points de todo el
// proyecto.
type TotalesAvance struct {
	Historias   int
	StoryPoints int
}

// AvanceHistorias es la vista de solo lectura de HU-18: tres grupos, siempre
// en el orden Nueva, En progreso, Completada, y los totales del proyecto.
type AvanceHistorias struct {
	Grupos []GrupoAvance
	Total  TotalesAvance
}

// AgruparAvance agrupa las historias por estado, en el orden Nueva, En
// progreso, Completada, y suma cantidad y Story Points de cada grupo y del
// proyecto. Conserva el orden recibido dentro de cada grupo.
func AgruparAvance(historias []BacklogItem, sprints []Sprint) (AvanceHistorias, error) {
	return AvanceHistorias{}, nil
}
