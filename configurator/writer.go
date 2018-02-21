package configurator

import (
	"github.com/pkg/errors"
	"github.com/BurntSushi/toml"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"io/ioutil"
	"bytes"
	"fmt"
	"os"
)

type writer struct {
}

func (r *writer) write(configPath string, data interface{}) (error) {
	ext := filepath.Ext(configPath)
	switch ext {
	case ".tml":
		fallthrough
	case ".toml":
		var buffer bytes.Buffer
		encoder := toml.NewEncoder(&buffer)
		err := encoder.Encode(data)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not encode with toml (%v)", configPath + ".tmp"))
		}
		err = ioutil.WriteFile(configPath + ".tmp", buffer.Bytes(), 0644)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not write file with toml (%v)", configPath + ".tmp"))
		}
	case ".yml":
		fallthrough
	case ".yaml":
		y, err := yaml.Marshal(data)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not encode with yaml (%v)", configPath + ".tmp"))
		}
		err = ioutil.WriteFile(configPath + ".tmp", y, 0644)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not write file with yaml (%v)", configPath + ".tmp"))
		}
	case ".jsn":
		fallthrough
	case ".json":
		j, err := json.Marshal(data)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not encode with json (%v)", configPath + ".tmp"))
		}
		err = ioutil.WriteFile(configPath + ".tmp", j, 0644)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not write file with json (%v)", configPath + ".tmp"))
		}
	default:
		return errors.Errorf("unexpected file extension (%v)", ext)
	}
	err := os.Rename(configPath + ".tmp", configPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not rename (%v)", configPath))
	}
	return nil
}

func newWriter() (*writer) {
	return &writer{}
}