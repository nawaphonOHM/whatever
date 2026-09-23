package server

// ProxyFields defines upstream proxy and forwarded header settings.
type ProxyFields struct {
	TrustedProxies      []string `env:"OHM9996_SERVER_TRUSTED_PROXIES"`
	RemoteIPHeaders     []string `env:"OHM9996_SERVER_REMOTE_IP_HEADERS"`
	ForwardedByClientIP bool     `env:"OHM9996_SERVER_FORWARDED_BY_CLIENT_IP" envDefault:"true"`
}

func defaultProxyFields() ProxyFields {
	return ProxyFields{
		ForwardedByClientIP: true,
		RemoteIPHeaders:     []string{"X-Forwarded-For", "X-Real-IP"},
		TrustedProxies:      nil,
	}
}
