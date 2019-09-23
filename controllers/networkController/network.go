package networkController

import (
	"strconv"

	"github.com/gin-gonic/gin"
	auth "github.com/ophum/foruka/controllers/authController"

	contmodel "github.com/ophum/foruka/models/containerModel"
	netmodel "github.com/ophum/foruka/models/networkModel"
)

func Index(c *gin.Context) {
	user := auth.Auth(c)

	endpoints := netmodel.GetEndPoints(user.ID)

	c.HTML(200, "networks/index.tmpl", gin.H{
		"EndPoints": endpoints,
	})
}

func Create(c *gin.Context) {
	user := auth.Auth(c)

	containers := contmodel.GetContainers(user.ID)

	c.HTML(200, "networks/create.tmpl", gin.H{
		"containers": containers,
	})
}

func Store(c *gin.Context) {
	user := auth.Auth(c)
	name := c.PostForm("name")
	cont := c.PostForm("container")
	port, _ := strconv.Atoi(c.PostForm("endpoint"))

	netmodel.Create(user.ID, name, cont, uint(port))
	c.Redirect(302, "/networks/")
}

func Show(c *gin.Context) {
	user := auth.Auth(c)

	var hash string
	endpoint := netmodel.GetEndPoint(user.ID, hash)
	c.HTML(200, "networks/show.tmpl", gin.H{
		"EndPoint": endpoint,
	})
}
