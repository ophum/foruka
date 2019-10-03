package main

import (
	"fmt"

	"github.com/ophum/foruka/core"
	"github.com/ophum/foruka/datastore"
	"github.com/urfave/cli"
)

func InitForuka(ctx *cli.Context) error {
	connection_type := ctx.String("connection-type")
	path := ctx.String("path")
	install_dir := ctx.String("install-dir")
	if install_dir == "" {
		return fmt.Errorf("Error: Empty `--install-dir`\n")
	}
	fmt.Printf("install-dir: %s\n", install_dir)
	if connection_type == "unix" {
		_, err := core.NewForukaUnix(path, install_dir)
		if err != nil {
			return err
		}
	} else if connection_type == "https" {
		//f = core.NewForuka()
	} else {
		return fmt.Errorf("Error: Not support connection type `%s`\n", connection_type)
	}

	fmt.Printf("%s://%s\n", connection_type, path)
	exbr := ctx.String("external-bridge")
	if exbr == "" {
		return fmt.Errorf("Error: Empty `--external-bridge`\n")
	}
	fmt.Printf("external-bridge: %s\n", exbr)

	ds := datastore.NewDataStore(install_dir)

	ds.Set("external-bridge", exbr)
	ds.Set("connection-type", connection_type)
	ds.Set("path", path)

	fmt.Println(ds.Get("external-bridge"))
	fmt.Println(ds.Get("connection-type"))
	fmt.Println(ds.Get("path"))
	return nil
}
