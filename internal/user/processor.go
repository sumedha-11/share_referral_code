package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ViewUser(g *gin.Context) {
	var id int
	var err error
	ids := g.Param("id")
	id, err = strconv.Atoi(ids)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"data":     nil,
			"debugMsg": "Error converting id to integer",
		})
		return
	}
	u := &User{}
	u.ID = uint(id)
	err = u.Get()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"data":     nil,
			"debugMsg": "Error finding user from id",
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"data":     u,
		"debugMsg": "",
	})
}

func CreateUser(g *gin.Context) {

	u := &User{}
	err := g.ShouldBindJSON(u)
	if err != nil {
		g.JSON(http.StatusUnprocessableEntity, gin.H{
			"data":     nil,
			"debugMsg": "Invalid json provided",
		})
		return
	}

	err = u.Create()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"data":     nil,
			"debugMsg": "Error in creating user",
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"data":     u,
		"debugMsg": "",
	})
	return
}
