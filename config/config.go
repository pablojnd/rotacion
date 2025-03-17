package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contiene todas las configuraciones de la aplicación
type Config struct {
	// SQL Server
	SQLServerHost     string
	SQLServerUser     string
	SQLServerPassword string
	SQLServerDatabase string
	SQLServerPort     int

	// MySQL
	MySQLHost     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string
	MySQLPort     int

	// Server
	ServerPort string
}

// Load carga la configuración desde el archivo .env o variables de entorno
func Load() (*Config, error) {
	// Intentar cargar .env, pero no fallar si no existe
	_ = godotenv.Load()

	// Convertir puertos a enteros
	sqlServerPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	if sqlServerPort == 0 {
		sqlServerPort = 1433 // Puerto por defecto de SQL Server
	}

	mysqlPort, _ := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if mysqlPort == 0 {
		mysqlPort = 3306 // Puerto por defecto de MySQL
	}

	// Obtener puerto del servidor
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080" // Puerto por defecto
	}

	cfg := &Config{
		// SQL Server
		SQLServerHost:     getEnv("DB_SERVER", "localhost"),
		SQLServerUser:     getEnv("DB_USER", "sa"),
		SQLServerPassword: getEnv("DB_PASSWORD", ""),
		SQLServerDatabase: getEnv("DB_NAME", ""),
		SQLServerPort:     sqlServerPort,

		// MySQL
		MySQLHost:     getEnv("MYSQL_HOST", "localhost"),
		MySQLUser:     getEnv("MYSQL_USER", "root"),
		MySQLPassword: getEnv("MYSQL_PASSWORD", ""),
		MySQLDatabase: getEnv("MYSQL_DATABASE", ""),
		MySQLPort:     mysqlPort,

		// Server
		ServerPort: serverPort,
	}

	return cfg, nil
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
