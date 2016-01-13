package comm

type Authenticator interface {
	Authenticate(Connection) error
}
