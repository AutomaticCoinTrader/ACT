package integrator

import (
	"github.com/AutomaticCoinTrader/ACT/exchange/zaif"
	"github.com/AutomaticCoinTrader/ACT/exchange/coincheck"
)

type exchangesConfig struct {
	Zaif *zaif.ExchangeConfig `json:"zaif" yaml:"zaif" toml:"zaif" config:"zaif"`
	Coincheck *coincheck.CoincheckExchangeConfig `json:"coincheck" yaml:"coincheck" config:"coincheck"`
}

