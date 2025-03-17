package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pablojnd/rotacion/config"
	"github.com/pablojnd/rotacion/db"
	"github.com/pablojnd/rotacion/server"
)

func main() {
	// Configurar logger
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Iniciando aplicación de rotación...")

	// Asegurar que la carpeta static/docs existe
	ensureStaticDirs()

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Error al cargar configuración: %v", err)
	}
	log.Println("✅ Configuración cargada correctamente")

	// Inicializar conexiones a bases de datos
	log.Println("📊 Conectando a SQL Server...")
	sqlServer, err := db.NewSQLServerConnection(cfg)
	if err != nil {
		log.Fatalf("❌ Error al conectar con SQL Server: %v", err)
	}
	log.Printf("✅ Conexión a SQL Server establecida: %s@%s:%d/%s\n",
		cfg.SQLServerUser, cfg.SQLServerHost, cfg.SQLServerPort, cfg.SQLServerDatabase)
	defer sqlServer.Close()

	log.Println("📊 Conectando a MySQL...")
	mysql, err := db.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("❌ Error al conectar con MySQL: %v", err)
	}
	log.Printf("✅ Conexión a MySQL establecida: %s@%s:%d/%s\n",
		cfg.MySQLUser, cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDatabase)
	defer mysql.Close()

	// Inicializar el servidor
	srv := server.New(cfg, sqlServer, mysql)

	// Manejar señales para cerrar gracefully
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Iniciar el servidor en una goroutine
	go func() {
		log.Printf("🚀 Servidor iniciado en http://localhost:%s\n", cfg.ServerPort)
		log.Printf("📚 Documentación disponible en http://localhost:%s/docs\n", cfg.ServerPort)
		if err := srv.Start(); err != nil {
			log.Fatalf("❌ Error en el servidor: %v", err)
		}
	}()

	// Esperar señal de interrupción
	<-stop
	log.Println("⚠️ Recibida señal de cierre. Finalizando aplicación...")

	log.Println("👋 Aplicación cerrada correctamente")
}

// ensureStaticDirs crea las carpetas necesarias para archivos estáticos
func ensureStaticDirs() {
	dirs := []string{
		"./static",
		"./static/docs",
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Printf("Creando directorio: %s", dir)
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("Error al crear directorio %s: %v", dir, err)
			}
		}
	}
}
