package comm

type ConnectionInfo interface {
}

type TelnetConnectionInfo struct {
	PoolSize int

	//TODO: Remove
	Login         string
	Password      string
	Authenticator Authenticator

	Network string
	Address string
}

type HTTPConnectionInfo struct {
	PoolSize int
	HostName string
	Port     string
}
