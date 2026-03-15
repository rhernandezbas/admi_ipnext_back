package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	appnomina "github.com/ipnext/admin-backend/internal/application/nomina"
	"github.com/ipnext/admin-backend/internal/domain/empleado"
	"github.com/ipnext/admin-backend/internal/infrastructure/http/middleware"
)

type NominaHandlerImpl struct {
	listEmpleados      *appnomina.ListEmpleadosUseCase
	getEmpleado        *appnomina.GetEmpleadoUseCase
	createEmpleado     *appnomina.CreateEmpleadoUseCase
	updateEmpleado     *appnomina.UpdateEmpleadoUseCase
	deleteEmpleado     *appnomina.DeleteEmpleadoUseCase
	getKPIs            *appnomina.GetEmpleadoKPIsUseCase
	listLiquidaciones  *appnomina.ListLiquidacionesUseCase
	createLiquidacion  *appnomina.CreateLiquidacionUseCase
	aprobarLiquidacion *appnomina.AprobarLiquidacionUseCase
	listGuardias       *appnomina.ListGuardiasUseCase
	createGuardia      *appnomina.CreateGuardiaUseCase
	listCompensaciones *appnomina.ListCompensacionesUseCase
	createCompensacion *appnomina.CreateCompensacionUseCase
}

func NewNominaHandler(
	listEmp *appnomina.ListEmpleadosUseCase,
	getEmp *appnomina.GetEmpleadoUseCase,
	createEmp *appnomina.CreateEmpleadoUseCase,
	updateEmp *appnomina.UpdateEmpleadoUseCase,
	deleteEmp *appnomina.DeleteEmpleadoUseCase,
	kpis *appnomina.GetEmpleadoKPIsUseCase,
	listLiq *appnomina.ListLiquidacionesUseCase,
	createLiq *appnomina.CreateLiquidacionUseCase,
	aprobarLiq *appnomina.AprobarLiquidacionUseCase,
	listGuard *appnomina.ListGuardiasUseCase,
	createGuard *appnomina.CreateGuardiaUseCase,
	listComp *appnomina.ListCompensacionesUseCase,
	createComp *appnomina.CreateCompensacionUseCase,
) *NominaHandlerImpl {
	return &NominaHandlerImpl{
		listEmpleados: listEmp, getEmpleado: getEmp, createEmpleado: createEmp,
		updateEmpleado: updateEmp, deleteEmpleado: deleteEmp, getKPIs: kpis,
		listLiquidaciones: listLiq, createLiquidacion: createLiq, aprobarLiquidacion: aprobarLiq,
		listGuardias: listGuard, createGuardia: createGuard,
		listCompensaciones: listComp, createCompensacion: createComp,
	}
}

// --- Empleados ---

func (h *NominaHandlerImpl) ListEmpleados(c *gin.Context) {
	items, err := h.listEmpleados.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *NominaHandlerImpl) GetEmpleado(c *gin.Context) {
	e, err := h.getEmpleado.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, appnomina.ErrEmpleadoNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Empleado no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": e})
}

type createEmpleadoRequest struct {
	Nombre       string  `json:"nombre" binding:"required"`
	Puesto       string  `json:"puesto" binding:"required"`
	Area         string  `json:"area" binding:"required"`
	Rol          string  `json:"rol" binding:"required"`
	SueldoBruto  float64 `json:"sueldoBruto" binding:"required,gt=0"`
	ObraSocial   string  `json:"obraSocial" binding:"required"`
	FechaIngreso string  `json:"fechaIngreso" binding:"required"`
	Avatar       *string `json:"avatar"`
}

func (h *NominaHandlerImpl) CreateEmpleado(c *gin.Context) {
	var req createEmpleadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	fecha, err := time.Parse("2006-01-02", req.FechaIngreso)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion("fechaIngreso debe ser YYYY-MM-DD"))
		return
	}
	e, err := h.createEmpleado.Execute(c.Request.Context(), appnomina.CreateEmpleadoRequest{
		Nombre: req.Nombre, Puesto: req.Puesto, Area: req.Area, Rol: req.Rol,
		SueldoBruto: req.SueldoBruto, ObraSocial: req.ObraSocial,
		FechaIngreso: fecha, Avatar: req.Avatar,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": e})
}

func (h *NominaHandlerImpl) UpdateEmpleado(c *gin.Context) {
	var req appnomina.UpdateEmpleadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	e, err := h.updateEmpleado.Execute(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, appnomina.ErrEmpleadoNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Empleado no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": e})
}

func (h *NominaHandlerImpl) DeleteEmpleado(c *gin.Context) {
	if err := h.deleteEmpleado.Execute(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, appnomina.ErrEmpleadoNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Empleado no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "Empleado dado de baja"}})
}

