package core

import (
	lxd "github.com/lxc/lxd/client"
	//"github.com/lxc/lxd/shared/api"
)

type Foruka struct {
	server	lxd.ContainerServer
}

func NewForuka(sockPath string) (*Foruka, error) {
	s, err := lxd.ConnectLXDUnix(sockPath, nil)
	if err != nil {
		return nil, err
	}
	f := &Foruka{
		server: s,
	}
	return f, nil
}
