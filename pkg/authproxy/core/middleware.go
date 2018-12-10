package core

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/glog"

	"github.com/syunkitada/goapp/pkg/authproxy/authproxy_model"
	"github.com/syunkitada/goapp/pkg/lib/logger"
)

func (authproxy *Authproxy) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		client := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		traceId := logger.NewTraceId()

		c.Set("TraceId", traceId)
		logger.TraceInfo(traceId, authproxy.host, authproxy.name, map[string]string{
			"Msg":    "Start",
			"Client": client,
			"Method": method,
			"Path":   path,
		})

		c.Next()
		end := time.Now()
		latency := end.Sub(start)

		statusCode := c.Writer.Status()

		logger.TraceInfo(traceId, authproxy.host, authproxy.name, map[string]string{
			"Msg":       "End",
			"Client":    client,
			"Method":    method,
			"Path":      path,
			"StausCode": strconv.Itoa(statusCode),
			"Latency":   strconv.FormatInt(latency.Nanoseconds()/1000000, 10),
		})
	}
}

// SecureHeaders adds secure headers to the API
// func (a *API) SecureHeaders(next http.Handler) http.Handler {
// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
func (authproxy *Authproxy) ValidateHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check AllowedHosts
		var err error
		if len(authproxy.AllowedHosts) > 0 {
			isGoodHost := false
			for _, allowedHost := range authproxy.AllowedHosts {
				if strings.EqualFold(allowedHost, c.Request.Host) {
					isGoodHost = true
					break
				}
			}
			if !isGoodHost {
				c.JSON(http.StatusForbidden, gin.H{
					"err": fmt.Sprintf("Bad host name: %s", c.Request.Host),
				})
				c.Abort()
				return
			}
		}
		// If there was an error, do not continue request
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": fmt.Sprintf("Failed to check allowed hosts"),
			})
			c.Abort()
			return
		}

		// Add X-XSS-Protection header
		// Enables XSS filtering. Rather than sanitizing the page, the browser will prevent rendering of the page if an attack is detected.
		c.Writer.Header().Add("X-XSS-Protection", "1; mode=block")

		// Add Content-Type header
		// Content type tells the browser what type of content you are sending. If you do not include it, the browser will try to guess the type and may get it wrong.
		// w.Header().Add("Content-Type", "application/json")

		// Add X-Content-Type-Options header
		// Content Sniffing is the inspecting the content of a byte stream to attempt to deduce the file format of the data within it.
		// Browsers will do this to try to guess at the content type you are sending.
		// By setting this header to “nosniff”, it prevents IE and Chrome from content sniffing a response away from its actual content type. This reduces exposure to drive-by download attacks.
		c.Writer.Header().Add("X-Content-Type-Options", "nosniff")

		// Prevent page from being displayed in an iframe
		c.Writer.Header().Add("X-Frame-Options", "DENY")

		// Allow Origin
		c.Writer.Header().Add("Access-Control-Allow-Origin", "http://192.168.10.103:3000")
		c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
	}
}

func (authproxy *Authproxy) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenAuthRequest authproxy_model.TokenAuthRequest

		if err := c.ShouldBindWith(&tokenAuthRequest, binding.JSON); err != nil {
			glog.Warningf("Invalid TokenAuthRequest: Failed ShouldBindJSON: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "Invalid AuthRequest",
			})
			c.Abort()
			return
		}

		value, cookieErr := c.Cookie("token")
		if cookieErr == nil {
			tokenAuthRequest.Token = value
			glog.Info(tokenAuthRequest.Token)
		}

		claims, err := authproxy.Token.ParseToken(tokenAuthRequest)
		if err != nil {
			glog.Warning("Invalid AuthRequest: Failed ParseToken")
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "Invalid AuthRequest",
			})
			c.Abort()
			return
		}

		username := claims["Username"].(string)
		userAuthority, getUserAuthorityErr := authproxy.AuthproxyModelApi.GetUserAuthority(username, &tokenAuthRequest.Action)
		if getUserAuthorityErr != nil {
			glog.Error(getUserAuthorityErr)
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "Invalid AuthRequest",
			})
			c.Abort()
			return
		}

		c.Set("Username", claims["Username"])
		c.Set("UserAuthority", userAuthority)
		c.Set("Action", tokenAuthRequest.Action)
	}
}
