package company

import (
	"github.com/jinzhu/gorm"
	"github.com/sumedha-11/share_referral_code/internal/code"
	"github.com/sumedha-11/share_referral_code/pkg/sql"
)

type Company struct {
	gorm.Model
	Name     string      `json:"name" gorm:"unique;not null"`
	Details  string      `json:"details"`
	Signup   string      `json:"signup"`
	Referral string      `json:"referral"`
	Codes    []code.Code `gorm:"foreignkey:UserID"`
}

func (c *Company) Get() error {
	//err := sql.DB.Model(c).First(c).Error
	//if err != nil {
	//	return err
	//}
	//return sql.DB.Model(c).Related(&c.Codes).Error
	return nil
}

func (c *Company) GetAll() ([]Company, error) {
	var cs []Company
	err := sql.DB.Model(c).Find(&cs).Error
	return cs, err
}

func (c *Company) Create() (err error) {
	err = sql.DB.Model(c).Create(c).Error
	return
}

func (c *Company) Update() (err error) {
	err = sql.DB.Model(c).Updates(c).Error
	return
}
