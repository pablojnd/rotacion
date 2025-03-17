package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yourusername/rotacion/api"
	"github.com/yourusername/rotacion/config"
	"github.com/yourusername/rotacion/db"
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
	// Crear handlers para la API
	handlers := api.NewHandlers(s.sqlServer, s.mysql)

	// Rutas de la API
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	// Consultas SQL Server
	apiRouter.HandleFunc("/sqlserver/query", handlers.SQLServerQuery).Methods("POST")

	// Consultas MySQL
	apiRouter.HandleFunc("/mysql/query", handlers.MySQLQuery).Methods("POST")

	// Exportar a Excel
	apiRouter.HandleFunc("/export/excel", handlers.ExportToExcel).Methods("POST")
}

// Start inicia el servidor HTTP
func (s *Server) Start() error {
	return http.ListenAndServe(":"+s.config.ServerPort, s.router)
}
