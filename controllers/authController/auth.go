package authController

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ophum/foruka/models/authModel"
)

func IsAuth(c *gin.Context) bool {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	if user_id != nil {
		return true
	}
	return false
}

func Auth(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	if user_id == nil {
		c.Redirect(301, "/login")
	}
}

func Verified(c *gin.Context) {

	if IsAuth(c) {
		session := sessions.Default(c)
		user_id := session.Get("user_id").(uint)
		user := authModel.GetUser(user_id)

		c.HTML(200, "verified.tmpl", gin.H{
			"Name": user.Name,
		})
	} else {
		c.Redirect(301, "/login")
	}
}
func Index(c *gin.Context) {
	if IsAuth(c) {
		c.Redirect(301, "/")
	} else {
		c.HTML(200, "login.tmpl", gin.H{})
	}
}

func Login(c *gin.Context) {
	id := c.PostForm("id")
	pass := c.PostForm("pass")

	v := authModel.Verify(id, pass)

	if v {
		user := authModel.GetUserWhereName(id)
		session := sessions.Default(c)

		session.Clear()
		session.Save()
		session.Set("user_id", user.ID)
		session.Save()
		c.Redirect(301, "/")
	} else {
		c.Redirect(301, "/login")
	}
}

func Create(c *gin.Context) {
	c.HTML(200, "register.tmpl", gin.H{})
}

func Register(c *gin.Context) {
	id := c.PostForm("id")
	pass := c.PostForm("pass")

	authModel.Create(id, pass)

	user := authModel.GetUserWhereName(id)
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()
	c.Redirect(301, "/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(301, "/")
}
