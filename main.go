package main

import (
	"fmt"

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
	err = frk.CreateRouter("hogeRouter", map[string]string{
		"eth0": "lxdbr0",
		"ens1": "testtesttest",
	})

	if err != nil {
		fmt.Println(err)
		//return
	}
	err = frk.StartRouter("hogeRouter")
	if err != nil {
		fmt.Println(err)
	}

	err = frk.StopRouter("hogeRouter")
	if err != nil {
		fmt.Println(err)
	}
}
