package validators

import (
	"fmt"
)

type Config struct {
	Pubkey string            `yaml:"pubkey"`
	Labels map[string]string `yaml:"labels"`
}

func (c *Config) Validate() error {
	if c.Pubkey == "" {
		return fmt.Errorf("pubkey is required")
	}

	return nil
}
