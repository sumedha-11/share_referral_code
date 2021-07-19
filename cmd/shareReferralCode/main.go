package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"

	"github.com/sumedha-11/share_referral_code/internal/company"
	"github.com/sumedha-11/share_referral_code/internal/config"
	"github.com/sumedha-11/share_referral_code/internal/login"
	"github.com/sumedha-11/share_referral_code/internal/user"
	"github.com/sumedha-11/share_referral_code/pkg/logger"
	"github.com/sumedha-11/share_referral_code/pkg/sql"
)

func main() {
	config.ReadConfig("common_lol.json")
	login.InitGoogleCred()

	err := sql.DBConn(config.Config.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	//defer sql.DB.Close()

	//sql.DB.SetLogger(logger.DBLogger)
	fmt.Println("DBCONNECTED")
	gin.DefaultWriter = logger.File
	router := gin.Default()
	token, err := login.RandToken(64)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}

	store := cookie.NewStore([]byte(token))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})

	router.Use(sessions.Sessions("refercodes", store))

	router.POST("/user/create", user.CreateUser)
	router.GET("/auth", login.AuthHandler)
	router.GET("/", company.ViewAllCompany)
	router.GET("/company/view/:id", company.ViewCompany)

	authorized := router.Group("/client")
	authorized.Use(login.AuthorizeRequest())
	{
		authorized.POST("/company/code/:id", company.AddCodes)
		authorized.GET("/user/view/:id", user.ViewUser)
		authorized.GET("/logout", login.Logout)
	}
	adminAuth := router.Group("/admin")
	adminAuth.Use(login.AdminAuthRequest())
	{
		//adminAuth.GET("/company/edit/:id", company.UpdateCompanyForm)
		adminAuth.POST("/company/edit/:id", company.UpdateCompany)
		adminAuth.GET("/company/create", company.CreateCompanyForm)
		adminAuth.POST("/company/create", company.CreateCompany)
	}

	router.LoadHTMLGlob("templates/*")
	router.Static("/static/", "./static/")

	logger.InfoLogger.Printf("msg:%v", "server starting....")

	err = router.Run(":" + config.Config.Common.ServerPort)
	if err != nil {
		log.Fatalf("error in starting server , err:%v", err)
	}
}
