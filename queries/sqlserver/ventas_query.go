package sqlserver

// GetVentasQuery devuelve la consulta SQL para obtener información de ventas detalladas
func GetVentasQuery() string {
	return `
WITH VentasBase AS (
    -- BOLETAS
    SELECT 
        VB.FECHA_VENTA AS FechaDocumento,
        VB.ID_SUCURSAL AS Sucursal,
        PR.CODIGO_INTERNO AS CodigoProducto,
        PR.NOMBRE_PRODUCTO AS NombreProducto,
        ROUND(z.COSTO_UNITARIO, 2) AS CostoUnitarioUSD,
        PR.PRECIO_VENTA AS PrecioBaseCLP,
        PR.PRECIO_OFERTA AS PrecioOfertaCLP,
        DB.VALORUNITARIO AS PrecioVentaCLP,
        DB.CANTIDAD AS CantidadVendida,
        CASE WHEN VB.NULA = 1 THEN 0 ELSE DB.TOTAL END AS TotalProductoCLP,
        CASE WHEN VB.NULA = 1 THEN 0 ELSE VB.TOTAL END AS TotalDocumentoCLP,
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
        PR.CODIGO_INTERNO AS CodigoProducto,
        PR.NOMBRE_PRODUCTO AS NombreProducto,
        ROUND(z.COSTO_UNITARIO, 2) AS CostoUnitarioUSD,
        PR.PRECIO_VENTA AS PrecioBaseCLP,
        PR.PRECIO_OFERTA AS PrecioOfertaCLP,
        DF.VALORUNITARIO AS PrecioVentaCLP,
        DF.CANTIDAD AS CantidadVendida,
        CASE WHEN VF.NULA = 1 THEN 0 ELSE DF.TOTAL END AS TotalProductoCLP,
        CASE WHEN VF.NULA = 1 THEN 0 ELSE VF.TOTAL END AS TotalDocumentoCLP,
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
CalculosProducto AS (
    SELECT 
        CodigoProducto,
        CAST(
            ROUND(
                SUM(PrecioVentaCLP * CantidadVendida) / NULLIF(SUM(CantidadVendida), 0),
                0
            ) AS INT
        ) AS PrecioPromedioCLP,
        COUNT(*) AS CantidadTransacciones
    FROM VentasBase
    GROUP BY CodigoProducto
)
SELECT 
    v.CodigoDocumento AS "Código Documento",
    v.FechaDocumento AS "Fecha Emisión",
    v.TipoDocumento AS "Tipo Documento",
    v.Cliente AS "Cliente",
    v.CodigoProducto AS "Código Producto",
    v.NombreProducto AS "Producto",
    v.CantidadVendida AS "Cantidad",
    v.PrecioVentaCLP AS "Precio Unitario (CLP)",
    v.TotalProductoCLP AS "Total Venta (CLP)",
    v.Sucursal AS "Sucursal",
    v.CostoUnitarioUSD AS "Costo Unitario (USD)",
    v.PrecioBaseCLP AS "Precio Base (CLP)",
    v.PrecioOfertaCLP AS "Precio Oferta (CLP)",
    c.PrecioPromedioCLP AS "Precio Promedio (CLP)",
    c.CantidadTransacciones AS "Cant. Transacciones"
FROM VentasBase v
INNER JOIN CalculosProducto c ON v.CodigoProducto = c.CodigoProducto
WHERE (? = '' OR v.CodigoProducto LIKE '%' + ? + '%')
ORDER BY v.FechaDocumento DESC
`
}

// GetVentasAgrupadasQuery devuelve la consulta SQL para obtener ventas agrupadas por producto
func GetVentasAgrupadasQuery() string {
	return `
WITH VentasBase AS (
    -- BOLETAS
    SELECT 
        VB.FECHA_VENTA AS FechaDocumento,
        VB.ID_SUCURSAL AS Sucursal,
        PR.CODIGO_INTERNO AS CodigoProducto,
        PR.NOMBRE_PRODUCTO AS NombreProducto,
        ROUND(z.COSTO_UNITARIO, 2) AS CostoUnitarioUSD, -- Asumiendo que el costo está en USD
        PR.PRECIO_VENTA AS PrecioBaseCLP, -- Precio base en CLP
        PR.PRECIO_OFERTA AS PrecioOfertaCLP, -- Precio de oferta en CLP
        DB.VALORUNITARIO AS PrecioVentaCLP, -- Precio real de venta en CLP
        DB.CANTIDAD AS CantidadVendida,
        CASE WHEN VB.NULA = 1 THEN 0 ELSE VB.TOTAL END AS TotalDocumentoCLP
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

    UNION ALL

    -- FACTURAS
    SELECT 
        VF.FECHA_EMISION AS FechaDocumento,
        VF.ID_SUCURSAL AS Sucursal,
        PR.CODIGO_INTERNO AS CodigoProducto,
        PR.NOMBRE_PRODUCTO AS NombreProducto,
        ROUND(z.COSTO_UNITARIO, 2) AS CostoUnitarioUSD,
        PR.PRECIO_VENTA AS PrecioBaseCLP,
        PR.PRECIO_OFERTA AS PrecioOfertaCLP,
        DF.VALORUNITARIO AS PrecioVentaCLP,
        DF.CANTIDAD AS CantidadVendida,
        CASE WHEN VF.NULA = 1 THEN 0 ELSE VF.TOTAL END AS TotalDocumentoCLP
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
),
Resumen AS (
    SELECT 
        CodigoProducto,
        MAX(NombreProducto) AS NombreProducto,
        MAX(CostoUnitarioUSD) AS CostoUnitarioUSD,
        MAX(PrecioBaseCLP) AS PrecioBaseCLP,
        MAX(PrecioOfertaCLP) AS PrecioOfertaCLP,
        SUM(CantidadVendida) AS CantidadTotalVendida,
        SUM(TotalDocumentoCLP) AS TotalVentasCLP,
        MAX(FechaDocumento) AS UltimaFechaVenta,
        CAST(
            ROUND(
                SUM(PrecioVentaCLP * CantidadVendida) / NULLIF(SUM(CantidadVendida), 0),
                0
            ) AS INT
        ) AS PrecioPromedioPonderadoCLP,
        MIN(PrecioVentaCLP) AS PrecioMinimoCLP,
        MAX(PrecioVentaCLP) AS PrecioMaximoCLP,
        COUNT(*) AS CantidadDeVentas
    FROM VentasBase
    GROUP BY CodigoProducto
)
SELECT 
    CodigoProducto AS "Código de Producto",
    NombreProducto AS "Nombre del Producto",
    CostoUnitarioUSD AS "Costo Unitario (USD)",
    PrecioBaseCLP AS "Precio Base (CLP)",
    PrecioOfertaCLP AS "Precio de Oferta (CLP)",
    CantidadTotalVendida AS "Cantidad Total Vendida",
    TotalVentasCLP AS "Total Ventas (CLP)",
    UltimaFechaVenta AS "Última Fecha de Venta",
    PrecioPromedioPonderadoCLP AS "Precio Promedio Ponderado (CLP)",
    PrecioMinimoCLP AS "Precio Mínimo (CLP)",
    PrecioMaximoCLP AS "Precio Máximo (CLP)",
    CantidadDeVentas AS "Cantidad de Ventas Registradas"
FROM Resumen
WHERE (? = '' OR CodigoProducto LIKE '%' + ? + '%')
ORDER BY CodigoProducto
`
}
