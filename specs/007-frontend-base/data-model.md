# Data Model: Base del frontend

La base no persiste datos ni define tablas; el esquema vive en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md) y no se toca. Estas son las formas que maneja el código del navegador.

## ApiError
Respuesta de error del backend ya interpretada.
- `status`: 400, 404, 409 u otro código recibido (por ejemplo, 500).
- `message`: texto para mostrar. Si el cuerpo no se puede interpretar: "La API respondió con estado N".
- `field`: nombre del campo del formulario; ausente si el error es general.

## ConnectionError
El backend no respondió (caído o red cortada). Tiene `message` fijo; no tiene `status`. Se distingue de `ApiError` con `instanceof`.

## EntradaDeRuta (registro de rutas)
Una línea por feature: `path` (o `index`) y la carga diferida de su pantalla. Va anidada bajo `projects/:id`.

## EntradaDeNavegacion (registro de navegación)
Una línea por feature: `ruta` relativa al proyecto y `texto` del menú.

## Respuestas de la API
Las usan las pantallas, no la base; los campos están en los DTO de `internal/http/*_dto.go`.
