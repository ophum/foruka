package containerController

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lxc/lxd/shared/api"
	auth "github.com/ophum/foruka/controllers/authController"
	contmodel "github.com/ophum/foruka/models/containerModel"
)

func Index(c *gin.Context) {
	user := auth.Auth(c)

	containers := contmodel.GetContainers(user.ID)
	c.HTML(200, "containers/index.tmpl", gin.H{
		"containers": containers,
	})
}

func Create(c *gin.Context) {
	auth.Auth(c)

	c.HTML(200, "containers/create.tmpl", gin.H{})
}

func Store(c *gin.Context) {
	user := auth.Auth(c)
	name := c.PostForm("name")
	image := c.PostForm("image")

	contmodel.Create(user.ID, name, image)
	c.Redirect(302, "/containers/")
}

func Show(c *gin.Context) {
	user := auth.Auth(c)
	hashId := c.Param("id")

	var cont contmodel.Container
	cont = contmodel.GetContainer(user.ID, hashId)
	status, _ := contmodel.Status(cont.Name)
	var addresses = []api.ContainerStateNetworkAddress{}
	if status.Status == "Running" {
		addresses = status.Network["eth0"].Addresses
	}

	c.HTML(200, "containers/show.tmpl", gin.H{
		"container": cont,
		"status":    status,
		"addresses": addresses,
	})
}

func Start(c *gin.Context) {
	user := auth.Auth(c)

	hashId := c.Param("id")

	cont := contmodel.GetContainer(user.ID, hashId)
	err := contmodel.LaunchContainer(cont.Name)
	if err != nil {
		fmt.Println("err: ", err)
	}

	c.Redirect(302, "/containers/show/"+hashId)
}

func Stop(c *gin.Context) {
	user := auth.Auth(c)
	hashId := c.Param("id")

	cont := contmodel.GetContainer(user.ID, hashId)
	err := contmodel.StopContainer(cont.Name)
	if err != nil {
		fmt.Println("err: ", err)
	}

	c.Redirect(302, "/containers/show/"+hashId)
}

func Delete(c *gin.Context) {
	user := auth.Auth(c)
	hashId := c.Param("id")

	cont := contmodel.GetContainer(user.ID, hashId)
	err := contmodel.DeleteContainer(cont.Name)
	if err != nil {
		fmt.Println("err: ", err)
	}

	contmodel.Delete(hashId)

	c.Redirect(302, "/containers/")
}
