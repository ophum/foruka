package container

import (
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

type ContainerIface map[string]string
type ContainerLimit map[string]string

type ContainerCreateRequest struct {
	Name   string         `json:"name"`
	Alias  string         `json:"alias"`
	Ifaces ContainerIface `json:"ifaces"`
	Limits ContainerLimit `json:"limits"`
}

func (a *ContainerAPI) Create(c *gin.Context) {
	ccr := ContainerCreateRequest{}
	c.BindJSON(&ccr)

	name := ccr.Name
	alias := ccr.Alias
	ifaces := ccr.Ifaces
	limits := ccr.Limits

	err := a.foruka.CreateContainer(name, alias, ifaces, limits)
	if err != nil {
		c.JSON(200, err)
	} else {
		container, _, err := a.foruka.GetContainer(name)
		if err != nil {
			c.JSON(200, err)
		} else {
			c.JSON(200, container)
		}
	}
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

type ContainerStateUpdateRequest struct {
	Name string `json:"name"`
}

func (a *ContainerAPI) Start(c *gin.Context) {
	csur := ContainerStateUpdateRequest{}
	c.BindJSON(&csur)
	name := csur.Name
	err := a.foruka.StartContainer(name)
	if err != nil {
		c.JSON(200, err)
	} else {
		container, _, err := a.foruka.GetContainer(name)
		if err != nil {
			c.JSON(200, err)
		} else {
			c.JSON(200, map[string]string{
				"status":      container.Status,
				"status_code": fmt.Sprintf("%d", container.StatusCode),
			})
		}
	}
}

func (a *ContainerAPI) Stop(c *gin.Context) {
	csur := ContainerStateUpdateRequest{}
	c.BindJSON(&csur)
	name := csur.Name
	err := a.foruka.StopContainer(name)
	if err != nil {
		c.JSON(200, err)
	} else {
		container, _, err := a.foruka.GetContainer(name)
		if err != nil {
			c.JSON(200, err)
		} else {
			c.JSON(200, map[string]string{
				"status":      container.Status,
				"status_code": fmt.Sprintf("%d", container.StatusCode),
			})
		}
	}
}

func (a *ContainerAPI) SetIP(c *gin.Context) {
	name := c.PostForm("name")
	address := c.PostForm("ipv4.address")
	prefix := c.PostForm("ipv4.prefix")
	device := c.PostForm("device")

	err := a.foruka.ExecContainer(
		name,
		[]string{"ip", "a", "add", fmt.Sprintf("%s/%s", address, prefix), "dev", device},
	)

	c.JSON(200, err)
}

func (a *ContainerAPI) GetState(c *gin.Context) {
	name := c.Param("name")

	state, err := a.foruka.GetContainerState(name)
	if err != nil {
		c.JSON(200, err)
	} else {
		c.JSON(200, state)
	}
}
