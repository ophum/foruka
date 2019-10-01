package core

import (
	"fmt"
	"net"
)

func (f *Foruka) CreateRouter(name string, iface map[string]string) error {
	if name == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	}

	err := f.CreateContainer(name, "router", iface, map[string]string{
		"cpu":    "1",
		"memory": "32MB",
	})

	if err != nil {
		return err
	}
	return nil
}

func (f *Foruka) StartRouter(name string) error {
	return f.StartContainer(name)
}

func (f *Foruka) StopRouter(name string) error {
	return f.StopContainer(name)
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
