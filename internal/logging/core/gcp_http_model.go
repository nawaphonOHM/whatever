package core

const (
	floatBitSize   = 64
	durationSuffix = "s"
)

// HttpRequest defines Google Cloud Logging HTTP request payload schema.
type HttpRequest struct {
	RequestURL    string `json:"requestUrl,omitempty"`
	RequestMethod string `json:"requestMethod,omitempty"`
	Latency       string `json:"latency,omitempty"`
	UserAgent     string `json:"userAgent,omitempty"`
	RemoteIP      string `json:"remoteIp,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	ResponseSize  int64  `json:"responseSize,omitempty"`
	Status        int    `json:"status,omitempty"`
}

func (r *HttpRequest) hasCoreFields() bool {
	return r.RequestURL != "" || r.RequestMethod != "" || r.Status != 0
}

func (r *HttpRequest) hasClientFields() bool {
	return r.UserAgent != "" || r.RemoteIP != ""
}

func (r *HttpRequest) hasNetworkFields() bool {
	return r.Latency != "" || r.Protocol != ""
}

func (r *HttpRequest) hasAuxFields() bool {
	return r.hasClientFields() || r.hasNetworkFields()
}

func (r *HttpRequest) hasAnyFields() bool {
	if r.ResponseSize != 0 {
		return true
	}
	return r.hasCoreFields() || r.hasAuxFields()
}

// IsEmpty reports whether all fields of HttpRequest are unset.
func (r *HttpRequest) IsEmpty() bool {
	if r == nil {
		return true
	}
	return !r.hasAnyFields()
}
