# Frontend compartido

Base React + TypeScript + Vite para las pantallas del Sprint 2.

## Ejecutar en desarrollo

Desde esta carpeta:

```bash
npm install
npm run dev
```

La aplicación queda disponible en `http://localhost:5173`. El backend debe
estar ejecutándose en `http://localhost:8080`; Vite redirige automáticamente
las rutas `/backlog`, `/projects` y `/sprints` mediante su proxy.

## Comandos

```bash
npm run build
npm run lint
npm test
```

## Decisión de ruteo

La base usa un router liviano sobre la History API para no agregar una
dependencia de ruteo antes de la decisión del equipo. El registro central
mantiene una línea por feature y expone `matchRoute`, por lo que React Router
puede incorporarse posteriormente sin cambiar las URLs públicas.

## Estructura

Cada historia vive en una carpeta propia dentro de `src/features/` y no
importa código de otra historia. La infraestructura compartida está en
`src/compartido/`:

```text
src/
  compartido/
    api/             Cliente HTTP y errores de la API
    componentes/     Button, Field, ErrorMessage y LoadingIndicator
    router/          Registro común de rutas
  features/
    hu01/
    hu04/
    hu05/
    hu13/
```

Para agregar una pantalla:

1. Crear su carpeta en `src/features/<nombre>/`.
2. Exportar el componente de la pantalla desde esa carpeta.
3. Agregar una única línea al registro de rutas en
   `src/compartido/router/routes.ts`.
4. Agregar sus tests junto a la feature.

Las rutas reciben el proyecto actual mediante `/projects/:id/...`; no se debe
guardar el identificador del proyecto únicamente en estado global o local.

## Cliente de API

`src/compartido/api/client.ts` expone `apiFetch`, que:

- decodifica respuestas JSON;
- envía `Accept: application/json`;
- lanza `ApiError` para estados HTTP no exitosos;
- conserva `status`, `message` y `field`;
- soporta errores con campo y errores generales sin campo.

El cliente usa rutas relativas para que el proxy funcione en desarrollo y la
misma aplicación pueda configurarse detrás de un servidor reverse proxy en
otros entornos.
