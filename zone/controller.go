package zone

type Controller string

const (
	ZCDefault           Controller = ""
	ZCFluxWIFI                     = "FluxWIFI"
	ZCWeMoInsightSwitch            = "WeMoInsightSwitch"
)

func ControllerFromString(c string) Controller {
	switch c {
	case "FluxWIFI":
		return ZCFluxWIFI
	case "WeMoInsightSwitch":
		return ZCWeMoInsightSwitch
	default:
		return ZCDefault
	}
}
