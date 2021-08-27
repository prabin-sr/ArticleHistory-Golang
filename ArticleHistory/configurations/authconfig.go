package configurations

// AuthConfiguration - Struct for Configurations
type AuthConfiguration struct {
	JWTCipherKey []byte
}

// GetAuthConfigurations - Retrieves the email data
func GetAuthConfigurations() AuthConfiguration {
	authconfiguration := AuthConfiguration{}

	// Here I am using JWTCipherKey as 1956437603246794 for demonstration purpose
	// You should change this key with your 16 digit random key
	authconfiguration.JWTCipherKey = []byte("1956437603246794")

	return authconfiguration
}
