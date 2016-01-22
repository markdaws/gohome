package comm

type ConnectionInfo interface {
}

//TODO: Delete?
type TelnetConnectionInfo struct {
	PoolSize int

	//TODO: Remove
	Login         string
	Password      string
	Authenticator Authenticator

	Network string
	Address string
}
