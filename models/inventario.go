package models

import (
	"time"
)

// Inventario representa un elemento de inventario
type Inventario struct {
	CodigoProducto            string     `json:"codigoProducto"`
	Marca                     string     `json:"marca"`
	Categoria                 string     `json:"categoria"`
	Dimensiones               string     `json:"dimensiones"`
	Zeta                      string     `json:"zeta"`
	AnioProduccion            string     `json:"anioProduccion"`
	NombreProducto            string     `json:"nombreProducto"`
	Packing                   float64    `json:"packing"`
	CantidadTotal             float64    `json:"cantidadTotal"`
	PromedioPonderadoCostoCIF *float64   `json:"promedioPonderadoCostoCIF"`
	CostoReal                 float64    `json:"costoReal"`
	FechaIngreso              *time.Time `json:"fechaIngreso"`
	DiasDesdeFechaIngreso     *int       `json:"diasDesdeFechaIngreso"`
	GlobalPromedioCostoCIF    *float64   `json:"globalPromedioCostoCIF"`
}

// InventarioFiltro define los filtros para consultar inventario
type InventarioFiltro struct {
	Anio int `json:"anio"`
}

// Validar valida los parámetros del filtro
func (f *InventarioFiltro) Validar() error {
	// Validar año
	if f.Anio <= 0 {
		f.Anio = time.Now().Year() // Usar año actual como predeterminado
	}

	return nil
}
