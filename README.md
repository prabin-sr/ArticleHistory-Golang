# ArticleHistory
Helps users to get the articles published in NewYork Times newspaper from **September, 1851** to **last month**. The news articles are extracted from NewYork Times through APIs will be used for data-processing and NLP (as per your requirements). This project is built using **Go (GoLang)** web framework **`Gin`** and **MySQL**.


Environment Setup
--------------------
1. Install Go programming language (Go or GoLang).
2. Install any code editor or IDE which has support for Go.
3. Install all dependency packages using `go get`.
    1. `go get -u github.com/gin-gonic/gin`
    2. `go get -u github.com/jinzhu/gorm`
    3. `go get -u github.com/go-sql-driver/mysql`
    4. `go get golang.org/x/crypto/bcrypt`
4. Install `MySQL` or `MariaDB` - (Both will act like same).
5. Configure MySQL/MariaDB with a username and password.
6. Create a database inside the MySQL which is needed to run go http server.
7. Configure NY API (nytimes) credentials in `configurations/static/nyapiconf.json` file.


Database Setup
--------------------
1. Download and install `MariaDB` or `MySQL`.
2. Create an user with a strong password.
3. Create a database with a suitable name.
4. Update the database configuration file with the username and password.


EmailServer Setup
----------------------
1. Create a email-id from any email service.
2. Grant permission to the email to use it from any 3rd party application.
2. Update the email configuration file.


Server Setup
----------------------
1. Run Application.
    1. Format: `go run <filename>.go`
    2. Example: `go run main.go`
2. Build Application with dependencies (Creates executables).
    1. Format: `go build <filename>.go` or `go build`
    2. Example: `go build main.go`


API List
---------------------
1. POST   /api/v1/user/register
2. POST   /api/v1/user/login
3. GET    /api/v1/user/logout
4. GET    /api/v1/user/details
5. GET    /api/v1/archive/details/:year/:month


API Instructions
-----------------------
1. Register an user with `/api/v1/user/register` API call.
2. Login with the registered user details, `/api/v1/user/login` API call.
3. You will get Auth Token when login.
4. Set the auth token request header `Authorization: Bearer YourToken` for all API calls, after login.
5. API call `/api/v1/user/details` is for getting the details of logged in user.
6. To get the list of all articles link and details for a particular year and month call `/api/v1/archive/details/:year/:month` API.
7. To logout from the application, use `/api/v1/user/logout` API call.


Help
----------------------
1. Don't fotgot to change the configurations for the application in `configurations` directory.
1. `https://golang.org/dl/`
2. `https://code.visualstudio.com/#alt-downloads`
3. `https://downloads.mariadb.org/`
4. `https://golang.org/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them`


Licence
-----------------------
This software is licenced under *GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007*.
