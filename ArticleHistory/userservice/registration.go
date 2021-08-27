package userservice

import (
	"net/http"
	"strings"
	"unicode"

	emailservice "../emailservice"
	schema "../schema"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword - Generates a Hash value for the passord
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CreateNewUser - Saves new user in DB
func CreateNewUser(firstname string, lastname string, email string, hash string) bool {
	user := schema.User{FirstName: firstname, LastName: lastname, UserName: email, Email: email, PasswordHash: hash}

	// => returns `true` as primary key is blank
	newUserFlag := schema.DBConn.NewRecord(user)

	if newUserFlag {
		// Check if returns RecordNotFound error
		if schema.DBConn.Where("email = ?", email).First(&user).RecordNotFound() {
			schema.DBConn.Create(&user)
			return true
		}
		return false
	}
	return newUserFlag
}

// IsLetter - checks if a string only contains alphabets or not.
func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// UserRegistration - Registers a user account
func UserRegistration(c *gin.Context) {
	var Status string
	var Message string

	statusCode := http.StatusBadRequest

	if c.Request.Method == "POST" {
		firstname := c.PostForm("firstname")
		lastname := c.PostForm("lastname")
		email := c.PostForm("email")
		password := c.PostForm("password")

		hash, _ := HashPassword(password)

		if len(email) == 0 {
			Status = "fail"
			Message = "Email expected."
		} else if len(password) == 0 {
			Status = "fail"
			Message = "Password expected."
		} else if len(email) < 8 || len(email) > 64 || strings.Count(email, "@") < 1 || strings.Count(email, ".") < 1 {
			Status = "fail"
			Message = "Invalid email address."
		} else if len(password) < 6 {
			Status = "fail"
			Message = "Password must be minimum of length 6."
		} else if len(password) > 15 {
			Status = "fail"
			Message = "Password must be maximum of length 15."
		} else if !IsLetter(firstname) {
			Status = "fail"
			Message = "Only Alphabets are allowed in Firstname."
		} else if !IsLetter(lastname) {
			Status = "fail"
			Message = "Only Alphabets are allowed in Lastname"
		} else if len(firstname) < 4 {
			Status = "fail"
			Message = "Firstname must be minimum of length 4."
		} else if len(firstname) > 15 {
			Status = "fail"
			Message = "Firstname must be maximum of length 15."
		} else if len(lastname) < 1 {
			Status = "fail"
			Message = "Lastname must be minimum of length 1."
		} else if len(lastname) > 15 {
			Status = "fail"
			Message = "Lastname must be maximum of length 15."
		} else {
			userFlag := CreateNewUser(firstname, lastname, email, hash)
			if userFlag {
				Status = "success"
				Message = "User registration successful."
				statusCode = http.StatusOK

				// Send Welcome Email
				emailservice.EmailSend(firstname+" "+lastname, email)

			} else {
				Status = "fail"
				Message = "User already exists."
			}
		}

	} else {
		Status = "fail"
		Message = "Invalid request."
		statusCode = http.StatusNotFound
	}
	c.JSON(statusCode, gin.H{"status": Status, "message": Message})
}
