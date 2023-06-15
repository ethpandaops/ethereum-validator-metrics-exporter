package api

type Config struct {
	Endpoint string `yaml:"endpoint" default:"https://beaconcha.in"`
	APIKey   string `yaml:"apikey"`
}

func (c *Config) Validate() error {
	return nil
}
