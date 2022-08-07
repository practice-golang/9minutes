package email

type BodyType int

const (
	_             = iota
	TEXT BodyType = 1
	HTML BodyType = 2
)

type Config struct {
	UseEmail   bool
	Domain     string
	SendDirect bool
	Service    Service
	SenderInfo From
}

type From struct {
	Name  string
	Email string
}
type To struct {
	Name  string
	Email string
}

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
	From             From
	To               To
	FromName         string
	ToName           string
	Subject          string
	Body             string
	BodyType         BodyType
}
