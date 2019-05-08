package database

import (
	"database/sql"
	"fmt"

	//"github.com/jmoiron/sqlx"

	//_ "github.com/go-sql-driver/mysql"

	_ "github.com/mattn/go-sqlite3"
)

const sqlTable = `
    CREATE TABLE IF NOT EXISTS userinfo(
  	id INTEGER PRIMARY KEY AUTOINCREMENT,
  	meeting_id INT NOT NULL,
  	enter_status TINYINT NOT NULL,
  	enter_channel TINYINT NOT NULL,
  	is_import TINYINT NOT NULL ,
  	sign_status TINYINT NOT NULL ,
  	sign_type TINYINT NOT NULL ,
  	user_name VARCHAR(60) NOT NULL ,
  	picture VARCHAR(255) NOT NULL ,
  	telephone VARCHAR(20) NOT NULL ,
  	create_time DATE  NULL,
  	update_time DATE  NULL
	);`

// create the connection
func NewMySQLConn(config *MySQLConfig) (*sql.DB, error) {
	// connect to the database
	config = mergeConfig(config)
	//par := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s", config.Username, config.Password, config.Host, config.Port, config.Database)
	//db, err := sqlx.Open("mysql", par) // 第一个参数为驱动名

	db, err  := sql.Open("sqlite3", "./bin/foo.db")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to log mysql: %s", err)
	}
	// ping the mysql


	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %s", err)
	}
	// set db

	// reuse the connection forever(Expired connections may be closed lazily before reuse)
	// If d <= 0, connections are reused forever.
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	// db.SetConnMaxLifetime(10*time.Second)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	// If n <= 0, no idle connections are retained.
	db.SetMaxIdleConns(config.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than MaxIdleConns, then MaxIdleConns will be reduced to match the new MaxOpenConns limit
	// If n <= 0, then there is no limit on the number of open connections. The default is 0 (unlimited).
	db.SetMaxOpenConns(config.MaxOpenConns)

	return db, nil
}
