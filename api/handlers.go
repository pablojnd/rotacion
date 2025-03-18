package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pablojnd/rotacion/db"
	"github.com/pablojnd/rotacion/models"
	"github.com/pablojnd/rotacion/services"
	"github.com/pablojnd/rotacion/utils"
)

// Handlers contiene los handlers para las rutas del API
type Handlers struct {
	sqlServer         *db.SQLServerDB
	mysql             *db.MySQLDB
	ventasService     *services.VentasService
	inventarioService *services.InventarioService
}

// QueryRequest representa una solicitud de consulta
type QueryRequest struct {
	Query string        `json:"query"`
	Args  []interface{} `json:"args,omitempty"`
}

// NewHandlers crea una nueva instancia de Handlers
func NewHandlers(sqlServer *db.SQLServerDB, mysql *db.MySQLDB) *Handlers {
	excelService := services.NewExcelService()
	return &Handlers{
		sqlServer:         sqlServer,
		mysql:             mysql,
		ventasService:     services.NewVentasService(sqlServer, excelService),
		inventarioService: services.NewInventarioService(mysql, excelService),
	}
}

// SQLServerQuery maneja las consultas a SQL Server
func (h *Handlers) SQLServerQuery(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := h.sqlServer.ExecuteQuery(req.Query, req.Args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result, err := utils.RowsToJSON(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// MySQLQuery maneja las consultas a MySQL
func (h *Handlers) MySQLQuery(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := h.mysql.ExecuteQuery(req.Query, req.Args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result, err := utils.RowsToJSON(rows)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetVentas obtiene las ventas según los filtros proporcionados
func (h *Handlers) GetVentas(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la consulta
	filtro := models.VentasFiltro{
		FechaInicio:    r.URL.Query().Get("fechaInicio"),
		FechaFin:       r.URL.Query().Get("fechaFin"),
		Sucursal:       parseIntParam(r.URL.Query().Get("sucursal"), 211),
		CodigoProducto: r.URL.Query().Get("codigo"),
	}

	// Usar el servicio para obtener los datos
	result, err := h.ventasService.GetVentas(filtro)
	if err != nil {
		log.Printf("Error al consultar ventas: %v", err)
		http.Error(w, fmt.Sprintf("Error al ejecutar consulta: %v", err), http.StatusInternalServerError)
		return
	}

	// Devolver resultados
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetInventario obtiene el inventario según los filtros proporcionados
func (h *Handlers) GetInventario(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la consulta
	anio := parseIntParam(r.URL.Query().Get("anio"), time.Now().Year())
	codigoProducto := r.URL.Query().Get("codigo")

	filtro := models.InventarioFiltro{
		Anio:           anio,
		CodigoProducto: codigoProducto,
	}

	// Usar el servicio para obtener los datos
	result, err := h.inventarioService.GetInventario(filtro)
	if err != nil {
		log.Printf("Error al consultar inventario: %v", err)
		http.Error(w, fmt.Sprintf("Error al ejecutar consulta: %v", err), http.StatusInternalServerError)
		return
	}

	// Devolver resultados
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// parseIntParam convierte un string a int con valor predeterminado
func parseIntParam(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
