package api

type Response[T any] struct {
	Data   T      `json:"data"`
	Status string `json:"status"`
}

type Validator struct {
	ActivationEligibilityEpoch int    `json:"activationeligibilityepoch"`
	ActivationEpoch            int    `json:"activationepoch"`
	Balance                    int    `json:"balance"`
	EffectiveBalance           int    `json:"effectivebalance"`
	ExitEpoch                  int    `json:"exitepoch"`
	LastAttestationSlot        int    `json:"lastattestationslot"`
	Name                       string `json:"name"`
	Pubkey                     string `json:"pubkey"`
	Slashed                    bool   `json:"slashed"`
	Status                     string `json:"status"`
	ValidatorIndex             int    `json:"validatorindex"`
	WithdrawableEpoch          int    `json:"withdrawableepoch"`
	WithdrawalCredentials      string `json:"withdrawalcredentials"`
	TotalWithdrawals           int    `json:"total_withdrawals"`
}
