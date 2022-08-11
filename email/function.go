package email

import (
	"log"
	"net/smtp"
)

func SendDirect(message Message) (err error) {
	from := message.From
	to := message.To
	subject := message.Subject
	body := message.Body

	headerFrom := from.Name + " <" + from.Email + ">"
	headerTo := to.Name + " <" + to.Email + ">"
	if !message.AppendFromToName {
		headerFrom = " <" + from.Email + ">"
		headerTo = " <" + to.Email + ">"
	}

	msg := []byte{}
	switch message.BodyType {
	case TEXT:
		msg = ComposeMimeMailTEXT(headerTo, headerFrom, subject, body)
	case HTML:
		msg = ComposeMimeMailHTML(headerTo, headerFrom, subject, body)
	}

	mx, err := GetMX(to.Email)
	if err != nil {
		log.Println("MX:", err)
		return err
	}

	domain, err := getDomain(from.Email)
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

	err = smtp.SendMail(mx+":25", nil, from.Email, []string{to.Email}, signedMessage)
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

	// from/to name will be set as email address when using gmail
	// Error caused when an email address using which different with id - The sender address is unauthorized
	from := id
	to := []string{message.To.Email}
	fromName := message.From.Name
	toName := message.To.Name
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
