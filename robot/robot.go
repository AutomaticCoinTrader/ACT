package robot

import (
	"github.com/pkg/errors"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"github.com/AutomaticCoinTrader/ACT/algorithm"
	"github.com/AutomaticCoinTrader/ACT/notifier"
	"log"
	"fmt"
	"path"
)

type Robot struct {
	config                  *Config
	configDir               string
	notifier                *notifier.Notifier
	internalTradeAlgorithms map[string][]algorithm.InternalTradeAlgorithm
	externalTradeAlgorithms []algorithm.ExternalTradeAlgorithm
}

func (r *Robot) CreateInternalTradeAlgorithms(tradeID string, tradeContext exchange.TradeContext) (error) {
	tradeAlgorithms := make([]algorithm.InternalTradeAlgorithm, 0)
	for name, registeredAlgorithm := range algorithm.GetRegisterdAlgoriths() {
		if registeredAlgorithm.InternalTradeAlgorithmNewFunc == nil {
			continue
		}
		log.Printf("create %v algorithm (trade id = %v)", name, tradeID)
		newTradeAlgoritm, err := registeredAlgorithm.InternalTradeAlgorithmNewFunc(path.Join(r.configDir, algorithm.AlgorithmConfigDir))
		if err != nil {
			log.Printf("can not create algorithm of %v (trade id = %v, reason = %v)", name, tradeID, err)
			continue
		}
		err = newTradeAlgoritm.Initialize(tradeContext, r.notifier)
		if err != nil {
			r.DestroyTradeAlgorithms(tradeID, tradeContext)
			return errors.Wrap(err, fmt.Sprintf("algorithm initialize error of %v (trade id = %v)", name, tradeID))
		}
		tradeAlgorithms = append(tradeAlgorithms, newTradeAlgoritm)
	}
	r.internalTradeAlgorithms[tradeID] = tradeAlgorithms
	return nil
}

func (r *Robot) UpdateTradeAlgorithms(tradeID string, tradeContext exchange.TradeContext) (error) {
	tradeAlgorithms := r.internalTradeAlgorithms[tradeID]
	for _, tradeAlgorithm := range tradeAlgorithms {
		err := tradeAlgorithm.Update(tradeContext, r.notifier)
		if err != nil {
			log.Printf("algorithm update error (name = %v, trade id = %v reason = %v)", tradeAlgorithm.GetName(), tradeID, err)
		}
	}
	return nil
}

func (r *Robot) DestroyTradeAlgorithms(tradeID string, tradeContext exchange.TradeContext) (error) {
	tradeAlgorithms := r.internalTradeAlgorithms[tradeID]
	for _, tradeAlgorithm := range tradeAlgorithms {
		err := tradeAlgorithm.Finalize(tradeContext, r.notifier)
		if err != nil {
			log.Printf("algorithm finalize error (name = %v, trade id = %v reason = %v)", tradeAlgorithm.GetName(), tradeID, err)
		}
	}
	return nil
}

func (r *Robot) CreateExternalTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for name, registeredAlgorithm := range algorithm.GetRegisterdAlgoriths() {
		if registeredAlgorithm.ExternalTradeAlgorithmNewFunc == nil {
			continue
		}
		log.Printf("create %v arbitrage algorithm", name)
		newArbitrageTradeAlgoritm, err := registeredAlgorithm.ExternalTradeAlgorithmNewFunc(path.Join(r.configDir, algorithm.AlgorithmConfigDir))
		if err != nil {
			log.Printf("can not create arbitrage algorithm of %v (reason = %v)", name, err)
			continue
		}
		err = newArbitrageTradeAlgoritm.Initialize(exchanges, r.notifier)
		if err != nil {
			r.DestroyArbitrageTradeAlgorithms(exchanges)
			return errors.Wrap(err, fmt.Sprintf("arbitrage algorithm initialize error of %v", name))
		}
		r.externalTradeAlgorithms = append(r.externalTradeAlgorithms, newArbitrageTradeAlgoritm)
	}
	return nil
}

func (r *Robot) UpdateArbitrageTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for _, arbitrageTradeAlgoritm := range r.externalTradeAlgorithms {
		err := arbitrageTradeAlgoritm.Update(exchanges, r.notifier)
		if err != nil {
			log.Printf("arbitrage algorithm update error (name = %v, reason = %v)", arbitrageTradeAlgoritm.GetName(), err)
		}
	}
	return nil
}

func (r *Robot) DestroyArbitrageTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for _, arbitrageTradeAlgoritm := range r.externalTradeAlgorithms {
		err := arbitrageTradeAlgoritm.Finalize(exchanges, r.notifier)
		if err != nil {
			log.Printf("arbitrage algorithm finalize error (name = %v, reason = %v)", arbitrageTradeAlgoritm.GetName(), err)
		}
	}
	return nil
}

type Config struct {
	AlgorithmPluginDir string `json:"algorithmPluginDir" yaml:"algorithmPluginDir" toml:"algorithmPluginDir"`
}

func NewRobot(config *Config, configDir string, notifier *notifier.Notifier) (*Robot, error) {
	r := &Robot{
		config:                  config,
		configDir:               configDir,
		notifier:                notifier,
		internalTradeAlgorithms: make(map[string][]algorithm.InternalTradeAlgorithm),
		externalTradeAlgorithms: make([]algorithm.ExternalTradeAlgorithm, 0),
	}
	r.loadPluginFiles()
	return r, nil
}
