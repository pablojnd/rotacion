package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yourusername/rotacion/config"
)

// MySQLDB representa una conexión a MySQL
type MySQLDB struct {
	*sql.DB
}

// NewMySQLConnection crea una nueva conexión a MySQL
func NewMySQLConnection(cfg *config.Config) (*MySQLDB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.MySQLUser, cfg.MySQLPassword, cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDatabase)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &MySQLDB{db}, nil
}

// ExecuteQuery ejecuta una consulta SQL y devuelve los resultados
func (db *MySQLDB) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

// ExecuteNonQuery ejecuta una consulta SQL que no devuelve filas
func (db *MySQLDB) ExecuteNonQuery(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}
