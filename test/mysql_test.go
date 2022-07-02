package test

import (
	"bluebell/dao/mysql"
	"bluebell/setting"
	"testing"
)

func TestSqlx(t *testing.T) {

	if err := setting.Init(); err != nil {
		t.Error(err)
	}

	if err := mysql.InitDB(setting.Conf.Db); err != nil {
		t.Error(err)
	}
}
