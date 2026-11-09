# Research: Autenticación de usuarios (HU-20)

## D1. Hash de contraseña

- **Decision**: PBKDF2-HMAC-SHA-256 con `crypto/pbkdf2` (`pbkdf2.Key(sha256.New, ...)`),
  sal aleatoria de 16 bytes por usuario, clave derivada de 32 bytes y 600 000 iteraciones
  (valor recomendado por OWASP para PBKDF2-SHA-256). Formato guardado:
  `pbkdf2-sha256$<iteraciones>$<sal base64>$<hash base64>`.
- **Rationale**: está en la biblioteca estándar (Go 1.24+), no suma dependencia y el
  formato autodescriptivo permite subir las iteraciones sin tocar a los usuarios
  existentes. Las iteraciones y el generador de sal se pasan al constructor; los tests usan
  pocas iteraciones.
- **Alternatives**: bcrypt/argon2 (dependencia nueva, prohibida por el Scrum Master);
  SHA-256 simple (no es un hash de contraseñas).

## D2. Comparación sin revelar cuál falló

- **Decision**: `Verificar` compara con `subtle.ConstantTimeCompare`. Si el usuario no
  existe o no tiene contraseña, `Ingresar` verifica igual contra un hash señuelo calculado
  al construir el servicio (mismos parámetros) y devuelve el mismo error.
- **Rationale**: iguala el costo y la respuesta de los tres casos (usuario inexistente,
  usuario previo sin contraseña, contraseña incorrecta), que es lo que exige FR-007.
- **Alternatives**: devolver antes si no hay usuario (filtra por tiempo si el usuario existe).

## D3. Token de sesión

- **Decision**: 32 bytes de `crypto/rand`, codificados en base64url sin relleno (43
  caracteres). Se guarda el SHA-256 en hexadecimal (`token_hash`, UNIQUE). Un token que no
  tenga esa forma se rechaza antes de consultar la base, con el mismo error.
- **Rationale**: con 256 bits de entropía basta un hash rápido (no hace falta PBKDF2 para
  un secreto que no es una contraseña humana); si se filtra la base, los tokens no sirven.
- **Alternatives**: JWT (no se puede revocar al cerrar sesión sin una lista, y el Scrum
  Master decidió sesiones en la base).

## D4. Vencimiento y vigencia

- **Decision**: `vence_en = creado_en + 7 días`, fijo, sin renovar. Una sesión es vigente si
  `cerrada_en IS NULL` y `ahora < vence_en`. La regla vive en `domain.Sesion.Vigente(ahora)`;
  el reloj se inyecta como `func() time.Time` (igual que `time.Now` en `CambiarEstadoHistoria`).
- **Rationale**: testeable sin dormir; una sola definición de «vigente».
- **Alternatives**: expiración deslizante (YAGNI, no pedida).

## D5. `vence_en` contra la convención de instantes (desvío documentado)

- **Decision**: `sesiones` lleva `creado_en` (auditoría, como el resto) y `vence_en`, que **sí
  entra en una regla de negocio**. Se agrega una nota en `MODELO-DE-DATOS.md` §2 y §4.11.
- **Rationale**: calcular `creado_en + 7 días` al leer ata la regla al momento de consulta y
  obliga a recalcular si cambia la duración; guardarlo congela el contrato de cada sesión.
- **Alternatives**: no guardar `vence_en` (cumple la letra de §2, pero cambiar la duración
  alteraría sesiones ya emitidas).

## D6. Cerrar sesión sin borrar (desvío evitado)

- **Decision**: `cerrada_en TEXT NULL`; cerrar escribe el instante, no borra la fila.
- **Rationale**: §2 «Borrado: ninguna historia borra datos». Además deja auditable cuándo
  se cerró.
- **Alternatives**: `DELETE` (más simple, pero contradice la convención del modelo).

## D7. Puertos (ISP)

- **Decision**: cinco puertos de un método cada uno, agrupados por archivo:
  `RegistroCredenciales.Registrar`, `BuscadorCredenciales.ObtenerPorNombreNormalizado`,
  `CreadorSesiones.Crear`, `BuscadorSesiones.ObtenerPorHuella`,
  `CerradorSesiones.Cerrar`. Cada caso de uso depende solo de los que usa; las
  estructuras SQLite implementan varios.
- **Rationale**: pedido del equipo; además los fakes de los tests quedan mínimos.
- **Alternatives**: ampliar `UserRepository` (mezcla identidad de integrantes con
  credenciales y obliga a todos sus consumidores a conocer contraseñas).

## D8. Registro atómico del usuario

- **Decision**: `Registrar` inserta usuario con `contrasena_hash` en una sola sentencia; la
  unicidad la garantiza `usuarios.nombre_normalizado UNIQUE`, y `SQLITE_CONSTRAINT_UNIQUE`
  se traduce a `ErrUsuarioDuplicado` (ya existe), que el servicio convierte en
  `domain.ErrNombreUsuarioEnUso`. La sesión se crea después; si fallara, la cuenta existe y
  la persona puede ingresar.
