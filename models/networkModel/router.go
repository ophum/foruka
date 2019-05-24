package networkModel

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	pipeline "github.com/mattn/go-pipeline"
)

type Interface struct {
	Name      string
	Peer_name string
	Addr      net.IP
	Prefix    uint
}

type DNat struct {
	Proto string
	Port  uint
	Daddr net.IP
	Dport uint
}

type Route struct {
	Dest       net.IP
	DestPrefix uint
	Next       net.IP
	Adapter    string
}

type Router struct {
	Name     string
	Adapters map[string]*Interface
	Dnats    map[uint]*DNat
	Routes   map[string]map[uint]*Route
	Running  *Router
}

func (r *Router) FetchRunning() error {
	// is not exists...
	if !r.Exists() {
		r.Running = nil
		return fmt.Errorf("%s is not exists", r.Name)
	}

	r.Adapters = map[string]*Interface{}
	r.Dnats = map[uint]*DNat{}
	r.Routes = map[string]map[uint]*Route{}

	r.Running = &Router{
		Name:     r.Name,
		Adapters: map[string]*Interface{},
		Dnats:    map[uint]*DNat{},
		Routes:   map[string]map[uint]*Route{},
	}

	out, err := pipeline.Output(
		[]string{"ip", "netns", "exec", r.Name, "ip", "-4", "-o", "a"},
		[]string{"xargs", "-L", "1", "echo"},
		[]string{"awk", "{print $2 \" \" $4}"},
	)
	if err != nil {
		return err
	}

	slice := strings.Split(string(out), "\n")

	for _, str := range slice {
		if str == "" {
			continue
		}

		str = strings.Replace(str, "/", " ", 1)
		sl := strings.Split(str, " ")
		if len(sl) == 3 {
			name := sl[0]
			addr := net.ParseIP(sl[1])
			prefix, err := strconv.Atoi(sl[2])
			fmt.Println("err:", err)
			r.Adapters[sl[0]] = &Interface{
				Name:   name,
				Addr:   addr,
				Prefix: uint(prefix),
			}
			r.Running.Adapters[sl[0]] = &Interface{
				Name:   name,
				Addr:   addr,
				Prefix: uint(prefix),
			}
		}
	}

	out, err = pipeline.Output(
		[]string{"ip", "netns", "exec", r.Name, "iptables", "-t", "nat", "-n", "-L"},
		[]string{"grep", "DNAT"},
		[]string{"awk", "{print $2 \" \" $7 \" \" $8}"},
	)
	if err != nil {
		return err
	}

	slice = strings.Split(string(out), "\n")
	for _, str := range slice {
		if str == "" {
			continue
		}

		sl := strings.Split(str, " ")
		proto := sl[0]
		port, _ := strconv.Atoi(strings.Split(sl[1], ":")[1])
		daddr := net.ParseIP(strings.Split(sl[2], ":")[1])
		dport, _ := strconv.Atoi(strings.Split(sl[2], ":")[2])
		r.Dnats[uint(port)] = &DNat{
			Proto: proto,
			Port:  uint(port),
			Daddr: daddr,
			Dport: uint(dport),
		}
		r.Running.Dnats[uint(port)] = &DNat{
			Proto: proto,
			Port:  uint(port),
			Daddr: daddr,
			Dport: uint(dport),
		}
	}

	out, err = pipeline.Output(
		[]string{"ip", "netns", "exec", r.Name, "ip", "route"},
		[]string{"grep", "-v", "scope link"},
	)
	if err != nil {
		return err
	}

	slice = strings.Split(string(out), "\n")
	for _, str := range slice {
		if str == "" {
			continue
		}
		sl := strings.Split(str, " ")
		op := 0

		var dest, next net.IP
		var destPrefix int
		var adapter string

		for _, c := range sl {
			switch op {
			case -1:
				switch c {
				case "via":
					op = 1
					break
				case "dev":
					op = 2
					break
				default:
					break
				}
				break
			case 0:
				if c == "default" {
					dest = net.ParseIP("0.0.0.0")
					destPrefix = 0
				} else {
					dest = net.ParseIP(strings.Split(c, "/")[0])
					p, _ := strconv.Atoi(strings.Split(c, "/")[1])
					destPrefix = p
				}
				op = -1
				break
			case 1:
				next = net.ParseIP(c)
				op = -1
				break
			case 2:
				adapter = c
				op = -1
				break
			}

		}
		r.Routes[dest.String()] = map[uint]*Route{}
		r.Routes[dest.String()][uint(destPrefix)] = &Route{
			Dest:       dest,
			DestPrefix: uint(destPrefix),
			Next:       next,
			Adapter:    adapter,
		}
		r.Running.Routes[dest.String()] = map[uint]*Route{}
		r.Running.Routes[dest.String()][uint(destPrefix)] = &Route{
			Dest:       dest,
			DestPrefix: uint(destPrefix),
			Next:       next,
			Adapter:    adapter,
		}
	}
	return nil
}

