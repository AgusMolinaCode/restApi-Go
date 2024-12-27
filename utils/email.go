package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	// Configuración del servidor SMTP
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Mensaje
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	// Autenticación
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Enviar correo
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	return err
}
