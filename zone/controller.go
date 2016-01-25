package zone

type Controller string

const (
	ZCDefault  Controller = ""
	ZCFluxWIFI            = "FluxWIFI"
)

func ControllerFromString(c string) Controller {
	switch c {
	case "FluxWIFI":
		return ZCFluxWIFI
	default:
		return ZCDefault
	}
}
