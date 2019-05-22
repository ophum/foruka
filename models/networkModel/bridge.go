package networkModel

import (
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
