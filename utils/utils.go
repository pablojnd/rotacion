package utils

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

// RowsToJSON convierte filas SQL a un slice de mapas
func RowsToJSON(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)

		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, column := range columns {
			value := values[i]

			// Convertir tipos especiales de SQL al formato adecuado
			if value != nil {
				switch columnTypes[i].DatabaseTypeName() {
				case "DECIMAL", "NUMERIC", "MONEY", "SMALLMONEY":
					if v, ok := value.([]byte); ok {
						row[column] = string(v)
					} else {
						row[column] = value
					}
				default:
					if b, ok := value.([]byte); ok {
						row[column] = string(b)
					} else {
						row[column] = value
					}
				}
			} else {
				row[column] = nil
			}
		}

		result = append(result, row)
	}

	return result, nil
}

// GenerateExcel genera un archivo Excel a partir de filas SQL
func GenerateExcel(rows *sql.Rows, filename string) ([]byte, error) {
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

	// Establecer encabezados
	for i, col := range columns {
		cell := fmt.Sprintf("%c1", 'A'+i) // A1, B1, C1, etc.
		f.SetCellValue(sheetName, cell, col)
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
				f.SetCellValue(sheetName, cell, values[i])
			}
		}
		rowIndex++
	}

	// Guardar en buffer
	var buffer bytes.Buffer
	if err := f.Write(&buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// ParseInt convierte un string a int con un valor por defecto si hay error
func ParseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// ParseFloat convierte un string a float64 con un valor por defecto si hay error
func ParseFloat(value string, defaultValue float64) float64 {
	if value == "" {
		return defaultValue
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return floatValue
}

// ParseDate convierte un string a time.Time con formato "2006-01-02"
func ParseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

// AnalyzarTipos analiza los tipos de datos de las columnas
func AnalyzarTipos(rows *sql.Rows) {
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		log.Printf("Error al obtener tipos de columnas: %v", err)
		return
	}

	log.Println("Análisis de tipos de columnas:")
	for _, ct := range columnTypes {
		precision, scale, ok := ct.DecimalSize()
		if ok {
			log.Printf("- Columna: %s, Tipo: %s, Precision: %d, Scale: %d",
				ct.Name(), ct.DatabaseTypeName(), precision, scale)
		} else {
			length, ok := ct.Length()
			if ok {
				log.Printf("- Columna: %s, Tipo: %s, Length: %d",
					ct.Name(), ct.DatabaseTypeName(), length)
			} else {
				log.Printf("- Columna: %s, Tipo: %s",
					ct.Name(), ct.DatabaseTypeName())
			}
		}
	}
}
