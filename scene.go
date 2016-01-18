package gohome

type Scene struct {
	//TODO: Why does scene need a local and global ID now? Map back to local device
	LocalID     string
	GlobalID    string
	Name        string
	Description string
	Commands    []Command
}
