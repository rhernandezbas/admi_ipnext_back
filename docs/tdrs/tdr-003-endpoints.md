# TDR-003: Catálogo de Endpoints

## Estado
Implementado — 48/48 endpoints verificados en producción

## Contexto
Catálogo de todos los endpoints de la API, agrupados por módulo. Incluye método, ruta, permiso requerido y descripción.

### Formato de fechas en request body
Todos los campos de fecha en el body de los requests deben enviarse como **`YYYY-MM-DD`** (ej: `"2026-03-15"`). Las respuestas devuelven fechas en **RFC3339** (ej: `"2026-03-15T00:00:00Z"`).

---

## Auth

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| POST | `/api/v1/auth/login` | — | Login. Body: `{"email","password"}`. Setea cookie `token` (HttpOnly, Secure) |
| POST | `/api/v1/auth/logout` | autenticado | Invalida cookie |
| GET | `/api/v1/auth/me` | autenticado | Usuario actual + permisos |

> **Nota:** La cookie tiene flag `Secure`. El frontend debe acceder por HTTPS, de lo contrario el browser descarta la cookie y el usuario queda deslogueado.

---

## Dashboard

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/dashboard/kpis` | dashboard | KPIs: total pagos mes, pendientes, vencidos, flujo |
| GET | `/api/v1/dashboard/pagos-urgentes` | dashboard | Pagos próximos a vencer (próximos 7 días) |
| GET | `/api/v1/dashboard/distribucion-egresos` | dashboard | Egresos agrupados por categoría |
| GET | `/api/v1/dashboard/actividad-reciente` | dashboard | Últimas acciones del sistema |

---

## Transferencias

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/transferencias` | transferencias:lectura | Lista con filtros |
| GET | `/api/v1/transferencias/:id` | transferencias:lectura | Detalle |
| POST | `/api/v1/transferencias` | transferencias:escritura | Crear transferencia |
| PATCH | `/api/v1/transferencias/:id` | transferencias:escritura | Actualizar |
| DELETE | `/api/v1/transferencias/:id` | admin_only | Eliminar (soft delete) |
| GET | `/api/v1/transferencias/calendario` | transferencias:lectura | Pagos agrupados por fecha del mes |
| GET | `/api/v1/transferencias/recurrentes` | transferencias:lectura | Solo frecuencia != manual |
| PATCH | `/api/v1/transferencias/:id/estado` | transferencias:escritura | Cambiar estado |

**Body POST/PATCH:**
```json
{
  "beneficiario": "string (required)",
  "cbu": "string (optional)",
  "alias": "string (optional)",
  "categoria": "string (required)",
  "monto": 5000,
  "moneda": "ARS|USD (required)",
  "fechaPago": "YYYY-MM-DD (required)",
  "frecuencia": "manual|mensual|... (required)",
  "metodoPago": "transferencia|debito|... (required)",
  "estado": "pendiente|pagado|... (optional en PATCH)",
  "notas": "string (optional)",
  "proveedorId": "uuid (optional)"
}
```

Query params para lista: `?page`, `?per_page`, `?estado`, `?categoria`, `?desde`, `?hasta`, `?q`

---

## Nóminas

### Empleados

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/nominas/empleados/kpis` | nominas:lectura | KPIs: total empleados, costo mensual, etc. |
| GET | `/api/v1/nominas/empleados` | nominas:lectura | Lista empleados |
| GET | `/api/v1/nominas/empleados/:id` | nominas:lectura | Detalle empleado |
| POST | `/api/v1/nominas/empleados` | nominas:escritura | Crear empleado |
| PATCH | `/api/v1/nominas/empleados/:id` | nominas:escritura | Actualizar empleado |
| DELETE | `/api/v1/nominas/empleados/:id` | admin_only | Dar de baja (activo=false) |

**Body POST empleado:**
```json
{
  "nombre": "string (required)",
  "puesto": "string (required)",
  "area": "string (required)",
  "rol": "string (required)",
  "sueldoBruto": 180000,
  "obraSocial": "string (required)",
  "fechaIngreso": "YYYY-MM-DD (required)",
  "avatar": "url (optional)"
}
```

### Liquidaciones

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/nominas/liquidaciones` | nominas:lectura | Lista. Query: `?periodo=YYYY-MM` |
| POST | `/api/v1/nominas/liquidaciones` | nominas:escritura | Crear liquidación |
| POST | `/api/v1/nominas/liquidaciones/:id/aprobar` | admin_only | Aprobar liquidación |

