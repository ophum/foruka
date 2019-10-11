package core

import (
	"fmt"
	//lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

func (f *Foruka) CreateNetwork(name string, config map[string]string) error {

	if name == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	}

	req := api.NetworksPost{
		NetworkPut: api.NetworkPut{
			Config: config,
		},
		Name: name,
	}

	err := f.server.CreateNetwork(req)
	if err != nil {
		return err
	}
	return nil
}

func (f *Foruka) DeleteNetwork(name string) error {
	if name == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	}
	err := f.server.DeleteNetwork(name)
	return err
}

func (f *Foruka) GetNetworks() ([]api.Network, error) {
	networks, err := f.server.GetNetworks()
	return networks, err
}

func (f *Foruka) GetNetwork(name string) (*api.Network, string, error) {
	network, etag, err := f.server.GetNetwork(name)
	return network, etag, err
}
