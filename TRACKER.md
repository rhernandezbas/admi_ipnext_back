# IPNEXT Backend — Tracker de Implementación

## Estado general

| Módulo | Estado | Iteraciones |
|--------|--------|-------------|
| Setup inicial | ✅ completo | 1–8 |
| Auth | ✅ completo | 19–23 |
| Dashboard | ✅ completo | 24–27 |
| Transferencias | ✅ completo | 28–30 |
| Nóminas | ✅ completo | 31–35 |
| Proveedores | ✅ completo | 36–38 |
| Servicios | ✅ completo | 39 |
| Alquileres | ✅ completo | 40 |
| Tesorería | ✅ completo | 41 |
| Reportes | ✅ completo | 42 |
| Usuarios | ✅ completo | 42 |

---

## Última iteración

> (Ralph completa esta sección al final de cada iteración)

- **Iteración**: 42
- **Tarea completada**: 12.1–12.6 + 13.1–13.2 + 14.1–14.2 — Reportes, Usuarios y Migraciones completos
- **Archivos creados/modificados**: `financiero.go`, `nomina.go`, `proveedores.go`, `inmuebles.go`, `csv.go`, `reporte_handler.go`, `usuario_usecases.go`, `usuario_handler.go`, `001_initial_schema.sql`, `main.go`
- **Estado de compilación**: ✅ OK (`go build ./...`)
- **Próxima tarea**: BACKEND COMPLETO ✅

---

## Tareas

### Fase 1: Setup del proyecto

- [x] **1.1** Inicializar módulo Go: `go mod init github.com/ipnext/admin-backend`
- [x] **1.2** Crear estructura de carpetas base (`cmd/`, `internal/domain/`, `internal/application/`, `internal/infrastructure/`, `config/`, `migrations/`)
- [x] **1.3** Agregar dependencias: Gin, GORM, GORM MySQL driver, golang-jwt, google/uuid, godotenv, zap
- [x] **1.4** Crear `config/config.go` — struct Config + carga desde .env
- [x] **1.5** Crear `internal/infrastructure/config/database.go` — conexión MySQL con GORM
- [x] **1.6** Crear `cmd/server/main.go` — entrypoint básico que levanta el servidor en el puerto configurado
- [x] **1.7** Crear `docker-compose.yml` con servicio MySQL + el backend
- [x] **1.8** Crear `.env.example` con todas las variables necesarias

### Fase 2: Dominio — Entidades e Interfaces

- [x] **2.1** `internal/domain/usuario/` — entity.go + repository.go (interface)
- [x] **2.2** `internal/domain/transferencia/` — entity.go + repository.go
- [x] **2.3** `internal/domain/empleado/` — entity.go (Empleado + Liquidacion + Guardia + Compensacion) + repository.go
- [x] **2.4** `internal/domain/proveedor/` — entity.go (Proveedor + ContratoProveedor) + repository.go
- [x] **2.5** `internal/domain/servicio/` — entity.go + repository.go
- [x] **2.6** `internal/domain/alquiler/` — entity.go (Inmueble + ContratoAlquiler + PagoAlquiler) + repository.go
- [x] **2.7** `internal/domain/tesoreria/` — entity.go (CuentaBancaria + MovimientoBancario) + repository.go

### Fase 3: Infraestructura — Middleware y Router base

- [x] **3.1** `internal/infrastructure/http/middleware/auth.go` — JWT parsing + inyección de claims en context
- [x] **3.2** `internal/infrastructure/http/middleware/permiso.go` — RequirePermiso(modulo, nivel)
- [x] **3.3** `internal/infrastructure/http/router.go` — router base con grupos de rutas y middleware aplicado

### Fase 4: Auth

- [x] **4.1** `internal/infrastructure/persistence/repository/usuario_repo.go` — implementación MySQL
- [x] **4.2** `internal/application/auth/login.go` — caso de uso Login (valida credenciales, firma JWT)
- [x] **4.3** `internal/application/auth/me.go` — caso de uso GetMe
- [x] **4.4** `internal/infrastructure/http/handler/auth_handler.go` — POST /login, POST /logout, GET /me
- [x] **4.5** Registrar rutas de auth en router + wiring en main.go

### Fase 5: Dashboard

