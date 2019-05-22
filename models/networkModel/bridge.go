package networkModel

import (
	"fmt"
	"os/exec"
)

func AddEthAdapterBridge(bridge, adapter_name string) error {
	err := exec.Command(
		"brctl", "addif", bridge, adapter_name,
	).Run()
	return err
}

func DelEthAdapterBridge(bridge, adapter_name string) error {
	err := exec.Command(
		"brctl", "delif", bridge, adapter_name,
	).Run()
	return err
}

func UpEthAdapterBridge(adapter_name string) error {
	err := exec.Command(
		"ip", "link", "set", adapter_name, "up",
	).Run()
	return err
}

func DownEthAdapterBridge(adapter_name string) error {
	err := exec.Command(
		"ip", "link", "set", adapter_name, "down",
	).Run()
	return err
}

func AddVlan(adapter_name string, vid uint) error {
	err := exec.Command(
		"bridge", "vlan", "add",
		"vid", fmt.Sprintf("%d", vid),
		"dev", adapter_name,
		"pvid", "untagged",
	).Run()
	return err
}

func DelVlan(adapter_name string, vid uint) error {
	err := exec.Command(
		"bridge", "vlan", "del",
		"dev", adapter_name,
		"vid", fmt.Sprintf("%d", vid),
	).Run()
	return err
}
