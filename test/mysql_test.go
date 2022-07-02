package test

import (
	"bluebell/dao/mysql"
	"testing"
)

func TestSqlx(t *testing.T) {
	if err := mysql.InitDB(); err != nil {
		t.Error(err)
	}
}
