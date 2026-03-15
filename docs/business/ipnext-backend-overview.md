# IPNEXT Backend — Visión General

## Descripción

API REST que sirve al frontend de IPNEXT Admin. Expone los datos y operaciones de todos los módulos del sistema: transferencias, nóminas, proveedores, servicios, alquileres, tesorería, reportes y gestión de usuarios.

El backend es la única capa de seguridad real: valida JWT, verifica permisos por módulo y nivel, y persiste todos los datos en MySQL/MariaDB.

## Responsabilidades

- Autenticación y autorización (JWT + roles + permisos por módulo)
- CRUD completo de todas las entidades del negocio
- Lógica de negocio: cálculo de liquidaciones, proyección de flujo de caja, alertas de vencimiento
- Generación de reportes exportables (CSV)
- Validación de datos de entrada en todos los endpoints
- Migraciones automáticas al iniciar (goose)

## Stack

| Capa | Tecnología |
|------|-----------|
| Lenguaje | Go 1.25 |
| Framework HTTP | Gin |
| ORM | GORM |
| Base de datos | MySQL 8 / MariaDB (instancia en host, no containerizada) |
| Autenticación | JWT (httpOnly cookie, flag `Secure` activo) |
| Migraciones | goose (SQL puro, una sentencia por bloque StatementBegin) |
| Contenerización | Docker multi-stage (golang:1.25-alpine → alpine:3.19) |
| Deploy | Docker Compose + GitHub Actions (self-hosted runner en VPS) |

## Módulos / Dominios

| Dominio | Ruta base | Descripción |
|---------|-----------|-------------|
| Auth | `/api/v1/auth` | Login, logout, me |
| Dashboard | `/api/v1/dashboard` | KPIs, pagos urgentes, distribución de egresos, actividad reciente |
| Transferencias | `/api/v1/transferencias` | Pagos, recurrentes, calendario, cambio de estado |
| Nóminas | `/api/v1/nominas` | Empleados, liquidaciones, guardias, compensaciones |
| Proveedores | `/api/v1/proveedores` | Directorio, contratos, ranking por monto pagado |
| Servicios | `/api/v1/servicios` | Servicios y utilities por categoría |
| Alquileres | `/api/v1/alquileres` | Inmuebles, contratos, pagos, vencimientos |
| Tesorería | `/api/v1/tesoreria` | Flujo de caja, cuentas bancarias, movimientos, conciliación, proyecciones |
| Reportes | `/api/v1/reportes` | Financiero, nómina, proveedores, inmuebles, exportación CSV |
| Usuarios | `/api/v1/usuarios` | Gestión de sub-usuarios (solo admin) |

## Reglas de negocio globales

- Todos los endpoints (excepto `/api/v1/auth/login`) requieren JWT válido en cookie `token`.
- El backend valida el permiso del usuario sobre el módulo en cada request — no se confía en el frontend.
- Los sub-usuarios con nivel `lectura` solo pueden usar métodos `GET`.
- Los sub-usuarios con nivel `escritura` pueden usar `GET`, `POST`, `PATCH`.
- Solo el rol `admin` puede usar `DELETE`, aprobar liquidaciones y gestionar usuarios.
- Las respuestas de error siguen el formato estándar definido en ADR-002.
- Las fechas en request body van en formato `YYYY-MM-DD`; las respuestas devuelven RFC3339.
- Los montos son `float64` en JSON; la DB los almacena como `DECIMAL(15,2)`.

## Infraestructura en producción

- **VPS:** `190.7.234.37`, puerto SSH `2222`
- **Puerto del backend:** `8288`
- **Base de datos:** MariaDB corriendo en el host del VPS (no en Docker). Accedida desde el contenedor via `host.docker.internal`.
- **Runner CI/CD:** GitHub Actions self-hosted runner con usuario `github-runner`.
- **Env file de producción:** `/home/github-runner/ipnext.env` (fuera del repo).
- **Secret en GitHub Actions:** `ENV_FILE_PATH` apunta al .env de producción.

## Usuario administrador seed

| Campo | Valor |
|-------|-------|
| Email | `admin@ipnext.com` |
| Password | `password` |
| Rol | `admin` |

## Nota sobre HTTPS y cookies

La cookie JWT tiene el flag `Secure` activo. Para que el login funcione correctamente desde el frontend, **el backend debe servirse por HTTPS**. Con HTTP, el browser descarta la cookie y el usuario queda sin sesión.

## Relación con el frontend

El backend sirve exclusivamente al frontend de IPNEXT Admin (SPA React). Los contratos de datos (tipos de respuesta JSON) deben mantenerse alineados con los tipos TypeScript del frontend. Cualquier cambio de schema debe coordinarse.
