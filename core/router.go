package core

import (
	"fmt"

	"github.com/lxc/lxd/shared/api"
)

func (f *Foruka) CreateRouterProfile(name string, iface map[string]string) error {
	if name == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	}

	device := map[string]map[string]string{}

	for k, v := range iface {
		device[k] = map[string]string{}
		device[k]["name"] = k
		device[k]["nictype"] = "bridged"
		device[k]["parent"] = v
		device[k]["type"] = "nic"
	}

	device["root"] = map[string]string{}
	device["root"]["path"] = "/"
	device["root"]["pool"] = "default"
	device["root"]["type"] = "disk"

	pp := api.ProfilesPost{
		ProfilePut: api.ProfilePut{
			Devices: device,
		},
		Name: name,
	}

	return f.server.CreateProfile(pp)
}

func (f *Foruka) CreateRouter(name string, iface map[string]string) error {
	if name == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	}

	device := map[string]map[string]string{}

	for k, v := range iface {
		device[k] = map[string]string{}
		device[k]["name"] = k
		device[k]["nictype"] = "bridged"
		device[k]["parent"] = v
		device[k]["type"] = "nic"
	}

	device["root"] = map[string]string{}
	device["root"]["path"] = "/"
	device["root"]["pool"] = "default"
	device["root"]["type"] = "disk"

	// 仮想1コア, 32MBメモリ
	config := map[string]string{
		"limits.cpu":    "1",
		"limits.memory": "32MB",
	}

	// alias routerで登録されているイメージを使用(後々変更できるようにしたい)
	req := api.ContainersPost{
		Name: name,
		ContainerPut: api.ContainerPut{
			Devices: device,
			Config:  config,
		},
		Source: api.ContainerSource{
			Type:  "image",
			Alias: "router",
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

func (f *Foruka) CreateRouterFromProfile(name, profile string) error {

	if name == "" && profile == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`, `profile`\n")
	} else if name == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `name`\n")
	} else if profile == "" {
		return fmt.Errorf("Error: Invalid Arguments, Empty `profile`\n")
	}

	req := api.ContainersPost{
		Name: name,
		ContainerPut: api.ContainerPut{
			Profiles: []string{
				profile,
			},
		},
		Source: api.ContainerSource{
			Type:  "image",
			Alias: "router",
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

func (f *Foruka) StartRouter(name string) error {
	return f.updateRouterState(name, "start")
}

func (f *Foruka) StopRouter(name string) error {
	return f.updateRouterState(name, "stop")
}

func (f *Foruka) updateRouterState(name, action string) error {

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
