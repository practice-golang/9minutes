package email

var Info = Config{
	UseEmail:   false,
	Domain:     "http://localhost:1234",
	SendDirect: true,
}

func SendVerificationEmail(message Message) error {
	err := SendDirect(message)
	if err != nil {
		return err
	}

	return nil
}
