package config

// Protocol constants for MongoDB connection schemes.
const (
	ProtocolMongoDB    = "mongodb"
	ProtocolMongoDBSrv = "mongodb+srv"
)

// UUID representation constants for MongoDB binary UUID encoding.
const (
	UUIDRepresentationUnspecified  = "unspecified"
	UUIDRepresentationStandard     = "standard"
	UUIDRepresentationCSharpLegacy = "csharpLegacy"
	UUIDRepresentationJavaLegacy   = "javaLegacy"
	UUIDRepresentationPythonLegacy = "pythonLegacy"
)

// BaseFields defines core connection endpoint, authentication, and application metadata.
type BaseFields struct {
	Host               string `env:"OHM9996_MONGODB_HOST"`
	Protocol           string `env:"OHM9996_MONGODB_PROTOCOL" envDefault:"mongodb"`
	Database           string `env:"OHM9996_MONGODB_DATABASE" envDefault:""`
	Username           string `env:"OHM9996_MONGODB_USERNAME" envDefault:""`
	Password           string `env:"OHM9996_MONGODB_PASSWORD" envDefault:""`
	AuthSource         string `env:"OHM9996_MONGODB_AUTH_SOURCE" envDefault:""`
	AppName            string `env:"OHM9996_MONGODB_APP_NAME" envDefault:""`
	UUIDRepresentation string `env:"OHM9996_MONGODB_UUID_REPRESENTATION" envDefault:"unspecified"`
	Port               int    `env:"OHM9996_MONGODB_PORT" envDefault:"27017"`
}

// Config defines configuration options for connecting to MongoDB.
type Config struct {
	BaseFields
	TimeoutFields
	PoolFields
}
