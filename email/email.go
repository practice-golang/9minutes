package email

var Info = Config{
	UseEmail:   false,
	Domain:     "http://localhost:1234",
	SendDirect: true,
}

func SendVerificationEmail(message Message) error {
	if Info.SendDirect {
		err := SendDirect(message)
		if err != nil {
			return err
		}
	} else {
		err := SendViaService(message)
		if err != nil {
			return err
		}
	}

	return nil
}
