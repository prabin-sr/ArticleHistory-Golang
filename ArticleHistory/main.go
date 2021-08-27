package main

import (
	"fmt"
	"strconv"

	nyapiservice "./nyapiservice"
	archivespider "./nyapiservice/archiveapi/spiders"
	archivewrangler "./nyapiservice/archiveapi/wranglers"
	userservice "./userservice"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initiates Archive API Meta data crawler
	// Run as a GoRoutine - Cuncurrent process
	go archivespider.NYMetaDataGrabber()

	// Initiates Archive Corpus Cleaner
	// Run as a GoRoutine - Cuncurrent process
	go archivewrangler.ArchiveMetaCleaner()

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apiVersion := router.Group("api/v1")

	users := apiVersion.Group("/user")
	archives := apiVersion.Group("/archive")

	// User related APIs
	users.POST("/register", userservice.UserRegistration)
	users.POST("/login", userservice.UserLogin)
	users.GET("/logout", userservice.UserLogout)
	users.Use(userservice.UserAuthentication())
	{
		users.GET("/details", userservice.UserDetails)
	}

	// NY Archive API related APIs
	archives.Use(userservice.UserAuthentication())
	{
		archives.GET("/details/:year/:month", nyapiservice.GetArchiveMetaData)
	}

	/*
		// Asking admin to enter the port number to run the server
		scanner := bufio.NewScanner(os.Stdin)

		var port string
		fmt.Print("\n\nEnter PORT Number: ")
		if scanner.Scan() {
			port = scanner.Text()
		}
	*/
	port := "3000"

	_, err := strconv.Atoi(port)
	if err == nil && len(port) > 2 && len(port) < 5 {
		fmt.Printf("\n%q is the selected port.\n\n", port)
	} else {
		fmt.Print("\nYou have entered an invalid port number.\nTrying to run Golang server in default port.\n\n")
		port = "3000"
	}

	router.Run("127.0.0.1:" + port)
}
