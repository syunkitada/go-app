package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func (authproxy *Authproxy) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Health",
	})
}

func (authproxy *Authproxy) AuthHealth(c *gin.Context) {
	userName, userNameOk := c.Get("UserName")
	roleName, roleNameOk := c.Get("RoleName")
	projectName, projectNameOk := c.Get("ProjectName")
	projectRoleName, projectRoleNameOk := c.Get("ProjectRoleName")
	if !userNameOk || !roleNameOk || !projectNameOk || !projectRoleNameOk {
		glog.Error("Success AuthHealth: userName(%v), roleName(%v), projectName(%v), projectRoleName(%v)", userNameOk, roleNameOk, projectNameOk, projectRoleNameOk)
		c.JSON(500, gin.H{
			"message": "Invalid request",
		})
		return
	}

	glog.Info("Success AuthHealth: userName(%v), roleName(%v), projectName(%v), projectRoleName(%v)", userName, roleName, projectName, projectRoleName)

	c.JSON(200, gin.H{
		"message": "Health",
	})
}

func (authproxy *Authproxy) HealthGrpc(c *gin.Context) {
	status, err := authproxy.HealthClient.Status()
	if err != nil {
		glog.Error("Failed HealthClient.Status", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid AuthRequest",
		})
		c.Abort()
	}
	glog.Info(status)

	c.JSON(200, gin.H{
		"message": "Health",
	})
}
