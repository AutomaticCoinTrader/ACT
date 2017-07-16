package integrator

import (
	"github.com/pkg/errors"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"log"
	"time"
	"fmt"
	"reflect"
)

type StartStreamingCallback func(tradeContext exchange.TradeContext, userCallbackData interface{}) (error)
type UpdateStreamingCallback func(tradeContext exchange.TradeContext, userCallbackData interface{}) (error)
type StopStreamingCallback func(tradeContext exchange.TradeContext, userCallbackData interface{}) (error)
type StartArbitrageCallback func(exchanges map[string]exchange.Exchange, userCallbackData interface{}) (error)
type UpdateArbitrageCallback func(exchanges map[string]exchange.Exchange, userCallbackData interface{}) (error)
type StopArbitrageCallback func(exchanges map[string]exchange.Exchange, userCallbackData interface{}) (error)

type Integrator struct {
	config                  *Config
	exchanges               map[string]exchange.Exchange
	arbitrageLoopFinishChan chan bool
	startStreamingCallback  StartStreamingCallback
	updateStreamingCallback UpdateStreamingCallback
	stopStreamingCallback   StopStreamingCallback
	startArbitrageCallback  StartArbitrageCallback
	updateArbitrageCallback UpdateArbitrageCallback
	stopArbitrageCallback   StopArbitrageCallback
	userCallbackData        interface{}
}

func (i *Integrator) streamingCallback(tradeContext exchange.TradeContext, userCallbackData interface{}) (error) {
	// bypassするだけ
	return i.updateStreamingCallback(tradeContext, i.userCallbackData)
}

func (i *Integrator) Initialize() (error) {
	for name, exchangeNewFunc := range exchange.GetRegisterdExchanges() {
		t := reflect.TypeOf(i.config).Elem()
		for idx := 0; idx < t.NumField(); idx++ {
			f := t.Field(idx)
			if f.Tag.Get("config") != name {
				continue
			}
			v := reflect.ValueOf(i.config)
			if v.IsNil() {
				continue
			}
			v = v.Elem()
			fv := v.FieldByName(f.Name)
			if fv.IsNil() {
				continue
			}
			conf := fv.Interface()
			if exchangeNewFunc == nil {
				continue
			}
			log.Printf("%v exchange create", name)
			ex, err :=  exchangeNewFunc(conf)
			if err != nil {
				i.Finalize()
				return errors.Wrap(err, fmt.Sprintf("can not create exchange of %v", name))
			}
			ex.Initialize(i.streamingCallback, nil)
			// 作った取引所を保存しておく
			i.exchanges[name] = ex
		}
	}
	return nil
}

func (i *Integrator) Finalize() (error) {
	return nil
}

func (i *Integrator) StartStreaming() (error) {
	for _, ex := range i.exchanges {
		tradeContextCursor := ex.GetTradeContextCursor()
		for {
			tradeContext, ok := tradeContextCursor.Next()
			if !ok {
				break
			}
			// callbackを呼び出す
			// streamingを始める前の前処理を期待
			err := i.startStreamingCallback(tradeContext, i.userCallbackData)
			if err != nil {
				i.StopStreaming()
				return errors.Wrap(err, fmt.Sprintf("start straming callback error (name = %v)", ex.GetName()))
			}
			// ストリーミングを開始
			err = ex.StartStreaming(tradeContext)
			if err != nil {
				i.StopStreaming()
				return errors.Wrap(err, fmt.Sprintf("can not start streaming (name = %v)", ex.GetName()))
			}
		}
	}

	return nil
}

func (i *Integrator) StopStreaming() (error) {
	// 取引所を停止する処理
	for _, ex := range i.exchanges {
		tradeContextCursor := ex.GetTradeContextCursor()
		for {
			tradeContext, ok := tradeContextCursor.Next()
			if !ok {
				break
			}
			// streamingを停止
			err := ex.StopStreaming(tradeContext)
			if err != nil {
				log.Printf("can not stop streaming (name = %v)", ex.GetName())
			}
			// callbackを呼び出す
			// straming止めた後の終了処理を期待
			err = i.stopStreamingCallback(tradeContext, i.userCallbackData)
			if err != nil {
				log.Printf("stop straming callback error (name = %v)", ex.GetName())
			}
		}
	}
	return nil
}

func (i *Integrator) ArbitrageLoop (){
	for {
		select {
		case <- i.arbitrageLoopFinishChan:
			return
		case <- time.After(500 * time.Millisecond):
			i.updateArbitrageCallback(i.exchanges, i.userCallbackData)
		}
	}
}

func (i *Integrator) StartArbitrage() (error) {
	err := i.startArbitrageCallback(i.exchanges, i.userCallbackData )
	if err != nil {
		return errors.Wrap(err,"start arbitrage callback error")
	}
	go i.ArbitrageLoop()
	return nil
}

func (i *Integrator) StopArbitrageTrade() (error) {
	close(i.arbitrageLoopFinishChan)
	err := i.stopArbitrageCallback(i.exchanges, i.userCallbackData )
	if err != nil {
		log.Printf("stop arbitrage callback error")
	}
	return nil
}

func NewIntegrator(config *Config,
	startStreamingCallback StartStreamingCallback,
	updateStreamingCallback UpdateStreamingCallback,
	stopStreamingCallback StopStreamingCallback,
	startArbitrageCallback StartArbitrageCallback,
	updateArbitrageCallback UpdateArbitrageCallback,
	stopArbitrageCallback StopArbitrageCallback,
	userCallbackData interface{}) (*Integrator, error) {
	return &Integrator{
		config: config,
		exchanges: make(map[string]exchange.Exchange),
		arbitrageLoopFinishChan: make(chan bool),
		startStreamingCallback: startStreamingCallback,
		updateStreamingCallback: updateStreamingCallback,
		stopStreamingCallback: stopStreamingCallback,
		startArbitrageCallback: startArbitrageCallback,
		updateArbitrageCallback: updateArbitrageCallback,
		stopArbitrageCallback: stopArbitrageCallback,
		userCallbackData: userCallbackData,
	}, nil
}



