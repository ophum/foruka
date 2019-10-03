package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "foruka - container hosting system"
	app.Usage = "foruka command"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "initialize foruka",
			Action: InitForuka,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "external-bridge",
				},
				cli.StringFlag{
					Name: "connection-type",
				},
				cli.StringFlag{
					Name: "path",
				},
				cli.StringFlag{
					Name: "install-dir",
				},
			},
		},
		{
			Name:   "serve",
			Usage:  "serve foruka server",
			Action: ServeForuka,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "install-dir",
				},
			},
		},
	}

	app.Run(os.Args)
	//	frk, err := core.NewForukaUnix("/var/snap/lxd/common/lxd/unix.socket")
	//	//frk, err := core.NewForuka("https://10.55.37.84:8443")
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//	//	err = frk.CreateRouterProfile("testtestprofile", map[string]string{
	//	//		"eth0": "lxdbr0",
	//	//		"eth1": "testtesttest",
	//	//	})
	//	err = frk.CreateNetwork("test_network")
	//	//
	//
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	err = frk.CreateRouter("hogeRouter", map[string]string{
	//		"eth0": "lxdbr0",
	//		"ens1": "test_network",
	//	})
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		//return
	//	}
	//	err = frk.StartRouter("hogeRouter")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	err = frk.ConfigureRouterInterface("hogeRouter", core.RouterInterface{
	//		Name:        "ens1",
	//		Ipv4Address: net.ParseIP("192.168.10.254"),
	//		Ipv4Prefix:  24,
	//	})
	//
	//	err = frk.ConfigureRouterPortForwarding("hogeRouter", []core.RouterPortForwardTable{
	//		core.RouterPortForwardTable{
	//			Iface:      "eth0",
	//			Proto:      "tcp",
	//			Dport:      80,
	//			ToDestIP:   net.ParseIP("192.168.10.1"),
	//			ToDestPort: 80,
	//		},
	//		core.RouterPortForwardTable{
	//			Iface:      "eth0",
	//			Proto:      "tcp",
	//			Dport:      443,
	//			ToDestIP:   net.ParseIP("192.168.10.1"),
	//			ToDestPort: 80,
	//		},
	//	})
	//
	//	err = frk.ConfigureRouterNat(core.RouterNat{
	//		RouterName: "hogeRouter",
	//		SrcCidr:    "192.168.10.0/24",
	//		OutIface:   "eth0",
	//	})
	//
	//	err = frk.CreateContainer("t1", "router", map[string]string{
	//		"eth0": "test_network",
	//	}, map[string]string{
	//		"cpu":    "2",
	//		"memory": "64MB",
	//	})
	//
	//	err = frk.StartContainer("t1")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = frk.ExecContainer("t1", []string{
	//		"ip", "a", "add", "192.168.10.1/24", "dev", "eth0",
	//	})
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = frk.ExecContainer("t1", []string{
	//		"ip", "route", "add", "default", "via", "192.168.10.254",
	//	})
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = frk.CreateContainer("t2", "router", map[string]string{
	//		"eth0": "test_network",
	//	}, map[string]string{
	//		"cpu":    "4",
	//		"memory": "64MB",
	//	})
	//
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = frk.StartContainer("t2")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = frk.ExecContainer("t2", []string{
	//		"ip", "a", "add", "192.168.10.2/24", "dev", "eth0",
	//	})
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = frk.ExecContainer("t2", []string{
	//		"ip", "route", "add", "default", "via", "192.168.10.254",
	//	})
	//	if err != nil {
	//		fmt.Println(err)
	//	}
}
