package main

import (
	"fmt"
	"log"

	"github.com/yourusername/rotacion/config"
	"github.com/yourusername/rotacion/db"
	"github.com/yourusername/rotacion/server"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Inicializar conexiones a bases de datos
	sqlServer, err := db.NewSQLServerConnection(cfg)
	if err != nil {
		log.Fatalf("Error al conectar con SQL Server: %v", err)
	}
	defer sqlServer.Close()

	mysql, err := db.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Error al conectar con MySQL: %v", err)
	}
	defer mysql.Close()

	// Inicializar el servidor
	srv := server.New(cfg, sqlServer, mysql)
	fmt.Printf("Servidor iniciado en puerto %s\n", cfg.ServerPort)
	if err := srv.Start(); err != nil {
		log.Fatalf("Error en el servidor: %v", err)
	}
}
