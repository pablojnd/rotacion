package models

// ReporteCombinado representa la unión de datos de inventario y ventas
type ReporteCombinado struct {
	// Datos generales del producto
	CodigoProducto string  `json:"Codigo_Producto"`
	Marca          string  `json:"MARCA"`
	Categoria      string  `json:"CATEGORIA"`
	Dimensiones    string  `json:"DIMENSIONES"`
	Nombre         string  `json:"NOMBRE"`
	Packing        float64 `json:"PACKING"`

	// Datos de inventario (MySQL)
	CifPromedioUsd     float64 `json:"CIF_PROMEDIO_USD"`
	CifPromedioClp     float64 `json:"CIF_PROMEDIO_CLP"`
	CantidadIngresada  float64 `json:"CANTIDAD_INGRESADA"`
	DiasEnInventario   int     `json:"CANTIDAD_DE_DIAS_EN_INVENTARIO"`
	FechaPrimerIngreso string  `json:"FECHA_PRIMER_INGRESO"`
	FechaUltimoIngreso string  `json:"FECHA_ULTIMO_INGRESO"`
	CantidadIngresos   int     `json:"CANTIDAD_INGRESOS"`
	HistorialIngresos  string  `json:"HISTORIAL_INGRESOS"`

	// Datos de ventas (SQL Server)
	PrecioProductoClp      int     `json:"PRECIO_PRODUCTO_CLP"`
	PrecioOfertaClp        int     `json:"PRECIO_OFERTA_CLP"`
	PrecioVentaPromedioClp int     `json:"PROMEDIO_DEL_PRECIO_VENTA_CLP"`
	CantidadVendida        float64 `json:"CANTIDAD_VENDIDA"`
	RankingCantidad        int     `json:"RANKING_POR_CANTIDAD_VENDIDA"`
	VentaNetaTotalClp      int     `json:"VENTA_NETA_TOTAL_CLP"`
	UltimaFechaVenta       string  `json:"ULTIMA_FECHA_VENTA"`
	CantidadTransacciones  int     `json:"CANTIDAD_TRANSACCIONES"`

	// Campos calculados
	PorcentajeVendido float64 `json:"PORCENTAJE_VENDIDO"`
	UtilidadClp       float64 `json:"UTILIDAD_CLP"`
	RankingVenta      int     `json:"RANKING_VENTA"`
}

// ReporteFiltro define los filtros para consultar el reporte combinado
type ReporteFiltro struct {
	Anio           int    `json:"anio"`
	FechaInicio    string `json:"fechaInicio"`
	FechaFin       string `json:"fechaFin"`
	Sucursal       int    `json:"sucursal"`
	CodigoProducto string `json:"codigoProducto"`
}
