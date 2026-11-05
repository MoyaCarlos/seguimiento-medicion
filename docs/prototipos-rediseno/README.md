# Prototipos del rediseño (referencia visual)

Prototipos HTML estáticos del rediseño del frontend (issue #230). Se abren con el navegador (doble clic), sin servidor y sin dependencias.

| Archivo | Qué muestra |
|---|---|
| `remake-con-login.html` | **El más completo y el de referencia:** pantallas de ingreso y registro (HU-20), «Mis proyectos», panel, Product Backlog, Sprints, Estado y votación, con la paleta aprobada (azul marino y teal) |
| `remake.html` | La misma propuesta sin las pantallas de login |
| `crear-proyecto.html` | La pantalla de crear proyecto y la lista de proyectos |

## Cómo usarlos

- **Son solo una referencia de estructura y de estilo.** Los datos son inventados: no copiar contenido que la API no entrega (por ejemplo estado, avance o Scrum Master en las tarjetas de proyecto).
- **Por ahora no se implementa el rediseño.** Cada historia nueva (login, permisos, dashboard, defectos, etc.) arma su pantalla básica con el estilo que ya tiene la aplicación desplegada, y usa estos prototipos solo para ubicar los elementos (qué va en cada pantalla y dónde).
- **El rediseño moderno se aplica después,** cuando las pantallas nuevas estén listas, en PR aparte y por pantallas (#230).
- Decisiones ya tomadas: fechas en formato dd/mm/aaaa, ancho completo de pantalla y sin librerías de estilos nuevas.
- Los archivos no son código de la aplicación y no se ejecutan en el CI ni en Render (esta carpeta está bajo `docs/`, ruta ignorada por el despliegue).
