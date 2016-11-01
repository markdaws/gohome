package belkin

// BinaryState represents the BinaryState information returned from a device
type BinaryState struct {
	OnOff   int
	OnSince int64
}

// TODO: decipher all the parameters from the response
// <BinaryState>1|1477978435|0|0|0|1168438|0|100|0|0</BinaryState>
