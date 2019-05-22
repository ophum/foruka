package networkModel

import (
	"fmt"
	"os/exec"
)

func AddRouter(name string) error {
	err := exec.Command("ip", "netns", "add", name).Run()
	return err
}

func DelRouter(name string) error {
	err := exec.Command("ip", "netns", "del", name).Run()
	return err
}

func AddEthAdapter(router_name, adapter_name, peer_name string) error {
	err := exec.Command("ip", "link", "add", adapter_name, "type", "veth", "peer", "name", peer_name).Run()
	if err != nil {
		return err
	}

	err = exec.Command("ip", "link", "set", peer_name, "netns", router_name, "up").Run()
	return err
}

func DelEthAdapter(adapter_name string) error {
	err := exec.Command("ip", "link", "del", adapter_name).Run()
	return err
}

func AddDNatHost(proto string, port uint, daddr string, dport uint) error {

	err := exec.Command(
		"iptables",
		"-t", "nat",
		"-A", "PREROUTING",
		"-p", proto, "--dport", fmt.Sprintf("%d", port),
		"-j", "DNAT",
		"--to-destination", fmt.Sprintf("%s:%d", daddr, dport),
	).Run()

	return err
}

func DelDNatHost(proto string, port uint, daddr string, dport uint) error {

	err := exec.Command(
		"iptables",
		"-t", "nat",
		"-D", "PREROUTING",
		"-p", proto, "--dport", fmt.Sprintf("%d", port),
		"-j", "DNAT",
		"--to-destination", fmt.Sprintf("%s:%d", daddr, dport),
	).Run()

	return err
}
