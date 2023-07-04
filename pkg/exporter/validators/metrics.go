package validators

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	balance             *prometheus.GaugeVec
	exited              *prometheus.GaugeVec
	credentialsCode     *prometheus.GaugeVec
	lastAttestationSlot *prometheus.GaugeVec
	totalWithdrawals    *prometheus.GaugeVec
}

func NewMetrics(namespace string, constLabels map[string]string, labels []string) Metrics {
	m := Metrics{
		balance: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "balance",
				Help:        "The balance of the validator.",
				ConstLabels: constLabels,
			},
			labels,
		),
		exited: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "exited",
				Help:        "The exited status of the validator.",
				ConstLabels: constLabels,
			},
			labels,
		),
		credentialsCode: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "credentials_code",
				Help:        "The withdrawal credentials code of the validator.",
				ConstLabels: constLabels,
			},
			labels,
		),
		lastAttestationSlot: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "last_attestation_slot",
				Help:        "The last attestation slot of the validator.",
				ConstLabels: constLabels,
			},
			labels,
		),
		totalWithdrawals: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "total_withdrawals",
				Help:        "The total withdrawals of the validator.",
				ConstLabels: constLabels,
			},
			labels,
		),
	}

	prometheus.MustRegister(m.balance)
	prometheus.MustRegister(m.exited)
	prometheus.MustRegister(m.credentialsCode)
	prometheus.MustRegister(m.lastAttestationSlot)
	prometheus.MustRegister(m.totalWithdrawals)

	return m
}

func (m Metrics) UpdateBalance(balance float64, labels []string) {
	m.balance.WithLabelValues(labels...).Set(balance)
}

func (m Metrics) UpdateExited(exited float64, labels []string) {
	m.exited.WithLabelValues(labels...).Set(exited)
}

func (m Metrics) UpdateCredentialsCode(code float64, labels []string) {
	m.credentialsCode.WithLabelValues(labels...).Set(code)
}

func (m Metrics) UpdateLastAttestationSlot(slot float64, labels []string) {
	m.lastAttestationSlot.WithLabelValues(labels...).Set(slot)
}

func (m Metrics) UpdateTotalWithdrawals(total float64, labels []string) {
	m.totalWithdrawals.WithLabelValues(labels...).Set(total)
}
