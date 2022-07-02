package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

func Register(param models.ParamRegister) (err error) {

	err = mysql.IsAvailableUsername(param.Username)
	if err != nil {
		return err
	}

	user := models.User{
		Uid:      snowflake.GenerateID(),
		Username: param.Username,
		Password: param.Password,
	}

	return mysql.InsertUser(&user)
}