- **Rationale**: cubre la carrera de dos registros simultáneos sin consultar antes.

## D9. Reglas de nombre y contraseña (dominio puro)

- **Decision**: nombre: 3 a 30 caracteres, solo `[A-Za-z0-9._-]` (ASCII, sin recortar
  espacios); identidad con `NormalizarNombre` (ya existe). Contraseña: de 8 a 128 caracteres
  (se cuentan runas). Errores `ValidationError` con `Campo` `nombre` y `contrasena`.
- **Nota**: el registro nunca crea nombres con espacios, pero el ingreso normaliza el nombre con `NormalizarNombre` (recorta y pasa a minúsculas) antes de buscarlo; por eso « Ana.Perez » ingresa. Es inofensivo: la identidad es el nombre normalizado.
- **Rationale**: lo definió el Scrum Master; la tabla de casos cubre los bordes (2, 3, 30, 31,
  espacios, tildes, `@`, vacía, 7, 8, 128, 129).

## D10. Errores y códigos HTTP

- **Decision**: `ErrCredencialesInvalidas` → 401 «Usuario o contraseña incorrectos»;
  `ErrNoAutenticado` → 401 con un único mensaje para token ausente, mal formado, inexistente,
  vencido o cerrado; `ErrNombreUsuarioEnUso` → 409 «Ese nombre de usuario ya existe»;
  `ValidationError` sigue siendo 400 con `campo`. Se agregan a `escribirError`.
- **Rationale**: reutiliza el formato `{campo, mensaje}` que el cliente ya sabe leer.

## D11. Logs

- **Decision**: ningún handler ni servicio registra cuerpos de petición ni cabeceras. Un test
  captura el log durante registro, ingreso y cierre y verifica que no aparezcan la
  contraseña ni el token.

## D12. Alcance de las rutas protegidas en el frontend

- **Decision (confirmada por el Scrum Master)**: todas las rutas de la interfaz, salvo
  `/ingresar` y `/registro`, exigen sesión. Eso incluye la raíz (`/`), todas las rutas de
  `projects/:id/*` y la ruta «no encontrada» (`*`). Sin sesión se redirige a `/ingresar`
  conservando la ruta a la que volver (`state.desde`, con ruta y búsqueda), y tras ingresar o
  registrarse se vuelve a ella. El backend sigue sin exigir sesión; solo la interfaz pide
  ingresar (escenario 6).
- **Rationale**: es lo que piden el escenario 6 («se vuelve a la pantalla en la que estaba») y
  FR-022.
- **Consecuencia aceptada** por el Scrum Master (decisión 3): los usuarios previos sin
  contraseña no pueden entrar a la interfaz hasta registrarse con otro nombre.
- **Implementación**: un solo archivo (`rutas.tsx`); `rutas.test.tsx` usa una sesión simulada y
  verifica que cada ruta protegida redirige a `/ingresar` sin sesión, y que `/ingresar` y
  `/registro` siguen abiertas. Ningún test del backend cambia.

## D13. Enganche del cliente HTTP con la sesión

- **Decision**: `compartido/api/sesionHttp.ts` guarda dos funciones registradas por la
  feature (`obtenerToken`, `alNoAutenticado`). `client.ts` agrega `Authorization: Bearer` si
  hay token y llama a `alNoAutenticado` ante un 401 de una petición que llevó token y no es
  de `/auth/ingreso` ni `/auth/registro`.
- **Rationale**: `compartido` no puede importar una feature (test de aislamiento).

## D14. Pantalla de origen tras ingresar

- **Decision**: se pasa en el estado de navegación (`state.desde`, con ruta y búsqueda) y el
  aviso «Tu sesión venció» en `state.vencida`. Tras ingresar o registrarse se navega a
  `desde`, o a `/` si no hay.
- **Rationale**: sin parámetros en la URL (evita redirecciones abiertas) y sin estado global
  extra.

## D15. Token en `localStorage`

- **Decision**: clave única, todas las lecturas y escrituras dentro de `try/catch`; si el
  almacenamiento falla, la sesión vive solo en memoria hasta recargar.
- **Rationale**: pedido del equipo; el navegador puede bloquear el almacenamiento.

## D16. Vigilancia de la migración

- **Decision**: `ALTER TABLE usuarios ADD COLUMN contrasena_hash TEXT` (NULL) más
  `CREATE TABLE sesiones`, ambos en `006_autenticacion.sql`. Test de esquema con y sin
  usuarios previos.
- **Rationale**: `ADD COLUMN` nulo no toca las filas existentes, así que no puede fallar
  por usuarios previos.
