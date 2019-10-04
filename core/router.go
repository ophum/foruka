package core

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/rs/xid"
)

type Router struct {
	Id    string            `json:"id"`
	Cid   string            `json:"cid"`
	Name  string            `json:"name"`
	Iface map[string]string `json:"iface"`
}

func (f *Foruka) CreateRouter(name string, iface map[string]string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	}

	cid, err := f.CreateContainer(name, "router", iface, map[string]string{
		"cpu":    "1",
		"memory": "32MB",
	})

	if err != nil {
		return "", err
	}

	id := xid.New()
	c := &Router{
		Id:    id.String(),
		Cid:   cid,
		Name:  name,
		Iface: iface,
	}

	s, _ := json.Marshal(c)
	f.ds.Set(fmt.Sprintf("router/%s", id.String()), string(s))
	return id.String(), nil
}

func (f *Foruka) StartRouter(id string) error {
	s := f.ds.Get(fmt.Sprintf("router/%s", id))
	r := &Router{}
	_ = json.Unmarshal([]byte(s), &r)

	return f.StartContainer(r.Cid)
}

func (f *Foruka) StopRouter(id string) error {
	s := f.ds.Get(fmt.Sprintf("router/%s", id))
	r := &Router{}
	_ = json.Unmarshal([]byte(s), &r)

	return f.StopContainer(r.Cid)
}

func (f *Foruka) GetRouter(id string) Router {
	s := f.ds.Get(fmt.Sprintf("router/%s", id))
	r := Router{}
	_ = json.Unmarshal([]byte(s), &r)
	return r
}

func (f *Foruka) GetRouters() map[string]Router {
	rts := f.ds.Gets("router/%")
	rs := map[string]Router{}
	for _, v := range rts {
		r := &Router{}
		_ = json.Unmarshal([]byte(v), r)
		rs[r.Id] = *r
	}
	return rs
}
func (f *Foruka) GetRouterStatus(id string) string {
	s := f.ds.Get(fmt.Sprintf("router/%s", id))
	r := &Router{}
	_ = json.Unmarshal([]byte(s), &r)

	return f.GetContainerStatus(r.Cid)
}

type RouterInterface struct {
	Name        string
	Ipv4Address net.IP
	Ipv4Prefix  int
}

type RouterPortForwardTable struct {
	Iface      string
	Proto      string
	Dport      int
	ToDestIP   net.IP
	ToDestPort int
}

type RouterNat struct {
	RouterName string
	SrcCidr    string
	OutIface   string
}

func (f *Foruka) ConfigureRouterInterface(name string, iface RouterInterface) error {

	command := []string{
		"ip", "a", "add", fmt.Sprintf("%s/%d", iface.Ipv4Address, iface.Ipv4Prefix),
		"dev", fmt.Sprintf("%s", iface.Name),
	}

	err := f.ExecContainer(name, command)
	if err != nil {
		return err
	}
	return nil
}

func (f *Foruka) ConfigureRouterPortForwarding(routerName string, tables []RouterPortForwardTable) error {
	{
		command := []string{
			"iptables", "-t", "nat", "-F", "PREROUTING",
		}
		f.ExecContainer(routerName, command)
	}
	for _, v := range tables {
		command := []string{
			"iptables", "-t", "nat",
			"-A", "PREROUTING", "-p", v.Proto,
			"-i", v.Iface,
			"--dport", fmt.Sprintf("%d", v.Dport),
			"-j", "DNAT", "--to-destination",
			fmt.Sprintf("%s:%d", v.ToDestIP, v.ToDestPort),
		}
		f.ExecContainer(routerName, command)
	}
	return nil
}

func (f *Foruka) ConfigureRouterNat(rnat RouterNat) error {
	{
		command := []string{
			"iptables", "-t", "nat", "-F", "POSTROUTING",
		}
		f.ExecContainer(rnat.RouterName, command)
	}
	{
		command := []string{
			"iptables", "-t", "nat", "-A", "POSTROUTING",
			"-o", rnat.OutIface, "-s", rnat.SrcCidr,
			"-j", "MASQUERADE",
		}
		f.ExecContainer(rnat.RouterName, command)
	}

	return nil
}
