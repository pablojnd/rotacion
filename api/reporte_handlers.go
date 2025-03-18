package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pablojnd/rotacion/models"
	"github.com/pablojnd/rotacion/services"
)

// ReporteHandlers contiene handlers para las operaciones de reportes combinados
type ReporteHandlers struct {
	reporteService *services.ReporteService
}

// NewReporteHandlers crea una nueva instancia de ReporteHandlers
func NewReporteHandlers(reporteService *services.ReporteService) *ReporteHandlers {
	return &ReporteHandlers{reporteService: reporteService}
}

// ObtenerReporteCombinado obtiene un reporte que combina datos de inventario y ventas
func (h *ReporteHandlers) ObtenerReporteCombinado(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de consulta
	anio := parseIntParam(r.URL.Query().Get("anio"), time.Now().Year())
	fechaInicio := r.URL.Query().Get("fechaInicio")
	fechaFin := r.URL.Query().Get("fechaFin")
	sucursal := parseIntParam(r.URL.Query().Get("sucursal"), 211)
	codigoProducto := r.URL.Query().Get("codigo")

	// Si no se proporcionaron fechas, usar el año actual
	if fechaInicio == "" {
		fechaInicio = time.Date(anio, 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	}
	if fechaFin == "" {
		fechaFin = time.Date(anio, 12, 31, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	}

	// Crear filtro
	filtro := models.ReporteFiltro{
		Anio:           anio,
		FechaInicio:    fechaInicio,
		FechaFin:       fechaFin,
		Sucursal:       sucursal,
		CodigoProducto: codigoProducto,
	}

	// Obtener reporte
	reportesCoincidentes, reportesSinCoincidencia, err := h.reporteService.GenerarReporteCombinado(filtro)
	if err != nil {
		log.Printf("Error al generar reporte combinado: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver respuesta combinada
	result := struct {
		ReportesCoincidentes    []models.ReporteCombinado `json:"reportesCoincidentes"`
		ReportesSinCoincidencia []models.ReporteCombinado `json:"reportesSinCoincidencia"`
	}{
		ReportesCoincidentes:    reportesCoincidentes,
		ReportesSinCoincidencia: reportesSinCoincidencia,
	}

	// Devolver respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ExportarReporteCombinado exporta un reporte combinado a Excel
func (h *ReporteHandlers) ExportarReporteCombinado(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de consulta
	anio := parseIntParam(r.URL.Query().Get("anio"), time.Now().Year())
	fechaInicio := r.URL.Query().Get("fechaInicio")
	fechaFin := r.URL.Query().Get("fechaFin")
	sucursal := parseIntParam(r.URL.Query().Get("sucursal"), 211)
	codigoProducto := r.URL.Query().Get("codigo")

	// Si no se proporcionaron fechas, usar el año actual
	if fechaInicio == "" {
		fechaInicio = time.Date(anio, 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	}
	if fechaFin == "" {
		fechaFin = time.Date(anio, 12, 31, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	}

	// Crear filtro
	filtro := models.ReporteFiltro{
		Anio:           anio,
		FechaInicio:    fechaInicio,
		FechaFin:       fechaFin,
		Sucursal:       sucursal,
		CodigoProducto: codigoProducto,
	}

	// Generar Excel
	excelBytes, filename, err := h.reporteService.ExportarReporteCombinado(filtro)
	if err != nil {
		log.Printf("Error al exportar reporte combinado: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enviar respuesta
	services.SendExcelResponse(w, excelBytes, filename)
}
