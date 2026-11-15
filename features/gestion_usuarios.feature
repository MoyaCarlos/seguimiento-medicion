# language: es
@HU-17
Característica: Renombrar el propio usuario
  Como usuario
  Quiero cambiar mi nombre de usuario
  Para corregir un error de tipeo sin perder mi historial en los proyectos

  Antecedentes:
    Dado que "ana" tiene una cuenta con la contraseña "clave-segura-1"

  Escenario: Cambio exitoso del propio nombre
    Dado que "ana" creó el proyecto de prueba "Tienda"
    Cuando cambio mi nombre de usuario a "ana.perez"
    Entonces el sistema responde 200 con el nombre "ana.perez"
    Y el equipo de "Tienda" muestra a "ana.perez"
    Y puedo ingresar con "ana.perez" y la contraseña "clave-segura-1"
    Y el ingreso con "ana" es rechazado

  Escenario: Nombre en uso
    Dado que existe otra cuenta "bruno" con la contraseña "clave-segura-2"
    Cuando intento cambiar mi nombre a "BRUNO"
    Entonces el cambio se rechaza con 409 y el mensaje "Ese nombre de usuario ya existe"
    Y mi nombre sigue siendo "ana"

  Escenario: Sin sesión
    Cuando intento cambiar el nombre sin sesión a "ana.perez"
    Entonces la petición se rechaza con 401 y el mensaje "Tu sesión no es válida"

  Escenario: Nombre inválido
    Cuando intento cambiar mi nombre a "ab"
    Entonces la petición se rechaza con 400 en el campo "nombre"
