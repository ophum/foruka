package main

import (
	_ "net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	auth "github.com/ophum/foruka/controllers/authController"
	home "github.com/ophum/foruka/controllers/homeController"
)

func main() {
	r := gin.Default()

	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	r.LoadHTMLGlob("views/templates/*/**")

	r.GET("/", home.Index)

	r.GET("/register", auth.Create)
	r.POST("/register", auth.Register)

	r.GET("/login", auth.Index)
	r.POST("/login", auth.Login)

	r.GET("/verified", auth.Verified)

	r.GET("/logout", auth.Logout)

	r.GET("/containers/", func(c *gin.Context) {
		c.HTML(200, "containers/index.tmpl", gin.H{})
	})
	r.Run()
}
