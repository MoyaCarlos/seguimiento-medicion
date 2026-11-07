package domain

import "fmt"

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

var estadosDelAvance = []Estado{EstadoNueva, EstadoEnProgreso, EstadoCompletada}

// AgruparAvance agrupa las historias por estado, en el orden Nueva, En
// progreso, Completada, y suma cantidad y Story Points de cada grupo y del
// proyecto. Conserva el orden recibido dentro de cada grupo. Devuelve error si
// una historia tiene un estado desconocido.
func AgruparAvance(historias []BacklogItem, sprints []Sprint) (AvanceHistorias, error) {
	grupos := make([]GrupoAvance, len(estadosDelAvance))
	posicion := make(map[Estado]int, len(estadosDelAvance))
	for i, estado := range estadosDelAvance {
		grupos[i] = GrupoAvance{Estado: estado, Historias: []HistoriaAvance{}}
		posicion[estado] = i
	}

	var total TotalesAvance
	for _, h := range historias {
		i, ok := posicion[h.Estado]
		if !ok {
			return AvanceHistorias{}, fmt.Errorf("la historia %d tiene un estado desconocido: %q", h.ID, h.Estado)
		}
		puntos := 0
		if h.EstimacionSP != nil {
			puntos = *h.EstimacionSP
		}
		grupos[i].Historias = append(grupos[i].Historias, HistoriaAvance{
			ID: h.ID, Titulo: h.Titulo, EstimacionSP: h.EstimacionSP, Estado: h.Estado,
		})
		grupos[i].Cantidad++
		grupos[i].StoryPoints += puntos
		total.Historias++
		total.StoryPoints += puntos
	}
	return AvanceHistorias{Grupos: grupos, Total: total}, nil
}
