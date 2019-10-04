package main

import (
	"encoding/json"
	"fmt"

	"github.com/ophum/foruka/datastore"

	"github.com/ophum/foruka/core"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

func ServeForuka(ctx *cli.Context) error {
	install_dir := ctx.String("install-dir")
	if install_dir == "" {
		return fmt.Errorf("Error: Empty install-dir\n")
	}
	ds := datastore.NewDataStore(install_dir)
	fr, _ := core.NewForukaUnix(ds.Get("path"), install_dir)
	r := gin.Default()
	r.GET("/external-bridge", func(c *gin.Context) {
		c.JSON(200, ds.Get("external-bridge"))
	})
	r.GET("/create-network/:name", func(c *gin.Context) {
		nn := c.Param("name")
		id, err := fr.CreateNetwork(nn)
		if err != nil {
			c.JSON(200, err)
		} else {
			network := fr.GetNetwork(id)

			c.JSON(200, network)
		}
	})

	r.GET("/delete-network/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := fr.DeleteNetwork(id)
		if err != nil {
			c.JSON(200, err)
		} else {
			c.JSON(200, "success")
		}
	})
	r.GET("/create-router/:name/:netid", func(c *gin.Context) {
		name := c.Param("name")
		netid := c.Param("netid")

		s := ds.Get(fmt.Sprintf("networks/bridge/%s", netid))
		net := &core.Network{}
		_ = json.Unmarshal([]byte(s), &net)

		id, _ := fr.CreateRouter(name, map[string]string{
			"eth0": ds.Get("external-bridge"),
			"eth1": net.Name,
		})

		c.JSON(200, id)
	})

	r.GET("/routers", func(c *gin.Context) {
		rts := fr.GetRouters()
		c.JSON(200, rts)
	})

	r.GET("/routers/:id/status", func(c *gin.Context) {
		status := fr.GetRouterStatus(c.Param("id"))
		c.JSON(200, status)
	})

	r.GET("/routers/:id/start", func(c *gin.Context) {
		id := c.Param("id")
		fr.StartRouter(id)
		status := fr.GetRouterStatus(id)
		c.JSON(200, status)
	})
	r.GET("/routers/:id/stop", func(c *gin.Context) {
		id := c.Param("id")
		fr.StopRouter(id)
		status := fr.GetRouterStatus(id)
		c.JSON(200, status)
	})
	r.GET("/networks", func(c *gin.Context) {
		nws := fr.GetNetworks()
		c.JSON(200, nws)
	})

	r.GET("/networks/:id", func(c *gin.Context) {
		network := fr.GetNetwork(c.Param("id"))
		c.JSON(200, network)
	})
	r.Run()

	return nil
}
