# TDR-001: Estructura del Proyecto (Arquitectura Hexagonal)

## Estado
Implementado

## Contexto
El backend usa arquitectura hexagonal (ports & adapters). Este documento define la estructura de carpetas implementada, los roles de cada capa y las convenciones de nombrado.

## Estructura de carpetas (implementada)

```
administracion-backend/
├── cmd/
│   └── server/
│       └── main.go                    ← entrypoint: wiring manual de dependencias + runMigrations
├── internal/
│   ├── domain/                        ← capa de dominio (sin dependencias externas)
│   │   ├── transferencia/
│   │   │   ├── entity.go              ← struct Transferencia + métodos de dominio
│   │   │   └── repository.go          ← interface Repository (port)
│   │   ├── empleado/
│   │   ├── proveedor/
│   │   ├── servicio/
│   │   ├── alquiler/
│   │   ├── tesoreria/
│   │   └── usuario/
│   ├── application/                   ← casos de uso (orquestan dominio)
│   │   ├── transferencia/
│   │   │   └── transferencia_usecases.go   ← todos los use cases del módulo en un archivo
│   │   ├── nomina/
│   │   │   ├── empleado_usecases.go
│   │   │   ├── liquidacion_usecases.go
│   │   │   └── guardia_compensacion_usecases.go
│   │   ├── proveedor/
│   │   │   └── proveedor_usecases.go
│   │   ├── servicio/
│   │   ├── alquiler/
│   │   │   └── alquiler_usecases.go
│   │   ├── tesoreria/
│   │   │   └── tesoreria_usecases.go
│   │   ├── reporte/
│   │   │   ├── financiero.go
│   │   │   ├── nomina.go
│   │   │   ├── proveedores.go
│   │   │   ├── inmuebles.go
│   │   │   └── csv.go
│   │   └── usuario/
│   │       └── usuario_usecases.go
│   └── infrastructure/                ← adaptadores
│       ├── http/
│       │   ├── router.go              ← registro de todas las rutas (SetupRouter)
│       │   ├── middleware/
│       │   │   ├── auth.go            ← AuthMiddleware (valida JWT de cookie)
│       │   │   └── permiso.go         ← RequirePermiso(modulo, nivel)
│       │   └── handler/               ← un handler por módulo
│       │       ├── auth_handler.go
│       │       ├── dashboard_handler.go
│       │       ├── transferencia_handler.go
│       │       ├── nomina_handler.go
│       │       ├── proveedor_handler.go
│       │       ├── servicio_handler.go
│       │       ├── alquiler_handler.go
│       │       ├── tesoreria_handler.go
│       │       ├── reporte_handler.go
│       │       └── usuario_handler.go
│       └── persistence/
│           └── repository/            ← modelos GORM + implementaciones de repos (en el mismo paquete)
│               ├── transferencia_repo.go
│               ├── empleado_repo.go
│               ├── proveedor_repo.go
│               ├── servicio_repo.go
│               ├── alquiler_repo.go
│               ├── tesoreria_repo.go
│               └── usuario_repo.go
├── config/
│   └── config.go                      ← struct Config + carga desde env (godotenv)
├── migrations/
│   └── 001_initial_schema.sql         ← schema completo, versionado con goose
├── .github/
│   └── workflows/
│       └── deploy.yml                 ← CI/CD: push a main → build Docker → deploy VPS
├── docker-compose.yml
├── Dockerfile                         ← multi-stage: golang:1.25-alpine → alpine:3.19
└── go.mod
```

## Reglas de capas

| Capa | Puede importar | No puede importar |
|------|---------------|-------------------|
| `domain/` | solo stdlib Go | `application/`, `infrastructure/` |
| `application/` | `domain/`, stdlib | `infrastructure/` |
| `infrastructure/` | `domain/`, `application/`, libs externas | — |

## Convenciones implementadas

- **Modelos GORM:** definidos en el mismo archivo que el repo que los usa (ej: `empleadoModel` en `empleado_repo.go`). Si hay discrepancia entre nombre Go y columna DB, se usa tag `gorm:"column:..."`.
- **Use cases:** todos los use cases de un módulo en un archivo (o pocos archivos por sub-dominio). Cada use case es una struct con método `Execute`.
- **Handlers:** un archivo por módulo. Los structs de request se definen localmente en el handler. Las fechas se reciben como `string` y se parsean con `time.Parse("2006-01-02", ...)`.
- **Router:** todas las rutas en `router.go`. Los grupos de rutas protegidas usan `middleware.AuthMiddleware` y `middleware.RequirePermiso`.

## Inyección de dependencias

Wiring manual en `cmd/server/main.go`. Sin frameworks DI. Orden: repos → use cases → handlers → router.

```go
// Repos
transRepo := repository.NewMySQLTransferenciaRepository(db)

// Use cases
listUC := transferencia.NewListUseCase(transRepo)
createUC := transferencia.NewCreateUseCase(transRepo)

// Handler
transHandler := handler.NewTransferenciaHandler(listUC, createUC, ...)

// Router
r := infrahttp.SetupRouter(handlers, cfg.JWT.Secret)
```

## Migraciones

Las migraciones se ejecutan automáticamente al iniciar el servidor:

```go
func runMigrations(dsn, migrationsDir string) {
    // goose.Up con el directorio de migrations
}
```

La ruta al directorio `migrations/` se resuelve en tiempo de ejecución con `runtime.Caller(0)` para ser compatible tanto en desarrollo como dentro del contenedor Docker.

## Consecuencias
- Positivo: dominio sin dependencias externas = testeable con mocks puros.
- Positivo: cambiar DB (ej: a PostgreSQL) = solo tocar `infrastructure/persistence/`.
- Positivo: agregar un módulo nuevo = seguir el mismo patrón sin afectar otros.
- A tener en cuenta: más archivos que MVC clásico — compensado por claridad y mantenibilidad.
