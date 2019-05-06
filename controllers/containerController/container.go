package containerController

import (
	"strconv"

	"github.com/gin-gonic/gin"
	auth "github.com/ophum/foruka/controllers/authController"
	contmodel "github.com/ophum/foruka/models/containerModel"
)

func Index(c *gin.Context) {
	auth.Auth(c)

	containers := contmodel.GetContainers(1)
	c.HTML(200, "containers/index.tmpl", gin.H{
		"containers": containers,
		})
}

func Create(c *gin.Context) {
	auth.Auth(c)

	c.HTML(200, "containers/create.tmpl", gin.H{})
}

func Store(c *gin.Context) {
	auth.Auth(c)
	userId := 1
	name := c.PostForm("name")
	image := c.PostForm("image")

	contmodel.Create(userId, name, image)
	c.Redirect(301, "/containers/")
}

func Show(c *gin.Context) {
	auth.Auth(c)
	userId := 1
	contId, _ := strconv.Atoi(c.Param("id"))

	var cont contmodel.Container
	cont = contmodel.GetContainer(userId, contId)

	c.HTML(200, "containers/show.tmpl", gin.H{
		"container": cont,
	})
}
