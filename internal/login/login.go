package login

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sumedha-11/share_referral_code/internal/config"
	"github.com/sumedha-11/share_referral_code/internal/user"
	"github.com/sumedha-11/share_referral_code/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"math/rand"
	"net/http"
)

var cred config.Credentials
var conf *oauth2.Config

// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func InitGoogleCred() {
	cred = config.Credentials{
		Cid:     config.Config.GoogleCred.Cid,
		Csecret: config.Config.GoogleCred.Csecret,
	}

	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  config.Config.GoogleCred.Redirect,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

func HandleGoogleLogin(g *gin.Context) string {
	state, err := RandToken(32)
	if err != nil {
		logger.ErrorLogger.Printf("error in getting randToken, err:%v", err)
		return ""
	}
	session := sessions.Default(g)
	session.Set("state", state)
	err = session.Save()
	if err != nil {
		logger.ErrorLogger.Printf("error in saving session, err:%v", err)
		return ""
	}
	link := getLoginURL(state)
	return link
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func AuthHandler(g *gin.Context) {

	session := sessions.Default(g)
	retrievedState := session.Get("state")
	queryState := g.Request.URL.Query().Get("state")
	if retrievedState != queryState {
		logger.ErrorLogger.Printf("invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		g.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Invalid session state."})
		return
	}
	code := g.Request.URL.Query().Get("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		logger.ErrorLogger.Printf("error in token auth, error:%v", err)
		g.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Login failed. Please try again."})
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		logger.ErrorLogger.Printf("error in getting userinfo from google, error:%v", err)
		g.AbortWithStatus(http.StatusBadRequest)
		return
	}

	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	guser := user.GUser{}
	if err = json.Unmarshal(data, &guser); err != nil {
		logger.ErrorLogger.Printf("error in unmarshall, error:%v", err)
		g.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try again."})
		return
	}
	session.Set("user-id", guser.Email)
	err = session.Save()
	if err != nil {
		logger.ErrorLogger.Printf("error in saving from session, error:%v", err)
		g.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
		return
	}

	//TODO: search using force indexing (email) to reduce time complexity
	dbUser := user.User{}
	dbUser.Email = guser.Email

	err = dbUser.GetByEmail()
	if err != nil {
		if err.Error() == "record not found" {
			dbUser.Name = guser.Name
			dbUser.Active = true
			err = dbUser.Create()
			if err != nil {
				logger.ErrorLogger.Printf("error in creating new user, error:%v", err)
				g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{})
				return
			}

		} else {
			logger.ErrorLogger.Printf("error in db, error:%v", err)
			g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{})
			return
		}
	}
	if dbUser.Active == false {
		logger.ErrorLogger.Printf("user is deactivated, id:%d", dbUser.ID)
		g.HTML(http.StatusLocked, "error.tmpl", gin.H{"message": "user is deactivated"})
		return
	}
	g.Redirect(http.StatusFound, "/")
}

func AuthorizeRequest() gin.HandlerFunc {
	return func(g *gin.Context) {
		session := sessions.Default(g)
		v := session.Get("user-id")
		if v == nil {
			logger.ErrorLogger.Printf("user is not valid, email:%s", v)
			g.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Please login."})
			g.Abort()
		}
		g.Next()
	}
}

func AdminAuthRequest() gin.HandlerFunc {
	return func(g *gin.Context) {
		session := sessions.Default(g)
		v := session.Get("user-id")
		if v == nil || (v != "rajsumedha11@gmail.com") {
			logger.ErrorLogger.Printf("user is not admin, email:%s", v)
			g.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Admin login required."})
			g.Abort()
		}
		g.Next()
	}
}

func Logout(g *gin.Context) {
	session := sessions.Default(g)
	session.Clear()
	session.Save()
	g.Redirect(http.StatusFound, "/")
}
