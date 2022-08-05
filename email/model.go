package email

type BodyType int

const (
	_             = iota
	TEXT BodyType = 1
	HTML BodyType = 2
)

type Service struct {
	Host     string
	Port     string
	ID       string
	Password string
	KeyDKIM  string
}

type Message struct {
	Service          Service
	AppendFromToName bool
	From             string
	To               string
	FromName         string
	ToName           string
	Subject          string
	Body             string
	BodyType         BodyType
}
