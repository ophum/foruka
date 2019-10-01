package main

import (
	"fmt"
	"net"

	"github.com/ophum/foruka/core"
)

func main() {
	frk, err := core.NewForuka("/var/snap/lxd/common/lxd/unix.socket")
	if err != nil {
		fmt.Println(err)
		return
	}

	//	err = frk.CreateRouterProfile("testtestprofile", map[string]string{
	//		"eth0": "lxdbr0",
	//		"eth1": "testtesttest",
	//	})
	//

	err = frk.CreateNetwork("test_network")
	if err != nil {
		fmt.Println(err)
	}

	err = frk.CreateRouter("hogeRouter", map[string]string{
		"eth0": "lxdbr0",
		"ens1": "test_network",
	})

	if err != nil {
		fmt.Println(err)
		//return
	}
	err = frk.StartRouter("hogeRouter")
	if err != nil {
		fmt.Println(err)
	}

	err = frk.ConfigureRouterInterface("hogeRouter", core.RouterInterface{
		Name:        "ens1",
		Ipv4Address: net.ParseIP("192.168.10.1"),
		Ipv4Prefix:  24,
	})

	err = frk.ConfigureRouterPortForwarding("hogeRouter", []core.RouterPortForwardTable{
		core.RouterPortForwardTable{
			Iface:      "eth0",
			Proto:      "tcp",
			Dport:      80,
			ToDestIP:   net.ParseIP("192.168.10.12"),
			ToDestPort: 80,
		},
		core.RouterPortForwardTable{
			Iface:      "eth0",
			Proto:      "tcp",
			Dport:      442,
			ToDestIP:   net.ParseIP("192.168.10.12"),
			ToDestPort: 80,
		},
	})

	err = frk.ConfigureRouterNat(core.RouterNat{
		RouterName: "hogeRouter",
		SrcCidr:    "192.168.10.0/24",
		OutIface:   "eth0",
	})

	if err != nil {
		fmt.Println(err)
	}
}
