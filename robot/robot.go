package robot

import (
	"github.com/pkg/errors"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"github.com/AutomaticCoinTrader/ACT/algorithm"
	"log"
	"fmt"
	"github.com/AutomaticCoinTrader/ACT/notifier"
	"reflect"
)

type Robot struct {
	config                   *Config
	notifier                 *notifier.Notifier
	tradeAlgorithms          map[string][]algorithm.TradeAlgorithm
	arbitrageTradeAlgorithms []algorithm.ArbitrageTradeAlgorithm
}

func (r *Robot) CreateTradeAlgorithms(tradeID string, tradeContext exchange.TradeContext) (error) {
	tradeAlgorithms := make([]algorithm.TradeAlgorithm, 0)
	for name, registeredAlgorithm := range algorithm.GetRegisterdAlgoriths() {
		t := reflect.TypeOf(r.config.Trade).Elem()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Tag.Get("config") != name {
				continue
			}
			v := reflect.ValueOf(r.config.Trade)
			if v.IsNil() {
				continue
			}
			v = v.Elem()
			fv := v.FieldByName(f.Name)
			if fv.IsNil() {
				continue
			}
			conf := fv.Interface()
			if registeredAlgorithm.TradeAlgorithmNewFunc == nil {
				continue
			}
			log.Printf("create %v algorithm (trade id = %v)", name, tradeID)
			newTradeAlgoritm, err := registeredAlgorithm.TradeAlgorithmNewFunc(conf)
			if err != nil {
				r.DestroyTradeAlgorithms(tradeID, tradeContext)
				return errors.Wrap(err, fmt.Sprintf("can not create algorithm of %v (trade id = %v)", name, tradeID))
			}
			err = newTradeAlgoritm.Initialize(tradeContext, r.notifier)
			if err != nil {
				r.DestroyTradeAlgorithms(tradeID, tradeContext)
				return errors.Wrap(err, fmt.Sprintf("algorithm initialize error of %v (trade id = %v)", name, tradeID))
			}
			tradeAlgorithms = append(tradeAlgorithms, newTradeAlgoritm)
		}
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
		t := reflect.TypeOf(r.config.ArbitrageTrade).Elem()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Tag.Get("config") != name {
				continue
			}
			v := reflect.ValueOf(r.config.ArbitrageTrade)
			if v.IsNil() {
				continue
			}
			v = v.Elem()
			fv := v.FieldByName(f.Name)
			if fv.IsNil() {
				continue
			}
			conf := fv.Interface()
			if registeredAlgorithm.ArbitrageTradeAlgorithmNewFunc == nil {
				continue
			}
			log.Printf("create %v arbitrage algorithm", name)
			newArbitrageTradeAlgoritm, err := registeredAlgorithm.ArbitrageTradeAlgorithmNewFunc(conf)
			if err != nil {
				r.DestroyArbitrageTradeAlgorithms(exchanges)
				return errors.Wrap(err, fmt.Sprintf("can not create arbitrage algorithm of %v", name))
			}
			err = newArbitrageTradeAlgoritm.Initialize(exchanges, r.notifier)
			if err != nil {
				r.DestroyArbitrageTradeAlgorithms(exchanges)
				return errors.Wrap(err, fmt.Sprintf("arbitrage algorithm initialize error of %v", name))
			}
			r.arbitrageTradeAlgorithms = append(r.arbitrageTradeAlgorithms, newArbitrageTradeAlgoritm)
		}
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
	Trade          *TradeConfig          `json:"trade" yaml:"trade" toml:"trade"`
	ArbitrageTrade *ArbitrageTradeConfig `json:"arbitrageTrade" yaml:"arbitrageTrade" toml:"arbitrageTrade"`
}

func NewRobot(config *Config, notifier *notifier.Notifier) (*Robot, error) {
	return &Robot{
		config:                   config,
		notifier:                 notifier,
		tradeAlgorithms:          make(map[string][]algorithm.TradeAlgorithm),
		arbitrageTradeAlgorithms: make([]algorithm.ArbitrageTradeAlgorithm, 0),
	}, nil
}
