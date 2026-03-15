# TDR-001: Estructura del Proyecto (Arquitectura Hexagonal)

## Estado
Propuesto

## Contexto
El backend usa arquitectura hexagonal (ports & adapters). Necesitamos definir la estructura de carpetas, los roles de cada capa y las convenciones de nombrado para que el proyecto sea mantenible y coherente.

## Decisión

### Estructura de carpetas

```
ipnext-backend/
├── cmd/
│   └── server/
│       └── main.go                    ← entrypoint: wiring de dependencias, arranque
├── internal/
│   ├── domain/                        ← capa de dominio (sin dependencias externas)
│   │   ├── transferencia/
│   │   │   ├── entity.go              ← struct Transferencia + métodos de dominio
│   │   │   ├── repository.go          ← interface TransferenciaRepository (port)
│   │   │   └── service.go             ← interface TransferenciaService (port)
│   │   ├── empleado/
│   │   ├── proveedor/
│   │   ├── servicio/
│   │   ├── alquiler/
│   │   ├── tesoreria/
│   │   └── usuario/
│   ├── application/                   ← casos de uso (orquestan dominio)
│   │   ├── transferencia/
│   │   │   ├── list_transferencias.go
│   │   │   ├── create_transferencia.go
│   │   │   ├── update_transferencia.go
│   │   │   └── delete_transferencia.go
│   │   ├── nomina/
│   │   ├── proveedor/
│   │   └── ...
│   └── infrastructure/                ← adaptadores (implementan los puertos)
│       ├── http/                      ← adaptador HTTP (Gin)
│       │   ├── router.go              ← registro de todas las rutas
│       │   ├── middleware/
│       │   │   ├── auth.go
│       │   │   └── permiso.go
│       │   └── handler/               ← un handler por módulo
│       │       ├── auth_handler.go
│       │       ├── transferencia_handler.go
│       │       ├── nomina_handler.go
│       │       ├── proveedor_handler.go
│       │       ├── servicio_handler.go
│       │       ├── alquiler_handler.go
│       │       ├── tesoreria_handler.go
│       │       ├── reporte_handler.go
│       │       └── dashboard_handler.go
│       ├── persistence/               ← adaptador DB (GORM + MySQL)
│       │   ├── model/                 ← structs GORM (pueden diferir de las entidades)
│       │   │   ├── transferencia_model.go
│       │   │   └── ...
│       │   └── repository/            ← implementaciones de los ports
│       │       ├── transferencia_repo.go
│       │       └── ...
│       └── config/
│           └── database.go            ← conexión MySQL
├── config/
│   ├── config.go                      ← struct Config + carga desde env
│   └── .env.example
├── migrations/
│   ├── 001_create_usuarios.sql
│   ├── 002_create_transferencias.sql
│   └── ...
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

### Reglas de capas

| Capa | Puede importar | No puede importar |
|------|---------------|-------------------|
| `domain/` | solo stdlib Go | `application/`, `infrastructure/` |
| `application/` | `domain/` | `infrastructure/` (solo interfaces) |
| `infrastructure/` | `domain/`, `application/`, libs externas | nada prohibido |

### Convenciones de nombrado

- Entidades de dominio: `PascalCase` struct, archivo `entity.go`
- Ports (interfaces): sufijo `Repository` o `Service` → `TransferenciaRepository`
- Adaptadores (implementaciones): prefijo `MySQL` o `GORM` → `MySQLTransferenciaRepository`
- Handlers: sufijo `Handler` → `TransferenciaHandler`
- Casos de uso: archivo descriptivo → `create_transferencia.go`, struct `CreateTransferenciaUseCase`
- DTOs de request/response: en el handler, sufijo `Request` / `Response`

### Inyección de dependencias

Wiring manual en `cmd/server/main.go`:
```go
// Repositorios
transRepo := persistence.NewMySQLTransferenciaRepository(db)

// Casos de uso
createTransferencia := application.NewCreateTransferenciaUseCase(transRepo)

// Handlers
transHandler := handler.NewTransferenciaHandler(createTransferencia, ...)

// Router
router.SetupRoutes(transHandler, ...)
```

Sin frameworks DI — Go idiomático con constructores.

## Consecuencias
- Positivo: dominio sin dependencias externas = testeable con mocks puros.
- Positivo: cambiar DB (ej: a PostgreSQL) = solo tocar `infrastructure/persistence/`.
- Positivo: agregar un módulo nuevo = seguir el mismo patrón sin afectar otros.
- A tener en cuenta: más archivos que MVC clásico — compensado por claridad y testabilidad.
