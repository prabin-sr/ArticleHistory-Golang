package userservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserDetails - Returns the basic user details
func UserDetails(c *gin.Context) {
	userData, statusID := c.Get("User")

	if statusID == false {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return

	}

	userStruct := userData.(JWTData)

	c.JSON(http.StatusOK, gin.H{"first_name": userStruct.FirstName, "last_name": userStruct.LastName, "email": userStruct.Email})

}
