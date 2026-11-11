# Quickstart: validar Permisos por rol (HU-21)

## Pruebas automatizadas

```bash
gofmt -l . && go vet ./...
go test ./internal/...            # TDD: servicios, middleware de sesión, test de rutas
go test ./features/...            # Godog: features/permisos_por_rol.feature + existentes ajustados
cd web && npm run lint && npm test && npm run build
```

Todo en verde es condición del check `backend` y `frontend` del CI.

## Validación manual con curl

Con el servidor en `:8080` y la base vacía:

```bash
# 1. Dos cuentas
SM=$(curl -s -XPOST localhost:8080/auth/registro -d '{"nombre":"ana","contrasena":"secreta123"}' | jq -r .token)
PB=$(curl -s -XPOST localhost:8080/auth/registro -d '{"nombre":"bruno","contrasena":"secreta123"}' | jq -r .token)

# 2. Sin sesión → 401 (escenario 3)
curl -i localhost:8080/projects

# 3. Ana crea un proyecto y queda como Scrum Master (escenario 7: "creador" se ignora)
curl -s -XPOST -H "Authorization: Bearer $SM" localhost:8080/projects \
  -d '{"nombre":"Demo","creador":"bruno"}'

# 4. Bruno no integra → 403 "No sos integrante de este proyecto" (escenario 4)
curl -i -H "Authorization: Bearer $PB" localhost:8080/projects/1

# 5. Ana elige a Bruno de la lista y lo suma como Product Builder (escenario 8)
curl -s -H "Authorization: Bearer $SM" localhost:8080/users
curl -s -XPOST -H "Authorization: Bearer $SM" localhost:8080/projects/1/members -d '{"nombre":"bruno"}'

# 6. Segundo Scrum Master → 409 (escenario 6)
curl -i -XPOST -H "Authorization: Bearer $SM" localhost:8080/projects/1/members -d '{"nombre":"bruno","rol":"scrum_master"}'

# 7. Bruno intenta crear un Sprint → 403 "Solo el Scrum Master puede…" (escenario 2)
curl -i -XPOST -H "Authorization: Bearer $PB" localhost:8080/sprints -d '{"proyecto_id":1,"objetivo":"x"}'

# 8. Mis proyectos de cada uno, con su rol (escenarios 5 y 9)
curl -s -H "Authorization: Bearer $PB" localhost:8080/projects
```

Para la votación (escenario 7), con una historia en el proyecto: Ana inicia la votación,
Bruno vota y Ana consulta `GET /backlog/{id}/votacion?usuario_id=2`: el voto de Bruno no se
ve hasta revelar.

## En la interfaz

1. Ingresar como Bruno (Product Builder): las acciones de Sprints, votación, edición y
   equipo aparecen deshabilitadas con «Solo el Scrum Master puede…».
2. No hay selector «Quién actúa» en la votación ni de integrante en las horas.
3. Abrir por URL un proyecto ajeno: pantalla «No sos integrante de este proyecto».
4. Cuenta nueva sin proyectos: «Mis proyectos» muestra el mensaje y el botón de crear.
