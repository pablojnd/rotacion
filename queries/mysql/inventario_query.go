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
    CASE
        WHEN MAX(IFNULL(NomDetSubFam, '')) <> '' AND MAX(IFNULL(NomDetSubFam, '')) <> 'POR ASIGNAR' AND MAX(IFNULL(NomDetSubFam, '')) <> 'Sin Asignar' THEN MAX(IFNULL(NomDetSubFam, ''))
        ELSE 
            -- Intentamos extraer dimensiones del nombre del producto
            CASE 
                -- Buscar patrón de dimensiones como "20X60", "20X60 CM", "20X60 CMS", etc.
                WHEN MAX(DES_ADU) REGEXP '[0-9]+[Xx][0-9]+'
                THEN 
                    SUBSTRING(
                        MAX(DES_ADU),
                        GREATEST(
                            LOCATE(' ', MAX(DES_ADU), LOCATE('[0-9]+[Xx][0-9]+', MAX(DES_ADU))),
                            LOCATE(' ', REVERSE(MAX(DES_ADU)), LOCATE('[0-9]+[Xx][0-9]+', REVERSE(MAX(DES_ADU))))
                        )
                    )
                -- Si no encontramos un patrón claro, indicamos que está por asignar
                ELSE 'POR ASIGNAR'
            END
    END AS DIMENSIONES,
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
    DATEDIFF(CURDATE(), MIN(fec_ing)) AS Dias_Desde_Primer_Ingreso,
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
    AND (? = '' OR COD_ART LIKE CONCAT('%', ?, '%'))
GROUP BY
    COD_ART,
    UNI_CAJ
ORDER BY Codigo_Producto
`
}

// GetExtractDimensionsFunction devuelve una función MySQL para extraer dimensiones de nombres de productos
func GetExtractDimensionsFunction() string {
	return `
CREATE FUNCTION IF NOT EXISTS ExtractDimensions(productName VARCHAR(255))
RETURNS VARCHAR(100)
DETERMINISTIC
BEGIN
    DECLARE dimensionStr VARCHAR(100);
    DECLARE startPos INT;
    DECLARE endPos INT;
    
    -- Buscar patrones comunes de dimensiones: números seguidos de X y más números
    -- Por ejemplo: 20X60, 20x60, 20X60 CMS
    SET startPos = REGEXP_INSTR(productName, '[0-9]+[Xx][0-9]+');
    
    -- Si encontramos el patrón
    IF startPos > 0 THEN
        -- Buscar hacia atrás hasta el último espacio
        WHILE startPos > 1 AND SUBSTRING(productName, startPos-1, 1) != ' ' DO
            SET startPos = startPos - 1;
        END WHILE;
        
        -- Buscar hacia adelante hasta el próximo espacio o final
        SET endPos = startPos;
        WHILE endPos < LENGTH(productName) AND SUBSTRING(productName, endPos, 1) != ' ' DO
            SET endPos = endPos + 1;
        END WHILE;
        
        -- Si termina con CMS, CM, etc., incluirlo
        IF endPos < LENGTH(productName) - 5 AND 
           (UPPER(SUBSTRING(productName, endPos+1, 3)) = 'CMS' OR 
            UPPER(SUBSTRING(productName, endPos+1, 2)) = 'CM') THEN
            SET endPos = endPos + 4;
        END IF;
        
        -- Extraer la dimensión
        SET dimensionStr = TRIM(SUBSTRING(productName, startPos, endPos - startPos + 1));
        
        RETURN dimensionStr;
    ELSE
        RETURN 'POR ASIGNAR';
    END IF;
END;
`
}

// Para MySQL 5.1 que no soporta bien REGEXP_INSTR, usaremos una versión simplificada:
func GetSimplifiedDimensionsQuery() string {
	return `
SELECT
    COD_ART AS Codigo_Producto,
    MAX(DES_ADU) AS Nombre_Producto,
    MAX(IFNULL(Marca, 'POR ASIGNAR')) AS MARCA,
    MAX(IFNULL(SubFamilia, 'POR ASIGNAR')) AS CATEGORIA,
    CASE
        WHEN MAX(IFNULL(NomDetSubFam, '')) <> '' AND MAX(IFNULL(NomDetSubFam, '')) <> 'POR ASIGNAR' AND MAX(IFNULL(NomDetSubFam, '')) <> 'Sin Asignar' THEN MAX(IFNULL(NomDetSubFam, ''))
        ELSE 
            -- Intento manual de extraer dimensiones con funciones básicas
            CASE 
                WHEN LOCATE('X', MAX(DES_ADU)) > 0 THEN
                    SUBSTRING(
                        MAX(DES_ADU),
                        GREATEST(
                            LOCATE(' ', MAX(DES_ADU), LOCATE('X', MAX(DES_ADU)) - 5),
                            1
                        ),
                        LEAST(
                            LOCATE(' ', MAX(DES_ADU), LOCATE('X', MAX(DES_ADU)) + 3) - 
                            GREATEST(LOCATE(' ', MAX(DES_ADU), LOCATE('X', MAX(DES_ADU)) - 5), 1),
                            20
                        )
                    )
                ELSE 'POR ASIGNAR'
            END
    END AS DIMENSIONES,
    UNI_CAJ AS Unidades_Por_Paquete,
    SUM(CAN_ING) AS Cantidad_Total_Ingresada,
    ROUND(IF(SUM(CAN_ING) = 0, 0, SUM(CIF_UNI * CAN_ING) / SUM(CAN_ING)), 2) AS Promedio_Costo_CIF,
    ROUND(IF(SUM(COS_UNI * CAN_ING) = 0, 0, SUM(COS_UNI * CAN_ING) / SUM(CAN_ING)), 2) AS Costo_Promedio_IKA,
    MIN(fec_ing) AS Fecha_Primer_Ingreso,
    MAX(fec_ing) AS Fecha_Ultimo_Ingreso,
    DATEDIFF(CURDATE(), MIN(fec_ing)) AS Dias_Desde_Primer_Ingreso,
    COUNT(*) AS Cantidad_De_Ingresos,
    CONCAT('[', GROUP_CONCAT(DISTINCT CONCAT(
        '{\'Zeta\':\'', ZET_ART, '\',',
        '\'Anio_Produccion\':\'', ANIO_PRO, '\',',
        '\'Cantidad_Ingresada\':', CAN_ING, ',',
        '\'Fecha_Ingreso\':\'', DATE_FORMAT(fec_ing, '%Y-%m-%d'), '\'}'
    ) SEPARATOR ','), ']') AS Metadatos_JSON
FROM saldos s
WHERE
    CAST(ANIO_PRO AS SIGNED) = ?
    AND (? = '' OR COD_ART LIKE CONCAT('%', ?, '%'))
GROUP BY
    COD_ART,
    UNI_CAJ
ORDER BY Codigo_Producto
`
}
