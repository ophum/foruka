package main

import (
	_ "net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	auth "github.com/ophum/foruka/controllers/authController"
	cont "github.com/ophum/foruka/controllers/containerController"
	home "github.com/ophum/foruka/controllers/homeController"
)

func main() {
	r := gin.Default()

	store := sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 86400})

	r.Use(sessions.Sessions("session", store))

	r.LoadHTMLGlob("views/templates/*/**")

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
	r.Run()
}
