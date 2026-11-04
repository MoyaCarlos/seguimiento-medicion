# Quickstart: validar la pantalla de Planning Poker (spec 030)

## Prerrequisitos

- Go (`/usr/local/go/bin`) y Node 24 (`nvm use`).
- Dependencias del frontend: `cd web && npm ci`.

## Pruebas automáticas

```sh
cd web && npm run lint && npm test && npm run build
```

## Validación manual

1. Backend con base temporal:
   `mkdir -p /tmp/hu02-datos && PORT=8080 DB_PATH=/tmp/hu02-datos/seguimiento.db DEMO_HABILITADO=true go run ./cmd/api`
2. Datos de prueba: `curl -X POST http://localhost:8080/demo/datos` (proyecto 25, tres integrantes, Beta Uno es el scrum_master, una historia con una ronda).
3. Frontend: `cd web && npm run dev`, abrir `/projects/25/backlog`.
4. Pulsar «Votar» en una historia y comprobar:
   - Sin integrante elegido no hay cartas ni acciones.
   - Como Beta Uno: «Iniciar votación» (historia sin ronda), votar, ver «Oculto» en los demás.
   - Como otro integrante: votar y cambiar el voto; no aparecen «Iniciar votación» ni «Revelar».
   - «Revelar» deshabilitado hasta que votan todos; luego muestra valores, promedio exacto y carta sugerida.
   - «Actualizar» refleja votos hechos desde otra pestaña.
5. Apagar backend y frontend al terminar.

Contrato de la API: [contracts/api-votacion.md](./contracts/api-votacion.md).
