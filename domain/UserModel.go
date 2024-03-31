package domain

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html"
	"sensor_iot/utils"
)

type UserModel struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Username string `gorm:"size:255;not null;unique" json:"username" validate:"nonzero"`
	Password string `gorm:"size:255;not null" validate:"nonzero"`
	Role     string `gorm:"column:rule" json:"role"`
}

func (user *UserModel) Save() {
	utils.MyDataBase.Create(user)
}

func (user *UserModel) BeforeSave(tx *gorm.DB) error {
	hashPassWord, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassWord)
	user.Username = html.EscapeString(user.Username)
	return nil
}
