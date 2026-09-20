package contracts

// CorsSetting configures Cross-Origin Resource Sharing policy options.
type CorsSetting struct {
	allowOrigin      []string
	allowHTTPMethods []HTTPMethod
}

// NewCorsSetting initializes a new CorsSetting builder.
func NewCorsSetting() *CorsSetting {
	return &CorsSetting{}
}

// WithAllowOrigin sets the allowed CORS origins.
func (c *CorsSetting) WithAllowOrigin(origins ...string) *CorsSetting {
	c.allowOrigin = append([]string(nil), origins...)
	return c
}

// WithAllowHTTPMethods sets the allowed CORS HTTP methods.
func (c *CorsSetting) WithAllowHTTPMethods(
	methods ...HTTPMethod,
) *CorsSetting {
	c.allowHTTPMethods = append([]HTTPMethod(nil), methods...)
	return c
}

// AllowOrigin returns a copy of the allowed CORS origins.
func (c *CorsSetting) AllowOrigin() []string {
	if c == nil || c.allowOrigin == nil {
		return nil
	}
	return append([]string(nil), c.allowOrigin...)
}

// AllowHTTPMethods returns a copy of the allowed CORS HTTP methods.
func (c *CorsSetting) AllowHTTPMethods() []HTTPMethod {
	if c == nil || c.allowHTTPMethods == nil {
		return nil
	}
	return append([]HTTPMethod(nil), c.allowHTTPMethods...)
}

// Meta holds server configuration metadata such as CORS settings.
type Meta struct {
	cors *CorsSetting
}

// NewMeta initializes a new Meta builder.
func NewMeta() *Meta {
	return &Meta{}
}

// WithCors sets the CORS configuration on the Meta instance.
func (m *Meta) WithCors(cors *CorsSetting) *Meta {
	m.cors = cors
	return m
}

// Cors returns the configured CorsSetting or nil if unset.
func (m *Meta) Cors() *CorsSetting {
	if m == nil {
		return nil
	}
	return m.cors
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
