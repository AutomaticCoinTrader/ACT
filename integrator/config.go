package integrator

import "github.com/AutomaticCoinTrader/ACT/exchange/zaif"

type exchangesConfig struct {
	Zaif *zaif.ExchangeConfig `json:"zaif" yaml:"zaif" toml:"zaif" config:"zaif"`
}

