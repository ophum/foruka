package networkModel

import (
	"fmt"
	"os/exec"
)

type Bridge struct {
	Name string
}

func (b *Bridge) AddEthAdapter(adapter_name string) error {
	err := exec.Command(
		"brctl", "addif", b.Name, adapter_name,
	).Run()
	return err
}

func (b *Bridge) DelEthAdapter(adapter_name string) error {
	err := exec.Command(
		"brctl", "delif", b.Name, adapter_name,
	).Run()
	return err
}

func (b *Bridge) UpEthAdapter(adapter_name string) error {
	err := exec.Command(
		"ip", "link", "set", adapter_name, "up",
	).Run()
	return err
}

func (b *Bridge) DownEthAdapter(adapter_name string) error {
	err := exec.Command(
		"ip", "link", "set", adapter_name, "down",
	).Run()
	return err
}

func (b *Bridge) AddVlan(adapter_name string, vid uint) error {
	err := exec.Command(
		"bridge", "vlan", "add",
		"vid", fmt.Sprintf("%d", vid),
		"dev", adapter_name,
		"pvid", "untagged",
	).Run()
	return err
}

func (b *Bridge) DelVlan(adapter_name string, vid uint) error {
	err := exec.Command(
		"bridge", "vlan", "del",
		"dev", adapter_name,
		"vid", fmt.Sprintf("%d", vid),
	).Run()
	return err
}
