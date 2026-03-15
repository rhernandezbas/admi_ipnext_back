# ADR-001: Tech Stack Backend

## Estado
Aceptado — Implementado

## Contexto
Definir el stack tecnológico del backend de IPNEXT Admin antes de iniciar el desarrollo. El backend debe servir a una SPA React como única consumidora.

## Decisión

| Categoría | Tecnología | Justificación |
|-----------|-----------|---------------|
| Lenguaje | **Go 1.25** | Rendimiento nativo, tipado estático, compilado, binario único fácil de desplegar |
| Framework HTTP | **Gin** | El más usado en Go, rápido, middleware ecosystem maduro |
| ORM | **GORM** | ORM más popular en Go, soporte MySQL/MariaDB, hooks |
| Base de datos | **MariaDB** (en host) | RDBMS relacional, corriendo en el VPS como servicio del sistema operativo |
| Migraciones | **goose** | Migraciones con SQL puro, versionadas. **Importante:** una sentencia por bloque `StatementBegin/StatementEnd` (requerimiento de MariaDB) |
| Auth | **golang-jwt/jwt v5** | JWT estándar, sin dependencias extra |
| Config | **godotenv** | Variables de entorno desde archivo `.env` |
| Contenerización | **Docker multi-stage + Docker Compose** | Build: `golang:1.25-alpine`; Runtime: `alpine:3.19`. Sin MySQL en Docker (usa host) |
| Deploy | **GitHub Actions + self-hosted runner** | Push a `main` dispara build y deploy automático en VPS |

## Estructura del proyecto

```
administracion-backend/
├── cmd/
│   └── server/
│       └── main.go              ← entrypoint + wiring completo de dependencias
├── internal/
│   ├── domain/                  ← entidades y ports (interfaces) — sin dependencias externas
│   │   ├── transferencia/
│   │   ├── empleado/
│   │   ├── proveedor/
│   │   ├── servicio/
│   │   ├── alquiler/
│   │   ├── tesoreria/
│   │   └── usuario/
│   ├── application/             ← casos de uso (orquestan dominio)
│   │   ├── transferencia/
│   │   ├── nomina/
│   │   ├── proveedor/
│   │   ├── servicio/
│   │   ├── alquiler/
│   │   ├── tesoreria/
│   │   ├── reporte/
│   │   └── usuario/
│   └── infrastructure/          ← adaptadores
│       ├── http/
│       │   ├── router.go
│       │   ├── middleware/
│       │   └── handler/
│       └── persistence/
│           └── repository/      ← modelos GORM + implementaciones de repos
├── config/
│   └── config.go
├── migrations/
│   └── 001_initial_schema.sql   ← schema completo en un único archivo versionado con goose
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

## Decisiones de implementación relevantes

- Los modelos GORM viven en el mismo paquete que los repos (`repository/`), no en un paquete `model/` separado.
- Las migraciones goose se ejecutan automáticamente al iniciar el servidor (`runMigrations` en `main.go`), usando `runtime.Caller(0)` para resolver la ruta relativa al binario compilado.
- El contenedor Docker conecta a la BD del host mediante `extra_hosts: host.docker.internal:host-gateway` en `docker-compose.yml`.

## Consecuencias
- Positivo: binario Go compilado = despliegue simple sin runtime externo.
- Positivo: goose con SQL puro = migraciones legibles y versionadas.
- Positivo: arquitectura hexagonal = cada módulo es independiente y testeable.
- A tener en cuenta: Go no tiene excepciones — el manejo de errores es explícito (`error` return).
- A tener en cuenta: la cookie JWT tiene flag `Secure`; el frontend debe acceder por HTTPS.
