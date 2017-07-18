package algorithm

var registeredAlgorithms map[string]*registeredAlgorithm = make(map[string]*registeredAlgorithm)

type TradeAlgorithmNewFunc func(configDir string) (TradeAlgorithm, error)
type ArbitrageTradeAlgorithmNewFunc func(configDir string) (ArbitrageTradeAlgorithm, error)

type registeredAlgorithm struct {
	TradeAlgorithmNewFunc TradeAlgorithmNewFunc
	ArbitrageTradeAlgorithmNewFunc ArbitrageTradeAlgorithmNewFunc
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
