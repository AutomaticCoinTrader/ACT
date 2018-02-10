package robot

import (
	"github.com/pkg/errors"
	"github.com/AutomaticCoinTrader/ACT/algorithm"
	"plugin"
	"io/ioutil"
	"path/filepath"
	"log"
	"fmt"
	"regexp"
	"os/user"
)

type GetRegistrationInfoType func() (string, algorithm.InternalTradeAlgorithmNewFunc, algorithm.ExternalTradeAlgorithmNewFunc)

func (r *Robot) registerAlgorithm(getRegistrationInfo GetRegistrationInfoType) {
	name, tradeAlgorithmNewFunc, arbitrageTradeAlgorithmNewFunc := getRegistrationInfo()
	algorithm.RegisterAlgorithm(name, tradeAlgorithmNewFunc, arbitrageTradeAlgorithmNewFunc)
}

func (r *Robot) checkPluginSymbole(p *plugin.Plugin) (GetRegistrationInfoType, error) {
	s, err := p.Lookup("GetRegistrationInfo")
	if err != nil {
		return nil, errors.Wrap(err, "not found GetRegistrationInfo symbole")
	}
	return s.(func() (string, algorithm.InternalTradeAlgorithmNewFunc, algorithm.ExternalTradeAlgorithmNewFunc)), nil
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

func (r *Robot) fixupAlgorithmPluginDir(algorithmPluginDir string) (string) {
	// shell表現 "~/" をなんとかする
	u, err := user.Current()
	if err != nil {
		log.Printf("can not get user info (reason = %v)", err)
		return algorithmPluginDir
	}
	re := regexp.MustCompile("^~/")
	return re.ReplaceAllString(r.config.AlgorithmPluginDir, u.HomeDir+"/")
}

func (r *Robot) loadPluginFiles() {
	if r.config == nil || r.config.AlgorithmPluginDir == "" {
		return
	}
	algorithmPluginDir := r.fixupAlgorithmPluginDir(r.config.AlgorithmPluginDir)
	filist, err := ioutil.ReadDir(algorithmPluginDir)
	if err != nil {
		log.Printf("can not readdir (dir = %v)", algorithmPluginDir)
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
		pluginPath := filepath.Join(algorithmPluginDir, fi.Name())
		err := r.loadPluginFile(pluginPath)
		if err != nil {
			log.Printf("can not load plugin file (plugin file = %v, reason = %v)", pluginPath, err)
			continue
		}
	}
}
