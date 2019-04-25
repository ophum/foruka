package homeController

import (
	"github.com/gin-gonic/gin"
	auth "github.com/ophum/foruka/controllers/authController"
	// "github.com/gin-gonic/contrib/sessions"
)

func Index(c *gin.Context) {
	if auth.IsAuth(c) {
		auth.Verified(c)
	} else {
		c.HTML(200, "home/index.tmpl", gin.H{"Title": "foruka"})
	}
}
