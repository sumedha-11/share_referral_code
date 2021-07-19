package main

import (
	"github.com/sumedha-11/share_referral_code/internal/code"
	"github.com/sumedha-11/share_referral_code/internal/company"
	"github.com/sumedha-11/share_referral_code/internal/config"
	"github.com/sumedha-11/share_referral_code/internal/user"
	"github.com/sumedha-11/share_referral_code/pkg/sql"
	"log"
)

func main() {
	err := sql.DBConn(config.Config.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer sql.DB.Close()
	sql.DB.AutoMigrate(&user.User{}, &company.Company{}, &code.Code{})
}