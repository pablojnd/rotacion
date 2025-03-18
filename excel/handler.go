package excel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pablojnd/rotacion/db"
	"github.com/pablojnd/rotacion/models"
	"github.com/pablojnd/rotacion/services"
)

// Handler gestiona las solicitudes relacionadas con la exportación a Excel
type Handler struct {
	sqlServer         *db.SQLServerDB
	mysql             *db.MySQLDB
	excelService      *services.ExcelService
	ventasService     *services.VentasService
	inventarioService *services.InventarioService
}

// ExportRequest representa una solicitud de exportación a Excel
type ExportRequest struct {
	Query    string        `json:"query"`
	Args     []interface{} `json:"args,omitempty"`
	Database string        `json:"database"` // "sqlserver" o "mysql"
	Filename string        `json:"filename"`
}

// NewHandler crea un nuevo manejador para operaciones Excel
func NewHandler(
	sqlServer *db.SQLServerDB,
	mysql *db.MySQLDB,
	excelService *services.ExcelService,
	ventasService *services.VentasService,
	inventarioService *services.InventarioService,
) *Handler {
	return &Handler{
		sqlServer:         sqlServer,
		mysql:             mysql,
		excelService:      excelService,
		ventasService:     ventasService,
		inventarioService: inventarioService,
	}
}

// ExportGeneric exporta datos a Excel desde una consulta genérica
func (h *Handler) ExportGeneric(w http.ResponseWriter, r *http.Request) {
	var req ExportRequest
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
	services.SendExcelResponse(w, excelBytes, req.Filename)
}

// ExportVentas exporta las ventas a Excel
func (h *Handler) ExportVentas(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la consulta
	filtro := models.VentasFiltro{
		FechaInicio:    r.URL.Query().Get("fechaInicio"),
		FechaFin:       r.URL.Query().Get("fechaFin"),
		Sucursal:       parseIntParam(r.URL.Query().Get("sucursal"), 211),
		CodigoProducto: r.URL.Query().Get("codigo"),
	}

	// Determinar el tipo de consulta
	tipo := models.ConsultaVentasDetallada
	if r.URL.Query().Get("tipo") == "agrupada" {
		tipo = models.ConsultaVentasAgrupada
	}

	// Usar el servicio para exportar a Excel
	excelBytes, filename, err := h.ventasService.ExportVentasToExcel(filtro, tipo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al generar Excel: %v", err), http.StatusInternalServerError)
		return
	}

	// Enviar respuesta
	services.SendExcelResponse(w, excelBytes, filename)
}

// ExportVentasAgrupadas exporta las ventas agrupadas a Excel
func (h *Handler) ExportVentasAgrupadas(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la consulta
	filtro := models.VentasFiltro{
		FechaInicio:    r.URL.Query().Get("fechaInicio"),
		FechaFin:       r.URL.Query().Get("fechaFin"),
		Sucursal:       parseIntParam(r.URL.Query().Get("sucursal"), 211),
		CodigoProducto: r.URL.Query().Get("codigo"),
	}

	// Usar el servicio para exportar a Excel
	excelBytes, filename, err := h.ventasService.ExportVentasAgrupadasToExcel(filtro)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al generar Excel: %v", err), http.StatusInternalServerError)
		return
	}

	// Enviar respuesta
	services.SendExcelResponse(w, excelBytes, filename)
}

// ExportInventario exporta el inventario a Excel
func (h *Handler) ExportInventario(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la consulta
	anio := parseIntParam(r.URL.Query().Get("anio"), time.Now().Year())
	codigoProducto := r.URL.Query().Get("codigo")

	filtro := models.InventarioFiltro{
		Anio:           anio,
		CodigoProducto: codigoProducto,
	}

	// Usar el servicio para exportar a Excel
	excelBytes, filename, err := h.inventarioService.ExportInventarioToExcel(filtro)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al generar Excel: %v", err), http.StatusInternalServerError)
		return
	}

	// Enviar respuesta
	services.SendExcelResponse(w, excelBytes, filename)
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
