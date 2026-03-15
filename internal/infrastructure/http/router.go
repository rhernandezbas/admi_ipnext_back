package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ipnext/admin-backend/internal/infrastructure/http/middleware"
)

type Handlers struct {
	Auth          AuthHandler
	Dashboard     DashboardHandler
	Transferencia TransferenciaHandler
	Nomina        NominaHandler
	Proveedor     ProveedorHandler
	Servicio      ServicioHandler
	Alquiler      AlquilerHandler
	Tesoreria     TesoreriaHandler
	Reporte       ReporteHandler
	Usuario       UsuarioHandler
}

type AuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Me(c *gin.Context)
}

type DashboardHandler interface {
	GetKPIs(c *gin.Context)
	GetPagosUrgentes(c *gin.Context)
	GetDistribucionEgresos(c *gin.Context)
	GetActividadReciente(c *gin.Context)
}

type TransferenciaHandler interface {
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Calendario(c *gin.Context)
	Recurrentes(c *gin.Context)
	CambiarEstado(c *gin.Context)
}

type NominaHandler interface {
	ListEmpleados(c *gin.Context)
	GetEmpleado(c *gin.Context)
	CreateEmpleado(c *gin.Context)
	UpdateEmpleado(c *gin.Context)
	DeleteEmpleado(c *gin.Context)
	GetKPIs(c *gin.Context)
	ListLiquidaciones(c *gin.Context)
	CreateLiquidacion(c *gin.Context)
	AprobarLiquidacion(c *gin.Context)
	ListGuardias(c *gin.Context)
	CreateGuardia(c *gin.Context)
	ListCompensaciones(c *gin.Context)
	CreateCompensacion(c *gin.Context)
}

type ProveedorHandler interface {
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	ListContratos(c *gin.Context)
	CreateContrato(c *gin.Context)
	UpdateContrato(c *gin.Context)
	GetRanking(c *gin.Context)
}

