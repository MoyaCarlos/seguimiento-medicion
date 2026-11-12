-- Baja lógica de integrantes: MODELO-DE-DATOS.md §4.3 (D15, spec 045, #175).
--
-- fecha_baja es una fecha de calendario YYYY-MM-DD con el CHECK de D12. NULL
-- significa que el integrante está activo, así que los integrantes que ya
-- existían quedan activos. Quitar a alguien guarda la fecha; reactivarlo vuelve
-- a NULL en la misma fila. ADD COLUMN admite el CHECK porque la columna es
-- anulable y sin valor por defecto: no hace falta recrear la tabla (§9 punto 4).

ALTER TABLE integrantes ADD COLUMN fecha_baja TEXT CHECK (fecha_baja IS NULL OR date(fecha_baja) IS fecha_baja);
