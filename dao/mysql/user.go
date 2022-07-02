package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

func IsAvailableUsername(username string) error {

	sqlStr := "select count(*) from user where username = ?"

	var count int
	err := db.Get(&count, sqlStr, username)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user already exist")
	}
	return nil
}

func InsertUser(user *models.User) (err error) {
	user.Password = encryptPassword(user.Password)

	sqlStr := "insert into user(uid, username, password) values(?, ?, ?)"
	_, err = db.Exec(sqlStr, user.Uid, user.Username, user.Password)
	return
}

func encryptPassword(password string) string {
	hash := md5.New()
	hash.Write([]byte("bluebell"))
	return hex.EncodeToString(hash.Sum([]byte(password)))
}
