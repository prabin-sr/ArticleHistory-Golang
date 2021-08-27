package configurations

// DatabaseConfiguration - Struct for Configurations
type DatabaseConfiguration struct {
	DBServer   string
	DBName     string
	DBUser     string
	DBPassword string
}

// GetDatabaseConfigurations - Retrieves the database data
func GetDatabaseConfigurations() DatabaseConfiguration {
	databaseconfiguration := DatabaseConfiguration{}

	// Change below fields as per your database setup
	databaseconfiguration.DBServer = "mysql"
	databaseconfiguration.DBName = "db_name"
	databaseconfiguration.DBUser = "db_user_name"
	databaseconfiguration.DBPassword = "db_password"

	return databaseconfiguration
}
