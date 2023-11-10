package library

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DB struct {
	*sql.DB
}

const (
	dataSourceNameFormat = "%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s"
	driverName           = "mysql"
)

func NewDB(conf *MysqlConfig) (db *DB, err error) {
	dsn := fmt.Sprintf(dataSourceNameFormat,
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName,
		conf.Charset,
		conf.ParseTime,
		conf.Loc,
	)
	odb, err := sql.Open(driverName, dsn)
	if err != nil {
		err = fmt.Errorf("database connection:[%s] sql open: %w", conf.ConnectionName, err)
		return
	}

	if err = odb.Ping(); err != nil {
		err = fmt.Errorf("database connection:[%s] db ping: %w", conf.ConnectionName, err)
		return
	}

	odb.SetConnMaxIdleTime(time.Duration(conf.MaxLifeTime) * time.Second)
	odb.SetMaxOpenConns(conf.MaxOpenConn)
	odb.SetMaxIdleConns(conf.MaxIdleConn)

	db = &DB{
		odb,
	}
	return
}
