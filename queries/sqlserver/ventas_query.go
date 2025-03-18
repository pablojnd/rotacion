package sqlserver

// GetVentasQuery devuelve la consulta SQL para obtener información de ventas
func GetVentasQuery() string {
	return `
WITH Documentos AS (
    -- BOLETAS
    SELECT 
        VB.FECHA_VENTA AS FechaDocumento,
        VB.ID_SUCURSAL AS Sucursal,
        PR.NOMBRE_PRODUCTO AS Producto,
        ROUND(z.COSTO_UNITARIO, 2) AS CostoUnitario,
        PR.PRECIO_VENTA AS PrecioProducto,
        PR.PRECIO_OFERTA AS PrecioOferta,
        DB.VALORUNITARIO AS PrecioVenta,
        DB.CANTIDAD AS CantidadVendida,
        CASE WHEN VB.NULA = 1 THEN 0 ELSE VB.TOTAL END AS TotalDocumento,
        PR.CODIGO_INTERNO AS CodigoProducto,
        'BOLETA' AS TipoDocumento,
        CAST(VB.CORRELATIVO AS VARCHAR(20)) AS CodigoDocumento,
        ISNULL(PCLI.NOMBRE_P + ' ' + PCLI.APELLIDOPATERNO_P, 'Sin Cliente') AS Cliente
    FROM VENTA_BOLETA VB
    INNER JOIN DETALLE_VENTA_BOLETA DB ON VB.ID_VENTA_BOLETA = DB.ID_VENTA_BOLETA
    INNER JOIN PRODUCTO PR ON DB.ID_PRODUCTO = PR.ID_PRODUCTO
    LEFT JOIN CLIENTE CL ON VB.RUT_CLI = CL.RUT_P
    LEFT JOIN PERSONA PCLI ON CL.RUT_P = PCLI.RUT_P
    LEFT JOIN STOCKS z ON 
        z.ID_SUCURSAL = VB.ID_SUCURSAL AND
        z.ID_PRODUCTO = DB.ID_PRODUCTO AND
        z.ZETA = DB.ZETA AND
        z.ANIO = YEAR(VB.FECHA_VENTA)
    WHERE 
        VB.FECHA_VENTA BETWEEN ? AND ? AND
        VB.ID_SUCURSAL = ?
    
    UNION ALL
    
    -- FACTURAS
    SELECT 
        VF.FECHA_EMISION AS FechaDocumento,
        VF.ID_SUCURSAL AS Sucursal,
        PR.NOMBRE_PRODUCTO AS Producto,
        ROUND(z.COSTO_UNITARIO, 2) AS CostoUnitario,
        PR.PRECIO_VENTA AS PrecioProducto,
        PR.PRECIO_OFERTA AS PrecioOferta,
        DF.VALORUNITARIO AS PrecioVenta,
        DF.CANTIDAD AS CantidadVendida,
        CASE WHEN VF.NULA = 1 THEN 0 ELSE VF.TOTAL END AS TotalDocumento,
        PR.CODIGO_INTERNO AS CodigoProducto,
        'FACTURA' AS TipoDocumento,
        CAST(VF.CORRELATIVO AS VARCHAR(20)) AS CodigoDocumento,
        ISNULL(PCLIF.NOMBRE_P + ' ' + PCLIF.APELLIDOPATERNO_P, 'Sin Cliente') AS Cliente
    FROM VENTA_FACTURA VF
    INNER JOIN DETALLE_FAC_E DF ON VF.ID_VENTA_FACTURA = DF.ID_VENTA_FACTURA
    INNER JOIN PRODUCTO PR ON DF.ID_PRODUCTO = PR.ID_PRODUCTO
    LEFT JOIN CLIENTE CLIF ON VF.RUT_CLI = CLIF.RUT_P
    LEFT JOIN PERSONA PCLIF ON CLIF.RUT_P = PCLIF.RUT_P
    LEFT JOIN STOCKS z ON 
        z.ID_SUCURSAL = VF.ID_SUCURSAL AND
        z.ID_PRODUCTO = DF.ID_PRODUCTO AND
        z.ZETA = DF.ZETA AND
        z.ANIO = YEAR(VF.FECHA_EMISION)
    WHERE 
        VF.FECHA_EMISION BETWEEN ? AND ? AND
        VF.ID_SUCURSAL = ?
),
Calculos AS (
    SELECT 
        CodigoProducto,
        ROUND(
            SUM(PrecioVenta * CantidadVendida) / NULLIF(SUM(CantidadVendida), 0),
            2
        ) AS PrecioPromedio,
        COUNT(*) AS CantidadDeVentas
    FROM Documentos
    GROUP BY CodigoProducto
)
SELECT 
    d.CodigoDocumento,
    d.FechaDocumento AS Fecha_Emision,
    d.TipoDocumento AS Tipo_Documento,
    d.Cliente,
    d.CodigoProducto AS Codigo_Producto,
    d.Producto,
    d.CantidadVendida AS Cantidad,
    d.PrecioVenta AS Precio_Unitario,
    d.TotalDocumento AS Total_Venta,
    d.Sucursal,
    d.CostoUnitario,
    d.PrecioProducto,
    d.PrecioOferta,
    c.PrecioPromedio,
    c.CantidadDeVentas
FROM Documentos d
INNER JOIN Calculos c ON d.CodigoProducto = c.CodigoProducto
WHERE (? = '' OR d.CodigoProducto LIKE '%' + ? + '%')
ORDER BY d.FechaDocumento DESC
`
}

