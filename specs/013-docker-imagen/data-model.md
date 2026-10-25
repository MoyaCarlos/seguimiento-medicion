# Data Model: Imagen Docker

Esta historia **no crea ni modifica tablas, columnas ni migraciones**. La base es la descrita
en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md); solo cambia dónde vive el archivo
(`DB_PATH`, en la imagen `/data/seguimiento.db`).

Configuración de la app (no es modelo de datos):

| Variable | Default | En la imagen |
|---|---|---|
| `PORT` | `8080` | `8080` |
| `DB_PATH` | `seguimiento.db` | `/data/seguimiento.db` |
| `WEB_DIR` | `web/dist` | `/app/web/dist` |
