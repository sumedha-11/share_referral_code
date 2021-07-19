package user

import (
	"github.com/jinzhu/gorm"
	"github.com/sumedha-11/share_referral_code/internal/code"

	"github.com/sumedha-11/share_referral_code/pkg/sql"
)

type User struct {
	gorm.Model
	Name   string `json:"name"`
	Email  string `json:"email" gorm:"unique;not null"`
	Active bool
	Codes  []code.Code `gorm:"foreignkey:UserID"`
}

type GUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func (u *User) Get() (err error) {
	return sql.DB.Model(u).First(u).Error
}

func (u *User) GetByEmail() (err error) {
	return sql.DB.Model(u).Where("email = ?", u.Email).First(&u).Error
}

func (u *User) Create() (err error) {
	err = sql.DB.Model(u).Create(u).Error
	return
}
