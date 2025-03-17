package models

import (
	"time"
)

// Inventario representa un elemento de inventario
type Inventario struct {
	CodigoProducto         string    `json:"codigoProducto"`
	NombreProducto         string    `json:"nombreProducto"`
	Marca                  string    `json:"marca"`
	Categoria              string    `json:"categoria"`
	Dimensiones            string    `json:"dimensiones"`
	UnidadesPorPaquete     float64   `json:"unidadesPorPaquete"`
	CantidadTotalIngresada float64   `json:"cantidadTotalIngresada"`
	PromedioCostoCIF       float64   `json:"promedioCostoCIF"` // Cambiado el nombre
	CostoPromedioIKA       float64   `json:"costoPromedioIKA"` // Campo nuevo
	FechaPrimerIngreso     time.Time `json:"fechaPrimerIngreso"`
	FechaUltimoIngreso     time.Time `json:"fechaUltimoIngreso"`
	CantidadDeIngresos     int       `json:"cantidadDeIngresos"`
	MetadatosJSON          string    `json:"metadatosJSON"`
}

// MetadatoProducto representa un elemento dentro del array MetadatosJSON
type MetadatoProducto struct {
	Zeta              string    `json:"Zeta"`
	AnioProduccion    string    `json:"Anio_Produccion"`
	CantidadIngresada float64   `json:"Cantidad_Ingresada"`
	FechaIngreso      time.Time `json:"Fecha_Ingreso"`
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
