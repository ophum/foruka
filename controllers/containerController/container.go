package containerController

import (
	"github.com/gin-gonic/gin"
	auth "github.com/ophum/foruka/controllers/authController"
)

func Index(c *gin.Context) {
  auth.Auth(c)

  c.HTML(200, "containers/index.tmpl", gin.H{"Title": "foruka", "C": "hoge"})
}
