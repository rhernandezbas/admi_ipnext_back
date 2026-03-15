# ADR-003: Estrategia de Autenticación y Autorización

## Estado
Aceptado

## Contexto
El sistema tiene dos roles (Admin / Sub-usuario) con permisos granulares por módulo y nivel (lectura / escritura / ninguno). La autenticación debe ser segura contra XSS y fácil de implementar desde el frontend React.

## Decisión

### Autenticación — JWT en httpOnly Cookie

- Al hacer login el backend firma un JWT y lo establece como `httpOnly` + `Secure` + `SameSite=Strict` cookie.
- El cliente nunca accede al token por JS — protección contra XSS.
- El JWT expira en **8 horas**; el frontend redirige a `/login` al recibir 401.
- No hay refresh token en v1 — al vencer se re-loguea.

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

### Middleware de autenticación (Go)

```go
// AuthMiddleware — valida JWT, inyecta claims en context
func AuthMiddleware() gin.HandlerFunc

// RequirePermiso — valida módulo + nivel mínimo
func RequirePermiso(modulo string, nivelMinimo string) gin.HandlerFunc
```

Uso en rutas:
```go
transferencias := r.Group("/api/v1/transferencias", AuthMiddleware())
transferencias.GET("", RequirePermiso("transferencias", "lectura"), handler.List)
transferencias.POST("", RequirePermiso("transferencias", "escritura"), handler.Create)
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
- `POST /api/v1/usuarios` (crear sub-usuarios)

### Endpoints de autenticación

```
POST /api/v1/auth/login
  Body: { "email": string, "password": string }
  Response: { "data": { "usuario": Usuario } }
  Cookie: Set-Cookie: token=<jwt>; HttpOnly; Secure; SameSite=Strict

POST /api/v1/auth/logout
  Response: 200 OK
  Cookie: token="" con Max-Age=0 (invalida la cookie)

GET /api/v1/auth/me
  Response: { "data": { "usuario": Usuario } }
```

### Almacenamiento de contraseñas

- Hash con **bcrypt** (cost factor 12).
- Nunca almacenar ni loguear contraseña en texto plano.

## Consecuencias
- Positivo: httpOnly cookie = sin riesgo XSS sobre el token.
- Positivo: permisos en el JWT = el middleware no necesita ir a la DB en cada request.
- A tener en cuenta: si se cambian permisos de un sub-usuario, el JWT viejo sigue siendo válido hasta expirar (8h). Aceptable para v1.
- A tener en cuenta: en producción el flag `Secure` requiere HTTPS.
