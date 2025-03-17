package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pablojnd/rotacion/api"
	"github.com/pablojnd/rotacion/config"
	"github.com/pablojnd/rotacion/db"
	"github.com/pablojnd/rotacion/excel"
	"github.com/pablojnd/rotacion/services"
)

// Server representa el servidor HTTP
type Server struct {
	config    *config.Config
	router    *mux.Router
	sqlServer *db.SQLServerDB
	mysql     *db.MySQLDB
}

// New crea una nueva instancia del servidor
func New(cfg *config.Config, sqlServer *db.SQLServerDB, mysql *db.MySQLDB) *Server {
	s := &Server{
		config:    cfg,
		router:    mux.NewRouter(),
		sqlServer: sqlServer,
		mysql:     mysql,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configura las rutas del API
func (s *Server) setupRoutes() {
	// Crear servicios compartidos
	excelService := services.NewExcelService()
	ventasService := services.NewVentasService(s.sqlServer, excelService)
	inventarioService := services.NewInventarioService(s.mysql, excelService)

	// Crear handlers para la API
	handlers := api.NewHandlers(s.sqlServer, s.mysql)

	// Crear handler para Excel
	excelHandler := excel.NewHandler(
		s.sqlServer,
		s.mysql,
		excelService,
		ventasService,
		inventarioService,
	)

	// Ruta de estado del servidor
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API en funcionamiento"))
	}).Methods("GET")

	// Rutas de la API
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	// Consultas SQL Server
	apiRouter.HandleFunc("/sqlserver/query", handlers.SQLServerQuery).Methods("POST")

	// Consulta específica de ventas
	apiRouter.HandleFunc("/ventas", handlers.GetVentas).Methods("GET")
	apiRouter.HandleFunc("/ventas/excel", excelHandler.ExportVentas).Methods("GET")

	// Consultas MySQL
	apiRouter.HandleFunc("/mysql/query", handlers.MySQLQuery).Methods("POST")

	// Ruta para inventario
	apiRouter.HandleFunc("/inventario", handlers.GetInventario).Methods("GET")
	apiRouter.HandleFunc("/inventario/excel", excelHandler.ExportInventario).Methods("GET")

	// Exportar a Excel
	apiRouter.HandleFunc("/export/excel", excelHandler.ExportGeneric).Methods("POST")

	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("./static"))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Redirección simple para /docs
	s.router.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/docs/index.html", http.StatusMovedPermanently)
	})
}

// Start inicia el servidor HTTP
func (s *Server) Start() error {
	return http.ListenAndServe(":"+s.config.ServerPort, s.router)
}
