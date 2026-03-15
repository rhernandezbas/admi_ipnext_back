package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"

	"github.com/ipnext/admin-backend/config"
	appalq "github.com/ipnext/admin-backend/internal/application/alquiler"
	appauth "github.com/ipnext/admin-backend/internal/application/auth"
	appdash "github.com/ipnext/admin-backend/internal/application/dashboard"
	appnom "github.com/ipnext/admin-backend/internal/application/nomina"
	appprov "github.com/ipnext/admin-backend/internal/application/proveedor"
	apprep "github.com/ipnext/admin-backend/internal/application/reporte"
	appsvc "github.com/ipnext/admin-backend/internal/application/servicio"
	apptes "github.com/ipnext/admin-backend/internal/application/tesoreria"
	apptrans "github.com/ipnext/admin-backend/internal/application/transferencia"
	appusr "github.com/ipnext/admin-backend/internal/application/usuario"
	infraconfig "github.com/ipnext/admin-backend/internal/infrastructure/config"
	infrahttp "github.com/ipnext/admin-backend/internal/infrastructure/http"
	"github.com/ipnext/admin-backend/internal/infrastructure/http/handler"
	"github.com/ipnext/admin-backend/internal/infrastructure/persistence/repository"
)

func runMigrations(dsn, migrationsDir string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}
	return goose.Up(db, migrationsDir)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error cargando configuración: %v", err)
	}

	_, thisFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(thisFile), "..", "..", "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = "migrations"
	}
	if err := runMigrations(cfg.Database.DSN(), migrationsDir); err != nil {
		log.Fatalf("error ejecutando migraciones: %v", err)
	}

	db, err := infraconfig.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("error conectando a la base de datos: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Repositorios
	usuarioRepo := repository.NewMySQLUsuarioRepository(db)
	transRepo := repository.NewMySQLTransferenciaRepository(db)
	empleadoRepo := repository.NewMySQLEmpleadoRepository(db)
	liquidacionRepo := repository.NewMySQLLiquidacionRepository(db)
	guardiaRepo := repository.NewMySQLGuardiaRepository(db)
	compensacionRepo := repository.NewMySQLCompensacionRepository(db)
	proveedorRepo := repository.NewMySQLProveedorRepository(db)
	contratoProvRepo := repository.NewMySQLContratoProveedorRepository(db)
	rankingRepo := repository.NewMySQLRankingRepository(db)
	servicioRepo := repository.NewMySQLServicioRepository(db)
	inmuebleRepo := repository.NewMySQLAlquilerRepository(db)
	contratoAlqRepo := repository.NewMySQLContratoAlquilerRepository(db)
	pagoAlqRepo := repository.NewMySQLPagoAlquilerRepository(db)
	cuentaRepo := repository.NewMySQLCuentaRepository(db)
	movimientoRepo := repository.NewMySQLMovimientoRepository(db)

	// Casos de uso — Auth
	loginUC := appauth.NewLoginUseCase(usuarioRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	getMeUC := appauth.NewGetMeUseCase(usuarioRepo)

	// Casos de uso — Dashboard
	getKPIsDashUC := appdash.NewGetKPIsUseCase(transRepo, empleadoRepo)
	getPagosUrgentesUC := appdash.NewGetPagosUrgentesUseCase(transRepo)
	getDistribUC := appdash.NewGetDistribucionEgresosUseCase(transRepo)

	// Casos de uso — Transferencias
	listTransUC := apptrans.NewListUseCase(transRepo)
	getTransUC := apptrans.NewGetUseCase(transRepo)
	createTransUC := apptrans.NewCreateUseCase(transRepo)
	updateTransUC := apptrans.NewUpdateUseCase(transRepo)
	deleteTransUC := apptrans.NewDeleteUseCase(transRepo)
	cambiarEstadoUC := apptrans.NewCambiarEstadoUseCase(transRepo)
	calendarioUC := apptrans.NewCalendarioUseCase(transRepo)
	recurrentesUC := apptrans.NewRecurrentesUseCase(transRepo)

	// Casos de uso — Nóminas
	listEmpleadosUC := appnom.NewListEmpleadosUseCase(empleadoRepo)
	getEmpleadoUC := appnom.NewGetEmpleadoUseCase(empleadoRepo)
	createEmpleadoUC := appnom.NewCreateEmpleadoUseCase(empleadoRepo)
	updateEmpleadoUC := appnom.NewUpdateEmpleadoUseCase(empleadoRepo)
	deleteEmpleadoUC := appnom.NewDeleteEmpleadoUseCase(empleadoRepo)
	getKPIsNomUC := appnom.NewGetEmpleadoKPIsUseCase(empleadoRepo)
	listLiqUC := appnom.NewListLiquidacionesUseCase(liquidacionRepo)
	createLiqUC := appnom.NewCreateLiquidacionUseCase(liquidacionRepo, empleadoRepo)
	aprobarLiqUC := appnom.NewAprobarLiquidacionUseCase(liquidacionRepo)
	listGuardiasUC := appnom.NewListGuardiasUseCase(guardiaRepo)
	createGuardiaUC := appnom.NewCreateGuardiaUseCase(guardiaRepo)
	listCompUC := appnom.NewListCompensacionesUseCase(compensacionRepo)
	createCompUC := appnom.NewCreateCompensacionUseCase(compensacionRepo)

	// Casos de uso — Proveedores
	listProvUC := appprov.NewListUseCase(proveedorRepo)
	getProvUC := appprov.NewGetUseCase(proveedorRepo)
	createProvUC := appprov.NewCreateUseCase(proveedorRepo)
	updateProvUC := appprov.NewUpdateUseCase(proveedorRepo)
	deleteProvUC := appprov.NewDeleteUseCase(proveedorRepo)
	listContratosProvUC := appprov.NewListContratosUseCase(contratoProvRepo)
	createContratoProvUC := appprov.NewCreateContratoUseCase(contratoProvRepo, proveedorRepo)
	updateContratoProvUC := appprov.NewUpdateContratoUseCase(contratoProvRepo)
	getRankingUC := appprov.NewGetRankingUseCase(rankingRepo)

	// Casos de uso — Servicios
	listSvcUC := appsvc.NewListUseCase(servicioRepo)
	getKPIsSvcUC := appsvc.NewGetKPIsUseCase(servicioRepo)
	getSvcUC := appsvc.NewGetUseCase(servicioRepo)
	createSvcUC := appsvc.NewCreateUseCase(servicioRepo)
	updateSvcUC := appsvc.NewUpdateUseCase(servicioRepo)
	deleteSvcUC := appsvc.NewDeleteUseCase(servicioRepo)

	// Casos de uso — Alquileres
	listInmueblesUC := appalq.NewListInmueblesUseCase(inmuebleRepo)
	getInmuebleUC := appalq.NewGetInmuebleUseCase(inmuebleRepo)
	createInmuebleUC := appalq.NewCreateInmuebleUseCase(inmuebleRepo)
	updateInmuebleUC := appalq.NewUpdateInmuebleUseCase(inmuebleRepo)
	deleteInmuebleUC := appalq.NewDeleteInmuebleUseCase(inmuebleRepo)
	listContratosAlqUC := appalq.NewListContratosUseCase(contratoAlqRepo)
	createContratoAlqUC := appalq.NewCreateContratoUseCase(contratoAlqRepo)
	vencimientosUC := appalq.NewVencimientosUseCase(contratoAlqRepo)
	listPagosUC := appalq.NewListPagosUseCase(pagoAlqRepo)
	createPagoUC := appalq.NewCreatePagoUseCase(pagoAlqRepo)

	// Casos de uso — Tesorería
	flujoCajaUC := apptes.NewGetFlujoCajaUseCase(movimientoRepo)
	proyeccionesUC := apptes.NewGetProyeccionesUseCase(movimientoRepo, cuentaRepo)
	listCuentasUC := apptes.NewListCuentasUseCase(cuentaRepo)
	createCuentaUC := apptes.NewCreateCuentaUseCase(cuentaRepo)
	updateCuentaUC := apptes.NewUpdateCuentaUseCase(cuentaRepo)
	getConciliacionUC := apptes.NewGetConciliacionUseCase(movimientoRepo)
	createMovimientoUC := apptes.NewCreateMovimientoUseCase(movimientoRepo)
	conciliarMovimientoUC := apptes.NewConciliarMovimientoUseCase(movimientoRepo)

	// Casos de uso — Reportes
	repFinancieroUC := apprep.NewReporteFinancieroUseCase(transRepo)
	repNominaUC := apprep.NewReporteNominaUseCase(liquidacionRepo, empleadoRepo)
	repProveedoresUC := apprep.NewReporteProveedoresUseCase(rankingRepo)
	repInmueblesUC := apprep.NewReporteInmueblesUseCase(inmuebleRepo, pagoAlqRepo)

	// Casos de uso — Usuarios
	listUsrUC := appusr.NewListUseCase(usuarioRepo)
	createUsrUC := appusr.NewCreateUseCase(usuarioRepo)
	updateUsrUC := appusr.NewUpdateUseCase(usuarioRepo)
	deleteUsrUC := appusr.NewDeleteUseCase(usuarioRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(loginUC, getMeUC)
	dashHandler := handler.NewDashboardHandler(getKPIsDashUC, getPagosUrgentesUC, getDistribUC)
	transHandler := handler.NewTransferenciaHandler(listTransUC, getTransUC, createTransUC, updateTransUC, deleteTransUC, cambiarEstadoUC, calendarioUC, recurrentesUC)
	nominaHandler := handler.NewNominaHandler(listEmpleadosUC, getEmpleadoUC, createEmpleadoUC, updateEmpleadoUC, deleteEmpleadoUC, getKPIsNomUC, listLiqUC, createLiqUC, aprobarLiqUC, listGuardiasUC, createGuardiaUC, listCompUC, createCompUC)
	proveedorHandler := handler.NewProveedorHandler(listProvUC, getProvUC, createProvUC, updateProvUC, deleteProvUC, listContratosProvUC, createContratoProvUC, updateContratoProvUC, getRankingUC)
	servicioHandler := handler.NewServicioHandler(listSvcUC, getKPIsSvcUC, getSvcUC, createSvcUC, updateSvcUC, deleteSvcUC)
	alquilerHandler := handler.NewAlquilerHandler(listInmueblesUC, getInmuebleUC, createInmuebleUC, updateInmuebleUC, deleteInmuebleUC, listContratosAlqUC, createContratoAlqUC, vencimientosUC, listPagosUC, createPagoUC)
	tesoreriaHandler := handler.NewTesoreriaHandler(flujoCajaUC, proyeccionesUC, listCuentasUC, createCuentaUC, updateCuentaUC, getConciliacionUC, createMovimientoUC, conciliarMovimientoUC)
	reporteHandler := handler.NewReporteHandler(repFinancieroUC, repNominaUC, repProveedoresUC, repInmueblesUC)
	usuarioHandler := handler.NewUsuarioHandler(listUsrUC, createUsrUC, updateUsrUC, deleteUsrUC)

	// Router
	handlers := infrahttp.Handlers{
		Auth:          authHandler,
		Dashboard:     dashHandler,
		Transferencia: transHandler,
		Nomina:        nominaHandler,
		Proveedor:     proveedorHandler,
		Servicio:      servicioHandler,
		Alquiler:      alquilerHandler,
		Tesoreria:     tesoreriaHandler,
		Reporte:       reporteHandler,
		Usuario:       usuarioHandler,
	}

	r := infrahttp.SetupRouter(handlers, cfg.JWT.Secret)

	fmt.Printf("🚀 IPNEXT Backend corriendo en :%s [%s]\n", cfg.Server.Port, cfg.Server.Env)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("error iniciando servidor: %v", err)
	}
}
