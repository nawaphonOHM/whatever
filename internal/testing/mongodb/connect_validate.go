package mongodb

const (
	minValidPort = 1
	maxValidPort = 65535
)

// isPortValid checks whether port is within the allowable TCP range.
func isPortValid(port int) bool {
	return port >= minValidPort && port <= maxValidPort
}

// isPortIgnored reports whether port validation should be skipped for this option set.
func isPortIgnored(o *Options) bool {
	return o.Protocol == "mongodb+srv" || o.Port == 0
}

// validatePort checks that the port is within the valid range (1-65535).
func validatePort(o *Options) error {
	if isPortIgnored(o) {
		return nil
	}
	if !isPortValid(o.Port) {
		return ErrInvalidPort
	}
	return nil
}

// validateOptions checks validity of connection options before attempting connection.
func validateOptions(o *Options) error {
	if o == nil {
		return ErrNilConfig
	}
	if o.URI != "" {
		return nil
	}
	return validatePort(o)
}
