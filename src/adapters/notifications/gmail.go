package notifications

import (
	"fmt"
	"net/smtp"
	"os"
)

// GmailNotifier se encarga de enviar correos usando SMTP
type GmailNotifier struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

// NewGmailNotifier crea un nuevo notificador con variables de entorno
func NewGmailNotifier() *GmailNotifier {
	return &GmailNotifier{
		from:     os.Getenv("GMAIL_FROM"),     // tu-email@gmail.com
		password: os.Getenv("GMAIL_PASSWORD"), // Contraseña de aplicación
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
	}
}

// SendNotification envía un correo electrónico al destinatario
func (g *GmailNotifier) SendNotification(to, subject, body string) error {
	// Autenticación para el servidor SMTP de Gmail
	auth := smtp.PlainAuth("", g.from, g.password, g.smtpHost)

	// Construir el mensaje
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	// Enviar el correo
	addr := fmt.Sprintf("%s:%s", g.smtpHost, g.smtpPort)
	err := smtp.SendMail(addr, auth, g.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("error al enviar correo: %w", err)
	}
	return nil
}
