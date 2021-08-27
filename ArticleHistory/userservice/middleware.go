package userservice

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"
	configurations "../configurations"
	"github.com/gin-gonic/gin"
)

// EncryptToken - Encrypts the user token
func EncryptToken(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

// DecryptToken - Decrypts the user token
func DecryptToken(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return
}

// RespondWithError - Error response returner
func RespondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

// UserAuthentication - middleware to check user authentication
func UserAuthentication() gin.HandlerFunc {

	return func(c *gin.Context) {

		authorization := c.Request.Header.Get("Authorization")

		if authorization == "" {

			RespondWithError(c, 401, "Authorization header required.")
			return

		}

		bearerAuth := strings.Fields(authorization)

		if len(bearerAuth) <= 1 || len(bearerAuth) > 2 {

			RespondWithError(c, 401, "Bearer token required.")
			return

		}

		if bearerAuth[0] != "Bearer" {

			RespondWithError(c, 401, "Bearer token invalid.")
			return

		}

		userStruct := JWTData{}
		authConfiguration := configurations.GetAuthConfigurations()

		cipherKey := authConfiguration.JWTCipherKey
		jwtToken := bearerAuth[1]

		// Decrypt the token to get the user details
		decryptedUserString, err := DecryptToken(cipherKey, jwtToken)

		if err != nil {

			RespondWithError(c, 401, "Invalid authentication token.")
			return

		}

		decryptedUserByte := []byte(decryptedUserString)

		err = json.Unmarshal(decryptedUserByte, &userStruct)
		if err != nil {
			panic(err)
		}

		t1Year, t1Month, t1Day := time.Now().Date()
		t2Year, t2Month, t2Day := userStruct.TokenGeneratedAt.Date()

		t1 := time.Date(t1Year, t1Month, t1Day, 0, 0, 0, 0, time.Local)
		t2 := time.Date(t2Year, t2Month, t2Day, 0, 0, 0, 0, time.Local)

		days := t2.Sub(t1).Hours() / 24

		// Remove login credentials after 7 days.
		// Unauthorize.
		if days >= 7 {

			RespondWithError(c, 401, "Invalid authentication token.")
			return

		}

		c.Set("User", userStruct)

		c.Next()

	}

}
