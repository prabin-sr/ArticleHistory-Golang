package emailservice

import (
	"fmt"
	"net/smtp"

	configurations "../configurations"
)

// EmailSend - Sends an email to the email address
func EmailSend(name string, toAddress string) {
	EmailObj := configurations.GetEmailConfigurations()
	from := EmailObj.FromEmail
	pass := EmailObj.Password

	msg := "From: " + from + "\n" +
		"To: " + toAddress + "\n" +
		"Subject: Welcome to ArticleHistory\n\n" +
		"Hello " + name + ",\n\n" + "Thank you for registering with ArticleHistory."

	err := smtp.SendMail(EmailObj.MailServer,
		smtp.PlainAuth("", from, pass, EmailObj.SMTPType),
		from, []string{toAddress}, []byte(msg))

	if err != nil {
		fmt.Println(err)
		fmt.Println("Error in sending email to " + toAddress)
	}

	fmt.Println("Email has been sent to " + toAddress)
}
