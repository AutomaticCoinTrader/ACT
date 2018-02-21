package configurator

import (
	"github.com/pkg/errors"
	"github.com/theckman/go-flock"
	"os"
	"path"
	"os/user"
	"log"
	"regexp"
	"time"
)

// Configurator is struct of configurator
type Configurator struct {
	configPath string
	reader     *reader
	writer     *writer
	lock       *flock.Flock
}

// GetConfigPath is get config path
func (c *Configurator) GetConfigPath() (string) {
	return c.configPath
}

// Lock is lock config file
func (c *Configurator) Lock() (error) {
	for {
		if !c.lock.Locked() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return c.lock.Lock()
}

// Unlock is unlock config file
func (c *Configurator) Unlock() (error) {
	return c.lock.Unlock()
}

// Load is load config
func (c *Configurator) Load(data interface{}) (error) {
	return c.reader.read(c.configPath, data)
}

// Save is save config
func (c *Configurator) Save(data interface{}) (error) {
	return c.writer.write(c.configPath, data)
}

func fixupConfigFilePathPrefix(configFilePathPrefix string) (string) {
	// shell表現 "~/" をなんとかする
	u, err := user.Current()
	if err != nil {
		log.Printf("can not get user info (reason = %v)", err)
		return configFilePathPrefix
	}
	re := regexp.MustCompile("^~/")
	return re.ReplaceAllString(configFilePathPrefix, u.HomeDir+"/")
}

func checkConfigFilePath(configFilePathPrefix string) (string, error) {
	for _, extension := range []string{".yaml", ".yml", ".toml", ".tml", ".json", ".jsn"} {
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
	configFilePathPrefix = fixupConfigFilePathPrefix(configFilePathPrefix)
	configFilePath, err := checkConfigFilePath(configFilePathPrefix)
	if err != nil {
		return nil, err
	}
	return &Configurator{
		reader:  newReader(),
		writer : newWriter(),
		lock: flock.NewFlock(configFilePath + ".lock"),
		configPath: configFilePath,
	}, nil
}