type ServicioHandler interface {
	List(c *gin.Context)
	GetKPIs(c *gin.Context)
	ListByTipo(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type AlquilerHandler interface {
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	ListContratos(c *gin.Context)
	CreateContrato(c *gin.Context)
	ListPagos(c *gin.Context)
	CreatePago(c *gin.Context)
	GetVencimientos(c *gin.Context)
}

type TesoreriaHandler interface {
	GetFlujoCaja(c *gin.Context)
	ListCuentas(c *gin.Context)
	CreateCuenta(c *gin.Context)
	UpdateCuenta(c *gin.Context)
	GetConciliacion(c *gin.Context)
	CreateMovimiento(c *gin.Context)
	ConciliarMovimiento(c *gin.Context)
	GetProyecciones(c *gin.Context)
}

type ReporteHandler interface {
	Financiero(c *gin.Context)
	Nomina(c *gin.Context)
	Proveedores(c *gin.Context)
	Inmuebles(c *gin.Context)
	Exportar(c *gin.Context)
}

type UsuarioHandler interface {
	List(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func SetupRouter(h Handlers, jwtSecret string) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	// Auth (sin middleware)
	auth := api.Group("/auth")
	auth.POST("/login", h.Auth.Login)
	auth.POST("/logout", h.Auth.Logout)
	auth.GET("/me", middleware.AuthMiddleware(jwtSecret), h.Auth.Me)

	// Rutas protegidas
	protected := api.Group("", middleware.AuthMiddleware(jwtSecret))

	// Dashboard
	dash := protected.Group("/dashboard")
	dash.GET("/kpis", middleware.RequirePermiso("dashboard", "lectura"), h.Dashboard.GetKPIs)
	dash.GET("/pagos-urgentes", middleware.RequirePermiso("dashboard", "lectura"), h.Dashboard.GetPagosUrgentes)
	dash.GET("/distribucion-egresos", middleware.RequirePermiso("dashboard", "lectura"), h.Dashboard.GetDistribucionEgresos)
	dash.GET("/actividad-reciente", middleware.RequirePermiso("dashboard", "lectura"), h.Dashboard.GetActividadReciente)

	// Transferencias
	trans := protected.Group("/transferencias")
	trans.GET("", middleware.RequirePermiso("transferencias", "lectura"), h.Transferencia.List)
	trans.GET("/calendario", middleware.RequirePermiso("transferencias", "lectura"), h.Transferencia.Calendario)
	trans.GET("/recurrentes", middleware.RequirePermiso("transferencias", "lectura"), h.Transferencia.Recurrentes)
	trans.GET("/:id", middleware.RequirePermiso("transferencias", "lectura"), h.Transferencia.Get)
	trans.POST("", middleware.RequirePermiso("transferencias", "escritura"), h.Transferencia.Create)
	trans.PATCH("/:id", middleware.RequirePermiso("transferencias", "escritura"), h.Transferencia.Update)
	trans.PATCH("/:id/estado", middleware.RequirePermiso("transferencias", "escritura"), h.Transferencia.CambiarEstado)
	trans.DELETE("/:id", middleware.RequirePermiso("transferencias", "admin_only"), h.Transferencia.Delete)

	// Nóminas
	nom := protected.Group("/nominas")
	nom.GET("/empleados/kpis", middleware.RequirePermiso("nominas", "lectura"), h.Nomina.GetKPIs)
	nom.GET("/empleados", middleware.RequirePermiso("nominas", "lectura"), h.Nomina.ListEmpleados)
	nom.GET("/empleados/:id", middleware.RequirePermiso("nominas", "lectura"), h.Nomina.GetEmpleado)
	nom.POST("/empleados", middleware.RequirePermiso("nominas", "escritura"), h.Nomina.CreateEmpleado)
	nom.PATCH("/empleados/:id", middleware.RequirePermiso("nominas", "escritura"), h.Nomina.UpdateEmpleado)
	nom.DELETE("/empleados/:id", middleware.RequirePermiso("nominas", "admin_only"), h.Nomina.DeleteEmpleado)
	nom.GET("/liquidaciones", middleware.RequirePermiso("nominas", "lectura"), h.Nomina.ListLiquidaciones)
	nom.POST("/liquidaciones", middleware.RequirePermiso("nominas", "escritura"), h.Nomina.CreateLiquidacion)
	nom.POST("/liquidaciones/:id/aprobar", middleware.RequirePermiso("nominas", "admin_only"), h.Nomina.AprobarLiquidacion)
	nom.GET("/guardias", middleware.RequirePermiso("nominas", "lectura"), h.Nomina.ListGuardias)
	nom.POST("/guardias", middleware.RequirePermiso("nominas", "escritura"), h.Nomina.CreateGuardia)
	nom.GET("/compensaciones", middleware.RequirePermiso("nominas", "lectura"), h.Nomina.ListCompensaciones)
	nom.POST("/compensaciones", middleware.RequirePermiso("nominas", "escritura"), h.Nomina.CreateCompensacion)

	// Proveedores
	prov := protected.Group("/proveedores")
	prov.GET("", middleware.RequirePermiso("proveedores", "lectura"), h.Proveedor.List)
	prov.GET("/contratos", middleware.RequirePermiso("proveedores", "lectura"), h.Proveedor.ListContratos)
	prov.GET("/ranking", middleware.RequirePermiso("proveedores", "lectura"), h.Proveedor.GetRanking)
	prov.GET("/:id", middleware.RequirePermiso("proveedores", "lectura"), h.Proveedor.Get)
	prov.POST("", middleware.RequirePermiso("proveedores", "escritura"), h.Proveedor.Create)
	prov.PATCH("/:id", middleware.RequirePermiso("proveedores", "escritura"), h.Proveedor.Update)
	prov.DELETE("/:id", middleware.RequirePermiso("proveedores", "admin_only"), h.Proveedor.Delete)
	prov.POST("/contratos", middleware.RequirePermiso("proveedores", "escritura"), h.Proveedor.CreateContrato)
	prov.PATCH("/contratos/:id", middleware.RequirePermiso("proveedores", "escritura"), h.Proveedor.UpdateContrato)

	// Servicios
	svc := protected.Group("/servicios")
	svc.GET("", middleware.RequirePermiso("servicios", "lectura"), h.Servicio.List)
	svc.GET("/kpis", middleware.RequirePermiso("servicios", "lectura"), h.Servicio.GetKPIs)
	svc.GET("/:tipo", middleware.RequirePermiso("servicios", "lectura"), h.Servicio.ListByTipo)
	svc.GET("/item/:id", middleware.RequirePermiso("servicios", "lectura"), h.Servicio.Get)
	svc.POST("", middleware.RequirePermiso("servicios", "escritura"), h.Servicio.Create)
	svc.PATCH("/item/:id", middleware.RequirePermiso("servicios", "escritura"), h.Servicio.Update)
	svc.DELETE("/item/:id", middleware.RequirePermiso("servicios", "admin_only"), h.Servicio.Delete)

	// Alquileres
	alq := protected.Group("/alquileres")
	alq.GET("", middleware.RequirePermiso("alquileres", "lectura"), h.Alquiler.List)
	alq.GET("/contratos", middleware.RequirePermiso("alquileres", "lectura"), h.Alquiler.ListContratos)
	alq.GET("/pagos", middleware.RequirePermiso("alquileres", "lectura"), h.Alquiler.ListPagos)
	alq.GET("/vencimientos", middleware.RequirePermiso("alquileres", "lectura"), h.Alquiler.GetVencimientos)
	alq.GET("/:id", middleware.RequirePermiso("alquileres", "lectura"), h.Alquiler.Get)
	alq.POST("", middleware.RequirePermiso("alquileres", "escritura"), h.Alquiler.Create)
	alq.PATCH("/:id", middleware.RequirePermiso("alquileres", "escritura"), h.Alquiler.Update)
	alq.DELETE("/:id", middleware.RequirePermiso("alquileres", "admin_only"), h.Alquiler.Delete)
	alq.POST("/contratos", middleware.RequirePermiso("alquileres", "escritura"), h.Alquiler.CreateContrato)
	alq.POST("/pagos", middleware.RequirePermiso("alquileres", "escritura"), h.Alquiler.CreatePago)

	// Tesorería
	tes := protected.Group("/tesoreria")
	tes.GET("/flujo-caja", middleware.RequirePermiso("tesoreria", "lectura"), h.Tesoreria.GetFlujoCaja)
	tes.GET("/cuentas", middleware.RequirePermiso("tesoreria", "lectura"), h.Tesoreria.ListCuentas)
	tes.GET("/conciliacion", middleware.RequirePermiso("tesoreria", "lectura"), h.Tesoreria.GetConciliacion)
	tes.GET("/proyecciones", middleware.RequirePermiso("tesoreria", "lectura"), h.Tesoreria.GetProyecciones)
	tes.POST("/cuentas", middleware.RequirePermiso("tesoreria", "escritura"), h.Tesoreria.CreateCuenta)
	tes.PATCH("/cuentas/:id", middleware.RequirePermiso("tesoreria", "escritura"), h.Tesoreria.UpdateCuenta)
	tes.POST("/movimientos", middleware.RequirePermiso("tesoreria", "escritura"), h.Tesoreria.CreateMovimiento)
	tes.PATCH("/movimientos/:id/conciliar", middleware.RequirePermiso("tesoreria", "escritura"), h.Tesoreria.ConciliarMovimiento)

	// Reportes
	rep := protected.Group("/reportes")
	rep.GET("", middleware.RequirePermiso("reportes", "lectura"), h.Reporte.Financiero)
	rep.GET("/nomina", middleware.RequirePermiso("reportes", "lectura"), h.Reporte.Nomina)
	rep.GET("/proveedores", middleware.RequirePermiso("reportes", "lectura"), h.Reporte.Proveedores)
	rep.GET("/inmuebles", middleware.RequirePermiso("reportes", "lectura"), h.Reporte.Inmuebles)
	rep.GET("/exportar", middleware.RequirePermiso("reportes", "lectura"), h.Reporte.Exportar)

	// Usuarios
	usr := protected.Group("/usuarios")
	usr.GET("", middleware.RequirePermiso("usuarios", "admin_only"), h.Usuario.List)
	usr.POST("", middleware.RequirePermiso("usuarios", "admin_only"), h.Usuario.Create)
	usr.PATCH("/:id", middleware.RequirePermiso("usuarios", "admin_only"), h.Usuario.Update)
	usr.DELETE("/:id", middleware.RequirePermiso("usuarios", "admin_only"), h.Usuario.Delete)

	return r
}
