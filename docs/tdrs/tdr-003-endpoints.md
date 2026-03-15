# TDR-003: Catálogo de Endpoints

## Estado
Propuesto

## Contexto
Definir todos los endpoints de la API agrupados por módulo. Cada endpoint indica método, ruta, permiso requerido y qué pantalla del frontend lo consume.

---

## Auth

| Método | Ruta | Permiso | Descripción | Pantalla frontend |
|--------|------|---------|-------------|-------------------|
| POST | `/api/v1/auth/login` | — | Login con email/password | `/login` |
| POST | `/api/v1/auth/logout` | autenticado | Invalida cookie | cualquiera |
| GET | `/api/v1/auth/me` | autenticado | Usuario actual + permisos | app init |

---

## Dashboard

| Método | Ruta | Permiso | Descripción | Pantalla frontend |
|--------|------|---------|-------------|-------------------|
| GET | `/api/v1/dashboard/kpis` | dashboard | KPIs: total pagos mes, pendientes, vencidos, flujo | `/dashboard` |
| GET | `/api/v1/dashboard/pagos-urgentes` | dashboard | Pagos próximos a vencer (próximos 7 días) | `/dashboard` |
| GET | `/api/v1/dashboard/distribucion-egresos` | dashboard | Egresos agrupados por categoría (para gráfico) | `/dashboard` |
| GET | `/api/v1/dashboard/actividad-reciente` | dashboard | Últimas N acciones del sistema | `/dashboard` |

---

## Transferencias

| Método | Ruta | Permiso | Descripción | Pantalla frontend |
|--------|------|---------|-------------|-------------------|
| GET | `/api/v1/transferencias` | transferencias:lectura | Lista paginada con filtros | `/transferencias` |
| GET | `/api/v1/transferencias/:id` | transferencias:lectura | Detalle | panel lateral |
| POST | `/api/v1/transferencias` | transferencias:escritura | Crear transferencia | `/transferencias/nueva` |
| PATCH | `/api/v1/transferencias/:id` | transferencias:escritura | Actualizar parcial | formulario edición |
| DELETE | `/api/v1/transferencias/:id` | admin_only | Eliminar | — |
| GET | `/api/v1/transferencias/calendario` | transferencias:lectura | Pagos agrupados por fecha del mes | `/transferencias/calendario` |
| GET | `/api/v1/transferencias/recurrentes` | transferencias:lectura | Solo transferencias con frecuencia != manual | `/transferencias/recurrentes` |
| PATCH | `/api/v1/transferencias/:id/estado` | transferencias:escritura | Cambiar estado (ej: marcar como pagado) | acciones tabla |

Query params para lista: `?page`, `?per_page`, `?estado`, `?categoria`, `?desde`, `?hasta`, `?q`, `?orden`, `?dir`

---

## Nóminas

| Método | Ruta | Permiso | Descripción | Pantalla frontend |
|--------|------|---------|-------------|-------------------|
| GET | `/api/v1/nominas/empleados` | nominas:lectura | Lista de empleados activos | `/nominas` |
| GET | `/api/v1/nominas/empleados/:id` | nominas:lectura | Detalle empleado | panel detalle |
| POST | `/api/v1/nominas/empleados` | nominas:escritura | Crear empleado | modal nuevo empleado |
| PATCH | `/api/v1/nominas/empleados/:id` | nominas:escritura | Actualizar empleado | formulario edición |
| DELETE | `/api/v1/nominas/empleados/:id` | admin_only | Dar de baja | — |
| GET | `/api/v1/nominas/empleados/kpis` | nominas:lectura | KPIs nómina: total empleados, costo total, etc. | `/nominas` |
| GET | `/api/v1/nominas/liquidaciones` | nominas:lectura | Lista liquidaciones por período | `/nominas/liquidacion` |
| POST | `/api/v1/nominas/liquidaciones` | nominas:escritura | Generar liquidación del período | botón "Liquidar Nómina" |
| POST | `/api/v1/nominas/liquidaciones/:id/aprobar` | admin_only | Aprobar liquidación | acción tabla |
| GET | `/api/v1/nominas/guardias` | nominas:lectura | Lista de guardias | `/nominas/guardias` |
| POST | `/api/v1/nominas/guardias` | nominas:escritura | Registrar guardia | formulario |
| GET | `/api/v1/nominas/compensaciones` | nominas:lectura | Lista compensaciones adicionales | `/nominas/compensaciones` |
| POST | `/api/v1/nominas/compensaciones` | nominas:escritura | Crear compensación | formulario |

---

## Proveedores

| Método | Ruta | Permiso | Descripción | Pantalla frontend |
|--------|------|---------|-------------|-------------------|
| GET | `/api/v1/proveedores` | proveedores:lectura | Lista paginada con buscador | `/proveedores` |
| GET | `/api/v1/proveedores/:id` | proveedores:lectura | Detalle + historial pagos | panel lateral |
| POST | `/api/v1/proveedores` | proveedores:escritura | Crear proveedor | modal |
| PATCH | `/api/v1/proveedores/:id` | proveedores:escritura | Actualizar proveedor | formulario |
| DELETE | `/api/v1/proveedores/:id` | admin_only | Eliminar | — |
| GET | `/api/v1/proveedores/contratos` | proveedores:lectura | Lista contratos activos/vencidos | `/proveedores/contratos` |
| POST | `/api/v1/proveedores/contratos` | proveedores:escritura | Crear contrato | formulario |
| PATCH | `/api/v1/proveedores/contratos/:id` | proveedores:escritura | Actualizar contrato | formulario |
| GET | `/api/v1/proveedores/ranking` | proveedores:lectura | Proveedores ordenados por monto total pagado | `/proveedores/ranking` |

