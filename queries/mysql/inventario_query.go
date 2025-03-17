package mysql

// GetInventarioQuery devuelve la consulta SQL para obtener información de inventario
func GetInventarioQuery() string {
	return `
SELECT
    COD_ART AS Codigo_Producto,
    MAX(DES_ADU) AS Nombre_Producto,
    MAX(IFNULL(Marca, 'POR ASIGNAR')) AS MARCA,
    MAX(
        IFNULL(SubFamilia, 'POR ASIGNAR')
    ) AS CATEGORIA,
    MAX(
        IFNULL(NomDetSubFam, 'POR ASIGNAR')
    ) AS DIMENSIONES,
    UNI_CAJ AS Unidades_Por_Paquete,
    SUM(CAN_ING) AS Cantidad_Total_Ingresada,
    ROUND(
    IF(
        SUM(CAN_ING) = 0,
        0,
        SUM(CIF_UNI * CAN_ING) / SUM(CAN_ING)
    ),
    2
) AS Promedio_Costo_CIF,
ROUND(
    IF(
        SUM(COS_UNI) = 0,
        0,
        SUM(COS_UNI * CAN_ING) / SUM(CAN_ING)
    ),
    2
) AS Costo_Promedio_IKA,
    MIN(fec_ing) AS Fecha_Primer_Ingreso,
    MAX(fec_ing) AS Fecha_Ultimo_Ingreso,
    COUNT(*) AS Cantidad_De_Ingresos,
    -- Metadatos en formato JSON (compatible con MySQL 5.1)
    CONCAT(
        '[',
        GROUP_CONCAT(
            DISTINCT CONCAT(
                '{\'Zeta\':\'',
                ZET_ART,
                '\',',
                '\'Anio_Produccion\':\'',
                ANIO_PRO,
                '\',',
                '\'Cantidad_Ingresada\':',
                CAN_ING,
                ',',
                '\'Fecha_Ingreso\':\'',
                DATE_FORMAT(fec_ing, '%Y-%m-%d'),
                '\'}'
            ) SEPARATOR ','
        ),
        ']'
    ) AS Metadatos_JSON
FROM saldos s
WHERE
    CAST(ANIO_PRO AS SIGNED) = ?
GROUP BY
    COD_ART,
    UNI_CAJ
ORDER BY Codigo_Producto
`
}
