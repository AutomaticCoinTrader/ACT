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
	internalTradeAlgorithms []algorithm.InternalTradeAlgorithm
	externalTradeAlgorithms []algorithm.ExternalTradeAlgorithm
}

func (r *Robot) CreateInternalTradeAlgorithms(tradeContext exchange.TradeContext) (error) {

	for name, registeredAlgorithm := range algorithm.GetRegisterdAlgoriths() {
		if registeredAlgorithm.InternalTradeAlgorithmNewFunc == nil {
			continue
		}
		log.Printf("create %v algorithm", name)
		newInternalTradeAlgoritm, err := registeredAlgorithm.InternalTradeAlgorithmNewFunc(path.Join(r.configDir, algorithm.AlgorithmConfigDir))
		if err != nil {
			log.Printf("can not create algorithm of %v (reason = %v)", name, err)
			continue
		}
		err = newInternalTradeAlgoritm.Initialize(tradeContext, r.notifier)
		if err != nil {
			r.DestroyInternalTradeAlgorithms(tradeContext)
			return errors.Wrap(err, fmt.Sprintf("algorithm initialize error of %v", name))
		}
		r.internalTradeAlgorithms = append(r.internalTradeAlgorithms, newInternalTradeAlgoritm)
	}

	return nil
}

func (r *Robot) UpdateInternalTradeAlgorithms(currencyPair string, tradeContext exchange.TradeContext) (error) {
	for _, tradeAlgorithm := range r.internalTradeAlgorithms {
		err := tradeAlgorithm.Update(currencyPair, tradeContext, r.notifier)
		if err != nil {
			log.Printf("algorithm update error (name = %v, reason = %v)", tradeAlgorithm.GetName(), err)
		}
	}
	return nil
}

func (r *Robot) DestroyInternalTradeAlgorithms(tradeContext exchange.TradeContext) (error) {
	for _, tradeAlgorithm := range r.internalTradeAlgorithms {
		err := tradeAlgorithm.Finalize(tradeContext, r.notifier)
		if err != nil {
			log.Printf("algorithm finalize error (name = %v, reason = %v)", tradeAlgorithm.GetName(), err)
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
		newExternalTradeAlgoritm, err := registeredAlgorithm.ExternalTradeAlgorithmNewFunc(path.Join(r.configDir, algorithm.AlgorithmConfigDir))
		if err != nil {
			log.Printf("can not create arbitrage algorithm of %v (reason = %v)", name, err)
			continue
		}
		err = newExternalTradeAlgoritm.Initialize(exchanges, r.notifier)
		if err != nil {
			r.DestroyExternalTradeAlgorithms(exchanges)
			return errors.Wrap(err, fmt.Sprintf("arbitrage algorithm initialize error of %v", name))
		}
		r.externalTradeAlgorithms = append(r.externalTradeAlgorithms, newExternalTradeAlgoritm)
	}
	return nil
}

func (r *Robot) UpdateExternalTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for _, externalTradeAlgoritm := range r.externalTradeAlgorithms {
		err := externalTradeAlgoritm.Update(exchanges, r.notifier)
		if err != nil {
			log.Printf("arbitrage algorithm update error (name = %v, reason = %v)", externalTradeAlgoritm.GetName(), err)
		}
	}
	return nil
}

func (r *Robot) DestroyExternalTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for _, externalTradeAlgoritm := range r.externalTradeAlgorithms {
		err := externalTradeAlgoritm.Finalize(exchanges, r.notifier)
		if err != nil {
			log.Printf("arbitrage algorithm finalize error (name = %v, reason = %v)", externalTradeAlgoritm.GetName(), err)
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
		internalTradeAlgorithms: make([]algorithm.InternalTradeAlgorithm, 0),
		externalTradeAlgorithms: make([]algorithm.ExternalTradeAlgorithm, 0),
	}
	r.loadPluginFiles()
	return r, nil
}
