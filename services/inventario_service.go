package services

import (
	"encoding/json"

	"github.com/pablojnd/rotacion/db"
	"github.com/pablojnd/rotacion/models"
	"github.com/pablojnd/rotacion/queries/mysql"
	"github.com/pablojnd/rotacion/utils"
)

// InventarioService proporciona métodos para trabajar con datos de inventario
type InventarioService struct {
	mysql        *db.MySQLDB
	excelService *ExcelService
}

// NewInventarioService crea un nuevo servicio de inventario
func NewInventarioService(mysql *db.MySQLDB, excelService *ExcelService) *InventarioService {
	return &InventarioService{
		mysql:        mysql,
		excelService: excelService,
	}
}

// GetInventario obtiene el inventario según los filtros proporcionados
func (s *InventarioService) GetInventario(filtro models.InventarioFiltro) ([]map[string]interface{}, error) {
	// Validar filtros
	if err := filtro.Validar(); err != nil {
		return nil, err
	}

	// Obtener la consulta SQL
	query := mysql.GetInventarioQuery()

	// Ejecutar la consulta con los parámetros
	rows, err := s.mysql.ExecuteQuery(query,
		filtro.Anio,
		filtro.CodigoProducto,
		filtro.CodigoProducto) // pasamos dos veces para el CONCAT en la consulta
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Convertir resultados a JSON
	result, err := utils.RowsToJSON(rows)
	if err != nil {
		return nil, err
	}

	// Procesar los metadatos JSON para corregir las comillas
	for i := range result {
		if metadataJSON, ok := result[i]["Metadatos_JSON"].(string); ok {
			// Arreglar las comillas simples en el JSON
			fixedJSON := utils.FixJSONQuotes(metadataJSON)

			// Intentar analizar como JSON válido para verificar
			var metadata []interface{}
			if err := json.Unmarshal([]byte(fixedJSON), &metadata); err == nil {
				// Si el análisis tiene éxito, reemplazar el string con el objeto JSON analizado
				result[i]["Metadatos_JSON"] = metadata
			} else {
				// Si hay error, mantener el string arreglado
				result[i]["Metadatos_JSON"] = fixedJSON
			}
		}
	}

	return result, nil
}

// ExportInventarioToExcel exporta el inventario a un archivo Excel
func (s *InventarioService) ExportInventarioToExcel(filtro models.InventarioFiltro) ([]byte, string, error) {
	// Validar filtros
	if err := filtro.Validar(); err != nil {
		return nil, "", err
	}

	// Nombre del archivo
	filename := generateInventarioFilename(filtro)

	// Obtener y ejecutar la consulta
	query := mysql.GetInventarioQuery()
	rows, err := s.mysql.ExecuteQuery(query,
		filtro.Anio,
		filtro.CodigoProducto,
		filtro.CodigoProducto) // pasamos dos veces para el CONCAT en la consulta
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

// generateInventarioFilename genera un nombre de archivo para el reporte de inventario
func generateInventarioFilename(filtro models.InventarioFiltro) string {
	base := "Inventario_" + string(filtro.Anio)
	if filtro.CodigoProducto != "" {
		base += "_" + filtro.CodigoProducto
	}
	return base + ".xlsx"
}
