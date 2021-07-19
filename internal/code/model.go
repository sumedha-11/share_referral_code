package code

import (
	"github.com/jinzhu/gorm"
	"github.com/sumedha-11/share_referral_code/pkg/sql"
)

type Code struct {
	gorm.Model
	Link      string
	Text      string
	UserID    uint
	CompanyID uint
}

func (code *Code) Get() error {
	return sql.DB.Model(code).First(code).Error
}

func (code *Code) GetAll() ([]Code, error) {
	var codes []Code
	err := sql.DB.Model(code).Find(&codes).Error
	return codes, err
}

func (code *Code) GetFromWhere() (int, error) {
	var cnt int64
	err := sql.DB.Model(code).Where("user_id =? AND company_id =?", code.UserID, code.CompanyID).Count(&cnt).Error
	return int(cnt), err
}

func (code *Code) Create() (err error) {
	err = sql.DB.Model(code).Create(code).Error
	return
}

func (code *Code) Update() (err error) {
	err = sql.DB.Model(code).Updates(code).Error
	return
}
