package gohome

type SensorDataType string

const (
	SDTInt    SensorDataType = "int"
	SDTFloat  SensorDataType = "float"
	SDTBool   SensorDataType = "bool"
	SDTString SensorDataType = "string"
)

type SensorAttr struct {
	Name          string
	Value         string
	DataType      SensorDataType
	UnitOfMeasure string
}

type Sensor struct {
	Attrs []SensorAttr
}