**Body POST liquidación:**
```json
{
  "empleadoId": "uuid (required)",
  "periodo": "YYYY-MM (required)",
  "sueldoBruto": 180000,
  "deducciones": 18000
}
```
El `netoAPagar` se calcula automáticamente (`sueldoBruto - deducciones`).

### Guardias

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/nominas/guardias` | nominas:lectura | Lista. Query: `?empleado_id` |
| POST | `/api/v1/nominas/guardias` | nominas:escritura | Registrar guardia |

**Body POST guardia:**
```json
{
  "empleadoId": "uuid (required)",
  "fecha": "YYYY-MM-DD (required)",
  "horas": 8,
  "monto": 5000
}
```

### Compensaciones

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/nominas/compensaciones` | nominas:lectura | Lista. Query: `?empleado_id` |
| POST | `/api/v1/nominas/compensaciones` | nominas:escritura | Crear compensación |

**Body POST compensación:**
```json
{
  "empleadoId": "uuid (required)",
  "tipo": "bono|adelanto|extra|otro (required)",
  "monto": 10000,
  "descripcion": "string (optional)",
  "fecha": "YYYY-MM-DD (required)"
}
```

---

## Proveedores

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/proveedores` | proveedores:lectura | Lista. Query: `?q`, `?solo_activos` |
| GET | `/api/v1/proveedores/contratos` | proveedores:lectura | Todos los contratos. Query: `?proveedor_id` |
| GET | `/api/v1/proveedores/ranking` | proveedores:lectura | Top proveedores por monto pagado |
| GET | `/api/v1/proveedores/:id` | proveedores:lectura | Detalle proveedor |
| POST | `/api/v1/proveedores` | proveedores:escritura | Crear proveedor |
| PATCH | `/api/v1/proveedores/:id` | proveedores:escritura | Actualizar proveedor |
| DELETE | `/api/v1/proveedores/:id` | admin_only | Desactivar |
| POST | `/api/v1/proveedores/contratos` | proveedores:escritura | Crear contrato |
| PATCH | `/api/v1/proveedores/contratos/:id` | proveedores:escritura | Actualizar estado contrato |

**Body POST contrato proveedor:**
```json
{
  "proveedorId": "uuid (required)",
  "descripcion": "string (optional)",
  "vigenciaDesde": "YYYY-MM-DD (required)",
  "vigenciaHasta": "YYYY-MM-DD (required)",
  "montoAnual": 8000
}
```
El `codigo` se auto-genera con formato `CTR-{año}-{seq}`.

---

## Servicios

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/servicios` | servicios:lectura | Lista todos los servicios |
| GET | `/api/v1/servicios/kpis` | servicios:lectura | KPIs: costo total, por vencer, etc. |
| GET | `/api/v1/servicios/:tipo` | servicios:lectura | Servicios filtrados por tipo |
| GET | `/api/v1/servicios/item/:id` | servicios:lectura | Detalle servicio |
| POST | `/api/v1/servicios` | servicios:escritura | Crear servicio |
| PATCH | `/api/v1/servicios/item/:id` | servicios:escritura | Actualizar servicio |
| DELETE | `/api/v1/servicios/item/:id` | admin_only | Eliminar |

Tipos válidos: `internet`, `energia`, `seguridad`, `software`, `obra_social`, `seguro`

---

