package db

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pablojnd/rotacion/config"
)

// SQLServerDB representa una conexión a SQL Server
type SQLServerDB struct {
	*sql.DB
}

// NewSQLServerConnection crea una nueva conexión a SQL Server
func NewSQLServerConnection(cfg *config.Config) (*SQLServerDB, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		cfg.SQLServerHost, cfg.SQLServerUser, cfg.SQLServerPassword, cfg.SQLServerPort, cfg.SQLServerDatabase)

	db, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &SQLServerDB{db}, nil
}

// ExecuteQuery ejecuta una consulta SQL y devuelve los resultados
func (db *SQLServerDB) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

// ExecuteNonQuery ejecuta una consulta SQL que no devuelve filas
func (db *SQLServerDB) ExecuteNonQuery(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}
