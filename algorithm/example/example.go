package example

import (
	"github.com/AutomaticCoinTrader/ACT/algorithm"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"github.com/AutomaticCoinTrader/ACT/notifier"
	"github.com/AutomaticCoinTrader/ACT/configurator"
	"github.com/pkg/errors"
	"log"
	"path"
	"encoding/json"
)

const (
	algorithmName string = "example"
)

type internalTradeConfig struct {
}

type externalTradeConfig struct {
}

type config struct {
	InternalTrade *internalTradeConfig `json:"internalTrade" yaml:"internalTrade" toml:"internalTrade"`
	ExternalTrade *externalTradeConfig `json:"externalTrade" yaml:"externalTrade" toml:"externalTrade"`
}

type internalTradeExample struct {
	name   string
	config *internalTradeConfig
}

func (i *internalTradeExample) GetName() (string) {
	return i.name
}

func (i *internalTradeExample) Initialize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	return nil
}

func (i *internalTradeExample) Update(currencyPair string, tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	// trade
	log.Printf("currencyPair = %v, internal trade", currencyPair)

	boardCursor, err := tradeContext.GetBuyBoardCursor(currencyPair)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(boardCursor.PriceAll())
	if err != nil {
		return err
	}
	log.Printf("buy: %v\n", string(bytes))

	boardCursor, err = tradeContext.GetSellBoardCursor(currencyPair)
	if err != nil {
		return err
	}
	bytes, err = json.Marshal(boardCursor.PriceAll())
	if err != nil {
		return err
	}
	log.Printf("sell: %v\n", string(bytes))


	return nil
}

func (i *internalTradeExample) Finalize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	return nil
}

func newInternalTradeExample(configDir string) (algorithm.InternalTradeAlgorithm, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		return nil, errors.Errorf("can not create configurator (config file path prefix = %v)", configFilePathPrefix)
	}
	conf := new(config)
	err = cf.Load(conf)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &internalTradeExample{
		name:   algorithmName,
		config: conf.InternalTrade,
	}, nil
}

type externalTradeExample struct {
	name   string
	config *externalTradeConfig
}

func (e *externalTradeExample) GetName() (string) {
	return e.name
}

func (e *externalTradeExample) Initialize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	return nil
}

func (e *externalTradeExample) Update(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	// arbitrage trade
	log.Printf("external trade")
	return nil
}

func (e *externalTradeExample) Finalize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	return nil
}

func newExternalTradeExample(configDir string) (algorithm.ExternalTradeAlgorithm, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		return nil, errors.Errorf("can not create configurator (config file path prefix = %v)", configFilePathPrefix)
	}
	conf := new(config)
	err = cf.Load(conf)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &externalTradeExample{
		name:   algorithmName,
		config: conf.ExternalTrade,
	}, nil
}

func init() {
	algorithm.RegisterAlgorithm(algorithmName, newInternalTradeExample, newExternalTradeExample)
}
