# Data Model: Pantalla de Planning Poker (spec 029)

Sin cambios de esquema ni de base: esta spec no define tablas. Las rondas y votos están en [MODELO-DE-DATOS.md §4.7](../../MODELO-DE-DATOS.md#47-rondas_votacion-y-votos--hu-02-hu-14-tablas-nuevas-migración-003).

## Tipos del frontend (`votacion.ts`)

- `Rol`: `'scrum_master' | 'product_builder'`.
- `Carta`: `1 | 2 | 3 | 5 | 8 | 13 | 21`; `CARTAS` es la lista ordenada.
- `Integrante`: `{ id, nombre, rol }` (de `GET /projects/{id}/members`).
- `VotoDeIntegrante`: `{ usuario_id, nombre, rol, ya_voto, valor: number | null }`.
- `Votacion`: `{ historia_id, ronda, revelada_en: string | null, integrantes: VotoDeIntegrante[], promedio: { exacto, carta_sugerida } | null }`.

## Reglas puras de la vista

- `esScrumMaster(integrante)`: `rol === 'scrum_master'`.
- `estaRevelada(votacion)`: `revelada_en !== null`.
- `faltanPorVotar(votacion)`: integrantes con `ya_voto === false`.
- `puedeRevelar(votacion, actor)`: actor es Scrum Master, ronda sin revelar y `faltanPorVotar` vacío.
- `textoDeVoto(integrante, actorId, revelada)`: valor propio o revelado; «Oculto» si votó y es ajeno; «Falta» si no votó.