// GetVentasAgrupadasQuery devuelve la consulta SQL para obtener ventas agrupadas por producto
func GetVentasAgrupadasQuery() string {
	return `
WITH VentasAgrupadas AS (
    -- BOLETAS
    SELECT 
        PR.CODIGO_INTERNO AS CodigoProducto,
        MAX(PR.NOMBRE_PRODUCTO) AS Producto,
        MAX(ROUND(z.COSTO_UNITARIO, 2)) AS CostoUnitario,
        MAX(PR.PRECIO_VENTA) AS PrecioProducto,
        MAX(PR.PRECIO_OFERTA) AS PrecioOferta,
        SUM(DB.CANTIDAD) AS CantidadVendida,
        SUM(CASE WHEN VB.NULA = 1 THEN 0 ELSE VB.TOTAL END) AS TotalDocumento,
        MAX(VB.FECHA_VENTA) AS UltimaFechaVenta
    FROM VENTA_BOLETA VB
    INNER JOIN DETALLE_VENTA_BOLETA DB ON VB.ID_VENTA_BOLETA = DB.ID_VENTA_BOLETA
    INNER JOIN PRODUCTO PR ON DB.ID_PRODUCTO = PR.ID_PRODUCTO
    LEFT JOIN STOCKS z ON 
        z.ID_SUCURSAL = VB.ID_SUCURSAL AND
        z.ID_PRODUCTO = DB.ID_PRODUCTO AND
        z.ZETA = DB.ZETA AND
        z.ANIO = YEAR(VB.FECHA_VENTA)
    WHERE 
        VB.FECHA_VENTA BETWEEN ? AND ? AND
        VB.ID_SUCURSAL = ?
    GROUP BY PR.CODIGO_INTERNO

    UNION ALL

    -- FACTURAS
    SELECT 
        PR.CODIGO_INTERNO AS CodigoProducto,
        MAX(PR.NOMBRE_PRODUCTO) AS Producto,
        MAX(ROUND(z.COSTO_UNITARIO, 2)) AS CostoUnitario,
        MAX(PR.PRECIO_VENTA) AS PrecioProducto,
        MAX(PR.PRECIO_OFERTA) AS PrecioOferta,
        SUM(DF.CANTIDAD) AS CantidadVendida,
        SUM(CASE WHEN VF.NULA = 1 THEN 0 ELSE VF.TOTAL END) AS TotalDocumento,
        MAX(VF.FECHA_EMISION) AS UltimaFechaVenta
    FROM VENTA_FACTURA VF
    INNER JOIN DETALLE_FAC_E DF ON VF.ID_VENTA_FACTURA = DF.ID_VENTA_FACTURA
    INNER JOIN PRODUCTO PR ON DF.ID_PRODUCTO = PR.ID_PRODUCTO
    LEFT JOIN STOCKS z ON 
        z.ID_SUCURSAL = VF.ID_SUCURSAL AND
        z.ID_PRODUCTO = DF.ID_PRODUCTO AND
        z.ZETA = DF.ZETA AND
        z.ANIO = YEAR(VF.FECHA_EMISION)
    WHERE 
        VF.FECHA_EMISION BETWEEN ? AND ? AND
        VF.ID_SUCURSAL = ?
    GROUP BY PR.CODIGO_INTERNO
),
CalculosFinales AS (
    SELECT 
        CodigoProducto,
        MAX(Producto) AS Producto,
        MAX(CostoUnitario) AS CostoUnitario,
        MAX(PrecioProducto) AS PrecioProducto,
        MAX(PrecioOferta) AS PrecioOferta,
        SUM(CantidadVendida) AS CantidadVendidaTotal,
        SUM(TotalDocumento) AS TotalVentas,
        MAX(UltimaFechaVenta) AS UltimaFechaVenta,
        ROUND(
            SUM(PrecioProducto * CantidadVendida) / NULLIF(SUM(CantidadVendida), 0),
            2
        ) AS PrecioPromedio,
        COUNT(*) AS CantidadDeVentas
    FROM VentasAgrupadas
    GROUP BY CodigoProducto
)
SELECT 
    CodigoProducto,
    Producto,
    CostoUnitario,
    PrecioProducto,
    PrecioOferta,
    CantidadVendidaTotal,
    TotalVentas,
    PrecioPromedio,
    CantidadDeVentas,
    UltimaFechaVenta
FROM CalculosFinales
WHERE (? = '' OR CodigoProducto LIKE '%' + ? + '%')
ORDER BY CodigoProducto
`
}
