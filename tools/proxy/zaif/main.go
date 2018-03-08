package main

import (
	"runtime"
	"os"
	"path/filepath"
	"log"
	"flag"
	"path"
	"os/signal"
	"syscall"
	"log/syslog"
	actConfigurator "github.com/AutomaticCoinTrader/ACT/configurator"
	"github.com/AutomaticCoinTrader/ACT/tools/proxy/zaif/configurator"
	"github.com/AutomaticCoinTrader/ACT/tools/proxy/zaif/fetcher"
	"github.com/pkg/errors"
)

const (
	zaifProxyConfigPrefix string = "zaif-proxy"
)

func signalWait() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM)
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

func setupLogger(config *configurator.ZaifProxyConfig) (error) {
	log.SetFlags(log.Ldate|log.Ltime)
	if config.Logger != nil && config.Logger.Output == "syslog" {
		logger, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_USER, "ACT")
		if err != nil {
			return errors.Wrapf(err, "can not open syslog", err)
		}
		log.SetOutput(logger)
	}
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	wd, err := os.Getwd()
	if err == nil {
		abswd, err := filepath.Abs(wd)
		if err == nil {
			log.Printf("workdir: %v", abswd)
		} else {
			log.Printf("workdir: %v", wd)
		}
	}
	configDir := flag.String("confdir", "", "config directory")
	flag.Parse()
	cf, err := actConfigurator.NewConfigurator(path.Join(*configDir, zaifProxyConfigPrefix))
	if err != nil {
		log.Printf("can not create configurator (config dir = %v, reason = %v)", *configDir, err)
		return
	}
	newConfig := new(configurator.ZaifProxyConfig)
	err = cf.Load(newConfig)
	if err != nil {
		log.Printf("can not load config (config dir = %v, reason = %v)", *configDir, err)
		return

	}
	err = setupLogger(newConfig)
	if err != nil {
		log.Printf("can not setup logger (config dir = %v, reason = %v)", *configDir, err)
		return
	}
	f := fetcher.NewFetcher(newConfig)
	f.Start()
	signalWait()
	f.Stop()
}
