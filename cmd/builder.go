package cmd

type Builder interface {
	Build(Command) (*Func, error)
	ID() string
}
