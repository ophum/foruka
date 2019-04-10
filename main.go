package main

import(
  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/contrib/sessions"
  _ "net/http"

  auth "github.com/ophum/foruka/controllers/authController"
)

func main() {
  r := gin.Default()

  store := sessions.NewCookieStore([]byte("secret"))
  r.Use(sessions.Sessions("session", store))

  r.LoadHTMLGlob("views/templates/*")

  r.GET("/", func(c *gin.Context) {
    c.HTML(200, "index.tmpl", gin.H{
    })
  })

  r.GET("/register", auth.Create)
  r.POST("/register", auth.Register)

  r.GET("/login", auth.Index)
  r.POST("/login", auth.Login)

  r.GET("/verified", auth.Verified)
  r.Run()
}
