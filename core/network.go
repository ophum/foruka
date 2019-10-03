package core

import (
	"encoding/json"
	"fmt"

	//lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/rs/xid"
)

type Network struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (f *Foruka) CreateNetwork(name string) (string, error) {

	if name == "" {
		return "", fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	}

	req := api.NetworksPost{
		NetworkPut: api.NetworkPut{
			Config: map[string]string{
				"ipv4.address": "none",
				"ipv4.nat":     "false",
				"ipv6.address": "none",
				"ipv6.nat":     "false",
			},
		},
		Name: name,
	}

	err := f.server.CreateNetwork(req)
	if err != nil {
		return "", err
	}

	id := xid.New()
	network := Network{
		Id:   id.String(),
		Name: name,
	}
	s, _ := json.Marshal(network)
	f.ds.Set(fmt.Sprintf("networks/bridge/%s", id.String()), string(s))
	return id.String(), nil
}

func (f *Foruka) DeleteNetwork(id string) error {
	if id == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	}

	s := f.ds.Get(fmt.Sprintf("networks/bridge/%s", id))
	nw := &Network{}
	_ = json.Unmarshal([]byte(s), &nw)
	name := nw.Name
	err := f.server.DeleteNetwork(name)
	if err != nil {
		return err
	}
	f.ds.Del(fmt.Sprintf("networks/bridge/%s", nw.Id))
	return nil
}

func (f *Foruka) GetNetwork(id string) Network {
	s := f.ds.Get(fmt.Sprintf("networks/bridge/%s", id))
	fmt.Printf("GETNETWORK-->\n%s\n", s)
	var network Network
	_ = json.Unmarshal([]byte(s), &network)
	fmt.Println(network)
	return network
}

func (f *Foruka) GetNetworks() map[string]Network {
	nws := f.ds.Gets("networks/bridge/%")
	ns := map[string]Network{}
	for _, v := range nws {
		n := &Network{}
		_ = json.Unmarshal([]byte(v), n)
		ns[n.Id] = *n
	}
	return ns
}
