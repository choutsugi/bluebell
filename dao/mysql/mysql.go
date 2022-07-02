package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB() (err error) {
	dsn := "root:dangerous@tcp(127.0.0.1:3306)/bluebell?charset=utf8mb4&parseTime=True&loc=Local"
	//sqlx.Connect: Open and Ping
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("connect to mysql failed, err:", err)
		return
	}

	DB.SetMaxOpenConns(200)
	DB.SetConnMaxIdleTime(10)

	return
}
