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

func (r *Robot) CreateInternalTradeAlgorithms(ex exchange.Exchange) (error) {

	for name, registeredAlgorithm := range algorithm.GetRegisterdAlgoriths() {
		if registeredAlgorithm.InternalTradeAlgorithmNewFunc == nil {
			continue
		}
		log.Printf("create %v internal algorithm", name)
		newInternalTradeAlgoritm, err := registeredAlgorithm.InternalTradeAlgorithmNewFunc(path.Join(r.configDir, algorithm.AlgorithmConfigDir))
		if err != nil {
			log.Printf("can not create internal algorithm of %v (reason = %v)", name, err)
			continue
		}
		err = newInternalTradeAlgoritm.Initialize(ex, r.notifier)
		if err != nil {
			r.DestroyInternalTradeAlgorithms(ex)
			return errors.Wrap(err, fmt.Sprintf("internal algorithm initialize error of %v", name))
		}
		r.internalTradeAlgorithms = append(r.internalTradeAlgorithms, newInternalTradeAlgoritm)
	}

	return nil
}

func (r *Robot) UpdateInternalTradeAlgorithms(currencyPair string, ex exchange.Exchange) (error) {
	for _, tradeAlgorithm := range r.internalTradeAlgorithms {
		err := tradeAlgorithm.Update(currencyPair, ex, r.notifier)
		if err != nil {
			log.Printf("internal algorithm update error (name = %v, reason = %v)", tradeAlgorithm.GetName(), err)
		}
	}
	return nil
}

func (r *Robot) DestroyInternalTradeAlgorithms(ex exchange.Exchange) (error) {
	for _, tradeAlgorithm := range r.internalTradeAlgorithms {
		err := tradeAlgorithm.Finalize(ex, r.notifier)
		if err != nil {
			log.Printf("internal algorithm finalize error (name = %v, reason = %v)", tradeAlgorithm.GetName(), err)
		}
	}
	return nil
}

func (r *Robot) CreateExternalTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for name, registeredAlgorithm := range algorithm.GetRegisterdAlgoriths() {
		if registeredAlgorithm.ExternalTradeAlgorithmNewFunc == nil {
			continue
		}
		log.Printf("create %v external algorithm", name)
		newExternalTradeAlgoritm, err := registeredAlgorithm.ExternalTradeAlgorithmNewFunc(path.Join(r.configDir, algorithm.AlgorithmConfigDir))
		if err != nil {
			log.Printf("can not create external algorithm of %v (reason = %v)", name, err)
			continue
		}
		err = newExternalTradeAlgoritm.Initialize(exchanges, r.notifier)
		if err != nil {
			r.DestroyExternalTradeAlgorithms(exchanges)
			return errors.Wrap(err, fmt.Sprintf("external algorithm initialize error of %v", name))
		}
		r.externalTradeAlgorithms = append(r.externalTradeAlgorithms, newExternalTradeAlgoritm)
	}
	return nil
}

func (r *Robot) UpdateExternalTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for _, externalTradeAlgoritm := range r.externalTradeAlgorithms {
		err := externalTradeAlgoritm.Update(exchanges, r.notifier)
		if err != nil {
			log.Printf("external algorithm update error (name = %v, reason = %v)", externalTradeAlgoritm.GetName(), err)
		}
	}
	return nil
}

func (r *Robot) DestroyExternalTradeAlgorithms(exchanges map[string]exchange.Exchange) (error) {
	for _, externalTradeAlgoritm := range r.externalTradeAlgorithms {
		err := externalTradeAlgoritm.Finalize(exchanges, r.notifier)
		if err != nil {
			log.Printf("external algorithm finalize error (name = %v, reason = %v)", externalTradeAlgoritm.GetName(), err)
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