---

## Servicios

| Método | Ruta | Permiso | Descripción | Pantalla frontend |
|--------|------|---------|-------------|-------------------|
| GET | `/api/v1/servicios` | servicios:lectura | Resumen general: todos los servicios agrupados | `/servicios` |
| GET | `/api/v1/servicios/kpis` | servicios:lectura | KPIs: costo total mensual, por vencer, etc. | `/servicios` |
| GET | `/api/v1/servicios/:tipo` | servicios:lectura | Servicios filtrados por tipo | `/servicios/internet`, etc. |
| GET | `/api/v1/servicios/item/:id` | servicios:lectura | Detalle servicio | panel detalle |
| POST | `/api/v1/servicios` | servicios:escritura | Crear servicio | formulario |
| PATCH | `/api/v1/servicios/item/:id` | servicios:escritura | Actualizar servicio | formulario |
| DELETE | `/api/v1/servicios/item/:id` | admin_only | Eliminar | — |

Tipos válidos para `:tipo`: `internet`, `energia`, `seguridad`, `software`, `obra_social`, `seguro`

---

## Alquileres

| Método | Ruta | Permiso | Descripción | Pantalla frontend |
|--------|------|---------|-------------|-------------------|
| GET | `/api/v1/alquileres` | alquileres:lectura | Lista inmuebles activos | `/alquileres` |
| GET | `/api/v1/alquileres/:id` | alquileres:lectura | Detalle inmueble | panel detalle |
| POST | `/api/v1/alquileres` | alquileres:escritura | Crear inmueble | modal |
| PATCH | `/api/v1/alquileres/:id` | alquileres:escritura | Actualizar inmueble | formulario |
| DELETE | `/api/v1/alquileres/:id` | admin_only | Eliminar | — |
| GET | `/api/v1/alquileres/contratos` | alquileres:lectura | Contratos vigentes y vencidos | `/alquileres/contratos` |
| POST | `/api/v1/alquileres/contratos` | alquileres:escritura | Crear contrato | formulario |
| GET | `/api/v1/alquileres/pagos` | alquileres:lectura | Historial pagos con estado | `/alquileres/pagos` |
| POST | `/api/v1/alquileres/pagos` | alquileres:escritura | Registrar pago | botón "Registrar Pago" |
| GET | `/api/v1/alquileres/vencimientos` | alquileres:lectura | Contratos próximos a vencer | `/alquileres/vencimientos` |

---

## Tesorería

| Método | Ruta | Permiso | Descripción | Pantalla frontend |
|--------|------|---------|-------------|-------------------|
| GET | `/api/v1/tesoreria/flujo-caja` | tesoreria:lectura | Proyección ingresos/egresos por período | `/tesoreria` |
| GET | `/api/v1/tesoreria/cuentas` | tesoreria:lectura | Lista cuentas bancarias con saldo | `/tesoreria/cuentas` |
| POST | `/api/v1/tesoreria/cuentas` | tesoreria:escritura | Agregar cuenta bancaria | formulario |
| PATCH | `/api/v1/tesoreria/cuentas/:id` | tesoreria:escritura | Actualizar saldo o datos | formulario |
| GET | `/api/v1/tesoreria/conciliacion` | tesoreria:lectura | Movimientos pendientes de conciliar | `/tesoreria/conciliacion` |
| POST | `/api/v1/tesoreria/movimientos` | tesoreria:escritura | Registrar movimiento bancario | botón |
| PATCH | `/api/v1/tesoreria/movimientos/:id/conciliar` | tesoreria:escritura | Marcar movimiento como conciliado | acción |
| GET | `/api/v1/tesoreria/proyecciones` | tesoreria:lectura | Proyección de liquidez a N meses | `/tesoreria/proyecciones` |

---

## Reportes

| Método | Ruta | Permiso | Descripción | Pantalla frontend |
|--------|------|---------|-------------|-------------------|
| GET | `/api/v1/reportes/financiero` | reportes:lectura | Reporte egresos/ingresos del período | `/reportes` |
| GET | `/api/v1/reportes/nomina` | reportes:lectura + nominas:lectura | Reporte liquidaciones | `/reportes/nomina` |
| GET | `/api/v1/reportes/proveedores` | reportes:lectura + proveedores:lectura | Reporte pagos a proveedores | `/reportes/proveedores` |
| GET | `/api/v1/reportes/inmuebles` | reportes:lectura + alquileres:lectura | Reporte alquileres | `/reportes/inmuebles` |
| GET | `/api/v1/reportes/exportar` | reportes:lectura | Exportar datos en CSV o PDF | `/reportes/exportar` |

Query params de exportación: `?formato=csv` o `?formato=pdf`, `?desde=`, `?hasta=`, `?modulo=`

---

## Usuarios (gestión de sub-usuarios)

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/usuarios` | admin_only | Lista sub-usuarios |
| POST | `/api/v1/usuarios` | admin_only | Crear sub-usuario con permisos |
| PATCH | `/api/v1/usuarios/:id` | admin_only | Actualizar permisos |
| DELETE | `/api/v1/usuarios/:id` | admin_only | Desactivar usuario |
