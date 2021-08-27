package schema

import (
	configurations "../configurations"
	"github.com/jinzhu/gorm"

	// To connect mysql with the go: gin gonic
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DBConn - Global variable of DB connection to access from anywhere.
var DBConn *gorm.DB

func init() {
	DatabaseObj := configurations.GetDatabaseConfigurations()
	dbServer := DatabaseObj.DBServer
	dbName := DatabaseObj.DBName
	dbUser := DatabaseObj.DBUser
	dbPassword := DatabaseObj.DBPassword

	//open a db connection
	var err error
	DBConn, err = gorm.Open(dbServer, dbUser+":"+dbPassword+"@/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	// Disable table name's pluralization globally
	DBConn.SingularTable(true)

	//Migrate the schema
	DBConn.AutoMigrate(&User{})
	DBConn.AutoMigrate(&Archive{})
}
