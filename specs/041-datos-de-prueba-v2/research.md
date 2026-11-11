# Research: Datos de prueba v2

Las decisiones de la v1 (id fijo 25, transacción única, bloques, `DEMO_HABILITADO`) siguen en [`021/research.md`](../021-datos-de-prueba/research.md). Aquí solo lo nuevo.

## R1. Hash único de la contraseña de demostración
- **Decisión**: `Cargar` deriva `demo1234` una vez, fuera de la transacción, y lo reutiliza para los cuatro usuarios.
- **Razón**: PBKDF2 de 600 000 iteraciones es el costo dominante; cuatro derivaciones en el plan gratis de Render serían lentas, y dentro de la transacción tendrían tomada la escritura. La contraseña es pública, así que compartir sal no resta seguridad real.
- **Alternativas**: hash por usuario (4× el costo); hash precalculado como constante (frágil si cambian las iteraciones y deja un hash fijo en el código).

## R2. Interfaz `Hasheador` propia de `internal/demo`
- **Decisión**: `type Hasheador interface{ Hashear(string) (string, error) }` en `internal/demo`; `service.ServicioContrasenas` la satisface por tipado estructural.
- **Razón**: Dependency Inversion sin importar `service`; en pruebas se inyecta `service.NewServicioContrasenas(1)` o un fake.
- **Alternativa**: importar `service.HasheadorContrasenas`: acopla el paquete de carga a los casos de uso.

## R3. Clasificación de conflictos independiente de la restricción
- **Decisión**: tras el rollback se consulta: proyecto 25 con nombre reservado ⇒ `ya_cargado`; usuario de la demo existente ⇒ `usuario_existente` (nombre del primero); proyecto 25 de otro nombre ⇒ `id ocupado`.
- **Razón**: la v2 tiene dos fuentes de `UNIQUE` (usuario y proyecto); el orden de la consulta, no el orden de las escrituras, fija la prioridad (FR-006). Con una segunda carga, el usuario falla primero, pero el resultado correcto es `ya_cargado`.
- **Alternativa**: precomprobar con `SELECT` antes de escribir: reintroduce lectura→escritura (riesgo de `SQLITE_BUSY` entre procesos).

## R4. Datos históricos que no pasan por las reglas de la API
- **Decisión**: defectos Resueltos en Sprints finalizados y horas en Sprints cerrados se insertan directo (con constructores del dominio donde aplican), sin exigir Sprint activo.
- **Razón**: es el mismo criterio de la v1 para las horas; la regla «hay Sprint activo» es de la API (§6 de MODELO-DE-DATOS.md).

## R5. Ronda abierta
- **Decisión**: `revelada_en = NULL` con dos votos.
- **Razón**: ya lo permite §4.7 (`NULL` = votos ocultos); no hace falta cambio de esquema.
