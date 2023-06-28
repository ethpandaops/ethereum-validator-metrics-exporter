package api

type Config struct {
	Endpoint             string `yaml:"endpoint" default:"https://beaconcha.in"`
	APIKey               string `yaml:"apikey"`
	MaxRequestsPerMinute int    `yaml:"maxRequestsPerMinute" default:"10"`
	BatchSize            int    `yaml:"batchSize" default:"50"`
}

func (c *Config) Validate() error {
	return nil
}
