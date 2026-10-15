package config

// EndpointFields defines custom endpoint and addressing options for object storage.
type EndpointFields struct {
	Endpoint     string `env:"OHM9996_STORAGE_ENDPOINT" envDefault:""`
	BucketURL    string `env:"OHM9996_STORAGE_BUCKET_URL" envDefault:""`
	Port         int    `env:"OHM9996_STORAGE_PORT" envDefault:"0"`
	UsePathStyle bool   `env:"OHM9996_STORAGE_USE_PATH_STYLE" envDefault:"false"`
	Insecure     bool   `env:"OHM9996_STORAGE_INSECURE" envDefault:"false"`
	DisableSSL   bool   `env:"OHM9996_STORAGE_DISABLE_SSL" envDefault:"false"`
}
