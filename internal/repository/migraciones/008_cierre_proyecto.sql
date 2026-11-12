-- Cierre manual del proyecto: MODELO-DE-DATOS.md §4.1 (D16, spec 044, #300).
--
-- fecha_cierre es una fecha de calendario YYYY-MM-DD con el CHECK de D12.
-- NULL significa que el proyecto está abierto, así que los proyectos que ya
-- existían quedan abiertos. ADD COLUMN admite el CHECK porque la columna es
-- anulable y sin valor por defecto: no hace falta recrear la tabla (§9 punto 4).

ALTER TABLE proyectos ADD COLUMN fecha_cierre TEXT CHECK (date(fecha_cierre) IS fecha_cierre);
