package network

import (
	"fmt"

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
	Ipv6Address string `json:"ipv6.address"`
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

type NetworkDeleteRequest struct {
	Name string `json:"name"`
}

func (a *NetworkAPI) Delete(c *gin.Context) {
	ndr := NetworkDeleteRequest{}
	c.BindJSON(&ndr)

	err := a.foruka.DeleteNetwork(ndr.Name)
	if err != nil {
		c.JSON(200, err)
	} else {
		c.JSON(200, "success")
	}
}

type NetworkConfigProxyRequest struct {
	RouterName         string `json:"router_name"`
	EndpointPort       string `json:"endpoint_port"`
	DestinationPort    string `json:"destination_port"`
	DestinationAddress string `json:"destination_address"`
}

func (a *NetworkAPI) ConfigProxy(c *gin.Context) {
	ncp := NetworkConfigProxyRequest{}
	c.BindJSON(&ncp)

	router_name := ncp.RouterName
	endpoint_port := ncp.EndpointPort
	dport := ncp.DestinationPort
	daddr := ncp.DestinationAddress

	err := a.foruka.ExecContainer(
		router_name,
		[]string{
			"iptables", "-t", "nat",
			"-A", "PREROUTING",
			"-i", "eth0",
			"-p", "tcp",
			"--dport", endpoint_port,
			"-j", "DNAT",
			"--to-destination", fmt.Sprintf("%s:%s", daddr, dport),
		},
	)

	c.JSON(200, err)
}

type NetworkConfigMasqueradeRequest struct {
	RouterName string `json:"router_name"`
}

func (a *NetworkAPI) ConfigMasquerade(c *gin.Context) {
	cmr := NetworkConfigMasqueradeRequest{}
	c.BindJSON(&cmr)

	err := a.foruka.ExecContainer(
		cmr.RouterName,
		[]string{
			"iptables", "-t", "nat",
			"-A", "POSTROUTING",
			"-j", "MASQUERADE",
		},
	)

	c.JSON(200, err)
}
