package integrator

import (
	"github.com/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/braintree/manners"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"log"
	"time"
	"fmt"
	"reflect"
	"net/http"
)

type StartStreamingCallback func(tradeContext exchange.TradeContext, userCallbackData interface{}) (error)
type UpdateStreamingCallback func(tradeContext exchange.TradeContext, userCallbackData interface{}) (error)
type StopStreamingCallback func(tradeContext exchange.TradeContext, userCallbackData interface{}) (error)
type StartArbitrageCallback func(exchanges map[string]exchange.Exchange, userCallbackData interface{}) (error)
type UpdateArbitrageCallback func(exchanges map[string]exchange.Exchange, userCallbackData interface{}) (error)
type StopArbitrageCallback func(exchanges map[string]exchange.Exchange, userCallbackData interface{}) (error)

type gracefulServer struct {
	server    *manners.GracefulServer
	startChan chan error
}

type Integrator struct {
	config                  *Config
	gracefulServer          *gracefulServer
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

func (i *Integrator) setupRouting(engine *gin.Engine) {
	engine.HEAD( "/", i.index)
	engine.GET( "/", i.index)
}

func (i *Integrator) runHttpServer() {
	err := i.gracefulServer.server.ListenAndServe()
	if err != nil {
		i.gracefulServer.startChan <- err
	}
}

func (i *Integrator) initHttpServer() (error) {
	if i.config.AddrPort == "" {
		return nil
	}
	if !i.config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	i.setupRouting(engine)
	server := manners.NewWithServer(&http.Server{
		Addr:    i.config.AddrPort,
		Handler: engine,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	})
	i.gracefulServer = &gracefulServer{
		server: server,
		startChan: make(chan error),
	}
	go i.runHttpServer()
	select {
	case err := <- i.gracefulServer.startChan:
		return errors.Wrap(err, fmt.Sprintf("can not start http server (%s)", i.gracefulServer.server.Addr))
	case <-time.After(time.Second):
	}
	return nil
}

func (i *Integrator) streamingCallback(tradeContext exchange.TradeContext, userCallbackData interface{}) (error) {
	// bypassするだけ
	return i.updateStreamingCallback(tradeContext, i.userCallbackData)
}

func (i *Integrator) Initialize() (error) {
	err := i.initHttpServer()
	if err != nil {
		errors.Errorf("can not initalize of http server (reason = %v)", err)
	}
	for name, exchangeNewFunc := range exchange.GetRegisterdExchanges() {
		t := reflect.TypeOf(i.config.Exchanges).Elem()
		for idx := 0; idx < t.NumField(); idx++ {
			f := t.Field(idx)
			if f.Tag.Get("config") != name {
				continue
			}
			v := reflect.ValueOf(i.config.Exchanges)
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
	i.gracefulServer.server.BlockingClose()
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

type Config struct {
	Debug bool                 `json:"debug"     yaml:"debug"     toml:"debug"`
	AddrPort string            `json:"addrPort"  yaml:"addrPort"  toml:"addrPort"`
	Exchanges *ExchangesConfig `json:"exchanges" yaml:"exchanges" toml:"exchanges"`
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



