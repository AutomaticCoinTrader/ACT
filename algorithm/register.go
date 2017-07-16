package algorithm

var registeredAlgorithms map[string]*registeredAlgorithm = make(map[string]*registeredAlgorithm)

type TradeAlgorithmNewFunc func(config interface{}) (TradeAlgorithm, error)
type ArbitrageTradeAlgorithmNewFunc func(config interface{}) (ArbitrageTradeAlgorithm, error)

type registeredAlgorithm struct {
	TradeAlgorithmNewFunc func(config interface{}) (TradeAlgorithm, error)
	ArbitrageTradeAlgorithmNewFunc func(config interface{}) (ArbitrageTradeAlgorithm, error)
}

func RegisterAlgorithm(name string, tradeAlgorithmNewFunc TradeAlgorithmNewFunc, arbitrageTradeAlgorithmNewFunc ArbitrageTradeAlgorithmNewFunc) {
	registeredAlgorithms[name] = &registeredAlgorithm{
		TradeAlgorithmNewFunc: tradeAlgorithmNewFunc,
		ArbitrageTradeAlgorithmNewFunc: arbitrageTradeAlgorithmNewFunc,
	}
}

func GetRegisterdAlgoriths() (map[string]*registeredAlgorithm) {
	return registeredAlgorithms
}
