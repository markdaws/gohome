package cmd

type Builder interface {
	Build(Command) (*Func, error)

	//Why does this have an ID? Can't we just use ModelNumber
	ID() string
}
