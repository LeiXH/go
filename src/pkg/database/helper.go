package database

import (
	"database/sql"
	"fmt"
	"sync"

	//"github.com/jmoiron/sqlx"
)

var _mysql *sql.DB
var _mysqlOnce sync.Once

func InitMySQLConn(config *MySQLConfig) (err error) {
	_mysqlOnce.Do(func() {
		_mysql, err = NewMySQLConn(config)
		if err != nil {
			err = fmt.Errorf("mysql connection can not be established, because of %s", err)
		}
	})
	return
}

func MySQL() *sql.DB {
	return _mysql
}

func CloseConn() {
	_ = _mysql.Close()
}
