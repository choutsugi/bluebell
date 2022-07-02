package mysql

import (
	"bluebell/setting"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB(config *setting.DbConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DbName,
	)

	DB, err = sqlx.Connect(config.DriveName, dsn)
	if err != nil {
		fmt.Println("connect to mysql failed, err:", err)
		return
	}

	DB.SetMaxOpenConns(200)
	DB.SetConnMaxIdleTime(10)

	return
}
