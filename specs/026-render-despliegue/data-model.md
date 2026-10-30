# Data Model: Despliegue continuo en Render

Sin entidades ni cambios de esquema. Usa la base existente descrita en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md).

La base vive en `/data/seguimiento.db` dentro del contenedor, sin disco persistente: se recrea y migra en cada arranque.
