package core

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
)

func (f *Foruka) AddProxyRule(dport, rtip string) error {
	ipt, err := iptables.New()
	if err != nil {
		return err
	}

	err = ipt.AppendUnique(
		"nat", "PREROUTING",
		"-i", f.externalInterface, "-p", "tcp", "--dport", dport, "-j", "DNAT",
		"--to-destination", fmt.Sprintf("%s:%s", rtip, dport),
	)
	return err
}
