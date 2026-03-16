# TDR-004: Endpoints Faltantes — CRUD Edit/Delete

## Estado
Pendiente de implementación

## Contexto
El frontend necesita editar y eliminar registros en todos los módulos.
Este TDR documenta los endpoints que hay que agregar al backend Go+Gin.

---

## Endpoints a implementar

### Guardias

Agregar a la tabla de guardias en `tdr-003-endpoints.md`:

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| PATCH | `/api/v1/nominas/guardias/:id` | nominas:escritura | Actualizar guardia |
| DELETE | `/api/v1/nominas/guardias/:id` | admin_only | Eliminar guardia |

**Body PATCH guardia:**
```json
{
  "empleadoId": "uuid (optional)",
  "fecha": "YYYY-MM-DD (optional)",
  "horas": 8,
  "monto": 5000,
  "notas": "string (optional)"
}
```

---

### Compensaciones

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| PATCH | `/api/v1/nominas/compensaciones/:id` | nominas:escritura | Actualizar compensación |
| DELETE | `/api/v1/nominas/compensaciones/:id` | admin_only | Eliminar compensación |

**Body PATCH compensación:**
```json
{
  "tipo": "bono|adelanto|extra|otro (optional)",
  "monto": 10000,
  "descripcion": "string (optional)",
  "fecha": "YYYY-MM-DD (optional)",
  "estado": "aprobado|pendiente|rechazado (optional)"
}
```

---

### Alquileres — Contratos

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| PATCH | `/api/v1/alquileres/contratos/:id` | alquileres:escritura | Actualizar contrato |
| DELETE | `/api/v1/alquileres/contratos/:id` | admin_only | Eliminar contrato |

**Body PATCH contrato:**
```json
{
  "vigenciaDesde": "YYYY-MM-DD (optional)",
  "vigenciaHasta": "YYYY-MM-DD (optional)",
  "ajusteFrecuencia": "string (optional)",
  "montoMensual": 95000,
  "estado": "vigente|por_vencer|vencido (optional)"
}
```

---

### Alquileres — Pagos

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| PATCH | `/api/v1/alquileres/pagos/:id` | alquileres:escritura | Actualizar pago (estado, fecha, comprobante) |

**Body PATCH pago:**
```json
{
  "estado": "pagado|pendiente (optional)",
  "fechaPago": "YYYY-MM-DD (optional)",
  "monto": 95000,
  "comprobante": "string (optional)"
}
```

---

### Transferencias — DELETE

| Método | Ruta | Permiso | Descripción |
|--------|------|---------|-------------|
| DELETE | `/api/v1/transferencias/:id` | admin_only | Eliminar transferencia (ya está en docs, verificar si implementado) |

> Nota: `DELETE /transferencias/:id` ya figura en tdr-003 pero devuelve 404 en producción. Verificar si está implementado.

---

### Servicios — rutas correctas

Los endpoints de servicios YA están documentados en tdr-003 con rutas `/servicios/item/:id`.
Verificar en producción si están implementados:

```bash
PATCH /api/v1/servicios/item/:id   # devuelve 404 actualmente
DELETE /api/v1/servicios/item/:id  # devuelve 404 actualmente
```

Si no están implementados, agregarlos siguiendo el patrón del resto de handlers.

**Body PATCH servicio:**
```json
{
  "nombre": "string (optional)",
  "proveedor": "string (optional)",
  "costoMensual": 12000,
  "estado": "activo|proximo_vencer|vencido (optional)",
  "categoria": "internet|energia|seguridad|software (optional)",
  "extra": "string (optional)",
  "vtoFactura": "YYYY-MM-DD (optional)",
  "vigencia": "string (optional)",
  "renovacion": "YYYY-MM-DD (optional)"
}
```

---

## Response estándar para todos los endpoints

**PATCH exitoso:**
```json
{ "data": { <entidad completa actualizada> } }
```

**DELETE exitoso:**
```json
{ "data": { "message": "eliminado correctamente" } }
```

**Error 404:**
```json
{ "error": { "code": "NOT_FOUND", "message": "recurso no encontrado" } }
```

---

## Notas de implementación

- Seguir el patrón de handlers existentes en `internal/handlers/`
- Registrar rutas nuevas en `internal/routes/`
- Los DELETE son **soft delete** (marcar `activo=false` o similar) salvo indicación contraria
- Verificar permisos con el middleware de auth existente
