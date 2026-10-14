package contracts

// Meta holds server configuration metadata such as CORS settings.
type Meta struct {
	corsPath string
}

// NewMeta initializes a new Meta builder.
func NewMeta() *Meta {
	return &Meta{}
}

// WithCorsPath sets the CORS configuration file path on the Meta instance.
func (m *Meta) WithCorsPath(path string) *Meta {
	m.corsPath = path
	return m
}

// CorsPath returns the configured CORS configuration file path.
func (m *Meta) CorsPath() string {
	if m == nil {
		return ""
	}
	return m.corsPath
}

// BluePrint defines the application blueprint including metadata and APIs.
type BluePrint struct {
	meta *Meta
	apis []*RRestAPIRegistration
}

// NewBluePrint initializes a new BluePrint builder.
func NewBluePrint() *BluePrint {
	return &BluePrint{}
}

// WithMeta sets the server metadata for the blueprint.
func (b *BluePrint) WithMeta(meta *Meta) *BluePrint {
	b.meta = meta
	return b
}

// WithAPIs sets the registered API groups, replacing any existing ones.
func (b *BluePrint) WithAPIs(apis ...*RRestAPIRegistration) *BluePrint {
	b.apis = append([]*RRestAPIRegistration(nil), apis...)
	return b
}

// AddAPIs appends registered API groups to the blueprint.
func (b *BluePrint) AddAPIs(apis ...*RRestAPIRegistration) *BluePrint {
	b.apis = append(b.apis, apis...)
	return b
}

// Meta returns the configured Meta instance or nil if unset.
func (b *BluePrint) Meta() *Meta {
	if b == nil {
		return nil
	}
	return b.meta
}

// Apis returns a copy of the registered API groups.
func (b *BluePrint) Apis() []*RRestAPIRegistration {
	if b == nil || b.apis == nil {
		return nil
	}
	return append([]*RRestAPIRegistration(nil), b.apis...)
}
