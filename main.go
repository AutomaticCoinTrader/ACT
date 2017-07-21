package main

import (
	"github.com/pkg/errors"
	"github.com/AutomaticCoinTrader/ACT/notifier"
	"github.com/AutomaticCoinTrader/ACT/robot"
	"github.com/AutomaticCoinTrader/ACT/integrator"
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

func actStart(integrator *integrator.Integrator) (error){
	err := integrator.Initialize()
	if err != nil {
		return errors.Wrap(err, "initalize error of integrator")
	}
	err = integrator.Start()
	if err != nil {
		return errors.Wrap(err, "can not start intergrator")
	}
	return nil
}

func actStop(integrator *integrator.Integrator) (error) {
	err := integrator.Stop()
	if err != nil {
		log.Printf("can not stop integrator (reason = %v)", err)
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
	it, err := integrator.NewIntegrator(
		newConfig.Integrator,
		newConfig.Notifier,
		newConfig.Robot,
		*configDir)
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

