package common

type Params struct {
	MarketOrderEnabled   bool   `json:"market_order_enabled"`
	StakeEnabled         bool   `json:"stake_enabled"`
	ProcessOrdersEnabled bool   `json:"process_orders_enabled"`
	SwapEnabled          bool   `json:"swap_enabled"`
	PerpetualEnabled     bool   `json:"perpetual_enabled"`
	RewardEnabled        bool   `json:"reward_enabled"`
	LeverageEnabled      bool   `json:"leverage_enabled"`
	LimitProcessOrder    uint64 `json:"limit_process_order"`
	RewardPercentage     string `json:"reward_percentage"`
	MarginError          string `json:"margin_error"`
	MinimumDeposit       string `json:"minimum_deposit"`
}
