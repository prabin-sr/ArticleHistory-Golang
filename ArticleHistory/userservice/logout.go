package userservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserLogout - Logouts a user from the application
func UserLogout(c *gin.Context) {
	if c.Request.Method == "GET" {

		statusCode := http.StatusOK
		Status := "success"
		Message := "Logout successful."
		c.JSON(statusCode, gin.H{"status": Status, "message": Message})

	} else {

		statusCode := http.StatusNotFound
		Status := "fail"
		Message := "Logout failed."
		c.JSON(statusCode, gin.H{"status": Status, "message": Message})

	}

}
