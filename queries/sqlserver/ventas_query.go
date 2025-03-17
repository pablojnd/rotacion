package sqlserver

// GetVentasQuery devuelve la consulta SQL para obtener información de ventas
func GetVentasQuery() string {
	return `
WITH BoletasConsolidadas AS (
    SELECT DISTINCT
        'BOLETA' AS TipoDocumento,
        CAST(VB.CORRELATIVO AS VARCHAR(20)) AS NumeroDocumento,
        VB.FECHA_VENTA AS FechaDocumento,
        VB.ID_SUCURSAL AS Sucursal,
        CASE 
            WHEN VB.NULA = 1 THEN 'ANULADA'
            ELSE 'ACTIVA'
        END AS EstadoDocumento,
        CASE 
            WHEN VB.FORMA_PAGO IN ('Cheque', 'Cheques') THEN 'CHEQUE'
            WHEN VB.FORMA_PAGO = 'Contado' THEN 'CONTADO'
            WHEN VB.FORMA_PAGO = 'Cta Cte' THEN 'CTA CTE'
            WHEN VB.FORMA_PAGO = 'Convenios' THEN 'CONVENIO'
            WHEN VB.FORMA_PAGO IN ('Tarjeta Credito', 'Tarjeta de Credito', 'Tarjeta de Crédito') THEN 'TARJETA CRÉDITO'
            WHEN VB.FORMA_PAGO IN ('Tarjeta Débito', 'Tarjeta Debito', 'Tarjeta de Débito') THEN 'TARJETA DÉBITO'
            WHEN VB.FORMA_PAGO IN ('Transferencia', 'Transferencia Bancaria') THEN 'TRANSFERENCIA'
            ELSE 'OTROS'
        END AS FormaPagoUnificada,
        CASE 
            WHEN VB.NULA = 1 THEN 0
            ELSE VB.TOTAL
        END AS TotalDocumento,
        ISNULL(PTRA.NOMBRE_P + ' ' + PTRA.APELLIDOPATERNO_P, 'Sin Vendedor') AS Vendedor,
        ISNULL(PCLI.NOMBRE_P + ' ' + PCLI.APELLIDOPATERNO_P, 'Sin Cliente') AS Cliente,
        CASE 
            WHEN NV.CORRELATIVO IS NOT NULL THEN NV.CORRELATIVO
            ELSE NULL
        END AS NotaVenta,
        PR.CODIGO_INTERNO AS CodigoProducto,
        PR.NOMBRE_PRODUCTO AS Producto,
        PR.PRECIO_VENTA AS PrecioProducto,
        PR.PRECIO_OFERTA AS PrecioOferta,
        DB.CANTIDAD AS CantidadVendida,
        DB.VALORUNITARIO AS PrecioVenta,
        ROUND(((DB.TOTAL - DB.IMPUESTO - DB.iva) / T.valor) - DB.ley18219, 2) AS VentaNeta,
        ROUND(z.COSTO_UNITARIO, 2) AS CostoUnitario
    FROM VENTA_BOLETA VB
        LEFT JOIN DETALLE_VENTA_BOLETA DB ON VB.ID_VENTA_BOLETA = DB.ID_VENTA_BOLETA
        LEFT JOIN PRODUCTO PR ON DB.ID_PRODUCTO = PR.ID_PRODUCTO               
        LEFT JOIN TRABAJADOR TR ON VB.RUT_VEN = TR.RUT_P               
        LEFT JOIN PERSONA PTRA ON TR.RUT_P = PTRA.RUT_P               
        LEFT JOIN CLIENTE CL ON VB.RUT_CLI = CL.RUT_P               
        LEFT JOIN PERSONA PCLI ON CL.RUT_P = PCLI.RUT_P               
        LEFT JOIN NOTA_VENTA NV ON VB.NotaVenta = NV.ID_NOTA_VENTA OR VB.NotaVenta = NV.CORRELATIVO
        LEFT JOIN STOCKS z ON z.ID_SUCURSAL = VB.ID_SUCURSAL
                          AND z.ID_PRODUCTO = DB.ID_PRODUCTO
                          AND z.ZETA = DB.ZETA
                          AND z.ANIO = YEAR(VB.FECHA_VENTA)
        INNER JOIN TIPO_CAMBIO T ON T.fecha = VB.FECHA_VENTA
    WHERE VB.FECHA_VENTA >= ?
      AND VB.FECHA_VENTA <= ?
      AND VB.ID_SUCURSAL = ?
),
FacturasConsolidadas AS (
    SELECT DISTINCT
        'FACTURA' AS TipoDocumento,
        CAST(VF.CORRELATIVO AS VARCHAR(20)) AS NumeroDocumento,
        VF.FECHA_EMISION AS FechaDocumento,
        VF.ID_SUCURSAL AS Sucursal,
        CASE 
            WHEN VF.NULA = 1 THEN 'ANULADA'
            ELSE 'ACTIVA'
        END AS EstadoDocumento,
        CASE 
            WHEN VF.FORMA_PAGO IN ('Cheque', 'Cheques') THEN 'CHEQUE'
            WHEN VF.FORMA_PAGO = 'Contado' THEN 'CONTADO'
            WHEN VF.FORMA_PAGO = 'Cta Cte' THEN 'CTA CTE'
            WHEN VF.FORMA_PAGO = 'Convenios' THEN 'CONVENIO'
            WHEN VF.FORMA_PAGO IN ('Tarjeta Credito', 'Tarjeta de Credito', 'Tarjeta de Crédito') THEN 'TARJETA CRÉDITO'
            WHEN VF.FORMA_PAGO IN ('Tarjeta Débito', 'Tarjeta Debito', 'Tarjeta de Débito') THEN 'TARJETA DÉBITO'
            WHEN VF.FORMA_PAGO IN ('Transferencia', 'Transferencia Bancaria') THEN 'TRANSFERENCIA'
            ELSE 'OTROS'
        END AS FormaPagoUnificada,
        CASE 
            WHEN VF.NULA = 1 THEN 0
            ELSE VF.TOTAL
        END AS TotalDocumento,
        ISNULL(PTRAF.NOMBRE_P + ' ' + PTRAF.APELLIDOPATERNO_P, 'Sin Vendedor') AS Vendedor,
        ISNULL(PCLIF.NOMBRE_P + ' ' + PCLIF.APELLIDOPATERNO_P, 'Sin Cliente') AS Cliente,
        CASE 
            WHEN NV.CORRELATIVO IS NOT NULL THEN NV.CORRELATIVO
            ELSE NULL
        END AS NotaVenta,
        PR.CODIGO_INTERNO AS CodigoProducto,
        PR.NOMBRE_PRODUCTO AS Producto,
        PR.PRECIO_VENTA AS PrecioProducto,
        PR.PRECIO_OFERTA AS PrecioOferta,
        DF.CANTIDAD AS CantidadVendida,        
        DF.VALORUNITARIO AS PrecioVenta,
        ROUND(((DF.TOTAL - DF.IMPUESTO - DF.iva) / T.valor) - DF.ley18219, 2) AS VentaNeta,
        ROUND(z.COSTO_UNITARIO, 2) AS CostoUnitario
    FROM VENTA_FACTURA VF
        LEFT JOIN DETALLE_FAC_E DF ON VF.ID_VENTA_FACTURA = DF.ID_VENTA_FACTURA
        LEFT JOIN PRODUCTO PR ON DF.ID_PRODUCTO = PR.ID_PRODUCTO               
        LEFT JOIN TRABAJADOR TRAF ON VF.RUT_VEN = TRAF.RUT_P               
        LEFT JOIN PERSONA PTRAF ON TRAF.RUT_P = PTRAF.RUT_P               
        LEFT JOIN CLIENTE CLIF ON VF.RUT_CLI = CLIF.RUT_P               
        LEFT JOIN PERSONA PCLIF ON CLIF.RUT_P = PCLIF.RUT_P               
        LEFT JOIN NOTA_VENTA NV ON VF.NotaVenta = NV.ID_NOTA_VENTA OR VF.NotaVenta = NV.CORRELATIVO
        LEFT JOIN STOCKS z ON z.ID_SUCURSAL = VF.ID_SUCURSAL
                          AND z.ID_PRODUCTO = DF.ID_PRODUCTO
                          AND z.ZETA = DF.ZETA
                          AND z.ANIO = YEAR(VF.FECHA_EMISION)               
        INNER JOIN TIPO_CAMBIO T ON T.fecha = VF.FECHA_EMISION
    WHERE VF.FECHA_EMISION >= ?
      AND VF.FECHA_EMISION <= ?
      AND VF.ID_SUCURSAL = ?
),
DocumentosUnificados AS (
    SELECT * FROM BoletasConsolidadas
    UNION ALL
    SELECT * FROM FacturasConsolidadas
),
PromediosProducto AS (
    SELECT
        CodigoProducto,
        ROUND(SUM(PrecioVenta * CantidadVendida) / NULLIF(SUM(CantidadVendida), 0), 2) AS AvgPrecioVenta,
        ROUND(SUM(CostoUnitario * CantidadVendida) / NULLIF(SUM(CantidadVendida), 0), 2) AS AvgCostoUnitario,
        COUNT(DISTINCT NumeroDocumento) AS CantidadDeVentas
    FROM DocumentosUnificados
    GROUP BY CodigoProducto
)
SELECT
    CONVERT(VARCHAR(10), d.FechaDocumento, 120) AS FechaDocumento,
    d.Sucursal,
    d.EstadoDocumento,
    d.TotalDocumento,
    d.NotaVenta,
    d.CodigoProducto,
    d.Producto,
    d.PrecioProducto,
    d.PrecioOferta,
    d.CantidadVendida,
    p.AvgPrecioVenta AS PrecioVenta,
    p.AvgCostoUnitario AS CostoUnitario,
    p.CantidadDeVentas
FROM DocumentosUnificados d
JOIN PromediosProducto p ON d.CodigoProducto = p.CodigoProducto
WHERE d.Sucursal = ?
  AND d.FechaDocumento BETWEEN ? AND ?
ORDER BY d.NumeroDocumento DESC`
}
