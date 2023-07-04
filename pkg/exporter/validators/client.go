package validators

import (
	"context"
	"time"

	"github.com/ethpandaops/ethereum-validator-metrics-exporter/pkg/exporter/api"
	"github.com/sirupsen/logrus"
)

const (
	LabelPubkey       string = "pubkey"
	LabelDefaultValue string = ""
)

type Client struct {
	log                   logrus.FieldLogger
	checkInterval         time.Duration
	validators            map[string]*Config
	validatorChunks       [][]string
	validatorRequestDelay time.Duration
	labelsMap             map[string]int
	beaconchain           api.Client
	metrics               Metrics
}

// NewClient creates a new validators instance
func NewClient(log logrus.FieldLogger, validators []*Config, checkInterval time.Duration, namespace string, constLabels map[string]string, beaconchain api.Client) *Client {
	labelsMap := map[string]int{}
	labelsMap[LabelPubkey] = 0

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

	for i := 0; i < len(keys); i += beaconchain.GetBatchSize() {
		end := i + beaconchain.GetBatchSize()
		if end > len(keys) {
			end = len(keys)
		}

		chunks = append(chunks, keys[i:end])
	}

	instance := Client{
		log:                   log.WithField("module", "validators"),
		validators:            validatorsMap,
		validatorChunks:       chunks,
		validatorRequestDelay: time.Minute / time.Duration(beaconchain.GetMaxRequestsPerMinute()),
		beaconchain:           beaconchain,
		checkInterval:         checkInterval,
		labelsMap:             labelsMap,
		metrics:               NewMetrics(namespace, constLabels, labels),
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
	c.log.Debug("Starting validators update")

	ticker := time.NewTicker(c.validatorRequestDelay)
	defer ticker.Stop()

	for i, chunk := range c.validatorChunks {
		c.log.WithFields(logrus.Fields{
			"chunk":  i,
			"length": len(chunk),
		}).Debug("Processing validator pubkeys chunk")

		err := c.getValidators(ctx, chunk)

		if err != nil {
			c.log.WithError(err).WithField("pubkeys", chunk).Error("Error updating validators")
		}

		// Don't delay after the last chunk
		if i == len(c.validatorChunks)-1 {
			break
		}

		c.log.WithField("delay", c.validatorRequestDelay).Debug("Delaying request for next validator chunk")
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}

	c.log.Debug("Finished validators update")
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

	exited := float64(0)
	if data.IsExited() {
		exited = float64(1)
	}

	c.metrics.UpdateExited(exited, labels)

	credentialsCode, err := data.GetWithdrawalCredentialsCode()
	if err != nil {
		c.log.WithError(err).WithField("credentials", data.WithdrawalCredentials).Error("Error parsing withdrawal credentials")
	}

	code := float64(0)
	if credentialsCode != nil {
		code = float64(*credentialsCode)
	}

	c.metrics.UpdateCredentialsCode(code, labels)
	c.metrics.UpdateLastAttestationSlot(float64(data.LastAttestationSlot), labels)
	c.metrics.UpdateTotalWithdrawals(float64(data.TotalWithdrawals), labels)
}
