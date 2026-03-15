# IPNEXT Backend — Visión General

## Descripción

API REST que sirve al frontend de IPNEXT Admin. Expone los datos y operaciones de todos los módulos del sistema: transferencias, nóminas, proveedores, servicios, alquileres, tesorería y reportes.

El backend es la única capa de seguridad real: valida JWT, verifica permisos por módulo y nivel, y persiste todos los datos en MySQL.

## Responsabilidades

- Autenticación y autorización (JWT + roles + permisos por módulo)
- CRUD completo de todas las entidades del negocio
- Lógica de negocio: cálculo de liquidaciones, proyección de flujo de caja, alertas de vencimiento
- Generación de reportes exportables (PDF / CSV)
- Validación de datos de entrada en todos los endpoints

## Stack

| Capa | Tecnología |
|------|-----------|
| Lenguaje | Go 1.22+ |
| Framework HTTP | Gin |
| ORM | GORM |
| Base de datos | MySQL 8 |
| Autenticación | JWT (httpOnly cookie) |
| Migraciones | GORM AutoMigrate + goose |
| Contenerización | Docker + Docker Compose |

## Módulos / Dominios

| Dominio | Ruta base | Descripción |
|---------|-----------|-------------|
| Auth | `/api/auth` | Login, logout, me |
| Dashboard | `/api/dashboard` | KPIs y resumen ejecutivo |
| Transferencias | `/api/transferencias` | Pagos, recurrentes, calendario |
| Nóminas | `/api/nominas` | Empleados, liquidaciones, guardias, compensaciones |
| Proveedores | `/api/proveedores` | Directorio, contratos, ranking |
| Servicios | `/api/servicios` | Servicios y utilities por categoría |
| Alquileres | `/api/alquileres` | Inmuebles, contratos, pagos, vencimientos |
| Tesorería | `/api/tesoreria` | Flujo de caja, cuentas bancarias, conciliación |
| Reportes | `/api/reportes` | Generación y exportación de informes |

## Reglas de negocio globales

- Todos los endpoints (excepto `/api/auth/login`) requieren JWT válido.
- El backend valida el permiso del usuario sobre el módulo en cada request — no se confía en el frontend.
- Los sub-usuarios con nivel `lectura` solo pueden usar métodos `GET`.
- Los sub-usuarios con nivel `escritura` pueden usar `GET`, `POST`, `PUT`, `PATCH`.
- Solo el rol `admin` puede usar `DELETE` y aprobar liquidaciones.
- Las respuestas de error siguen el formato estándar definido en ADR-002.
- Todas las fechas viajan en formato ISO 8601 (`YYYY-MM-DDTHH:mm:ssZ`).
- Los montos son `float64` en JSON; la DB los almacena como `DECIMAL(15,2)`.

## Relación con el frontend

El backend sirve exclusivamente al frontend de IPNEXT Admin (SPA React). Los contratos de datos (tipos de respuesta JSON) deben mantenerse alineados con los tipos TypeScript definidos en `TDR-002` del frontend. Cualquier cambio de schema debe coordinarse.
