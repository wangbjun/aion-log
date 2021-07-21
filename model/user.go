package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Salt     string
}

func (u User) IsEmailExisted(email string) (bool, error) {
	var user User
	err := DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u User) GetUser(id int) (User, error) {
	var user User
	err := DB().Find(&user, "id = ?", id).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
