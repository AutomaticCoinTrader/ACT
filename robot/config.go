package robot

import (
	"github.com/AutomaticCoinTrader/ACT/algorithm/example"
)

type TradeConfig struct {
	Example *example.Config `json:"example" yaml:"example" toml:"example" config:"example"`
}

type ArbitrageTradeConfig struct {
	Example *example.ArbitrageConfig `json:"example" yaml:"example" toml:"example" config:"example"`
}


