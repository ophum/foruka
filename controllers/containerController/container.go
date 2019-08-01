package containerController

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lxc/lxd/shared/api"
	auth "github.com/ophum/foruka/controllers/authController"
	contmodel "github.com/ophum/foruka/models/containerModel"
)

func Index(c *gin.Context) {
	user := auth.Auth(c)

	containers := contmodel.GetContainers(user.ID)
	c.HTML(200, "containers/index.tmpl", gin.H{
		"containers": containers,
	})
}

func Create(c *gin.Context) {
	auth.Auth(c)

	c.HTML(200, "containers/create.tmpl", gin.H{})
}

func Store(c *gin.Context) {
	user := auth.Auth(c)
	name := c.PostForm("name")
	username := c.PostForm("user")
	sshkey := c.PostForm("ssh-authorized-key")
	image := c.PostForm("image")

	fmt.Println("user: ", user)
	fmt.Println("name: ", name)
	fmt.Println("username: ", username)
	fmt.Println("ssh-authorized-key: ", sshkey)
	image = strings.Replace(image, " ", "", -1)
	fmt.Println("image:", image)
	contmodel.Create(user.ID, name, image)

	contmodel.ExecContainer(user.Name+"-"+name, []string{"apt", "install", "-y", "openssh-server"})
	cmdlines := [][]string{
		[]string{"useradd", "-m", "-s", "/bin/bash", username},
		[]string{"mkdir", "-p", fmt.Sprintf("/home/%s/.ssh", username)},
		[]string{"bash", "-c", fmt.Sprintf("echo \"%s\" > /home/%s/.ssh/authorized_keys", sshkey, username)},
		[]string{"chown", "-R", fmt.Sprintf("%s.%s", username, username), fmt.Sprintf("/home/%s/.ssh", username)},
		[]string{"chmod", "700", fmt.Sprintf("/home/%s/.ssh", username)},
		[]string{"chmod", "600", fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)},
	}

	for _, cmd := range cmdlines {
		contmodel.ExecContainer(user.Name+"-"+name, cmd)
	}
	c.Redirect(302, "/containers/")
}

func Show(c *gin.Context) {
	user := auth.Auth(c)
	hashId := c.Param("id")

	var cont contmodel.Container
	cont = contmodel.GetContainer(user.ID, hashId)
	status, _ := contmodel.Status(cont.Name)
	var addresses = []api.ContainerStateNetworkAddress{}
	if status.Status == "Running" {
		addresses = status.Network["eth0"].Addresses
	}

	c.HTML(200, "containers/show.tmpl", gin.H{
		"container": cont,
		"status":    status,
		"addresses": addresses,
	})
}

func Start(c *gin.Context) {
	user := auth.Auth(c)

	hashId := c.Param("id")

	cont := contmodel.GetContainer(user.ID, hashId)
	err := contmodel.LaunchContainer(cont.Name)
	if err != nil {
		fmt.Println("err: ", err)
	}

	c.Redirect(302, "/containers/show/"+hashId)
}

func Stop(c *gin.Context) {
	user := auth.Auth(c)
	hashId := c.Param("id")

	cont := contmodel.GetContainer(user.ID, hashId)
	err := contmodel.StopContainer(cont.Name)
	if err != nil {
		fmt.Println("err: ", err)
	}

	c.Redirect(302, "/containers/show/"+hashId)
}

func Delete(c *gin.Context) {
	user := auth.Auth(c)
	hashId := c.Param("id")

	cont := contmodel.GetContainer(user.ID, hashId)
	err := contmodel.DeleteContainer(cont.Name)
	if err != nil {
		fmt.Println("err: ", err)
	}

	contmodel.Delete(hashId)

	c.Redirect(302, "/containers/")
}
