package container

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ophum/foruka/core"
)

type ContainerAPI struct {
	foruka *core.Foruka
}

func NewContainerAPI(f *core.Foruka) *ContainerAPI {
	api := &ContainerAPI{
		foruka: f,
	}
	return api
}

func (a *ContainerAPI) List(c *gin.Context) {
	containers, err := a.foruka.GetContainers()
	if err != nil {
		c.JSON(200, err)
	} else {
		c.JSON(200, containers)
	}

}
func (a *ContainerAPI) Create(c *gin.Context) {
	name := c.PostForm("name")
	alias := c.PostForm("alias")
	ifaces_json := c.PostForm("ifaces")
	limits_json := c.PostForm("limits")

	ifaces := map[string]string{}
	_ = json.Unmarshal([]byte(ifaces_json), &ifaces)

	limits := map[string]string{}
	_ = json.Unmarshal([]byte(limits_json), &limits)

	err := a.foruka.CreateContainer(name, alias, ifaces, limits)
	c.JSON(200, err)
}

func (a *ContainerAPI) Get(c *gin.Context) {
	name := c.Param("name")

	fmt.Println("container name: ", name)
	container, _, err := a.foruka.GetContainer(name)
	if err != nil {
		fmt.Println("error!!!")
		c.JSON(200, err)
	} else {
		c.JSON(200, container)
	}
}

func (a *ContainerAPI) ListNames(c *gin.Context) {
	names, err := a.foruka.GetContainerNames()
	if err != nil {
		c.JSON(200, err)
	} else {
		c.JSON(200, names)
	}
}
