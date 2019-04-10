package authController

import (
  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/contrib/sessions"
  "github.com/ophum/foruka/models/authModel"
)
func Verified(c *gin.Context) {
  session := sessions.Default(c)
  id := session.Get("id")
  if id != nil {
    c.HTML(200, "verified.tmpl", gin.H{
      "Id": id,
    })
  }else {
    c.Redirect(301, "/login")
  }
}
func Index(c *gin.Context) {
  c.HTML(200, "login.tmpl", gin.H{})
}

func Login(c *gin.Context) {
  id := c.PostForm("id")
  pass := c.PostForm("pass")

  v := authModel.Verify(id, pass)

  if v {
    session := sessions.Default(c)
    session.Set("id", id)
    session.Save()
    c.Redirect(301, "/verified")
  }else {
    c.Redirect(301, "/login")
  }
}

func Create(c *gin.Context) {
  c.HTML(200,"register.tmpl", gin.H{})
}

func Register(c *gin.Context) {
  id := c.PostForm("id")
  pass := c.PostForm("pass")

  authModel.Create(id, pass)

  c.Redirect(301, "/login")
}
