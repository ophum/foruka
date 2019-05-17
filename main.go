package main

import (
	_ "net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	auth "github.com/ophum/foruka/controllers/authController"
	cont "github.com/ophum/foruka/controllers/containerController"
	home "github.com/ophum/foruka/controllers/homeController"
	networks "github.com/ophum/foruka/controllers/networkController"
)

func main() {
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 86400})

	r.Use(sessions.Sessions("session", store))

	r.LoadHTMLGlob("views/templates/*/**")
	r.Static("/assets", "./assets")

	r.GET("/", home.Index)

	r.GET("/register", auth.Create)
	r.POST("/register", auth.Register)

	r.GET("/login", auth.Index)
	r.POST("/login", auth.Login)

	r.GET("/verified", auth.Verified)

	r.GET("/logout", auth.Logout)

	containers := r.Group("/containers")
	containers.GET("/", cont.Index)
	containers.GET("/create", cont.Create)
	containers.GET("/store", cont.Create)
	containers.POST("/store", cont.Store)
	containers.GET("/show/:id", cont.Show)
	containers.GET("/stop/:id", cont.Stop)
	containers.GET("/start/:id", cont.Start)
	containers.GET("/delete/:id", cont.Delete)

	net := r.Group("/networks")
	net.GET("/", networks.Index)
	net.GET("/create", networks.Create)

	r.Run()
}
