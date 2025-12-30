package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendEmail(email string, message string) {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD") // Для Gmail используйте app password
	to := []string{email}

	smtpHost := os.Getenv("SMTPHOST")
	smtpPort := os.Getenv("SMTPPORT")

	// Формируем сообщение
	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: Notify\r\n" +
		"\r\n" +
		message +
		"\r\n")

	// Авторизация
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Отправка
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Письмо отправлено!")
}
