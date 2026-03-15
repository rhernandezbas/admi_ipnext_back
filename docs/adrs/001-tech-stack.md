# ADR-001: Tech Stack Backend

## Estado
Aceptado

## Contexto
Definir el stack tecnológico del backend de IPNEXT Admin antes de iniciar el desarrollo. El backend debe servir a una SPA React como única consumidora.

## Decisión

| Categoría | Tecnología | Justificación |
|-----------|-----------|---------------|
| Lenguaje | **Go 1.22+** | Rendimiento nativo, tipado estático, compilado, manejo de concurrencia idiomático, binario único fácil de desplegar |
| Framework HTTP | **Gin** | El más usado en Go, rápido, middleware ecosystem maduro, buena DX |
| ORM | **GORM** | ORM más popular en Go, soporte MySQL, migraciones, hooks, associations |
| Base de datos | **MySQL 8** | RDBMS relacional, soporta JSON columns, bien soportado por GORM |
| Migraciones | **goose** | Migraciones con SQL puro o Go, versionadas, compatible con CI |
| Auth | **golang-jwt/jwt v5** | JWT estándar, sin dependencias extra, bien mantenido |
| Config | **godotenv + viper** | Variables de entorno + archivo de config con override por entorno |
| Logging | **zap (uber-go)** | Structured logging, alto rendimiento, niveles de log |
| Testing | **testify** | Assert/require/mock estándar en el ecosistema Go |
| Contenerización | **Docker + Docker Compose** | Entorno reproducible, fácil de desplegar en cualquier servidor |
| Hot reload (dev) | **air** | Live reload para Go en desarrollo |

## Estructura de módulos Go

```
go.mod
cmd/
  server/
    main.go          ← entrypoint
internal/            ← código privado de la app (hexagonal)
  domain/            ← entidades, puertos (interfaces)
  application/       ← casos de uso
  infrastructure/    ← adaptadores: DB, HTTP, externos
config/
  config.go
migrations/          ← archivos .sql de goose
```

## Consecuencias
- Positivo: binario Go compilado = despliegue simple sin runtime externo.
- Positivo: GORM + goose = migraciones versionadas y reproducibles.
- Positivo: Gin + zap = logs estructurados listos para observabilidad.
- A tener en cuenta: Go no tiene excepciones — el manejo de errores es explícito (`error` return).
- A tener en cuenta: GORM `AutoMigrate` solo en desarrollo; en producción usar goose.
