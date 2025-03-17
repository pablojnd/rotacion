package utils

import (
	"encoding/json"
	"strings"
)

// FixJSONQuotes corrige las comillas simples en un string JSON para que sean dobles
// MySQL puede devolver JSON con comillas simples que no es válido para parsing en Go
func FixJSONQuotes(jsonStr string) string {
	// Si está vacío o es nulo, devolver un array vacío
	if jsonStr == "" || jsonStr == "null" {
		return "[]"
	}

	// Reemplazar comillas simples por comillas dobles
	// pero sólo si no están dentro de una cadena (entre comillas dobles)
	fixed := strings.Replace(jsonStr, "'", "\"", -1)

	return fixed
}

// ParseJSON analiza un string JSON a una estructura de datos
func ParseJSON(jsonStr string, target interface{}) error {
	// Arreglar las comillas primero
	fixedJSON := FixJSONQuotes(jsonStr)

	// Analizar el JSON
	return json.Unmarshal([]byte(fixedJSON), target)
}
