package network

import (
	"github.com/gin-gonic/gin"
	"github.com/ophum/foruka/core"
)

type NetworkAPI struct {
	foruka *core.Foruka
}

func NewNetworkAPI(f *core.Foruka) *NetworkAPI {
	napi := &NetworkAPI{
		foruka: f,
	}

	return napi
}

func (a *NetworkAPI) List(c *gin.Context) {
	networks, err := a.foruka.GetNetworks()
	if err != nil {
		c.JSON(200, err)
	} else {
		c.JSON(200, networks)
	}
}

func (a *NetworkAPI) Get(c *gin.Context) {
	name := c.Param("name")
	network, _, err := a.foruka.GetNetwork(name)
	if err != nil {
		c.JSON(200, err)
	} else {
		c.JSON(200, network)
	}
}

type NetworkConfig struct {
	Ipv4Address string `json:"ipv4.address"`
	Ipv6Address string `json:"ipv4.address"`
}
type NetworkCreateRequest struct {
	Name   string        `json:"name"`
	Config NetworkConfig `json:"config"`
}

func (a *NetworkAPI) Create(c *gin.Context) {
	ncr := NetworkCreateRequest{}
	c.BindJSON(&ncr)
	name := ncr.Name
	config := map[string]string{
		"ipv4.address": ncr.Config.Ipv4Address,
		"ipv6.address": ncr.Config.Ipv6Address,
	}
	err := a.foruka.CreateNetwork(name, config)
	if err != nil {
		c.JSON(200, err)
	} else {
		c.JSON(200, "success")
	}
}

func (a *NetworkAPI) Delete(c *gin.Context) {
	name := c.PostForm("name")

	err := a.foruka.DeleteNetwork(name)
	if err != nil {
		c.JSON(200, err)
	} else {
		c.JSON(200, "success")
	}
}
