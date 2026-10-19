package config

import "time"

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

// dur shortens time.Duration lines for struct tags.
type dur = time.Duration

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

// TimeoutFields defines operation and lifecycle timeout settings for MongoDB.
type TimeoutFields struct {
	ConnectTimeout         dur `env:"OHM9996_MONGODB_CONNECT_TIMEOUT" envDefault:"10s"`
	ServerSelectionTimeout dur `env:"OHM9996_MONGODB_SERVER_SELECTION_TIMEOUT" envDefault:"5s"`
	SocketTimeout          dur `env:"OHM9996_MONGODB_SOCKET_TIMEOUT" envDefault:"10s"`
	MaxConnIdleTime        dur `env:"OHM9996_MONGODB_MAX_CONN_IDLE_TIME" envDefault:"10m"`
}

// PoolFields defines connection pooling capacity bounds for MongoDB.
type PoolFields struct {
	MaxPoolSize uint64 `env:"OHM9996_MONGODB_MAX_POOL_SIZE" envDefault:"100"`
	MinPoolSize uint64 `env:"OHM9996_MONGODB_MIN_POOL_SIZE" envDefault:"5"`
}
