package models

import (
	"time"
)

// Inventario representa un elemento de inventario con nombres en español
type Inventario struct {
	CodigoProducto          string    `json:"codigoProducto"`
	NombreAduanero          string    `json:"nombreAduanero"`
	MarcaProducto           string    `json:"marcaProducto"`
	CategoriaPrincipal      string    `json:"categoriaPrincipal"`
	SubcategoriaDimensiones string    `json:"subcategoriaDimensiones"`
	UnidadesPorCaja         float64   `json:"unidadesPorCaja"`
	TotalUnidadesIngresadas float64   `json:"totalUnidadesIngresadas"`
	CostoPromedioCIF        float64   `json:"costoPromedioCIF"`
	CostoPromedioUnitario   float64   `json:"costoPromedioUnitario"`
	FechaPrimerIngreso      time.Time `json:"fechaPrimerIngreso"`
	FechaUltimoIngreso      time.Time `json:"fechaUltimoIngreso"`
	DiasDesdePrimerIngreso  int       `json:"diasDesdePrimerIngreso"`
	CantidadIngresos        int       `json:"cantidadIngresos"`
	HistorialIngresos       string    `json:"historialIngresos"`
}

// MetadatoProducto representa un elemento dentro del array de historial de ingresos
type MetadatoProducto struct {
	Zeta               string    `json:"Zeta"`
	AnioProduccion     string    `json:"Año Producción"`
	UnidadesIngresadas float64   `json:"Unidades Ingresadas"`
	FechaIngreso       time.Time `json:"Fecha Ingreso"`
}

// InventarioFiltro define los filtros para consultar inventario
type InventarioFiltro struct {
	Anio           int    `json:"anio"`
	CodigoProducto string `json:"codigoProducto"` // Nuevo campo para filtrar por código
}

// Validar valida los parámetros del filtro
func (f *InventarioFiltro) Validar() error {
	// Validar año
	if f.Anio <= 0 {
		f.Anio = time.Now().Year() // Usar año actual como predeterminado
	}

	return nil
}