## Alquileres

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/alquileres` | alquileres:lectura | Lista inmuebles |
| GET | `/api/v1/alquileres/contratos` | alquileres:lectura | Contratos de alquiler |
| GET | `/api/v1/alquileres/pagos` | alquileres:lectura | Historial de pagos. Query: `?inmueble_id` |
| GET | `/api/v1/alquileres/vencimientos` | alquileres:lectura | Contratos próximos a vencer. Query: `?dias=30` |
| GET | `/api/v1/alquileres/:id` | alquileres:lectura | Detalle inmueble |
| POST | `/api/v1/alquileres` | alquileres:escritura | Crear inmueble |
| PATCH | `/api/v1/alquileres/:id` | alquileres:escritura | Actualizar inmueble |
| DELETE | `/api/v1/alquileres/:id` | admin_only | Eliminar |
| POST | `/api/v1/alquileres/contratos` | alquileres:escritura | Crear contrato de alquiler |
| POST | `/api/v1/alquileres/pagos` | alquileres:escritura | Registrar pago |

**Body POST contrato alquiler:**
```json
{
  "inmuebleId": "uuid (required)",
  "vigenciaDesde": "YYYY-MM-DD (required)",
  "vigenciaHasta": "YYYY-MM-DD (required)",
  "ajusteFrecuencia": "string (required)",
  "montoMensual": 95000
}
```

**Body POST pago alquiler:**
```json
{
  "inmuebleId": "uuid (required)",
  "periodo": "YYYY-MM (required)",
  "monto": 95000,
  "fechaPago": "YYYY-MM-DD (optional)"
}
```

---

## Tesorería

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/tesoreria/flujo-caja` | tesoreria:lectura | Flujo de caja agrupado por día |
| GET | `/api/v1/tesoreria/proyecciones` | tesoreria:lectura | Proyección de liquidez |
| GET | `/api/v1/tesoreria/cuentas` | tesoreria:lectura | Lista cuentas bancarias |
| POST | `/api/v1/tesoreria/cuentas` | tesoreria:escritura | Crear cuenta bancaria |
| PATCH | `/api/v1/tesoreria/cuentas/:id` | tesoreria:escritura | Actualizar cuenta |
| GET | `/api/v1/tesoreria/conciliacion` | tesoreria:lectura | Movimientos pendientes de conciliar |
| POST | `/api/v1/tesoreria/movimientos` | tesoreria:escritura | Registrar movimiento bancario |
| PATCH | `/api/v1/tesoreria/movimientos/:id/conciliar` | tesoreria:escritura | Marcar como conciliado |

**Body POST movimiento:**
```json
{
  "cuentaId": "uuid (required)",
  "tipo": "ingreso|egreso (required)",
  "monto": 30000,
  "descripcion": "string (required)",
  "fecha": "YYYY-MM-DD (required)"
}
```

---

## Reportes

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/reportes` | reportes:lectura | Reporte financiero (egresos/ingresos por período) |
| GET | `/api/v1/reportes/nomina` | reportes:lectura | Reporte liquidaciones del período |
| GET | `/api/v1/reportes/proveedores` | reportes:lectura | Ranking pagos a proveedores |
| GET | `/api/v1/reportes/inmuebles` | reportes:lectura | Reporte alquileres con totales |
| GET | `/api/v1/reportes/exportar` | reportes:lectura | Exportar en CSV. Query: `?modulo=financiero|nomina|...` |

---

## Usuarios (gestión de sub-usuarios)

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| GET | `/api/v1/usuarios` | admin_only | Lista sub-usuarios |
| POST | `/api/v1/usuarios` | admin_only | Crear sub-usuario con permisos |
| PATCH | `/api/v1/usuarios/:id` | admin_only | Actualizar permisos / activar / desactivar |
| DELETE | `/api/v1/usuarios/:id` | admin_only | Desactivar usuario |

**Body POST usuario:**
```json
{
  "nombre": "string (required)",
  "email": "string (required)",
  "password": "string (required)",
  "rol": "sub-usuario",
  "permisos": {
    "dashboard": true,
    "transferencias": "lectura|escritura",
    "nominas": "lectura|escritura",
    "proveedores": "lectura|escritura",
    "servicios": "lectura|escritura",
    "alquileres": "lectura|escritura",
    "tesoreria": "lectura|escritura",
    "reportes": "lectura|escritura"
  }
}
```
