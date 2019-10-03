package core

import (
	lxd "github.com/lxc/lxd/client"
	"github.com/ophum/foruka/datastore"
	//"github.com/lxc/lxd/shared/api"
)

type Foruka struct {
	server    lxd.ContainerServer
	ds        *datastore.DataStore
	IsCluster bool
}

func NewForukaUnix(sockPath, install_dir string) (*Foruka, error) {
	s, err := lxd.ConnectLXDUnix(sockPath, nil)
	if err != nil {
		return nil, err
	}
	f := &Foruka{
		server:    s,
		IsCluster: s.IsClustered(),
		ds:        datastore.NewDataStore(install_dir),
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
