package models

// VentaDetallada representa una venta individual con detalles completos
type VentaDetallada struct {
	CodigoDocumento       string  `json:"codigoDocumento"`
	FechaEmision          string  `json:"fechaEmision"`
	TipoDocumento         string  `json:"tipoDocumento"`
	Cliente               string  `json:"cliente"`
	CodigoProducto        string  `json:"codigoProducto"`
	NombreProducto        string  `json:"nombreProducto"`
	Cantidad              float64 `json:"cantidad"`
	PrecioUnitarioCLP     int     `json:"precioUnitarioCLP"`
	TotalVentaCLP         int     `json:"totalVentaCLP"`
	Sucursal              int     `json:"sucursal"`
	CostoUnitarioUSD      float64 `json:"costoUnitarioUSD"`
	PrecioBaseCLP         int     `json:"precioBaseCLP"`
	PrecioOfertaCLP       int     `json:"precioOfertaCLP"`
	PrecioPromedioCLP     int     `json:"precioPromedioCLP"`
	CantidadTransacciones int     `json:"cantidadTransacciones"`
}
