package timezone

const (
	CtxTimezoneKey    = "Timezone"
	defaultHeaderName = "Timezone"
)

type Config struct {
	// HeaderName name with timezone data
	HeaderName string
}

var DefaultConfig = Config{
	HeaderName: defaultHeaderName,
}

func defaultConfig(config ...Config) Config {
	if len(config) < 1 {
		return DefaultConfig
	}

	cfg := config[0]

	if cfg.HeaderName == "" {
		cfg.HeaderName = defaultHeaderName
	}

	return cfg
}
