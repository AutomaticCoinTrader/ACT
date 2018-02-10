package exchange

var registeredExchanges = make(map[string]ExchangeNewFunc)

type ExchangeNewFunc func(config interface{}) (Exchange, error)

func RegisterExchange(name string, exchangeNewFunc ExchangeNewFunc) {
	registeredExchanges[name] = exchangeNewFunc
}

func GetRegisterdExchanges() (map[string]ExchangeNewFunc) {
	return registeredExchanges
}

