# Contrato: render.yaml

Referencia versionada de la configuración del servicio. No se crea un Blueprint.

```yaml
services:
  - type: web
    name: seguimiento-medicion-sprintometro
    runtime: docker
    plan: free
    region: oregon
    branch: main
    dockerfilePath: ./Dockerfile
    dockerContext: .
    healthCheckPath: /health
    autoDeployTrigger: checksPass
    envVars:
      - key: DEMO_HABILITADO
        value: "true"
    buildFilter:
      ignoredPaths:
        - docs/**
        - specs/**
        - features/**
        - "**/*.md"
        - .github/**
        - .specify/**
        - render.yaml
```

Correspondencia con el dashboard:

| Dashboard | Clave |
|---|---|
| Language: Docker | `runtime: docker` |
| Instance Type: Free | `plan: free` |
| Region: Oregon (US West) | `region: oregon` |
| Branch: main | `branch: main` |
| Dockerfile Path / Docker Build Context Directory | `dockerfilePath`, `dockerContext` |
| Health Check Path | `healthCheckPath` |
| Auto-Deploy: After CI Checks Pass | `autoDeployTrigger: checksPass` |
| Environment: `DEMO_HABILITADO` | `envVars` |
| Build Filters: Included Paths (vacío) / Ignored Paths | `buildFilter.ignoredPaths` (sin `paths`) |
