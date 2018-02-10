package integrator

import (
	"github.com/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/braintree/manners"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"github.com/AutomaticCoinTrader/ACT/robot"
	"log"
	"time"
	"fmt"
	"reflect"
	"net/http"
	"github.com/AutomaticCoinTrader/ACT/notifier"
)

type StartStreamingCallback func(ex exchange.Exchange, userCallbackData interface{}) (error)
type UpdateStreamingCallback func(ex exchange.Exchange, userCallbackData interface{}) (error)
type StopStreamingCallback func(ex exchange.Exchange, userCallbackData interface{}) (error)
type StartArbitrageCallback func(exs map[string]exchange.Exchange, userCallbackData interface{}) (error)
type UpdateArbitrageCallback func(exs map[string]exchange.Exchange, userCallbackData interface{}) (error)
type StopArbitrageCallback func(exs map[string]exchange.Exchange, userCallbackData interface{}) (error)

type gracefulServer struct {
	server    *manners.GracefulServer
	startChan chan error
}

type Integrator struct {
	config                  *Config
	gracefulServer          *gracefulServer
	exchanges               map[string]exchange.Exchange
	arbitrageLoopFinishChan chan bool
	notifier                *notifier.Notifier
	robot                   *robot.Robot
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
	if i.config.Server.AddrPort == "" {
		return nil
	}
	if !i.config.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	i.setupRouting(engine)
	server := manners.NewWithServer(&http.Server{
		Addr:    i.config.Server.AddrPort,
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

func (i *Integrator) streamingCallback(currencyPair string, ex exchange.Exchange) (error) {
	// トレード処理を期待
	err := i.robot.UpdateInternalTradeAlgorithms(currencyPair, ex)
	if err != nil {
		log.Printf("can not run algorithm (reason = %v)", err)
	}
	return nil
}

func (i *Integrator) Initialize() (error) {
	err := i.initHttpServer()
	if err != nil {
		errors.Errorf("can not initalize of http server (reason = %v)", err)
	}
	for name, exchangeNewFunc := range exchange.GetRegisterdExchanges() {
		log.Printf("listing exchange: %s\n", name)
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
			ex, err :=  exchangeNewFunc(conf, )
			if err != nil {
				i.Finalize()
				return errors.Wrap(err, fmt.Sprintf("can not create exchange of %v", name))
			}
			ex.Initialize()
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

func (i *Integrator) startInternalTrade() (error) {
	for _, ex := range i.exchanges {

		// streamingを始める前の前処理を期待
		err := i.robot.CreateInternalTradeAlgorithms(ex)
		if err != nil {
			i.stopInternalTrade()
			return errors.Wrap(err, fmt.Sprintf("can not create algorithm  (name = %v)", ex.GetName()))
		}
		// ストリーミングを開始
		err = ex.StartStreamings(i.streamingCallback)
		if err != nil {
			i.stopInternalTrade()
			return errors.Wrap(err, fmt.Sprintf("can not start streaming (name = %v)", ex.GetName()))
		}
	}

	return nil
}

func (i *Integrator) stopInternalTrade() (error) {
	// 取引所を停止する処理
	for _, ex := range i.exchanges {

		// streamingを停止
		err := ex.StopStreamings()
		if err != nil {
			log.Printf("can not stop streaming (name = %v)", ex.GetName())
		}
		// straming止めた後の終了処理を期待
		err = i.robot.DestroyInternalTradeAlgorithms(ex)
		if err != nil {
			log.Printf("can not destroy algorithm (name = %v, reason = %v)", ex.GetName(), err)
		}
	}
	return nil
}

func (i *Integrator) externalTradeEventLoop(){
	for {
		select {
		case <- i.arbitrageLoopFinishChan:
			return
		case <- time.After(500 * time.Millisecond):
			err := i.robot.UpdateExternalTradeAlgorithms(i.exchanges)
			if err != nil {
				log.Printf("can not update external trade algorithm (reason = %v)", err)
			}
		}
	}
}

func (i *Integrator) startExternalTrade() (error) {
	err := i.robot.CreateExternalTradeAlgorithms(i.exchanges)
	if err != nil {
		return errors.Wrap(err,"can not create external trade algorithm")
	}
	go i.externalTradeEventLoop()
	return nil
}

func (i *Integrator) stopExternalTrade() (error) {
	close(i.arbitrageLoopFinishChan)
	err := i.robot.DestroyExternalTradeAlgorithms(i.exchanges)
	if err != nil {
		log.Printf("can not destroy external trade algorithm (reason = %v)", err)
	}
	return nil
}

func (i *Integrator) Start() (error) {
	err := i.startInternalTrade()
	if err != nil {
		return errors.Wrap(err, "can not start streaming")
	}
	err = i.startExternalTrade()
	if err != nil {
		return errors.Wrap(err, "can not start external trade")
	}
	return nil
}

func (i *Integrator) Stop() (error) {
	err := i.stopExternalTrade()
	if err != nil {
		log.Printf("can not stop extenal trade (reason = %v)", err)
	}
	err = i.stopInternalTrade()
	if err != nil {
		log.Printf("can not stop internal trade (reason = %v)", err)
	}
	return nil
}

type serverConfig struct {
	Debug bool                 `json:"debug"     yaml:"debug"     toml:"debug"`
	AddrPort string            `json:"addrPort"  yaml:"addrPort"  toml:"addrPort"`
}

type Config struct {
	Server    *serverConfig     `json:"server"    yaml:"server"    toml:"server"`
	Exchanges *exchangesConfig  `json:"exchanges" yaml:"exchanges" toml:"exchanges"`
	Robot      *robot.Config    `json:"robot"     yaml:"robot"     toml:"robot"`
	Notifier   *notifier.Config `json:"notifier"  yaml:"notifier"  toml:"notifier"`
}

func NewIntegrator(config *Config, configDir string) (*Integrator, error) {
	ntf, err := notifier.NewMailNotifier(config.Notifier)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can not create notifier (config dir = %v, reason = %v)", configDir, err))
	}
	rbt, err := robot.NewRobot(config.Robot, configDir, ntf)
	if err != nil {
		return nil, errors.Wrap(err,fmt.Sprintf("can not create robot (config dir = %v, reason = %v)", configDir, err))
	}
	return &Integrator{
		config: config,
		exchanges: make(map[string]exchange.Exchange),
		arbitrageLoopFinishChan: make(chan bool),
		notifier: ntf,
		robot: rbt,
	}, nil
}



