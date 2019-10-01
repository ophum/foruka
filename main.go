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

	err = frk.DeleteNetwork("test1234")
	if err != nil {
		fmt.Println(err)
		return
	}

}
