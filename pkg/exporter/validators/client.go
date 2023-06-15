package validators

import (
	"context"
	"time"

	"github.com/ethpandaops/ethereum-validator-metrics-exporter/pkg/exporter/api"
	"github.com/sirupsen/logrus"
)

const (
	ChunkSize                         = 2
	LabelPubkey                string = "pubkey"
	LabelWithdrawalCredentials string = "withdrawal_credentials"
	LabelStatus                string = "status"
	LabelDefaultValue          string = ""
)

type Client struct {
	log             logrus.FieldLogger
	checkInterval   time.Duration
	validators      map[string]*Config
	validatorChunks [][]string
	labelsMap       map[string]int
	beaconchain     api.Client
	metrics         Metrics
}

// NewClient creates a new validators instance
func NewClient(log logrus.FieldLogger, validators []*Config, checkInterval time.Duration, namespace string, constLabels map[string]string, beaconchain api.Client) *Client {
	labelsMap := map[string]int{}
	labelsMap[LabelPubkey] = 0
	labelsMap[LabelWithdrawalCredentials] = 1
	labelsMap[LabelStatus] = 2

	validatorsMap := map[string]*Config{}

	var chunks [][]string

	keys := make([]string, 0, len(validators))

	for validator := range validators {
		validatorsMap[validators[validator].Pubkey] = validators[validator]
		keys = append(keys, validators[validator].Pubkey)

		for label := range validators[validator].Labels {
			if _, ok := labelsMap[label]; !ok {
				labelsMap[label] = len(labelsMap)
			}
		}
	}

	labels := make([]string, len(labelsMap))
	for label, index := range labelsMap {
		labels[index] = label
	}

	for i := 0; i < len(keys); i += ChunkSize {
		end := i + ChunkSize
		if end > len(keys) {
			end = len(keys)
		}

		chunks = append(chunks, keys[i:end])
	}

	instance := Client{
		log:             log.WithField("module", "validators"),
		validators:      validatorsMap,
		validatorChunks: chunks,
		beaconchain:     beaconchain,
		checkInterval:   checkInterval,
		labelsMap:       labelsMap,
		metrics:         NewMetrics(namespace, constLabels, labels),
	}

	return &instance
}

func (c *Client) Start(ctx context.Context) {
	c.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.checkInterval):
			c.log.WithField("check interval", c.checkInterval).Debug("Tick")
			c.tick(ctx)
		}
	}
}

func (c *Client) tick(ctx context.Context) {
	for i, chunk := range c.validatorChunks {
		c.log.WithField("chunk", i).Debug("Processing pubkeys")
		err := c.getValidators(ctx, chunk)

		if err != nil {
			c.log.WithError(err).WithField("pubkeys", chunk).Error("Error updating validators")
		}
	}
}

func (c *Client) getLabelValues(data *api.Validator) []string {
	values := make([]string, len(c.labelsMap))
	validator := c.validators[data.Pubkey]

	for label, index := range c.labelsMap {
		if validator.Labels != nil && validator.Labels[label] != "" {
			values[index] = validator.Labels[label]
		} else {
			switch label {
			case LabelPubkey:
				values[index] = data.Pubkey
			case LabelWithdrawalCredentials:
				values[index] = data.WithdrawalCredentials[:4]
			case LabelStatus:
				values[index] = data.Status
			default:
				values[index] = LabelDefaultValue
			}
		}
	}

	return values
}

func (c *Client) getValidators(ctx context.Context, validators []string) error {
	if len(validators) == 0 {
		return nil
	}

	if len(validators) == 1 {
		response, err := c.beaconchain.GetValidator(ctx, validators[0])
		if err != nil {
			return err
		}

		c.updateValidatorMetrics(response)
	} else {
		response, err := c.beaconchain.GetValidators(ctx, validators)
		if err != nil {
			return err
		}

		for _, validator := range response {
			if validator != nil {
				c.updateValidatorMetrics(validator)
			}
		}
	}

	return nil
}

func (c *Client) updateValidatorMetrics(data *api.Validator) {
	if data == nil {
		return
	}

	labels := c.getLabelValues(data)
	c.metrics.UpdateBalance(float64(data.Balance), labels)
	c.metrics.UpdateLastAttestationSlot(float64(data.LastAttestationSlot), labels)
	c.metrics.UpdateTotalWithdrawals(float64(data.TotalWithdrawals), labels)
}
