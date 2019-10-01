package core

import (
	"fmt"
	"net"
	"os"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

func (f *Foruka) CreateRouterProfile(name string, iface map[string]string) error {
	if name == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	}

	device := map[string]map[string]string{}

	for k, v := range iface {
		device[k] = map[string]string{}
		device[k]["name"] = k
		device[k]["nictype"] = "bridged"
		device[k]["parent"] = v
		device[k]["type"] = "nic"
	}

	device["root"] = map[string]string{}
	device["root"]["path"] = "/"
	device["root"]["pool"] = "default"
	device["root"]["type"] = "disk"

	pp := api.ProfilesPost{
		ProfilePut: api.ProfilePut{
			Devices: device,
		},
		Name: name,
	}

	return f.server.CreateProfile(pp)
}

func (f *Foruka) CreateRouter(name string, iface map[string]string) error {
	if name == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	}

	device := map[string]map[string]string{}

	for k, v := range iface {
		device[k] = map[string]string{}
		device[k]["name"] = k
		device[k]["nictype"] = "bridged"
		device[k]["parent"] = v
		device[k]["type"] = "nic"
	}

	device["root"] = map[string]string{}
	device["root"]["path"] = "/"
	device["root"]["pool"] = "default"
	device["root"]["type"] = "disk"

	// 仮想1コア, 32MBメモリ
	config := map[string]string{
		"limits.cpu":    "1",
		"limits.memory": "32MB",
	}

	// alias routerで登録されているイメージを使用(後々変更できるようにしたい)
	req := api.ContainersPost{
		Name: name,
		ContainerPut: api.ContainerPut{
			Devices: device,
			Config:  config,
		},
		Source: api.ContainerSource{
			Type:  "image",
			Alias: "router",
		},
	}

	op, err := f.server.CreateContainer(req)
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (f *Foruka) CreateRouterFromProfile(name, profile string) error {

	if name == "" && profile == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`, `profile`\n")
	} else if name == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	} else if profile == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `profile`\n")
	}

	req := api.ContainersPost{
		Name: name,
		ContainerPut: api.ContainerPut{
			Profiles: []string{
				profile,
			},
		},
		Source: api.ContainerSource{
			Type:  "image",
			Alias: "router",
		},
	}

	op, err := f.server.CreateContainer(req)
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (f *Foruka) StartRouter(name string) error {
	return f.updateRouterState(name, "start")
}

func (f *Foruka) StopRouter(name string) error {
	return f.updateRouterState(name, "stop")
}

func (f *Foruka) updateRouterState(name, action string) error {

	if action != "start" && action != "stop" && action != "freeze" && action != "unfreeze" {
		return fmt.Errorf("Error: Invalid Arguments, undefined action `%s`\n", action)
	}

	req := api.ContainerStatePut{
		Action:  action,
		Timeout: -1,
	}
	op, err := f.server.UpdateContainerState(name, req, "")
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}

	return nil
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

	err := f.execRouter(name, command)
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
		f.execRouter(routerName, command)
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
		f.execRouter(routerName, command)
	}
	return nil
}

func (f *Foruka) ConfigureRouterNat(rnat RouterNat) error {
	{
		command := []string{
			"iptables", "-t", "nat", "-F", "POSTROUTING",
		}
		f.execRouter(rnat.RouterName, command)
	}
	{
		command := []string{
			"iptables", "-t", "nat", "-A", "POSTROUTING",
			"-o", rnat.OutIface, "-s", rnat.SrcCidr,
			"-j", "MASQUERADE",
		}
		f.execRouter(rnat.RouterName, command)
	}

	return nil
}
func (f *Foruka) execRouter(name string, command []string) error {
	req := api.ContainerExecPost{
		Command:     command,
		WaitForWS:   true,
		Interactive: false,
	}

	args := lxd.ContainerExecArgs{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	op, err := f.server.ExecContainer(name, req, &args)
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}
