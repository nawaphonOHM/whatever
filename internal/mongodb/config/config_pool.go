package config

// PoolFields defines connection pooling capacity bounds for MongoDB.
type PoolFields struct {
	MaxPoolSize uint64 `env:"OHM9996_MONGODB_MAX_POOL_SIZE" envDefault:"100"`
	MinPoolSize uint64 `env:"OHM9996_MONGODB_MIN_POOL_SIZE" envDefault:"5"`
}
