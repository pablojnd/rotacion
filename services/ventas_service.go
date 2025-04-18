package services

import (
	"github.com/pablojnd/rotacion/db"
	"github.com/pablojnd/rotacion/models"
	"github.com/pablojnd/rotacion/queries/sqlserver"
	"github.com/pablojnd/rotacion/utils"
)

// VentasService proporciona métodos para trabajar con datos de ventas
type VentasService struct {
	sqlServer    *db.SQLServerDB
	excelService *ExcelService
}

// NewVentasService crea un nuevo servicio de ventas
func NewVentasService(sqlServer *db.SQLServerDB, excelService *ExcelService) *VentasService {
	return &VentasService{
		sqlServer:    sqlServer,
		excelService: excelService,
	}
}

// GetVentas obtiene las ventas según los filtros proporcionados
func (s *VentasService) GetVentas(filtro models.VentasFiltro, tipo models.TipoConsultaVentas) ([]map[string]interface{}, error) {
	// Validar filtros
	if err := filtro.Validar(); err != nil {
		return nil, err
	}

	var query string
	var args []interface{}

	// Seleccionar la consulta según el tipo
	if tipo == models.ConsultaVentasAgrupada {
		query = sqlserver.GetVentasAgrupadasQuery()
		args = []interface{}{
			filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para Boletas
			filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para Facturas
			filtro.CodigoProducto, filtro.CodigoProducto, // Para filtrado por código
		}
	} else {
		query = sqlserver.GetVentasQuery()
		args = []interface{}{
			filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para Boletas
			filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para Facturas
			filtro.CodigoProducto, filtro.CodigoProducto, // Para filtrado por código
		}
	}

	// Ejecutar la consulta con los parámetros
	rows, err := s.sqlServer.ExecuteQuery(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Convertir resultados a JSON
	result, err := utils.RowsToJSON(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetVentasAgrupadas obtiene las ventas agrupadas por producto según los filtros proporcionados
func (s *VentasService) GetVentasAgrupadas(filtro models.VentasFiltro) ([]map[string]interface{}, error) {
	return s.GetVentas(filtro, models.ConsultaVentasAgrupada)
}

// ExportVentasToExcel exporta ventas a un archivo Excel
func (s *VentasService) ExportVentasToExcel(filtro models.VentasFiltro, tipo models.TipoConsultaVentas) ([]byte, string, error) {
	// Validar filtros
	if err := filtro.Validar(); err != nil {
		return nil, "", err
	}

	// Nombre del archivo
	filename := generateVentasFilename(filtro)

	var query string
	var args []interface{}

	// Seleccionar la consulta según el tipo
	if tipo == models.ConsultaVentasAgrupada {
		query = sqlserver.GetVentasAgrupadasQuery()
		args = []interface{}{
			filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para Boletas
			filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para Facturas
			filtro.CodigoProducto, filtro.CodigoProducto, // Para filtrado por código
		}
	} else {
		query = sqlserver.GetVentasQuery()
		args = []interface{}{
			filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para Boletas
			filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para Facturas
			filtro.CodigoProducto, filtro.CodigoProducto, // Para filtrado por código
		}
	}

	// Obtener y ejecutar la consulta
	rows, err := s.sqlServer.ExecuteQuery(query, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	// Generar Excel
	excelBytes, err := s.excelService.GenerateExcel(rows, filename)
	if err != nil {
		return nil, "", err
	}

	return excelBytes, filename, nil
}

// ExportVentasAgrupadasToExcel exporta ventas agrupadas a un archivo Excel
func (s *VentasService) ExportVentasAgrupadasToExcel(filtro models.VentasFiltro) ([]byte, string, error) {
	return s.ExportVentasToExcel(filtro, models.ConsultaVentasAgrupada)
}

// generateVentasFilename genera un nombre de archivo para el reporte de ventas
func generateVentasFilename(filtro models.VentasFiltro) string {
	base := "Ventas_Sucursal_" + string(filtro.Sucursal) + "_" + filtro.FechaInicio + "_al_" + filtro.FechaFin
	if filtro.CodigoProducto != "" {
		base += "_Producto_" + filtro.CodigoProducto
	}
	return base + ".xlsx"
}
