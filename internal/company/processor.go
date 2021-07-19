package company

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sumedha-11/share_referral_code/internal/code"
	"github.com/sumedha-11/share_referral_code/internal/login"
	"github.com/sumedha-11/share_referral_code/internal/user"
	"github.com/sumedha-11/share_referral_code/pkg/logger"
	"net/http"
	"strconv"
)

func ViewCompany(g *gin.Context) {
	var id int
	var err error
	ids := g.Param("id")
	id, err = strconv.Atoi(ids)
	if err != nil {
		logger.ErrorLogger.Printf("error in converting param id to integer, ids:%v err:%v", ids, err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error."})
		return
	}
	c := &Company{}
	c.ID = uint(id)
	err = c.Get()
	if err != nil {
		logger.ErrorLogger.Printf("error in company Get method, cId:%d err:%v", id, err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error."})
		return
	}

	isLoggedIn := false
	session := sessions.Default(g)
	email := session.Get("user-id")
	if email != nil {
		isLoggedIn = true
	}
	link := login.HandleGoogleLogin(g)
	//TODO: hide form post url which include company id
	g.HTML(http.StatusOK, "companyView.html", gin.H{"Company": c, "Codes": c.Codes, "LoggedIn": isLoggedIn, "Link": link})
}

func ViewAllCompany(g *gin.Context) {
	c := Company{}
	cs, err := c.GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("error in getting all company get method, err:%v", err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error."})
		return
	}
	g.HTML(http.StatusOK, "companyAll.tmpl", gin.H{"List": cs})
}

func CreateCompanyForm(g *gin.Context) {
	g.HTML(http.StatusOK, "companyAdd.tmpl", gin.H{})
}

func UpdateCompany(g *gin.Context) {
	var id int
	var err error
	ids := g.Param("id")
	id, err = strconv.Atoi(ids)
	if err != nil {
		logger.ErrorLogger.Printf("error in converting parmas ids to integer, Id:%d, err:%v", id, err)
		g.JSON(http.StatusInternalServerError, gin.H{
			"data":     nil,
			"debugMsg": "Error converting id to integer",
		})
		return
	}
	c := &Company{}
	err = g.ShouldBindJSON(c)
	if err != nil {
		logger.ErrorLogger.Printf("invalid json provided, Id:%d, err:%v", id, err)
		g.JSON(http.StatusUnprocessableEntity, gin.H{
			"data":     nil,
			"debugMsg": "Invalid json provided",
		})
		return
	}
	c.ID = uint(id)
	err = c.Update()
	if err != nil {
		logger.ErrorLogger.Printf("error in updating company, Id:%d, err:%v", id, err)
		g.JSON(http.StatusInternalServerError, gin.H{
			"data":     nil,
			"debugMsg": "Error in updating company",
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"data":     c,
		"debugMsg": "",
	})
}

func CreateCompany(g *gin.Context) {

	c := &Company{
		Name:     g.PostForm("Name"),
		Details:  g.PostForm("Details"),
		Signup:   g.PostForm("Signup"),
		Referral: g.PostForm("Referral"),
	}

	err := c.Create()
	if err != nil {
		logger.ErrorLogger.Printf("error in creating company, err:%v", err)
		g.JSON(http.StatusInternalServerError, gin.H{
			"data":     nil,
			"debugMsg": "Error in creating company",
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"data":     c,
		"debugMsg": "",
	})
	return
}

//TODO: find a better place for this function
//TODO: Single code for unique(userid, companyid)
func AddCodes(g *gin.Context) {
	var id int
	var err error
	ids := g.Param("id")
	id, err = strconv.Atoi(ids)
	if err != nil {
		logger.ErrorLogger.Printf("error in converting parmas ids to integer, Id:%d, err:%v", id, err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error."})
		return
	}
	c := &Company{}
	c.ID = uint(id)
	session := sessions.Default(g)
	email := session.Get("user-id")
	u := user.User{}
	u.Email = fmt.Sprintf("%v", email)
	err = u.GetByEmail()
	if err != nil {
		logger.ErrorLogger.Printf("error in finding user from user-email:%s, err:%v", u.Email, err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error, User not valid, Please login again"})
		return
	}
	// TODO: verify if company exits
	code2 := code.Code{
		Link:      g.PostForm("Link"),
		Text:      g.PostForm("Text"),
		UserID:    u.ID,
		CompanyID: c.ID,
	}
	err = c.Get()
	if err != nil {
		logger.ErrorLogger.Printf("error in finding company from id:%d, err:%v", c.ID, err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Company not found with this Id"})
		return
	}
	cnt, err := code2.GetFromWhere()

	if err != nil || cnt != 0 {
		logger.ErrorLogger.Printf("err:%v in find code , count of code2:%d per company", err, cnt)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Only one code per userId can be created "})
		return
	}
	err = code2.Create()
	if err != nil {
		logger.ErrorLogger.Printf("error in creating code:%v, err:%v", code2, err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error."})
		return
	}
	g.Redirect(http.StatusFound, "/company/view/"+strconv.Itoa(int(c.ID)))
}
