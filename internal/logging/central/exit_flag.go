package central

// ExitFlag determines post-logging termination behavior.
type ExitFlag int

const (
	ExitNone ExitFlag = iota
	ExitGraceful
	ExitAbnormal
)

var exitFlagNames = map[ExitFlag]string{
	ExitNone:     "ExitNone",
	ExitGraceful: "ExitGraceful",
	ExitAbnormal: "ExitAbnormal",
}

func (f ExitFlag) String() string {
	if name, ok := exitFlagNames[f]; ok {
		return name
	}
	return "ExitUnknown"
}
