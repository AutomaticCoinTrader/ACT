package exchange

type CurrencyPair string

type BasicTradeRequester interface {
	LimitOrder(order_type string, pair CurrencyPair, rate float64, amount float64) error
	MarketOrder(order_type string, pair CurrencyPair, amount float64) error
	GetOpenOrders() error
	GetExecutionHistory() error
	GetBalance() error
}

type LeveragedTradeRequester interface {
	LeveragedLimitOrder(order_type string, pair CurrencyPair, rate float64, amount float64) error
	LeveragedMarketOrder(order_type string, pair CurrencyPair, amount float64) error
}

type AccountChargeRequester interface {
	Transfer(currency string, destination_address string, amount float64, fee float64) error
	// Withdraw(currency string, destination_account BankAccount, amount int64) error
	// Deposit(currency string, source_account BankAccount, amount int64) error
}