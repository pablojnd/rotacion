package models

import (
	"database/sql"
	"time"
)

// Venta representa el modelo de datos para una venta
type Venta struct {
	FechaDocumento   string  `json:"fechaDocumento"`
	Sucursal         int     `json:"sucursal"`
	EstadoDocumento  string  `json:"estadoDocumento"`
	TotalDocumento   float64 `json:"totalDocumento"`
	NotaVenta        int     `json:"notaVenta"`
	CodigoProducto   string  `json:"codigoProducto"`
	Producto         string  `json:"producto"`
	PrecioProducto   float64 `json:"precioProducto"`
	PrecioOferta     float64 `json:"precioOferta"`
	CantidadVendida  float64 `json:"cantidadVendida"`
	PrecioVenta      float64 `json:"precioVenta"`
	CostoUnitario    float64 `json:"costoUnitario"`
	CantidadDeVentas int     `json:"cantidadDeVentas"`
}

// VentasFiltro define los filtros para la consulta de ventas
type VentasFiltro struct {
	FechaInicio string `json:"fechaInicio"`
	FechaFin    string `json:"fechaFin"`
	Sucursal    int    `json:"sucursal"`
}

// Validar valida los parámetros del filtro
func (f *VentasFiltro) Validar() error {
	// Validar fechas
	if f.FechaInicio == "" || f.FechaFin == "" {
		return sql.ErrNoRows
	}

	// Validar formato de fechas
	_, err := time.Parse("2006-01-02", f.FechaInicio)
	if err != nil {
		return err
	}

	_, err = time.Parse("2006-01-02", f.FechaFin)
	if err != nil {
		return err
	}

	// Validar sucursal
	if f.Sucursal <= 0 {
		f.Sucursal = 211 // Valor predeterminado
	}

	return nil
}
