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
	excelService      *services.ExcelService
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
		excelService:      excelService,
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

// ExportToExcel exporta datos a Excel
func (h *Handlers) ExportToExcel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query    string        `json:"query"`
		Args     []interface{} `json:"args,omitempty"`
		Database string        `json:"database"` // "sqlserver" o "mysql"
		Filename string        `json:"filename"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var db interface{}

	// Seleccionar la base de datos correcta
	if req.Database == "sqlserver" {
		db = h.sqlServer.DB
	} else {
		db = h.mysql.DB
	}

	// Generar el Excel usando el servicio
	excelBytes, err := h.excelService.GenerateExcelFromQuery(db, req.Query, req.Args, req.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar respuesta para descarga
	sendExcelResponse(w, excelBytes, req.Filename)
}

// GetVentas obtiene las ventas según los filtros proporcionados
func (h *Handlers) GetVentas(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la consulta
	filtro := models.VentasFiltro{
		FechaInicio: r.URL.Query().Get("fechaInicio"),
		FechaFin:    r.URL.Query().Get("fechaFin"),
		Sucursal:    parseIntParam(r.URL.Query().Get("sucursal"), 211),
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

// ExportVentas exporta las ventas a Excel
func (h *Handlers) ExportVentas(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la consulta
	filtro := models.VentasFiltro{
		FechaInicio: r.URL.Query().Get("fechaInicio"),
		FechaFin:    r.URL.Query().Get("fechaFin"),
		Sucursal:    parseIntParam(r.URL.Query().Get("sucursal"), 211),
	}

	// Usar el servicio para exportar a Excel
	excelBytes, filename, err := h.ventasService.ExportVentasToExcel(filtro)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al generar Excel: %v", err), http.StatusInternalServerError)
		return
	}

	// Enviar respuesta
	sendExcelResponse(w, excelBytes, filename)
}

// GetInventario obtiene el inventario según los filtros proporcionados
func (h *Handlers) GetInventario(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la consulta
	anio := parseIntParam(r.URL.Query().Get("anio"), time.Now().Year())

	filtro := models.InventarioFiltro{
		Anio: anio,
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

// ExportInventario exporta el inventario a Excel
func (h *Handlers) ExportInventario(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la consulta
	anio := parseIntParam(r.URL.Query().Get("anio"), time.Now().Year())

	filtro := models.InventarioFiltro{
		Anio: anio,
	}

	// Usar el servicio para exportar a Excel
	excelBytes, filename, err := h.inventarioService.ExportInventarioToExcel(filtro)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al generar Excel: %v", err), http.StatusInternalServerError)
		return
	}

	// Enviar respuesta
	sendExcelResponse(w, excelBytes, filename)
}

// sendExcelResponse envía un archivo Excel como respuesta HTTP
func sendExcelResponse(w http.ResponseWriter, excelBytes []byte, filename string) {
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")

	w.Write(excelBytes)
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
