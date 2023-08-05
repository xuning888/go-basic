package domain

import "golang.org/x/crypto/bcrypt"

// User 领域对象
type User struct {
	Id           int64
	Email        string
	Password     string
	NickName     string
	Birthday     string
	Introduction string
}

func (u User) Compare(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// TODO log
		return false
	}
	return true
}
