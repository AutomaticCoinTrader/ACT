package main

import (
	"runtime"
	"os"
	"path/filepath"
	"log"
	"flag"
	"path"
	actConfigurator "github.com/AutomaticCoinTrader/ACT/configurator"
	"github.com/AutomaticCoinTrader/ACT/tools/proxy/zaif/configurator"
	"os/signal"
	"syscall"
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



	signalWait()


}