- [x] **5.1** `internal/application/dashboard/get_kpis.go` — agrega KPIs desde transferencias, empleados, etc.
- [x] **5.2** `internal/application/dashboard/get_pagos_urgentes.go` — pagos próximos 7 días
- [x] **5.3** `internal/application/dashboard/get_distribucion_egresos.go` — agrupación por categoría
- [x] **5.4** `internal/infrastructure/http/handler/dashboard_handler.go` + rutas

### Fase 6: Transferencias

- [x] **6.1** `internal/infrastructure/persistence/repository/transferencia_repo.go`
- [x] **6.2** Casos de uso: list, get, create, update, delete, cambiar estado
- [x] **6.3** `internal/infrastructure/http/handler/transferencia_handler.go` + rutas
- [x] **6.4** Endpoint calendario (agrupado por fecha)
- [x] **6.5** Endpoint recurrentes (frecuencia != manual)

### Fase 7: Nóminas

- [x] **7.1** `internal/infrastructure/persistence/repository/empleado_repo.go`
- [x] **7.2** Casos de uso empleados: list, get, create, update, delete, kpis
- [x] **7.3** Casos de uso liquidaciones: list, create, aprobar
- [x] **7.4** Casos de uso guardias: list, create
- [x] **7.5** Casos de uso compensaciones: list, create
- [x] **7.6** `internal/infrastructure/http/handler/nomina_handler.go` + rutas

### Fase 8: Proveedores

- [x] **8.1** `internal/infrastructure/persistence/repository/proveedor_repo.go`
- [x] **8.2** Casos de uso proveedores: list, get, create, update, delete
- [x] **8.3** Casos de uso contratos: list, create, update
- [x] **8.4** Endpoint ranking (agrupado por monto total)
- [x] **8.5** `internal/infrastructure/http/handler/proveedor_handler.go` + rutas

### Fase 9: Servicios

- [x] **9.1** `internal/infrastructure/persistence/repository/servicio_repo.go`
- [x] **9.2** Casos de uso: list (con filtro por tipo), get, create, update, delete, kpis
- [x] **9.3** `internal/infrastructure/http/handler/servicio_handler.go` + rutas

### Fase 10: Alquileres

- [x] **10.1** `internal/infrastructure/persistence/repository/alquiler_repo.go`
- [x] **10.2** Casos de uso inmuebles: list, get, create, update, delete
- [x] **10.3** Casos de uso contratos: list, create
- [x] **10.4** Casos de uso pagos: list, create
- [x] **10.5** Endpoint vencimientos (contratos próximos a vencer)
- [x] **10.6** `internal/infrastructure/http/handler/alquiler_handler.go` + rutas

### Fase 11: Tesorería

- [x] **11.1** `internal/infrastructure/persistence/repository/tesoreria_repo.go`
- [x] **11.2** Caso de uso flujo de caja (proyección ingresos/egresos)
- [x] **11.3** Casos de uso cuentas: list, create, update
- [x] **11.4** Casos de uso movimientos: list conciliación, create, conciliar
- [x] **11.5** Endpoint proyecciones (liquidez futura)
- [x] **11.6** `internal/infrastructure/http/handler/tesoreria_handler.go` + rutas

### Fase 12: Reportes

- [x] **12.1** `internal/application/reporte/financiero.go`
- [x] **12.2** `internal/application/reporte/nomina.go`
- [x] **12.3** `internal/application/reporte/proveedores.go`
- [x] **12.4** `internal/application/reporte/inmuebles.go`
- [x] **12.5** Exportación CSV (función utilitaria)
- [x] **12.6** `internal/infrastructure/http/handler/reporte_handler.go` + rutas

### Fase 13: Usuarios (gestión de sub-usuarios)

- [x] **13.1** Casos de uso: list, create, update permisos, delete
- [x] **13.2** `internal/infrastructure/http/handler/usuario_handler.go` + rutas (todos admin_only)

### Fase 14: Migraciones

- [x] **14.1** Crear archivos SQL de migración en `migrations/` para todas las tablas (basado en TDR-002)
- [x] **14.2** Setup de goose en main.go para correr migraciones al iniciar

---

## Notas de implementación

- UUID generation: `github.com/google/uuid`
- Password hashing: `golang.org/x/crypto/bcrypt`
- JWT: `github.com/golang-jwt/jwt/v5`
- Fechas: siempre `time.Time` en Go, ISO 8601 en JSON
- Montos: `float64` en JSON, `DECIMAL(15,2)` en MySQL
- Errores: siempre retornar `error`, no panic
- Logs: usar `zap.Logger` inyectado por dependencia
