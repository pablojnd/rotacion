package db

import (
	"database/sql"
	"io"
)

// MockRows implementa una versión mock de sql.Rows para uso con datos en memoria
type MockRows struct {
	columns []string
	data    [][]interface{}
	current int
	closed  bool
}

// NewMockRows crea un nuevo objeto MockRows con los datos proporcionados
func NewMockRows(columns []string, data [][]interface{}) *MockRows {
	return &MockRows{
		columns: columns,
		data:    data,
		current: -1,
		closed:  false,
	}
}

// Columns devuelve los nombres de las columnas
func (m *MockRows) Columns() ([]string, error) {
	if m.closed {
		return nil, sql.ErrNoRows
	}
	return m.columns, nil
}

// Next avanza al siguiente registro
func (m *MockRows) Next() bool {
	if m.closed {
		return false
	}
	m.current++
	return m.current < len(m.data)
}

// Scan copia los valores del registro actual a los destinos proporcionados
func (m *MockRows) Scan(dest ...interface{}) error {
	if m.closed {
		return sql.ErrNoRows
	}
	if m.current < 0 || m.current >= len(m.data) {
		return io.EOF
	}

	row := m.data[m.current]
	if len(dest) != len(row) {
		return sql.ErrNoRows
	}

	for i, val := range row {
		d := dest[i].(*interface{})
		*d = val
	}

	return nil
}

// Close cierra las filas
func (m *MockRows) Close() error {
	m.closed = true
	return nil
}

// ColumnTypes devuelve información sobre los tipos de columna
func (m *MockRows) ColumnTypes() ([]*sql.ColumnType, error) {
	if m.closed {
		return nil, sql.ErrNoRows
	}
	// Esta es una implementación simplificada, en realidad deberíamos obtener
	// los tipos reales de las columnas
	types := make([]*sql.ColumnType, len(m.columns))
	// Note: esta parte es compleja de implementar completamente
	return types, nil
}

// Err devuelve el último error que ocurrió durante la iteración
func (m *MockRows) Err() error {
	return nil
}

// CreateMockRowsFromMaps crea un objeto MockRows a partir de una lista de mapas
func CreateMockRowsFromMaps(data []map[string]interface{}) (*MockRows, error) {
	if len(data) == 0 {
		return nil, sql.ErrNoRows
	}

	// Extraer las columnas del primer mapa
	columns := make([]string, 0, len(data[0]))
	for k := range data[0] {
		columns = append(columns, k)
	}

	// Construir los datos en formato de filas
	rows := make([][]interface{}, len(data))
	for i, item := range data {
		row := make([]interface{}, len(columns))
		for j, col := range columns {
			row[j] = item[col]
		}
		rows[i] = row
	}

	return NewMockRows(columns, rows), nil
}
