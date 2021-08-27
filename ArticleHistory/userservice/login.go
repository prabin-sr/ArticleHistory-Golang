package userservice

import (
	"encoding/json"
	"net/http"
	"time"

	configurations "../configurations"
	schema "../schema"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash - Returns true if hash and password matches else returms false
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWTData - Struct to store data passed through auth request
type JWTData struct {
	ID               uint
	FirstName        string
	LastName         string
	Email            string
	IsActive         bool
	IsAdmin          bool
	TokenGeneratedAt *time.Time
}

// UserLogin - Logins the user if the user is authorised.
func UserLogin(c *gin.Context) {
	var Status string
	var Message string

	var statusCode int

	if c.Request.Method == "POST" {

		email := c.PostForm("username")
		password := c.PostForm("password")

		var jwtToken string

		if email == "" {

			Status = "fail"
			Message = "Username field required."
			statusCode = http.StatusUnauthorized

		} else if password == "" {

			Status = "fail"
			Message = "Password field required."
			statusCode = http.StatusUnauthorized

		} else {

			userStruct := schema.User{}

			if schema.DBConn.Where(schema.User{Email: email}).First(&userStruct).RecordNotFound() {

				Status = "fail"
				Message = "User account not found. Please register."
				statusCode = http.StatusUnauthorized

			} else {

				schema.DBConn.Where(schema.User{Email: email}).First(&userStruct)

				match := CheckPasswordHash(password, userStruct.PasswordHash)

				if match {
					authConfiguration := configurations.GetAuthConfigurations()

					cipherKey := authConfiguration.JWTCipherKey
					currentTime := time.Now()

					jwtStruct := JWTData{}
					jwtStruct.ID = userStruct.ID
					jwtStruct.FirstName = userStruct.FirstName
					jwtStruct.LastName = userStruct.LastName
					jwtStruct.Email = userStruct.Email
					jwtStruct.IsActive = userStruct.IsActive
					jwtStruct.IsAdmin = userStruct.IsAdmin
					jwtStruct.TokenGeneratedAt = &currentTime

					messageBytes, err := json.Marshal(jwtStruct)

					if err != nil {

						panic(err)

					}

					messageString := string(messageBytes)

					// Encrypt the user details
					encrypted, err := EncryptToken(cipherKey, messageString)

					if err != nil {

						Status = "fail"
						Message = "Something went wrong."
						statusCode = http.StatusUnauthorized

					} else {

						jwtToken = encrypted

						Status = "success"
						Message = "Login Success."
						statusCode = http.StatusOK

					}

				} else {

					Status = "fail"
					Message = "Invalid password. Please enter the actual password."
					statusCode = http.StatusUnauthorized

				}
			}
		}

		if jwtToken != "" {

			c.JSON(statusCode, gin.H{"status": Status, "message": Message, "token": jwtToken})

		} else {

			c.JSON(statusCode, gin.H{"status": Status, "message": Message})

		}

	} else {

		Status = "fail"
		Message = "Invalid request."
		statusCode = http.StatusNotFound
		c.JSON(statusCode, gin.H{"status": Status, "message": Message})

	}
}
