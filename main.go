package main

import (
	"github.com/pkg/errors"
	"github.com/AutomaticCoinTrader/ACT/notifier"
	"github.com/AutomaticCoinTrader/ACT/robot"
	"github.com/AutomaticCoinTrader/ACT/integrator"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"github.com/AutomaticCoinTrader/ACT/configurator"
	"log"
	"runtime"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"path"
)

const (
	actConfigPrefix string = "act"
)

type config struct {
	Integrator *integrator.Config `json:"integrator" yaml:"integrator" toml:"integrator"`
	Robot      *robot.Config      `json:"robot"      yaml:"robot"      toml:"robot"`
	Notifier   *notifier.Config   `json:"notifier"   yaml:"notifier"   toml:"notifier"`
}

func signalWait() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
Loop:
	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGQUIT:
			fallthrough
		case syscall.SIGTERM:
			break Loop
		default:
			log.Printf("unexpected signal (sig = %v)", sig)
		}
	}
}

type actCallbackData struct  {
	robot *robot.Robot
}

func startStreamingCallback(tradeContext exchange.TradeContext, userCallbackData interface{}) (error) {
	actCallbackData := userCallbackData.(*actCallbackData)
	tradeID := tradeContext.GetID()
	err := actCallbackData.robot.CreateTradeAlgorithms(tradeID, tradeContext)
	if err != nil {
		log.Printf("can not create algorithm (reason = %v)", err)
	}
	return nil
}

func updateStreamingCallback(tradeContext exchange.TradeContext, userCallbackData interface{}) (error) {
	actCallbackData := userCallbackData.(*actCallbackData)
	tradeID := tradeContext.GetID()
	err := actCallbackData.robot.UpdateTradeAlgorithms(tradeID, tradeContext)
	if err != nil {
		log.Printf("can not run algorithm (reason = %v)", err)
	}
	return nil
}

func stopStreamingCallback(tradeContext exchange.TradeContext, userCallbackData interface{}) (error) {
	actCallbackData := userCallbackData.(*actCallbackData)
	tradeID := tradeContext.GetID()
	err := actCallbackData.robot.DestroyTradeAlgorithms(tradeID, tradeContext)
	if err != nil {
		log.Printf("can not destroy algorithm (reason = %v)", err)
	}
	return nil
}

func startArbitrageCallback(exchanges map[string]exchange.Exchange, userCallbackData interface{}) (error) {
	actCallbackData := userCallbackData.(*actCallbackData)
	err := actCallbackData.robot.CreateArbitrageTradeAlgorithms(exchanges)
	if err != nil {
		log.Printf("can not create algorithm (reason = %v)", err)
	}
	return nil
}

func updateArbitrageCallback(exchanges map[string]exchange.Exchange, userCallbackData interface{}) (error) {
	actCallbackData := userCallbackData.(*actCallbackData)
	err := actCallbackData.robot.UpdateArbitrageTradeAlgorithms(exchanges)
	if err != nil {
		log.Printf("can not run algorithm (reason = %v)", err)
	}
	return nil
}

func stopArbitrageCallback(exchanges map[string]exchange.Exchange, userCallbackData interface{}) (error) {
	actCallbackData := userCallbackData.(*actCallbackData)
	err := actCallbackData.robot.DestroyArbitrageTradeAlgorithms(exchanges)
	if err != nil {
		log.Printf("can not destroy algorithm (reason = %v)", err)
	}
	return nil
}

func actStart(integrator *integrator.Integrator) (error){
	err := integrator.Initialize()
	if err != nil {
		return errors.Wrap(err, "initalize error of integrator")
	}
	err = integrator.StartArbitrage()
	if err != nil {
		return errors.Wrap(err, "can not start arbitarage")
	}
	err = integrator.StartStreaming()
	if err != nil {
		return errors.Wrap(err, "can not start streaming")
	}
	return nil
}

func actStop(integrator *integrator.Integrator) (error) {
	err := integrator.StopStreaming()
	if err != nil {
		log.Printf("can not stop streaming (reason = %v)", err)
	}
	err = integrator.StopArbitrageTrade()
	if err != nil {
		log.Printf("can not stop arbitarage (reason = %v)", err)
	}
	err = integrator.Finalize()
	if err != nil {
		log.Printf("finalize error of integrator(reazon = %v)", err)
	}
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	configDir := flag.String("confdir", "", "config directory")
	flag.Parse()
	cf, err := configurator.NewConfigurator(path.Join(*configDir, actConfigPrefix))
	if err != nil {
		log.Printf("can not create configurator (config dir = %v, reason = %v)", *configDir, err)
		return
	}
	newConfig := new(config)
	err = cf.Load(newConfig)
	if err != nil {
		log.Printf("can not load config (config dir = %v, reason = %v)", *configDir, err)
		return

	}
	ntf, err := notifier.NewNotifier(newConfig.Notifier)
	if err != nil {
		log.Printf("can not create notifier (config dir = %v, reason = %v)", *configDir, err)
		return
	}
	rbt, err := robot.NewRobot(newConfig.Robot, *configDir, ntf)
	if err != nil {
		log.Printf("can not create robot (config dir = %v, reason = %v)", *configDir, err)
		return
	}
	actCallbackData := &actCallbackData{
		robot: rbt,
	}
	it, err := integrator.NewIntegrator(
		newConfig.Integrator,
		startStreamingCallback,
		updateStreamingCallback,
		stopStreamingCallback,
		startArbitrageCallback,
		updateArbitrageCallback,
		stopArbitrageCallback,
		actCallbackData)
	if err != nil {
		log.Printf("can not create exchangers (config dir = %v, reason = %v)", *configDir, err)
		return
	}
	err = actStart(it)
	if err != nil {
		log.Printf("can not start act (reason = %v)", err)
		return
	}
	signalWait()
	err = actStop(it)
	if err != nil {
		log.Printf("can not stop act (reason = %v)", err)
	}
}

