package networkModel

import (
	"errors"
	"fmt"
	"net"
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

func AddIpRouter(router, adapter, addr string) error {
	err := exec.Command(
		"ip", "netns", "exec", router,
		"ip", "addr", "add", addr, "dev", adapter,
	).Run()
	return err
}

func DelIpRouter(router, adapter, addr string) error {
	err := exec.Command(
		"ip", "netns", "exec", router,
		"ip", "addr", "del", addr, "dev", adapter,
	).Run()
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

func UpEthAdapter(router, adapter_name string) error {
	err := exec.Command(
		"ip", "netns", "exec", router,
		"ip", "link", "set", adapter_name, "up",
	).Run()

	return err
}

func DownEthAdapter(router, adapter_name string) error {
	err := exec.Command(
		"ip", "netns", "exec", router,
		"ip", "link", "set", adapter_name, "down",
	).Run()

	return err
}

func checkProto(proto string) error {
	switch proto {
	case "tcp":
	case "udp":
		break
	default:
		return errors.New(fmt.Sprintf("Invalid proto -> %s", proto))
	}
	return nil
}

func checkPort(port uint) error {
	if port == 0 || port > 65535 {
		return errors.New(fmt.Sprintf("Invalid port -> %d", port))
	}
	return nil
}

func AddDNatHost(proto string, port uint, daddr string, dport uint) error {

	err := checkProto(proto)
	if err != nil {
		return err
	}

	err = checkPort(port)
	if err != nil {
		return err
	}

	err = checkPort(dport)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid dport -> %d", dport))
	}

	if net.ParseIP(daddr) == nil {
		return errors.New(fmt.Sprintf("Invalid daddr -> %s", daddr))
	}

	err = exec.Command(
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

	err := checkProto(proto)
	if err != nil {
		return err
	}

	err = checkPort(port)
	if err != nil {
		return err
	}

	err = checkPort(dport)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid dport -> %d", dport))
	}

	if net.ParseIP(daddr) == nil {
		return errors.New(fmt.Sprintf("Invalid daddr -> %s", daddr))
	}

	err = exec.Command(
		"iptables",
		"-t", "nat",
		"-D", "PREROUTING",
		"-p", proto, "--dport", fmt.Sprintf("%d", port),
		"-j", "DNAT",
		"--to-destination", fmt.Sprintf("%s:%d", daddr, dport),
	).Run()

	return err
}

func AddDNatRouter(router string, proto string, port uint, daddr string, dport uint) error {

	err := checkProto(proto)
	if err != nil {
		return err
	}

	err = checkPort(port)
	if err != nil {
		return err
	}

	err = checkPort(dport)
	if err != nil {
		return fmt.Errorf("Invalid dport -> %d", dport)
	}

	if net.ParseIP(daddr) == nil {
		return fmt.Errorf("Invalid daddr -> %s", daddr)
	}

	err = exec.Command(
		"ip", "netns", "exec", router,
		"iptables",
		"-t", "nat",
		"-A", "PREROUTING",
		"-p", proto, "--dport", fmt.Sprintf("%d", port),
		"-j", "DNAT",
		"--to-destination", fmt.Sprintf("%s:%d", daddr, dport),
	).Run()

	return err
}

func DelDNatRouter(router string, proto string, port uint, daddr string, dport uint) error {

	err := checkProto(proto)
	if err != nil {
		return err
	}

	err = checkPort(port)
	if err != nil {
		return err
	}

	err = checkPort(dport)
	if err != nil {
		return fmt.Errorf("Invalid dport -> %d", dport)
	}

	if net.ParseIP(daddr) == nil {
		return fmt.Errorf("Invalid daddr -> %s", daddr)
	}

	err = exec.Command(
		"ip", "netns", "exec", router,
		"iptables",
		"-t", "nat",
		"-D", "PREROUTING",
		"-p", proto, "--dport", fmt.Sprintf("%d", port),
		"-j", "DNAT",
		"--to-destination", fmt.Sprintf("%s:%d", daddr, dport),
	).Run()

	return err
}
