package comm

type Authenticator interface {
	Authenticate(Connection) error
}

type Auth struct {
	Login         string
	Password      string
	Token         string
	Authenticator Authenticator
}
