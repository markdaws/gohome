package cmd

// Builder is an interface indicating a type supports taking an abstract Command
// like ZoneSetLevel and returning a Func that knows how to execute that command
// on a specific type of hardware
type Builder interface {
	Build(Command) (*Func, error)
}