func (h *NominaHandlerImpl) GetKPIs(c *gin.Context) {
	kpis, err := h.getKPIs.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpis})
}

// --- Liquidaciones ---

func (h *NominaHandlerImpl) ListLiquidaciones(c *gin.Context) {
	periodo := c.DefaultQuery("periodo", time.Now().Format("2006-01"))
	items, err := h.listLiquidaciones.Execute(c.Request.Context(), periodo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *NominaHandlerImpl) CreateLiquidacion(c *gin.Context) {
	var req appnomina.CreateLiquidacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	l, err := h.createLiquidacion.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, appnomina.ErrEmpleadoNoEncontrado) {
			c.JSON(http.StatusNotFound, errNotFound("Empleado no encontrado"))
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": l})
}

func (h *NominaHandlerImpl) AprobarLiquidacion(c *gin.Context) {
	claims := middleware.GetClaims(c)
	l, err := h.aprobarLiquidacion.Execute(c.Request.Context(), c.Param("id"), claims.Sub)
	if err != nil {
		if errors.Is(err, appnomina.ErrLiquidacionNoEncontrada) {
			c.JSON(http.StatusNotFound, errNotFound("Liquidación no encontrada"))
			return
		}
		if errors.Is(err, appnomina.ErrLiquidacionYaAprobada) {
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "CONFLICT", "message": err.Error()}})
			return
		}
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": l})
}

// --- Guardias ---

func (h *NominaHandlerImpl) ListGuardias(c *gin.Context) {
	var empID *string
	if id := c.Query("empleado_id"); id != "" {
		empID = &id
	}
	items, err := h.listGuardias.Execute(c.Request.Context(), empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *NominaHandlerImpl) CreateGuardia(c *gin.Context) {
	var body struct {
		EmpleadoID  string  `json:"empleadoId" binding:"required"`
		Fecha       string  `json:"fecha" binding:"required"`
		Horas       float64 `json:"horas" binding:"required,gt=0"`
		Monto       float64 `json:"monto" binding:"required,gt=0"`
		Descripcion *string `json:"descripcion"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	fecha, err := time.Parse("2006-01-02", body.Fecha)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion("fecha debe ser YYYY-MM-DD"))
		return
	}
	g, err := h.createGuardia.Execute(c.Request.Context(), appnomina.CreateGuardiaRequest{
		EmpleadoID: body.EmpleadoID, Fecha: fecha,
		Horas: body.Horas, Monto: body.Monto,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": g})
}

// --- Compensaciones ---

func (h *NominaHandlerImpl) ListCompensaciones(c *gin.Context) {
	var empID *string
	if id := c.Query("empleado_id"); id != "" {
		empID = &id
	}
	items, err := h.listCompensaciones.Execute(c.Request.Context(), empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *NominaHandlerImpl) CreateCompensacion(c *gin.Context) {
	var body struct {
		EmpleadoID  string                   `json:"empleadoId" binding:"required"`
		Tipo        empleado.TipoCompensacion `json:"tipo" binding:"required"`
		Monto       float64                  `json:"monto" binding:"required,gt=0"`
		Fecha       string                   `json:"fecha" binding:"required"`
		Descripcion *string                  `json:"descripcion"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion(err.Error()))
		return
	}
	fecha, err := time.Parse("2006-01-02", body.Fecha)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, errValidacion("fecha debe ser YYYY-MM-DD"))
		return
	}
	comp, err := h.createCompensacion.Execute(c.Request.Context(), appnomina.CreateCompensacionRequest{
		EmpleadoID: body.EmpleadoID, Tipo: body.Tipo, Monto: body.Monto,
		Fecha: fecha, Descripcion: body.Descripcion,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errInterno())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": comp})
}
