package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ipnext/admin-backend/internal/domain/usuario"
)

// RequirePermiso valida que el usuario tenga el nivel mínimo sobre el módulo.
// Usar "admin_only" como nivelMinimo para restringir solo a admins.
func RequirePermiso(modulo string, nivelMinimo string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "UNAUTHORIZED", "message": "No autenticado"},
			})
			return
		}

		if nivelMinimo == "admin_only" {
			if claims.Rol != usuario.RolAdmin {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": gin.H{"code": "FORBIDDEN", "message": "Acción reservada para administradores"},
				})
			} else {
				c.Next()
			}
			return
		}

		u := &usuario.Usuario{Rol: claims.Rol, Permisos: claims.Permisos}
		if !u.TienePermiso(modulo, usuario.NivelPermiso(nivelMinimo)) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Sin permisos suficientes para este módulo",
				},
			})
			return
		}

		c.Next()
	}
}
