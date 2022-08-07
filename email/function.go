package email

import (
	"log"
	"net/smtp"
)

func SendDirect(message Message) (err error) {
	from := message.From
	to := message.To
	fromName := message.FromName
	toName := message.ToName
	subject := message.Subject
	body := message.Body

	headerFrom := fromName + " <" + from + ">"
	headerTo := toName + " <" + to + ">"
	if !message.AppendFromToName {
		headerFrom = " <" + from + ">"
		headerTo = " <" + to + ">"
	}

	msg := []byte{}
	switch message.BodyType {
	case TEXT:
		msg = ComposeMimeMailTEXT(headerTo, headerFrom, subject, body)
	case HTML:
		msg = ComposeMimeMailHTML(headerTo, headerFrom, subject, body)
	}

	mx, err := GetMX(to)
	if err != nil {
		log.Println("MX:", err)
		return err
	}

	domain, err := getDomain(from)
	if err != nil {
		log.Println("Domain:", err)
		return err
	}

	signedMessage := msg
	if message.Service.KeyDKIM != "" {
		signedMessage, err = GetSignedMessage(msg, domain, message.Service.KeyDKIM)
		if err != nil {
			log.Println("SignedMessage:", err)
			return err
		}
	}

	err = smtp.SendMail(mx+":25", nil, from, []string{to}, signedMessage)
	if err != nil {
		log.Println("SendMail:", mx)
		log.Println("SendMail:", err)
		return err
	}

	// Message with DKIM validation
	// _, err = email.Verify(string(signedMessage))
	// if err != nil {
	// 	return err
	// }

	return nil
}

func SendViaService(message Message) (err error) {
	addr := message.Service.Host
	port := message.Service.Port
	id := message.Service.ID
	password := message.Service.Password

	// 계정과 다른 메일주소로는 변경 안된다
	// gmail은 걍 id로 변환, 네이버는 에러 - id와 다른 주소 사용시 : The sender address is unauthorized
	from := id
	to := []string{message.To}
	fromName := message.FromName
	toName := message.ToName
	subject := "Subject: " + message.Subject + "\r\n"
	body := message.Body + "\r\n"

	headerFrom := "From: " + fromName + " <" + from + ">\r\n"
	headerTo := "To: " + toName + " <" + to[0] + ">\r\n"
	if !message.AppendFromToName {
		headerFrom = "From: " + from + "\r\n"
		headerTo = "To: " + to[0] + "\r\n"
	}

	mime := ""
	switch message.BodyType {
	case TEXT:
		mime = "MIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n"
	case HTML:
		mime = "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"utf-8\"\r\n"
	}

	lineBreak := "\r\n"

	msg := []byte(headerFrom + headerTo + subject + mime + lineBreak + body)

	auth := smtp.PlainAuth("", id, password, addr)
	err = smtp.SendMail(addr+":"+port, auth, from, to, msg)
	if err != nil {
		return err
	}

	return nil
}

func SendVerificationEmail(message Message) error {
	log.Println(message.From, message.FromName)
	log.Println(message.To, message.ToName)
	err := SendDirect(message)
	if err != nil {
		log.Println("Sender:", err)
		return err
	}

	return nil
}
