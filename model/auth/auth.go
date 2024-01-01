package auth

import (
	"github.com/kazukiyo17/synergy_api_server/model"
	"gorm.io/gorm"
)

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CheckAuth checks if authentication information exists
func CheckAuth(username, password string) (bool, error) {
	var auth Auth
	err := model.DB.Select("id").Where(Auth{Username: username, Password: password}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if auth.ID > 0 {
		return true, nil
	}

	return false, nil
}

// CheckUsername checks if username exists
func CheckUsername(username string) (bool, error) {
	var auth Auth
	err := model.DB.Select("id").Where(Auth{Username: username}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if auth.ID > 0 {
		return true, nil
	}

	return false, nil
}

// AddAuth add a user
func AddAuth(username, password string) error {
	auth := Auth{
		Username: username,
		Password: password,
	}
	if err := model.DB.Create(&auth).Error; err != nil {
		return err
	}

	return nil
}
