package encrypt

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

func Encrypt(origin string) (hashed, salt []byte, err error) {
	salt = make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}
	salted := append([]byte(origin), salt...)
	hashed, err = bcrypt.GenerateFromPassword(salted, bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}
	return hashed, salt, nil
}

func Verify(origin string, hashed, salt []byte) bool {
	salted := append([]byte(origin), salt...)
	if err := bcrypt.CompareHashAndPassword(hashed, salted); err != nil {
		return false
	}
	return true
}
