# ADR-002: Diseño de la API REST

## Estado
Aceptado

## Contexto
Definir las convenciones de la API para que el frontend pueda integrarse de forma predecible y consistente.

## Decisión

### Base URL
```
/api/v1/...
```
Versión en la URL desde el inicio para no romper al frontend si hay cambios mayores.

### Convenciones de rutas

| Patrón | Descripción |
|--------|-------------|
| `GET    /api/v1/{recurso}` | Lista paginada |
| `GET    /api/v1/{recurso}/:id` | Detalle de uno |
| `POST   /api/v1/{recurso}` | Crear nuevo |
| `PUT    /api/v1/{recurso}/:id` | Reemplazar completo |
| `PATCH  /api/v1/{recurso}/:id` | Actualización parcial |
| `DELETE /api/v1/{recurso}/:id` | Eliminar |

### Paginación

Todos los endpoints de lista devuelven:

```json
{
  "data": [...],
  "meta": {
    "total": 120,
    "page": 1,
    "per_page": 20,
    "total_pages": 6
  }
}
```

Query params: `?page=1&per_page=20`

### Filtros y búsqueda

Query params estándar:
- `?q=texto` → búsqueda full-text en campos relevantes
- `?estado=pendiente` → filtro por campo específico
- `?desde=2024-01-01&hasta=2024-12-31` → rango de fechas
- `?orden=monto&dir=desc` → ordenamiento

### Formato de respuesta exitosa

```json
{
  "data": { ... }
}
```
o para listas:
```json
{
  "data": [...],
  "meta": { ... }
}
```

### Formato de error estándar

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "El campo 'monto' es requerido",
    "details": [
      { "field": "monto", "message": "required" }
    ]
  }
}
```

### Códigos de error de negocio

| Code | HTTP | Descripción |
|------|------|-------------|
| `UNAUTHORIZED` | 401 | Sin token o token inválido |
| `FORBIDDEN` | 403 | Token válido pero sin permiso para este recurso |
| `NOT_FOUND` | 404 | Recurso no existe |
| `VALIDATION_ERROR` | 422 | Datos de entrada inválidos |
| `CONFLICT` | 409 | Conflicto (ej: email ya registrado) |
| `INTERNAL_ERROR` | 500 | Error interno del servidor |

### Headers requeridos en requests autenticados

```
Cookie: token=<jwt>
Content-Type: application/json
```

### CORS

Origen permitido: solo el dominio del frontend (`FRONTEND_URL` en config).
Credentials: `true` (necesario para httpOnly cookie).

## Consecuencias
- Positivo: paginación y filtros consistentes = menos código en el frontend.
- Positivo: errores con `code` string = el frontend puede manejar casos específicos.
- A tener en cuenta: todos los endpoints deben respetar el mismo envelope `{ "data": ... }`.
