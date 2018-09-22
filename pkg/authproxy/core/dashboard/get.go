package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func (dashboard *Dashboard) GetState(c *gin.Context) {
	username, usernameOk := c.Get("Username")
	if !usernameOk {
		c.JSON(500, gin.H{
			"error": "Invalid request",
		})
		return
	}

	glog.Info("Success AuthHealth: username(%v)", username)

	c.JSON(200, gin.H{
		"username": username,
	})
}
