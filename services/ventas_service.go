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
func (s *VentasService) GetVentas(filtro models.VentasFiltro) ([]map[string]interface{}, error) {
	// Validar filtros
	if err := filtro.Validar(); err != nil {
		return nil, err
	}

	// Obtener la consulta SQL
	query := sqlserver.GetVentasQuery()

	// Ejecutar la consulta con los parámetros - adaptado para la nueva consulta
	rows, err := s.sqlServer.ExecuteQuery(query,
		filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para BoletasConsolidadas
		filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para FacturasConsolidadas
		filtro.CodigoProducto, filtro.CodigoProducto) // Para el filtrado por código

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

// ExportVentasToExcel exporta ventas a un archivo Excel
func (s *VentasService) ExportVentasToExcel(filtro models.VentasFiltro) ([]byte, string, error) {
	// Validar filtros
	if err := filtro.Validar(); err != nil {
		return nil, "", err
	}

	// Nombre del archivo
	filename := generateVentasFilename(filtro)

	// Obtener y ejecutar la consulta
	query := sqlserver.GetVentasQuery()
	rows, err := s.sqlServer.ExecuteQuery(query,
		filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para BoletasConsolidadas
		filtro.FechaInicio, filtro.FechaFin, filtro.Sucursal, // Para FacturasConsolidadas
		filtro.CodigoProducto, filtro.CodigoProducto) // Para el filtrado por código

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

// generateVentasFilename genera un nombre de archivo para el reporte de ventas
func generateVentasFilename(filtro models.VentasFiltro) string {
	base := "Ventas_Sucursal_" + string(filtro.Sucursal) + "_" + filtro.FechaInicio + "_al_" + filtro.FechaFin
	if filtro.CodigoProducto != "" {
		base += "_Producto_" + filtro.CodigoProducto
	}
	return base + ".xlsx"
}
