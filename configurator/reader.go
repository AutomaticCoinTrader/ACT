package configurator

import (
	"github.com/pkg/errors"
	"github.com/BurntSushi/toml"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type reader struct {
}

func (r *reader) read(configPath string, data interface{}) (error) {
	ext := filepath.Ext(configPath)
	switch ext {
	case ".tml":
		fallthrough
	case ".toml":
		_, err := toml.DecodeFile(configPath, data)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not decode file with toml (%v)", configPath))
		}
	case ".yml":
		fallthrough
	case ".yaml":
		buf, err := ioutil.ReadFile(configPath)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not read file with yaml (%v)", configPath))
		}
		err = yaml.Unmarshal(buf, data)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not decode file with yaml (%v)", configPath))
		}
	case ".jsn":
		fallthrough
	case ".json":
		buf, err := ioutil.ReadFile(configPath)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not read file with json (%v)", configPath))
		}
		err = json.Unmarshal(buf, data)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not decode file with json (%v)", configPath))
		}
	default:
		return errors.Errorf("unexpected file extension (%v)", ext)
	}
	return nil
}

func newReader() (*reader) {
	return &reader{}
}

