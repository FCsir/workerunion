package pkg

import (
	"log"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func SendEmail(content string, subject string, recepients []string) {
	var emailConfig = viper.GetStringMap("email")
	log.Println("jwt config", emailConfig)
	var server string = emailConfig["server"].(string)
	var port int = emailConfig["port"].(int)
	var username string = emailConfig["username"].(string)
	var password string = emailConfig["password"].(string)
	d := gomail.NewDialer(server, port, username, password)

	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", recepients...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	if err := d.DialAndSend(m); err != nil {
		log.Println("---------send error: ", err.Error(), subject)
	}
}
