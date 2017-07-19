package robot

import (
	"github.com/pkg/errors"
	"github.com/AutomaticCoinTrader/ACT/algorithm"
	"plugin"
	"io/ioutil"
	"path/filepath"
	"log"
	"fmt"
)

type GetRegistrationInfo func() (string, algorithm.TradeAlgorithmNewFunc, algorithm.ArbitrageTradeAlgorithmNewFunc)

func (r *Robot) registerAlgorithm(getRegistrationInfo GetRegistrationInfo) {
	name, tradeAlgorithmNewFunc, arbitrageTradeAlgorithmNewFunc := getRegistrationInfo()
	algorithm.RegisterAlgorithm(name, tradeAlgorithmNewFunc, arbitrageTradeAlgorithmNewFunc)
}

func (r *Robot) checkPluginSymbole(p *plugin.Plugin) (GetRegistrationInfo, error) {
	s, err := p.Lookup("GetRegistrationInfo")
	if err != nil {
		return nil, errors.Wrap(err, "not found GetRegistrationInfo symbole")
	}
	return s.(GetRegistrationInfo), nil
}

func (r *Robot) loadPluginFile(pluginFile string) (error) {
	p, err := plugin.Open(pluginFile)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not open plugin file (plugin file = %v)", pluginFile))
	}
	f, err := r.checkPluginSymbole(p)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("not plugin file (plugin file = %v)", pluginFile))
	}
	r.registerAlgorithm(f)
	return nil
}

func (r *Robot) loadPluginFiles() {
	if r.config == nil || r.config.AlgorithmPluginDir == "" {
		return
	}
	filist, err := ioutil.ReadDir(r.config.AlgorithmPluginDir)
	if err != nil {
		log.Printf("can not readdir (dir = %v)", r.config.AlgorithmPluginDir)
		return
	}
	for _, fi := range filist {
		if fi.IsDir() {
			continue
		}
		ext := filepath.Ext(fi.Name())
		if ext != ".so" && ext != ".dylib" {
			continue
		}
		pluginPath := filepath.Join(r.config.AlgorithmPluginDir, fi.Name())
		err := r.loadPluginFile(pluginPath)
		if err != nil {
			log.Printf("can not load plugin file (plugin file = %v)", pluginPath)
			continue
		}
	}
}

