package services

import (
	"encoding/json"
	"regexp"
	"strings"

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
	query := mysql.GetSimplifiedDimensionsQuery()

	// Ejecutar la consulta con los parámetros
	rows, err := s.mysql.ExecuteQuery(query,
		filtro.Anio,
		filtro.CodigoProducto,
		filtro.CodigoProducto)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Convertir resultados a JSON
	result, err := utils.RowsToJSON(rows)
	if err != nil {
		return nil, err
	}

	// Procesar los metadatos JSON
	for i := range result {
		// Procesar metadatos JSON (ahora llamado "Historial de Ingresos (JSON)")
		if metadataJSON, ok := result[i]["Historial de Ingresos (JSON)"].(string); ok {
			// Arreglar las comillas simples en el JSON
			fixedJSON := utils.FixJSONQuotes(metadataJSON)

			// Intentar analizar como JSON válido para verificar
			var metadata []interface{}
			if err := json.Unmarshal([]byte(fixedJSON), &metadata); err == nil {
				// Si el análisis tiene éxito, reemplazar el string con el objeto JSON analizado
				result[i]["Historial de Ingresos (JSON)"] = metadata
			} else {
				// Si hay error, mantener el string arreglado
				result[i]["Historial de Ingresos (JSON)"] = fixedJSON
			}
		}

		// Procesar dimensiones si es necesario (adaptado a los nuevos nombres)
		if dimensiones, ok := result[i]["Subcategoría/Dimensiones"].(string); ok {
			if dimensiones == "POR ASIGNAR" || dimensiones == "Sin Asignar" || dimensiones == "" {
				if nombreProducto, ok := result[i]["Nombre Aduanero"].(string); ok {
					// Extraer dimensiones del nombre del producto
					result[i]["Subcategoría/Dimensiones"] = extractDimensionsFromName(nombreProducto)
				}
			}
		}
	}

	return result, nil
}

// extractDimensionsFromName extrae las dimensiones del nombre del producto
func extractDimensionsFromName(name string) string {
	// Buscamos patrones comunes de dimensiones: NxN, NXN, N X N, etc.
	// seguidos opcionalmente por CM, CMS, etc.
	re := regexp.MustCompile(`(\d+\s*[Xx]\s*\d+(?:\s*(?:CM|CMS|MM)?)?)`)

	matches := re.FindStringSubmatch(name)
	if len(matches) > 0 {
		return strings.TrimSpace(matches[0])
	}

	return "POR ASIGNAR"
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
