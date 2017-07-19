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
	config                   *Config
	configDir                string
	notifier                 *notifier.Notifier
	tradeAlgorithms          map[string][]algorithm.TradeAlgorithm
	arbitrageTradeAlgorithms []algorithm.ArbitrageTradeAlgorithm
}

func (r *Robot) CreateTradeAlgorithms(tradeID string, tradeContext exchange.TradeContext) (error) {
	tradeAlgorithms := make([]algorithm.TradeAlgorithm, 0)
	for name, registeredAlgorithm := range algorithm.GetRegisterdAlgoriths() {
		if registeredAlgorithm.TradeAlgorithmNewFunc == nil {
			continue
		}
		log.Printf("create %v algorithm (trade id = %v)", name, tradeID)
		newTradeAlgoritm, err := registeredAlgorithm.TradeAlgorithmNewFunc(path.Join(r.configDir, algorithm.AlgorithmConfigDir))
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
	r.tradeAlgorithms[tradeID] = tradeAlgorithms
	return nil
}

func (r *Robot) UpdateTradeAlgorithms(tradeID string, tradeContext exchange.TradeContext) (error) {
	tradeAlgorithms := r.tradeAlgorithms[tradeID]
	for _, tradeAlgorithm := range tradeAlgorithms {
		err := tradeAlgorithm.Update(tradeContext, r.notifier)
		if err != nil {
			log.Printf("algorithm update error (name = %v, trade id = %v reason = %v)", tradeAlgorithm.GetName(), tradeID, err)
		}
	}
	return nil
}

func (r *Robot) DestroyTradeAlgorithms(tradeID string, tradeContext exchange.TradeContext) (error) {
	tradeAlgorithms := r.tradeAlgorithms[tradeID]
	for _, tradeAlgorithm := range tradeAlgorithms {
		err := tradeAlgorithm.Finalize(tradeContext, r.notifier)
		if err != nil {
			log.Printf("algorithm finalize error (name = %v, trade id = %v reason = %v)", tradeAlgorithm.GetName(), tradeID, err)
		}
	}
	return nil
}

func (r *Robot) CreateArbitrageTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for name, registeredAlgorithm := range algorithm.GetRegisterdAlgoriths() {
		if registeredAlgorithm.ArbitrageTradeAlgorithmNewFunc == nil {
			continue
		}
		log.Printf("create %v arbitrage algorithm", name)
		newArbitrageTradeAlgoritm, err := registeredAlgorithm.ArbitrageTradeAlgorithmNewFunc(path.Join(r.configDir, algorithm.AlgorithmConfigDir))
		if err != nil {
			log.Printf("can not create arbitrage algorithm of %v (reason = %v)", name, err)
			continue
		}
		err = newArbitrageTradeAlgoritm.Initialize(exchanges, r.notifier)
		if err != nil {
			r.DestroyArbitrageTradeAlgorithms(exchanges)
			return errors.Wrap(err, fmt.Sprintf("arbitrage algorithm initialize error of %v", name))
		}
		r.arbitrageTradeAlgorithms = append(r.arbitrageTradeAlgorithms, newArbitrageTradeAlgoritm)
	}
	return nil
}

func (r *Robot) UpdateArbitrageTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for _, arbitrageTradeAlgoritm := range r.arbitrageTradeAlgorithms {
		err := arbitrageTradeAlgoritm.Update(exchanges, r.notifier)
		if err != nil {
			log.Printf("arbitrage algorithm update error (name = %v, reason = %v)", arbitrageTradeAlgoritm.GetName(), err)
		}
	}
	return nil
}

func (r *Robot) DestroyArbitrageTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for _, arbitrageTradeAlgoritm := range r.arbitrageTradeAlgorithms {
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
		config:                   config,
		configDir:                configDir,
		notifier:                 notifier,
		tradeAlgorithms:          make(map[string][]algorithm.TradeAlgorithm),
		arbitrageTradeAlgorithms: make([]algorithm.ArbitrageTradeAlgorithm, 0),
	}
	r.loadPluginFiles()
	return r, nil
}
