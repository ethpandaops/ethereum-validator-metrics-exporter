package api

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Client exposes beaconchain client
type Client interface {
	// GetValidators returns a map of validators
	GetValidators(ctx context.Context, pubkeys []string) (map[string]*Validator, error)
	// GetValidator returns a validator
	GetValidator(ctx context.Context, pubkey string) (*Validator, error)
}

type client struct {
	log     logrus.FieldLogger
	url     string
	apikey  string
	metrics Metrics
}

// NewClient creates a new beaconchain instance
func NewClient(log logrus.FieldLogger, conf *Config, namespace string) Client {
	return &client{
		log:     log.WithField("module", "beaconchain"),
		url:     conf.Endpoint,
		apikey:  conf.APIKey,
		metrics: NewMetrics(fmt.Sprintf("%s_%s", namespace, "http")),
	}
}

func (c *client) GetValidators(ctx context.Context, pubkeys []string) (map[string]*Validator, error) {
	response, err := c.getValidators(ctx, pubkeys)
	if err != nil {
		return nil, err
	}

	if response.Status != "OK" {
		return nil, fmt.Errorf("error response from server: %s", response.Status)
	}

	validators := make(map[string]*Validator)

	if response.Data == nil {
		return validators, nil
	}

	for i := range response.Data {
		validator := &response.Data[i]
		validators[validator.Pubkey] = validator
	}

	return validators, nil
}

func (c *client) GetValidator(ctx context.Context, pubkey string) (*Validator, error) {
	response, err := c.getValidator(ctx, pubkey)
	if err != nil {
		return nil, err
	}

	if response.Status != "OK" {
		return nil, fmt.Errorf("error response from server: %s", response.Status)
	}

	return &response.Data, nil
}
