package main

import (
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
	r.POST("/create-router", func(c *gin.Context) {

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
