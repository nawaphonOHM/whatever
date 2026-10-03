package mongodb

// WithURI sets the direct connection URI string.
func WithURI(uri string) Option {
	return func(o *Options) {
		o.URI = uri
	}
}

// WithHost sets the MongoDB host address.
func WithHost(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

// WithPort sets the MongoDB connection port.
func WithPort(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}

// WithProtocol sets the connection scheme (e.g., mongodb or mongodb+srv).
func WithProtocol(protocol string) Option {
	return func(o *Options) {
		o.Protocol = protocol
	}
}

// WithDatabase sets the default database name for the test client.
func WithDatabase(database string) Option {
	return func(o *Options) {
		o.Database = database
	}
}

// WithUsername sets the authentication username.
func WithUsername(username string) Option {
	return func(o *Options) {
		o.Username = username
	}
}

// WithPassword sets the authentication password.
func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

// WithAuthSource sets the authentication source database.
func WithAuthSource(authSource string) Option {
	return func(o *Options) {
		o.AuthSource = authSource
	}
}

// WithAppName sets the client application name.
func WithAppName(appName string) Option {
	return func(o *Options) {
		o.AppName = appName
	}
}
