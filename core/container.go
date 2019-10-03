package core

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rs/xid"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

type Container struct {
	Id     string            `json:"id"`
	Name   string            `json:"name"`
	Iface  map[string]string `json:"iface"`
	Limits map[string]string `json:"limits"`
}

func (f *Foruka) CreateContainer(name, alias string, iface map[string]string, limits map[string]string) (string, error) {
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
		return "", err
	}

	err = op.Wait()
	if err != nil {
		return "", err
	}

	id := xid.New()
	c := &Container{
		Id:     id.String(),
		Name:   name,
		Iface:  iface,
		Limits: limits,
	}

	s, _ := json.Marshal(c)
	f.ds.Set(fmt.Sprintf("container/%s", id.String()), string(s))
	f.ds.Set(fmt.Sprintf("container/status/%s", id.String()), "stopped")
	return id.String(), nil
}

func (f *Foruka) GetContainerStatus(id string) string {
	s := f.ds.Get(fmt.Sprintf("container/status/%s", id))
	return s
}

func (f *Foruka) StartContainer(id string) error {
	err := f.updateContainerState(id, "start")
	if err != nil {
		return err
	}
	f.ds.Set(fmt.Sprintf("container/status/%s", id), "running")
	time.Sleep(time.Second)
	return nil
}

func (f *Foruka) StopContainer(id string) error {
	f.ds.Set(fmt.Sprintf("container/status/%s", id), "stopped")
	return f.updateContainerState(id, "stop")
}

func (f *Foruka) updateContainerState(id, action string) error {

	if action != "start" && action != "stop" && action != "freeze" && action != "unfreeze" {
		return fmt.Errorf("Error: Invalid Arguments, undefined action `%s`\n", action)
	}

	s := f.ds.Get(fmt.Sprintf("container/%s", id))
	con := &Container{}
	_ = json.Unmarshal([]byte(s), &con)
	name := con.Name
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

func (f *Foruka) ExecContainer(id string, command []string) error {
	s := f.ds.Get(fmt.Sprintf("container/%s", id))
	con := &Container{}
	_ = json.Unmarshal([]byte(s), &con)
	name := con.Name

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
