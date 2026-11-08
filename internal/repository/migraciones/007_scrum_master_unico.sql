-- Scrum Master único por proyecto: MODELO-DE-DATOS.md §4.3 (D4, spec 039, #236).
--
-- 1) Corrección previa: si un proyecto ya tiene más de un scrum_master, queda el
--    de creado_en más antiguo (si empatan, el de menor usuario_id) y los demás
--    pasan a product_builder. No se borra ninguna fila: siguen siendo integrantes.
--    Sin esto el índice no se puede crear y la aplicación no arranca. Los
--    proyectos con un solo Scrum Master, o con ninguno, no se tocan.
--    Criterio aprobado el 2026-11-08: la base de Render se vacía en cada deploy y
--    los datos de prueba crean un solo Scrum Master por proyecto, así que la
--    migración tiene que arrancar sola, sin un arreglo manual.
--    creado_en se ordena como texto porque son instantes RFC 3339 en UTC (§2).
-- 2) Índice único parcial: a lo sumo un scrum_master por proyecto.

UPDATE integrantes
SET rol = 'product_builder'
WHERE rol = 'scrum_master'
  AND usuario_id <> (
      SELECT otro.usuario_id
      FROM integrantes AS otro
      WHERE otro.proyecto_id = integrantes.proyecto_id AND otro.rol = 'scrum_master'
      ORDER BY otro.creado_en, otro.usuario_id
      LIMIT 1
  );

CREATE UNIQUE INDEX integrantes_un_scrum_master ON integrantes (proyecto_id) WHERE rol = 'scrum_master';
