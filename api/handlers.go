package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pablojnd/rotacion/db"
	"github.com/pablojnd/rotacion/models"
	"github.com/pablojnd/rotacion/queries/sqlserver"
	"github.com/pablojnd/rotacion/utils"
)

// Handlers contiene los handlers para las rutas del API
type Handlers struct {
	sqlServer *db.SQLServerDB
	mysql     *db.MySQLDB
}

// QueryRequest representa una solicitud de consulta
type QueryRequest struct {
	Query string        `json:"query"`
	Args  []interface{} `json:"args,omitempty"`
}

// NewHandlers crea una nueva instancia de Handlers
func NewHandlers(sqlServer *db.SQLServerDB, mysql *db.MySQLDB) *Handlers {
	return &Handlers{
		sqlServer: sqlServer,
		mysql:     mysql,
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

	var rows *sql.Rows
	var err error

	// Seleccionar la base de datos correcta
	if req.Database == "sqlserver" {
		rows, err = h.sqlServer.ExecuteQuery(req.Query, req.Args...)
	} else {
		rows, err = h.mysql.ExecuteQuery(req.Query, req.Args...)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Generar el Excel
	excelBytes, err := utils.GenerateExcel(rows, req.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar respuesta para descarga
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+req.Filename)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")

	w.Write(excelBytes)
}

// GetVentas obtiene las ventas según los filtros proporcionados
func (h *Handlers) GetVentas(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la consulta
	filtro := models.VentasFiltro{
		FechaInicio: r.URL.Query().Get("fechaInicio"),
		FechaFin:    r.URL.Query().Get("fechaFin"),
		Sucursal:    parseIntParam(r.URL.Query().Get("sucursal"), 211),
	}

	// Validar filtros
	if err := filtro.Validar(); err != nil {
		http.Error(w, fmt.Sprintf("Error en parámetros: %v", err), http.StatusBadRequest)
		return
	}

	// Obtener la consulta SQL
	query := sqlserver.GetVentasQuery()

	// Ejecutar la consulta con los parámetros
	log.Printf("Consultando ventas desde %s hasta %s para sucursal %d",
		filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal)

	// Usamos parámetros posicionales en el orden correcto
	rows, err := h.sqlServer.ExecuteQuery(query,
		filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Parámetros para BoletasConsolidadas
		filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Parámetros para FacturasConsolidadas
		filtro.Sucursal, filtro.FechaInicio, filtro.FechaFin) // Parámetros para la consulta final

	if err != nil {
		log.Printf("Error al consultar ventas: %v", err)
		http.Error(w, fmt.Sprintf("Error al ejecutar consulta: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Convertir resultados a JSON
	result, err := utils.RowsToJSON(rows)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al procesar resultados: %v", err), http.StatusInternalServerError)
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

	// Validar filtros
	if err := filtro.Validar(); err != nil {
		http.Error(w, fmt.Sprintf("Error en parámetros: %v", err), http.StatusBadRequest)
		return
	}

	// Nombre del archivo
	filename := fmt.Sprintf("Ventas_Sucursal_%d_%s_al_%s.xlsx",
		filtro.Sucursal, filtro.FechaInicio, filtro.FechaFin)

	// Obtener y ejecutar la consulta
	query := sqlserver.GetVentasQuery()
	rows, err := h.sqlServer.ExecuteQuery(query,
		filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal,
		filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal,
		filtro.Sucursal, filtro.FechaInicio, filtro.FechaFin)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al ejecutar consulta: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Generar Excel
	excelBytes, err := utils.GenerateExcel(rows, filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al generar Excel: %v", err), http.StatusInternalServerError)
		return
	}

	// Configurar respuesta para descarga
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
