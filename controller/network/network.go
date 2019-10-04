package network

import (
	"encoding/json"

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

func (a *NetworkAPI) Create(c *gin.Context) {
	name := c.PostForm("name")
	config_json := c.PostForm("config")

	config := map[string]string{}
	_ = json.Unmarshal([]byte(config_json), &config)
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
