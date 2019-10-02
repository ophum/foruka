package core

import (
	lxd "github.com/lxc/lxd/client"
	//"github.com/lxc/lxd/shared/api"
)

type Foruka struct {
	server    lxd.ContainerServer
	IsCluster bool
}

func NewForukaUnix(sockPath string) (*Foruka, error) {
	s, err := lxd.ConnectLXDUnix(sockPath, nil)
	if err != nil {
		return nil, err
	}
	f := &Foruka{
		server:    s,
		IsCluster: s.IsClustered(),
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
