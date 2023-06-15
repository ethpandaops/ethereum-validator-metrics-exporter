package exporter

import (
	"fmt"
	"time"

	"github.com/ethpandaops/ethereum-validator-metrics-exporter/pkg/exporter/api"
	"github.com/ethpandaops/ethereum-validator-metrics-exporter/pkg/exporter/validators"
)

// Config holds the configuration for the ethereum sync status tool.
type Config struct {
	GlobalConfig GlobalConfig `yaml:"global"`
	// API is the configuration for the beaconcha.in client.
	API api.Config `yaml:"beaconcha_in"`
	// Validators is the configuration for the validators.
	Validators []*validators.Config `yaml:"validators"`
}

type GlobalConfig struct {
	LoggingLevel  string            `yaml:"logging" default:"warn"`
	MetricsAddr   string            `yaml:"metricsAddr" default:":9090"`
	Namespace     string            `yaml:"namespace" default:"eth_validator"`
	CheckInterval time.Duration     `yaml:"checkInterval" default:"24h"`
	Labels        map[string]string `yaml:"labels"`
}

func (c *Config) Validate() error {
	err := c.API.Validate()
	if err != nil {
		return err
	}

	for index, validator := range c.Validators {
		err := validator.Validate()
		if err != nil {
			return fmt.Errorf("validator %d: %w", index, err)
		}
	}

	return nil
}
