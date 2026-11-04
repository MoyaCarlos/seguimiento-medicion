# language: es

# Un escenario por Acceptance Scenario de la User Story 1 de
# specs/034-hu-20-autenticacion/spec.md.
# Cobertura por categoría:
# - Normal: registro exitoso, la persona queda con la sesión iniciada (US1.1).
# - Alternativo: nombre repetido con otras mayúsculas (US1.2).
# - Límite: el largo de la contraseña en 7, 8, 128 y 129 caracteres y el formato del
#   nombre se prueban con tests unitarios del dominio.
# - Error: nombre de usuario ya existente (US1.2) y contraseña de menos de 8 caracteres (US1.3).

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
    Entonces el registro se rechaza con el mensaje "Ese nombre de usuario ya existe"
    Y no se crea ninguna cuenta nueva

  Escenario: Contraseña de menos de 8 caracteres
    Cuando me registro con el nombre de usuario "ana.perez" y la contraseña "corta12"
    Entonces el registro se rechaza indicando el campo "contrasena"
    Y no se crea ninguna cuenta nueva
