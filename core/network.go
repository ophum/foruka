package core

import (
	"fmt"
	//lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

func (f *Foruka) CreateNetwork(name string) error {

	if name == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`")
	}

	req := api.NetworksPost{
		NetworkPut: api.NetworkPut{
			Config: map[string]string{
				"ipv4.address": "none",
				"ipv4.nat": "false",
				"ipv6.address": "none",
				"ipv6.nat": "false",
			},
		},
		Name: name,
	}

	err := f.server.CreateNetwork(req)
	if err != nil {
		return err
	}
	return nil
}

