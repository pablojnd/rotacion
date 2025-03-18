package services

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExcelService proporciona métodos para generar archivos Excel
type ExcelService struct{}

// NewExcelService crea un nuevo servicio de Excel
func NewExcelService() *ExcelService {
	return &ExcelService{}
}

// GenerateExcelFromQuery ejecuta una consulta y genera un excel con los resultados
func (s *ExcelService) GenerateExcelFromQuery(db interface{}, query string, args []interface{}, filename string) ([]byte, error) {
	var rows *sql.Rows
	var err error

	// Ejecutar la consulta en la base de datos correspondiente
	switch dbInstance := db.(type) {
	case *sql.DB:
		rows, err = dbInstance.Query(query, args...)
	default:
		return nil, fmt.Errorf("tipo de base de datos no soportado")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.GenerateExcel(rows, filename)
}

// GenerateExcel genera un archivo Excel a partir de filas SQL
func (s *ExcelService) GenerateExcel(rows *sql.Rows, filename string) ([]byte, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Crear un nuevo archivo Excel
	f := excelize.NewFile()
	defer f.Close()

	// Crear hoja de trabajo
	sheetName := "Datos"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)

	// Establecer encabezados con un estilo para destacarlos
	styleHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#DCE6F1"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})

	if err != nil {
		return nil, err
	}

	// Aplicar encabezados
	for i, col := range columns {
		cell := fmt.Sprintf("%c1", 'A'+i) // A1, B1, C1, etc.
		f.SetCellValue(sheetName, cell, col)
		f.SetCellStyle(sheetName, cell, cell, styleHeader)
	}

	// Agregar datos
	rowIndex := 2 // comenzamos desde la fila 2
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))

		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		for i := range columns {
			cell := fmt.Sprintf("%c%d", 'A'+i, rowIndex) // A2, B2, C2, etc.

			if values[i] == nil {
				f.SetCellValue(sheetName, cell, "")
			} else {
				// Manejar formato específico para decimales
				switch v := values[i].(type) {
				case float32:
					// Para valores decimales, formateamos manualmente con coma
					if needsCommaDecimal(columns[i]) {
						// Convertir a string con el formato deseado (coma como separador decimal)
						strValue := formatFloatWithComma(float64(v))
						// Establecer como texto para preservar la coma
						f.SetCellStr(sheetName, cell, strValue)
					} else {
						f.SetCellValue(sheetName, cell, v)
					}
				case float64:
					// Para valores decimales, formateamos manualmente con coma
					if needsCommaDecimal(columns[i]) {
						// Convertir a string con el formato deseado (coma como separador decimal)
						strValue := formatFloatWithComma(v)
						// Establecer como texto para preservar la coma
						f.SetCellStr(sheetName, cell, strValue)
					} else {
						f.SetCellValue(sheetName, cell, v)
					}
				case []byte:
					// Verificar si es un número decimal en formato string
					str := string(v)
					if strings.Contains(str, ".") && isNumeric(str) && needsCommaDecimal(columns[i]) {
						// Intentar convertir a float64
						if floatVal, err := strconv.ParseFloat(str, 64); err == nil {
							// Convertir de nuevo a string con coma
							f.SetCellStr(sheetName, cell, formatFloatWithComma(floatVal))
						} else {
							f.SetCellValue(sheetName, cell, str)
						}
					} else {
						f.SetCellValue(sheetName, cell, str)
					}
				default:
					f.SetCellValue(sheetName, cell, v)
				}
			}
		}
		rowIndex++
	}

	// Auto-ajustar columnas
	for i := range columns {
		colName := string(rune('A' + i))
		f.SetColWidth(sheetName, colName, colName, 18) // Ancho inicial razonable
	}

	// Guardar en buffer
	var buffer bytes.Buffer
	if err := f.Write(&buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// needsCommaDecimal determina si la columna debería usar coma como separador decimal
func needsCommaDecimal(columnName string) bool {
	// Columnas que necesitan formato decimal con coma
	decimalColumns := []string{
		// Columnas existentes
		"Costo Unitario (USD)",
		"Cantidad Total Vendida",
		"Cantidad",
		"CostoUnitarioUSD",
		"CantidadTotalVendida",
		"CantidadVendida",
		"Unidades_Por_Paquete",
		"Promedio_Costo_CIF",
		"Costo_Promedio_IKA",

		// Nuevas columnas en español
		"Costo Promedio CIF (USD)",
		"Costo Promedio Unitario (CLP)",
		"Unidades por Caja",
		"Total Unidades Ingresadas",
	}

	for _, col := range decimalColumns {
		if columnName == col {
			return true
		}
	}
	return false
}

// formatFloatWithComma formatea un float con coma como separador decimal
func formatFloatWithComma(value float64) string {
	// Formatear con 2 decimales
	str := strconv.FormatFloat(value, 'f', 2, 64)
	// Reemplazar punto por coma
	return strings.Replace(str, ".", ",", 1)
}

// isNumeric verifica si una cadena es numérica
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// SendExcelResponse envía un archivo Excel como respuesta HTTP
func SendExcelResponse(w http.ResponseWriter, excelBytes []byte, filename string) {
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")

	w.Write(excelBytes)
}
