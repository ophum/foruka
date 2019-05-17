package networkController

import (
	"github.com/gin-gonic/gin"
	auth "github.com/ophum/foruka/controllers/authController"

	contmodel "github.com/ophum/foruka/models/containerModel"
)

func Index(c *gin.Context) {
	auth.Auth(c)
	c.HTML(200, "networks/index.tmpl", gin.H{})
}

func Create(c *gin.Context) {
	user := auth.Auth(c)

	containers := contmodel.GetContainers(user.ID)

	c.HTML(200, "networks/create.tmpl", gin.H{
		"containers": containers,
	})
}
