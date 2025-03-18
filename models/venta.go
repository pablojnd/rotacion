package models

import (
	"errors"
	"time"
)

// Venta representa una venta realizada
type Venta struct {
	CodigoDocumento  string  `json:"codigoDocumento"`
	FechaEmision     string  `json:"fechaEmision"`
	TipoDocumento    string  `json:"tipoDocumento"`
	Cliente          string  `json:"cliente"`
	CodigoProducto   string  `json:"codigoProducto"`
	Producto         string  `json:"producto"`
	Cantidad         float64 `json:"cantidad"`
	PrecioUnitario   float64 `json:"precioUnitario"`
	TotalVenta       float64 `json:"totalVenta"`
	Sucursal         int     `json:"sucursal"`
	CostoUnitario    float64 `json:"costoUnitario"`
	PrecioProducto   float64 `json:"precioProducto"`
	PrecioOferta     float64 `json:"precioOferta"`
	PrecioPromedio   float64 `json:"precioPromedio"`
	CantidadDeVentas int     `json:"cantidadDeVentas"`
}

// VentasFiltro define los filtros para consultar ventas
type VentasFiltro struct {
	FechaInicio    string `json:"fechaInicio"`
	FechaFin       string `json:"fechaFin"`
	Sucursal       int    `json:"sucursal"`
	CodigoProducto string `json:"codigoProducto"` // Nuevo campo para filtrar por código
}

// Validar valida los parámetros del filtro
func (f *VentasFiltro) Validar() error {
	// Validar fechas
	if f.FechaInicio == "" || f.FechaFin == "" {
		return errors.New("las fechas de inicio y fin son obligatorias")
	}

	// Validar formato de fechas
	_, err := time.Parse("2006-01-02", f.FechaInicio)
	if err != nil {
		return errors.New("fecha de inicio no válida, debe estar en formato YYYY-MM-DD")
	}
	_, err = time.Parse("2006-01-02", f.FechaFin)
	if err != nil {
		return errors.New("fecha de fin no válida, debe estar en formato YYYY-MM-DD")
	}

	// Validar sucursal
	if f.Sucursal <= 0 {
		f.Sucursal = 211 // Usando sucursal por defecto si no se especifica
	}

	return nil
}
