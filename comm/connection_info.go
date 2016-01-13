package comm

type ConnectionInfo struct {
	Network  string
	Address  string
	Login    string
	Password string
	Stream   bool
	PoolSize int
}
