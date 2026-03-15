# ADR-002: Diseño de la API REST

## Estado
Aceptado — Implementado

## Contexto
Definir las convenciones de la API para que el frontend pueda integrarse de forma predecible y consistente.

## Decisión

### Base URL
```
/api/v1/...
```

### Convenciones de rutas

| Patrón | Descripción |
|--------|-------------|
| `GET    /api/v1/{recurso}` | Lista |
| `GET    /api/v1/{recurso}/:id` | Detalle |
| `POST   /api/v1/{recurso}` | Crear |
| `PATCH  /api/v1/{recurso}/:id` | Actualización parcial |
| `DELETE /api/v1/{recurso}/:id` | Eliminar (soft delete, solo admin) |

> No se usa `PUT`. Todas las actualizaciones son `PATCH`.

### Formato de respuesta exitosa

Todos los endpoints devuelven el mismo envelope:

```json
{ "data": { ... } }
```

Para listas devuelven array dentro de `data`:

```json
{ "data": [ ... ] }
```

> **Importante para el frontend:** las listas NO incluyen objeto `meta` de paginación en v1. El array llega directamente en `data`. El frontend debe leer `response.data.data` (con axios) o `data.data` según el cliente HTTP.

### Formato de error estándar

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "fechaPago debe ser YYYY-MM-DD"
  }
}
```

### Códigos de error

| Code | HTTP | Descripción |
|------|------|-------------|
| `UNAUTHORIZED` | 401 | Sin token o token inválido |
| `FORBIDDEN` | 403 | Token válido pero sin permiso |
| `NOT_FOUND` | 404 | Recurso no existe |
| `VALIDATION_ERROR` | 422 | Datos de entrada inválidos |
| `INTERNAL_ERROR` | 500 | Error interno del servidor |

### Fechas

- **Request body:** `YYYY-MM-DD` (ej: `"2026-03-15"`)
- **Respuestas:** RFC3339 (ej: `"2026-03-15T00:00:00Z"`)

### Campos JSON en respuestas

Los campos en las respuestas usan **PascalCase** (nombre del struct Go sin tag `json`):

```json
{
  "data": {
    "ID": "uuid",
    "Nombre": "Acme Corp",
    "Activo": true,
    "CreatedAt": "2026-03-15T00:00:00Z"
  }
}
```

> El frontend debe leer `item.ID`, `item.Nombre`, etc. (PascalCase), no `item.id`, `item.nombre`.

### Headers requeridos en requests autenticados

```
Cookie: token=<jwt>
Content-Type: application/json
```

### CORS

Origen permitido configurado por variable de entorno `FRONTEND_URL`.
`credentials: true` requerido para que el browser envíe la cookie.

### Cookie de sesión

- Nombre: `token`
- HttpOnly: `true` (no accesible desde JS)
- Secure: `false` (permite HTTP — cambiar a `true` cuando se configure HTTPS)
- Max-Age: 28800 (8 horas)

## Consecuencias
- Las listas no tienen paginación en v1 — se devuelve todo. A considerar para módulos con muchos registros.
- Los campos PascalCase en las respuestas requieren que el frontend los lea así.
- El flag `Secure=false` en la cookie es temporal hasta configurar HTTPS en el VPS.
