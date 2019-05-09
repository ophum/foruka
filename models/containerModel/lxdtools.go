package containerModel

import (
	"fmt"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

var cont = connect("/var/snap/lxd/common/lxd/unix.socket")

func connect(path string) (container lxd.ContainerServer) {
	container, err := lxd.ConnectLXDUnix(path, nil)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect lxd")
	}
	return container
}

func Status(name string) (*api.ContainerState, error) {
	var stat *api.ContainerState

	stat, str, err := cont.GetContainerState(name)
	if err != nil {
		return nil, err
	}
	fmt.Println(str)
	return stat, nil
}

func CreateContainer(name string, alias string) (string, error) {
	req := api.ContainersPost{
		Name: name,
		Source: api.ContainerSource{
			Type:     "image",
			Alias:    "ubuntu/disco",
			Server:   "https://images.linuxcontainers.org",
			Protocol: "simplestreams",
		},
	}

	op, err := cont.CreateContainer(req)
	if err != nil {
		return "", err
	}

	err = op.Wait()
	if err != nil {
		return "", err
	}
	return "success", nil
}

func LaunchContainer(name string) error {
	req := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err := cont.UpdateContainerState(name, req, "")
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}
