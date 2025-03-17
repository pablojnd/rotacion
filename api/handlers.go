package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/yourusername/rotacion/db"
	"github.com/yourusername/rotacion/utils"
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
