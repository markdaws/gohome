package comm

type ConnectionInfo interface {
}

type TelnetConnectionInfo struct {
	PoolSize int

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
