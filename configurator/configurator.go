package configurator

import (
	"github.com/pkg/errors"
	"os"
	"path"
)

// Configurator is struct of configurator
type Configurator struct {
	configPath string
	reader     *reader
}

// Load is load config
func (c *Configurator) Load(data interface{}) (err error) {
	return c.reader.read(c.configPath, data)
}

func checkConfigFilePath(configFilePathPrefix string) (string, error) {
	for _, extension := range []string{".yaml", ".yml", "toml", "tml", ".json", "jsn"} {
		configFilePath := path.Join(configFilePathPrefix + extension)
		_, err := os.Stat(configFilePath)
		if err != nil {
			continue
		}
		return configFilePath, nil
	}
	return "", errors.Errorf("can not found config file (config file path prefix = %v)", configFilePathPrefix)
}

// NewConfigurator is create Configurator
func NewConfigurator(configFilePathPrefix string) (*Configurator, error) {
	configFilePath, err := checkConfigFilePath(configFilePathPrefix)
	if err != nil {
		return nil, err
	}
	return &Configurator{
		reader : newReader(),
		configPath : configFilePath,
	}, nil
}
