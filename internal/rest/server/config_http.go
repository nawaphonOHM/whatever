package server

// HTTPFields defines HTTP routing and path handling options.
type HTTPFields struct {
	RedirectTrailingSlash  bool `env:"OHM9996_SERVER_REDIRECT_TRAILING_SLASH" envDefault:"true"`
	RedirectFixedPath      bool `env:"OHM9996_SERVER_REDIRECT_FIXED_PATH" envDefault:"false"`
	HandleMethodNotAllowed bool `env:"OHM9996_SERVER_HANDLE_METHOD_NOT_ALLOWED" envDefault:"true"`
	UseRawPath             bool `env:"OHM9996_SERVER_USE_RAW_PATH" envDefault:"false"`
	UnescapePathValues     bool `env:"OHM9996_SERVER_UNESCAPE_PATH_VALUES" envDefault:"true"`
	RemoveExtraSlash       bool `env:"OHM9996_SERVER_REMOVE_EXTRA_SLASH" envDefault:"false"`
}

func defaultHTTPFields() HTTPFields {
	return HTTPFields{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: true,
		UseRawPath:             false,
		UnescapePathValues:     true,
		RemoveExtraSlash:       false,
	}
}
