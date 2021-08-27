package configurations

// EmailConfiguration - Struct for Configurations
type EmailConfiguration struct {
	FromEmail  string
	Password   string
	MailServer string
	SMTPType   string
}

// GetEmailConfigurations - Retrieves the email data
func GetEmailConfigurations() EmailConfiguration {
	emailconfiguration := EmailConfiguration{}

	// Change below fields as per your SMTP server setup
	emailconfiguration.FromEmail = "example@gmail.com"
	emailconfiguration.Password = "ExamplePassword"
	emailconfiguration.MailServer = "smtp.gmail.com:587"
	emailconfiguration.SMTPType = "smtp.gmail.com"

	return emailconfiguration
}
