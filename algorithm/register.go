package algorithm

var registeredAlgorithms map[string]*registeredAlgorithm = make(map[string]*registeredAlgorithm)

type InternalTradeAlgorithmNewFunc func(configDir string) (InternalTradeAlgorithm, error)
type ExternalTradeAlgorithmNewFunc func(configDir string) (ExternalTradeAlgorithm, error)

type registeredAlgorithm struct {
	InternalTradeAlgorithmNewFunc InternalTradeAlgorithmNewFunc
	ExternalTradeAlgorithmNewFunc ExternalTradeAlgorithmNewFunc
}

func RegisterAlgorithm(name string, internalTradeAlgorithmNewFunc InternalTradeAlgorithmNewFunc, externalTradeAlgorithmNewFunc ExternalTradeAlgorithmNewFunc) {
	registeredAlgorithms[name] = &registeredAlgorithm{
		InternalTradeAlgorithmNewFunc: internalTradeAlgorithmNewFunc,
		ExternalTradeAlgorithmNewFunc: externalTradeAlgorithmNewFunc,
	}
}

func GetRegisterdAlgoriths() (map[string]*registeredAlgorithm) {
	return registeredAlgorithms
}
