package mysql

// GetInventarioQuery devuelve la consulta SQL para obtener información de inventario
func GetInventarioQuery() string {
	return `
SELECT
    COD_ART AS Codigo_Producto,
    Marca AS MARCA,
    SubFamilia AS CATEGORIA,
    saldos.NomDetSubFam AS DIMENSIONES,
    ZET_ART AS Zeta,
    ANIO_PRO AS Anio_Produccion,
    DES_INT AS Nombre_Producto,
    UNI_CAJ AS Packing,
    SUM(CAN_ING) AS Cantidad_Total,
    SUM(CIF_UNI * CAN_ING) / SUM(CAN_ING) AS PromedioPonderadoCostoCIF,
    cos_uni AS Costo_Real,
    MAX(fec_ing) AS Fecha_Ingreso,
    DATEDIFF(CURDATE(), (
        SELECT MAX(fec_ing) 
        FROM saldos AS s2 
        WHERE s2.COD_ART = saldos.COD_ART 
          AND s2.ZET_ART = saldos.ZET_ART 
          AND s2.ANIO_PRO = saldos.ANIO_PRO
    )) AS Dias_Desde_Ingreso,
    (SELECT SUM(s2.CIF_UNI * s2.CAN_ING) / SUM(s2.CAN_ING)
     FROM saldos s2
     WHERE s2.COD_ART = saldos.COD_ART
       AND s2.ANIO_PRO = saldos.ANIO_PRO
    ) AS GlobalPromedioCostoCIF
FROM saldos
WHERE CAST(ANIO_PRO AS SIGNED) = ?
GROUP BY 
    COD_ART, Marca, SubFamilia, saldos.NomDetSubFam, ZET_ART, ANIO_PRO, DES_INT, UNI_CAJ, cos_uni
ORDER BY COD_ART, ANIO_PRO
`
}