func DispRouter(r Router) {
	fmt.Println(r.Name)
	fmt.Println("Adapter")
	for key, val := range r.Adapters {
		fmt.Println("-->", key)
		fmt.Printf("----> %s/%d\n", val.Addr.String(), val.Prefix)
	}
	fmt.Println("DNAT")
	for _, val := range r.Dnats {
		fmt.Printf("--> %s %d -> %s:%d\n", val.Proto, val.Port, val.Daddr.String(), val.Dport)
	}
	fmt.Println("Route")
	for _, val := range r.Routes {
		for _, v := range val {
			fmt.Printf("--> %s/%d -> %s dev %s\n", v.Dest.String(), v.DestPrefix, v.Next.String(), v.Adapter)
		}
	}

}
func (r *Router) Exists() bool {
	_, err := os.Stat(fmt.Sprintf("/var/run/netns/%s", r.Name))
	return err == nil
}

func (r *Router) Apply() {
	if !r.Exists() {
		r.Add(r.Name)
	}

	// adapters update...
	for k, v := range r.Adapters {
		rv, ok := r.Running.Adapters[k]
		if ok {
			// check ip and prefix
			if !(rv.Addr.Equal(v.Addr) && rv.Prefix == v.Prefix) {
				r.AddIp(v.Name)
				r.Running.DelIp(rv.Name)
			}
		} else { // not found, create eth adapter
			r.AddEthAdapter(v.Name)
			r.AddIp(v.Name)
		}
	}

	// dnats update...
	for k, v := range r.Dnats {
		rv, ok := r.Running.Dnats[k]
		if ok {
			// check dnat
			if !(rv.Proto == v.Proto && rv.Port == v.Port && rv.Daddr.Equal(v.Daddr) && rv.Dport == v.Dport) {
				r.AddDNat(v.Port)
				r.Running.DelDNat(rv.Port)
			}
		} else {
			r.AddDNat(v.Port)
		}
	}

	// routes update...

}
func (r *Router) Add(name string) error {
	err := exec.Command("ip", "netns", "add", name).Run()
	if err != nil {
		return err
	}

	err = exec.Command(
		"ip", "netns", "exec", name,
		"sysctl", "-w", "net.ipv4.ip_forward=1",
	).Run()
	return err
}

func DelRouter(name string) error {
	err := exec.Command("ip", "netns", "del", name).Run()
	return err
}

func (r *Router) AddIp(adapter string) error {
	adap := r.Adapters[adapter]
	ipv4 := fmt.Sprintf("%s/%d", adap.Addr.String(), adap.Prefix)
	err := exec.Command(
		"ip", "netns", "exec", r.Name,
		"ip", "addr", "add", ipv4, "dev", adapter,
	).Run()
	return err
}

func (r *Router) DelIp(adapter string) error {
	err := exec.Command(
		"ip", "netns", "exec", r.Name,
		"ip", "addr", "del", r.Adapters[adapter].Addr.String(), "dev", adapter,
	).Run()
	return err
}

func (r *Router) AddEthAdapter(adapter string) error {
	peer_name := r.Adapters[adapter].Peer_name
	err := exec.Command("ip", "link", "add", adapter, "type", "veth", "peer", "name", peer_name).Run()
	if err != nil {
		return err
	}

	err = exec.Command("ip", "link", "set", adapter, "netns", r.Name, "up").Run()
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

func (r *Router) AddDNat(port uint) error {
	dnat := r.Dnats[port]
	err := checkProto(dnat.Proto)
	if err != nil {
		return err
	}

	err = checkPort(dnat.Port)
	if err != nil {
		return err
	}

	err = checkPort(dnat.Dport)
	if err != nil {
		return fmt.Errorf("Invalid dport -> %d", dnat.Dport)
	}

	err = exec.Command(
		"ip", "netns", "exec", r.Name,
		"iptables",
		"-t", "nat",
		"-A", "PREROUTING",
		"-p", dnat.Proto, "--dport", fmt.Sprintf("%d", dnat.Port),
		"-j", "DNAT",
		"--to-destination", fmt.Sprintf("%s:%d", dnat.Daddr, dnat.Dport),
	).Run()

	return err
}

func (r *Router) DelDNat(port uint) error {
	dnat := r.Dnats[port]
	err := checkProto(dnat.Proto)
	if err != nil {
		return err
	}

	err = checkPort(dnat.Port)
	if err != nil {
		return err
	}

	err = checkPort(dnat.Dport)
	if err != nil {
		return fmt.Errorf("Invalid dport -> %d", dnat.Dport)
	}

	err = exec.Command(
		"ip", "netns", "exec", r.Name,
		"iptables",
		"-t", "nat",
		"-D", "PREROUTING",
		"-p", dnat.Proto, "--dport", fmt.Sprintf("%d", dnat.Port),
		"-j", "DNAT",
		"--to-destination", fmt.Sprintf("%s:%d", dnat.Daddr, dnat.Dport),
	).Run()

	return err
}
