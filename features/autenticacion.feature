# language: es

# Un escenario por Acceptance Scenario de las User Stories 1 y 2 de
# specs/034-hu-20-autenticacion/spec.md.
# Cobertura por categoría:
# - Normal: registro exitoso, la persona queda con la sesión iniciada (US1.1).
# - Alternativo: nombre repetido con otras mayúsculas (US1.2).
# - Límite: el largo de la contraseña en 7, 8, 128 y 129 caracteres y el formato del
#   nombre se prueban con tests unitarios del dominio.
# - Error: nombre de usuario ya existente (US1.2) y contraseña de menos de 8 caracteres (US1.3).
# - Ingreso: credenciales correctas con otras mayúsculas (US2.1), usuario previo de HU-04
#   sin contraseña, que no puede ingresar, y credenciales incorrectas con el mismo mensaje
#   en todos los casos, sin revelar cuál falló (US2.2).

@HU-20
Característica: Autenticación de usuarios

  Como integrante del equipo
  Quiero crear una cuenta con nombre de usuario y contraseña, ingresar y mantener mi sesión
  Para que el sistema sepa quién soy sin pedirme mis credenciales en cada pantalla

  Escenario: Registro exitoso deja la sesión iniciada
    Cuando me registro con el nombre de usuario "ana.perez" y la contraseña "una-clave-larga"
    Entonces la cuenta "ana.perez" queda creada
    Y la persona queda con una sesión válida a nombre de "ana.perez"
    Y la respuesta del registro no contiene la contraseña

  Escenario: Nombre de usuario repetido con otras mayúsculas
    Dado que ya existe la cuenta "ana.perez"
    Cuando me registro con el nombre de usuario "Ana.Perez" y la contraseña "otra-clave-larga"
    Entonces el registro de la cuenta se rechaza con el mensaje "Ese nombre de usuario ya existe"
    Y no se crea ninguna cuenta nueva

  Escenario: Contraseña de menos de 8 caracteres
    Cuando me registro con el nombre de usuario "ana.perez" y la contraseña "corta12"
    Entonces el registro de la cuenta se rechaza indicando el campo "contrasena"
    Y no se crea ninguna cuenta nueva

  Escenario: Ingreso exitoso con otras mayúsculas en el nombre
    Dado que ya existe la cuenta "ana.perez"
    Cuando ingreso con el nombre de usuario "ANA.PEREZ" y la contraseña "clave-previa-1"
    Entonces la persona queda con una sesión válida a nombre de "ana.perez"
    Y la respuesta del ingreso no contiene la contraseña

  Escenario: Un usuario previo creado solo con nombre no puede ingresar
    Dado que existe el usuario previo "luis.gomez" creado solo con nombre
    Cuando ingreso con el nombre de usuario "luis.gomez" y la contraseña "cualquier-clave-1"
    Entonces el ingreso se rechaza con el mensaje "Usuario o contraseña incorrectos"

  Esquema del escenario: Credenciales incorrectas no revelan cuál falló
    Dado que ya existe la cuenta "ana.perez"
    Cuando ingreso con el nombre de usuario "<nombre>" y la contraseña "<contraseña>"
    Entonces el ingreso se rechaza con el mensaje "Usuario o contraseña incorrectos"
    Y no hay ninguna sesión nueva

    Ejemplos:
      | nombre     | contraseña         |
      | ana.perez  | otra-clave-larga   |
      | fantasma   | clave-previa-1     |

  Escenario: Con una sesión vigente el sistema sabe quién es la persona
    Dado que "ana.perez" inició sesión
    Cuando consulto quién soy con mi token
    Entonces el sistema responde que soy "ana.perez"
    Y la respuesta de quién soy no contiene la contraseña

  Escenario: Una sesión vencida se rechaza
    Dado que "ana.perez" inició sesión
    Y pasan 7 días
    Cuando consulto quién soy con mi token
    Entonces la consulta de quién soy se rechaza con el mensaje "Tu sesión no es válida"

  Escenario: Un token inventado se rechaza
    Cuando consulto quién soy con el token "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
    Entonces la consulta de quién soy se rechaza con el mensaje "Tu sesión no es válida"

  Escenario: Una consulta sin token se rechaza
    Cuando consulto quién soy sin token
    Entonces la consulta de quién soy se rechaza con el mensaje "Tu sesión no es válida"

  # Desde HU-21 (spec 042, escenario 3) las rutas que no son públicas exigen sesión.
  Escenario: Una ruta protegida sin token se rechaza
    Cuando consulto la lista de proyectos sin token
    Entonces la lista de proyectos se rechaza con el mensaje "Tu sesión no es válida"

  Escenario: Al cerrar la sesión el token deja de ser válido
    Dado que "ana.perez" inició sesión
    Cuando cierro mi sesión
    Entonces el cierre de sesión responde sin contenido
    Y consulto quién soy con mi token
    Y la consulta de quién soy se rechaza con el mensaje "Tu sesión no es válida"

  Escenario: Cerrar una sesión no cierra la otra sesión de la misma persona
    Dado que "ana.perez" inició sesión
    Y "ana.perez" abre una segunda sesión
    Cuando cierro mi sesión
    Entonces el cierre de sesión responde sin contenido
    Y consulto quién soy con el token de mi segunda sesión
    Y el sistema responde que soy "ana.perez"
