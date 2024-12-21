package config

type Config struct {
	Host string
	Port string
}

func NewConfig(_ []string, getenv func(string) string) *Config {
	config := &Config{
		Host: getenv("HOST"),
		Port: getenv("PORT"),
	}
	config.UseDefaults()
	return config
}

func (c *Config) UseDefaults() {
	if c.Host == "" {
		c.Host = ""
	}
	if c.Port == "" {
		c.Port = "8000"
	}
}
