# ADR-003: Estrategia de Autenticación y Autorización

## Estado
Aceptado — Implementado

## Contexto
El sistema tiene dos roles (Admin / Sub-usuario) con permisos granulares por módulo y nivel (lectura / escritura / ninguno). La autenticación debe ser segura contra XSS y funcionar desde el frontend React.

## Decisión

### Autenticación — JWT en httpOnly Cookie

- Al hacer login el backend firma un JWT y lo establece como `httpOnly` + `Secure` + `SameSite=Strict` cookie.
- El cliente nunca accede al token por JS — protección contra XSS.
- El JWT expira en **8 horas**; el frontend redirige a `/login` al recibir 401.
- No hay refresh token en v1 — al vencer se re-loguea.

> **Importante en producción:** el flag `Secure` requiere que el frontend acceda por **HTTPS**. Si el frontend usa HTTP, el browser descarta la cookie y el usuario queda deslogueado inmediatamente tras el login. Solución: configurar HTTPS en el VPS (Traefik/Nginx + Let's Encrypt).

### Estructura del JWT payload

```json
{
  "sub": "user-uuid",
  "rol": "admin",
  "permisos": {
    "dashboard": true,
    "transferencias": "escritura",
    "nominas": "lectura",
    "proveedores": "ninguno",
    "servicios": "escritura",
    "alquileres": "lectura",
    "tesoreria": "ninguno",
    "reportes": "lectura"
  },
  "iat": 1700000000,
  "exp": 1700028800
}
```

### Middleware implementado

```go
// AuthMiddleware — lee cookie "token", valida JWT, inyecta claims en context Gin
func AuthMiddleware(jwtSecret string) gin.HandlerFunc

// RequirePermiso — verifica que el usuario tenga el nivel mínimo sobre el módulo
func RequirePermiso(modulo string, nivelMinimo string) gin.HandlerFunc
```

Uso en rutas:
```go
protected := api.Group("", middleware.AuthMiddleware(jwtSecret))
trans := protected.Group("/transferencias")
trans.GET("", middleware.RequirePermiso("transferencias", "lectura"), h.Transferencia.List)
trans.POST("", middleware.RequirePermiso("transferencias", "escritura"), h.Transferencia.Create)
trans.DELETE("/:id", middleware.RequirePermiso("transferencias", "admin_only"), h.Transferencia.Delete)
```

### Reglas de autorización

| Nivel requerido | Admin | Sub escritura | Sub lectura | Sin permiso |
|-----------------|-------|---------------|-------------|-------------|
| `lectura` | ✅ | ✅ | ✅ | ❌ 403 |
| `escritura` | ✅ | ✅ | ❌ 403 | ❌ 403 |
| `admin_only` | ✅ | ❌ 403 | ❌ 403 | ❌ 403 |

Acciones `admin_only`:
- `DELETE` en cualquier recurso
- `POST /api/v1/nominas/liquidaciones/:id/aprobar`
- `GET|POST|PATCH|DELETE /api/v1/usuarios`

### Endpoints de autenticación

```
POST /api/v1/auth/login
  Body: { "email": string, "password": string }
  Response: { "data": { "usuario": { id, nombre, email, rol, permisos, avatar, creadoEn } } }
  Cookie: Set-Cookie: token=<jwt>; Path=/; Max-Age=28800; HttpOnly; Secure

POST /api/v1/auth/logout
  Response: 200 OK
  Cookie: token="" con Max-Age=0 (invalida la cookie)

GET /api/v1/auth/me
  Response: { "data": { "usuario": {...} } }
```

### Almacenamiento de contraseñas

- Hash con **bcrypt** (cost factor 10).
- Nunca se almacena ni loguea la contraseña en texto plano.

### Usuario administrador por defecto (seed)

| Campo | Valor |
|-------|-------|
| Email | `admin@ipnext.com` |
| Password | `password` |
| Rol | `admin` |

Este usuario se inserta con `INSERT IGNORE` en la migración inicial, por lo que no falla si ya existe.

## Consecuencias
- Positivo: httpOnly cookie = sin riesgo XSS sobre el token.
- Positivo: permisos en el JWT = el middleware no necesita ir a la DB en cada request.
- A tener en cuenta: si se cambian permisos de un sub-usuario, el JWT viejo sigue válido hasta expirar (8h). Aceptable para v1.
- A tener en cuenta: el flag `Secure` requiere HTTPS en producción para que el login funcione correctamente.
