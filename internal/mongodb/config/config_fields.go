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

// Config defines configuration options for connecting to MongoDB.
type Config struct {
	Host               string `env:"OHM9996_MONGODB_HOST"`
	Username           string `env:"OHM9996_MONGODB_USERNAME"`
	Password           string `env:"OHM9996_MONGODB_PASSWORD"`
	Protocol           string `env:"OHM9996_MONGODB_PROTOCOL" envDefault:"mongodb"`
	UUIDRepresentation string `env:"OHM9996_MONGODB_UUID_REPRESENTATION" envDefault:"unspecified"`
	Port               int    `env:"OHM9996_MONGODB_PORT"`
}
