package core

import (
	"fmt"
	"os"
	"time"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

func (f *Foruka) CreateContainer(name, alias string, iface map[string]string, limits map[string]string) error {
	devices := map[string]map[string]string{}
	for k, v := range iface {
		devices[k] = map[string]string{
			"name":    k,
			"nictype": "bridged",
			"parent":  v,
			"type":    "nic",
		}
	}

	devices["root"] = map[string]string{
		"path": "/",
		"pool": "default",
		"type": "disk",
	}

	config := map[string]string{}

	for k, v := range limits {
		config[fmt.Sprintf("limits.%s", k)] = v
	}

	req := api.ContainersPost{
		Name: name,
		ContainerPut: api.ContainerPut{
			Devices: devices,
			Config:  config,
		},
		Source: api.ContainerSource{
			Type:  "image",
			Alias: alias,
		},
	}

	op, err := f.server.CreateContainer(req)
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (f *Foruka) StartContainer(name string) error {
	err := f.updateContainerState(name, "start")
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	return nil
}

func (f *Foruka) StopContainer(name string) error {
	return f.updateContainerState(name, "stop")
}

func (f *Foruka) updateContainerState(name, action string) error {

	if action != "start" && action != "stop" && action != "freeze" && action != "unfreeze" {
		return fmt.Errorf("Error: Invalid Arguments, undefined action `%s`\n", action)
	}

	req := api.ContainerStatePut{
		Action:  action,
		Timeout: -1,
	}
	op, err := f.server.UpdateContainerState(name, req, "")
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (f *Foruka) ExecContainer(name string, command []string) error {
	req := api.ContainerExecPost{
		Command:     command,
		WaitForWS:   true,
		Interactive: false,
	}

	args := lxd.ContainerExecArgs{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	op, err := f.server.ExecContainer(name, req, &args)
	if err != nil {
		return err
	}

	err = op.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (f *Foruka) GetContainers() ([]api.Container, error) {
	return f.server.GetContainers()
}

func (f *Foruka) GetContainerNames() ([]string, error) {
	return f.server.GetContainerNames()
}

func (f *Foruka) GetContainer(name string) (*api.Container, string, error) {

	container, etag, err := f.server.GetContainer(name)
	return container, etag, err
}

func (f *Foruka) GetContainerState(name string) (*api.ContainerState, error) {
	state, _, err := f.server.GetContainerState(name)
	return state, err
}
