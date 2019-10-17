package core

import (
	"net"

	lxd "github.com/lxc/lxd/client"
	//"github.com/lxc/lxd/shared/api"
)

type Foruka struct {
	server            lxd.ContainerServer
	externalInterface string
	externalIP        string
	IsCluster         bool
}

func (f *Foruka) GetExternalIP() string {
	return f.externalIP
}
func getInternalIP(iface string) string {
	itf, _ := net.InterfaceByName(iface)
	item, _ := itf.Addrs()
	var ip net.IP
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() && v.IP.To4() != nil {
				ip = v.IP
			}
		}
	}
	if ip != nil {
		return ip.String()
	}
	return ""
}
func NewForukaUnix(externalInterface, sockPath string) (*Foruka, error) {
	s, err := lxd.ConnectLXDUnix(sockPath, nil)
	if err != nil {
		return nil, err
	}

	f := &Foruka{
		server:            s,
		externalInterface: externalInterface,
		externalIP:        getInternalIP(externalInterface),
		IsCluster:         s.IsClustered(),
	}
	return f, nil
}

// lxc remote add <FQDN|IP|URL>
// ~/snap/lxd/
func NewForuka(url, serverCert, clientCert, clientKey string) (*Foruka, error) {
	args := lxd.ConnectionArgs{
		TLSServerCert: serverCert,
		TLSClientCert: clientCert,
		TLSClientKey:  clientKey,
	}
	s, err := lxd.ConnectLXD(url, &args)
	if err != nil {
		return nil, err
	}

	f := &Foruka{
		server:    s,
		IsCluster: s.IsClustered(),
	}
	return f, nil
}
