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

// Load carga la configuración desde el archivo .env
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	sqlServerPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	mysqlPort, _ := strconv.Atoi(os.Getenv("MYSQL_PORT"))

	cfg := &Config{
		// SQL Server
		SQLServerHost:     os.Getenv("DB_SERVER"),
		SQLServerUser:     os.Getenv("DB_USER"),
		SQLServerPassword: os.Getenv("DB_PASSWORD"),
		SQLServerDatabase: os.Getenv("DB_NAME"),
		SQLServerPort:     sqlServerPort,

		// MySQL
		MySQLHost:     os.Getenv("MYSQL_HOST"),
		MySQLUser:     os.Getenv("MYSQL_USER"),
		MySQLPassword: os.Getenv("MYSQL_PASSWORD"),
		MySQLDatabase: os.Getenv("MYSQL_DATABASE"),
		MySQLPort:     mysqlPort,

		// Server
		ServerPort: os.Getenv("SERVER_PORT"),
	}

	return cfg, nil
}
