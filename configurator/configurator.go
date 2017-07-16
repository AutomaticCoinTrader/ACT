package configurator

import (
	"github.com/pkg/errors"
	"os"
	"fmt"
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

// NewConfigurator is create Configurator
func NewConfigurator(configPath string) (configurator *Configurator, err error) {
	_, err = os.Stat(configPath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("not exists config file (%v)", configPath))
	}
	return &Configurator{
		reader : newReader(),
		configPath : configPath,
	}, nil
}
