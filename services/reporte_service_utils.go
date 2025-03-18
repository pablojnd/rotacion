package services

import (
	"github.com/pablojnd/rotacion/db"
)

// ConvertMapToExcelRows convierte un mapa a filas para Excel
func ConvertMapToExcelRows(data []map[string]interface{}) (*db.MockRows, []byte, error) {
	if len(data) == 0 {
		return nil, nil, nil
	}

	mockRows, err := db.CreateMockRowsFromMaps(data)
	if err != nil {
		return nil, nil, err
	}

	// El segundo valor ([]byte) se mantiene como nil porque no lo necesitamos realmente
	// pero es requerido por la firma de la función en reporte_service.go
	return mockRows, nil, nil
}
