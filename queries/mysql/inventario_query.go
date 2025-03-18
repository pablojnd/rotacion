package mysql

// GetInventarioQuery devuelve la consulta SQL para obtener información de inventario
func GetInventarioQuery() string {
	return `
SELECT
    -- Identificación del producto
    COD_ART AS "Código de Producto",
    MAX(DES_ADU) AS "Nombre Aduanero", 
    MAX(IFNULL(Marca, 'POR ASIGNAR')) AS "Marca del Producto",
    MAX(IFNULL(SubFamilia, 'POR ASIGNAR')) AS "Categoría Principal",
    MAX(IFNULL(NomDetSubFam, 'POR ASIGNAR')) AS "Subcategoría/Dimensiones",

    -- Especificaciones de empaque
    UNI_CAJ AS "Unidades por Caja",

    -- Métricas de ingresos
    SUM(CAN_ING) AS "Total Unidades Ingresadas",
    ROUND(
        IF(SUM(CAN_ING) = 0, 0, 
           SUM(CIF_UNI * CAN_ING) / SUM(CAN_ING)),
        2
    ) AS "Costo Promedio CIF (USD)", 
    ROUND(
        IF(SUM(COS_UNI) = 0, 0, 
           SUM(COS_UNI * CAN_ING) / SUM(CAN_ING)),
        2
    ) AS "Costo Promedio Unitario (CLP)", 

    -- Fechas relevantes
    MIN(fec_ing) AS "Fecha Primer Ingreso",
    MAX(fec_ing) AS "Fecha Último Ingreso",
    DATEDIFF(CURDATE(), MIN(fec_ing)) AS "Días Desde Primer Ingreso",

    -- Conteo de registros
    COUNT(*) AS "Cantidad de Ingresos",

    -- Metadatos en JSON
    CONCAT(
        '[',
        GROUP_CONCAT(
            DISTINCT CONCAT(
                '{\'Zeta\':\'', ZET_ART, '\',',
                '\'Año Producción\':\'', ANIO_PRO, '\',',
                '\'Unidades Ingresadas\':', CAN_ING, ',',
                '\'Fecha Ingreso\':\'', DATE_FORMAT(fec_ing, '%Y-%m-%d'), '\'}'
            ) SEPARATOR ','
        ),
        ']'
    ) AS "Historial de Ingresos (JSON)"
FROM saldos s
WHERE
    CAST(ANIO_PRO AS SIGNED) = ?
    AND (? = '' OR COD_ART LIKE CONCAT('%', ?, '%'))
GROUP BY 
    COD_ART, 
    UNI_CAJ
ORDER BY "Código de Producto"
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
    -- Identificación del producto
    COD_ART AS "Código de Producto",
    MAX(DES_ADU) AS "Nombre Aduanero", 
    MAX(IFNULL(Marca, 'POR ASIGNAR')) AS "Marca del Producto",
    MAX(IFNULL(SubFamilia, 'POR ASIGNAR')) AS "Categoría Principal",
    CASE
        WHEN MAX(IFNULL(NomDetSubFam, '')) <> '' AND MAX(IFNULL(NomDetSubFam, '')) <> 'POR ASIGNAR' AND MAX(IFNULL(NomDetSubFam, '')) <> 'Sin Asignar' 
        THEN MAX(IFNULL(NomDetSubFam, ''))
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
    END AS "Subcategoría/Dimensiones",

    -- Especificaciones de empaque
    UNI_CAJ AS "Unidades por Caja",

    -- Métricas de ingresos
    SUM(CAN_ING) AS "Total Unidades Ingresadas",
    ROUND(IF(SUM(CAN_ING) = 0, 0, SUM(CIF_UNI * CAN_ING) / SUM(CAN_ING)), 2) AS "Costo Promedio CIF (USD)",
    ROUND(IF(SUM(COS_UNI * CAN_ING) = 0, 0, SUM(COS_UNI * CAN_ING) / SUM(CAN_ING)), 2) AS "Costo Promedio Unitario (CLP)",

    -- Fechas relevantes
    MIN(fec_ing) AS "Fecha Primer Ingreso",
    MAX(fec_ing) AS "Fecha Último Ingreso",
    DATEDIFF(CURDATE(), MIN(fec_ing)) AS "Días Desde Primer Ingreso",

    -- Conteo de registros
    COUNT(*) AS "Cantidad de Ingresos",

    -- Metadatos en JSON
    CONCAT('[', GROUP_CONCAT(DISTINCT CONCAT(
        '{\'Zeta\':\'', ZET_ART, '\',',
        '\'Año Producción\':\'', ANIO_PRO, '\',',
        '\'Unidades Ingresadas\':', CAN_ING, ',',
        '\'Fecha Ingreso\':\'', DATE_FORMAT(fec_ing, '%Y-%m-%d'), '\'}'
    ) SEPARATOR ','), ']') AS "Historial de Ingresos (JSON)"
FROM saldos s
WHERE
    CAST(ANIO_PRO AS SIGNED) = ?
    AND (? = '' OR COD_ART LIKE CONCAT('%', ?, '%'))
GROUP BY
    COD_ART,
    UNI_CAJ
ORDER BY "Código de Producto"
`
}
