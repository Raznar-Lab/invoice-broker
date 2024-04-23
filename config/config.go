package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	filepath      string
	Web           WebConfig           `yaml:"web" validate:"required"`
	Gateway       GatewayConfig       `yaml:"gateway" validate:"required"`
	// MessageBroker MessageBrokerConfig `yaml:"message_broker" validate:"required"`
}

var validate = validator.New()
var templateFile = "config.example.yml"

func New(filePath string) (config *Config, err error) {
	config = &Config{
		filepath: filePath,
	}

	if err = config.Reload(); err != nil {
		config = nil
		return
	}

	err = validate.Struct(config)
	return
}

func (c *Config) Parse() (err error) {
	if c.Web.Bind == "" {
		c.Web.Bind = "0.0.0.0"
	}

	return
}

func (c *Config) Reload() (err error) {
	file, err := os.OpenFile(c.filepath, os.O_RDONLY, os.ModePerm)
	if os.IsNotExist(err) {
		tContent, err := os.ReadFile(templateFile)
		if err != nil {
			return err
		}

		err = os.WriteFile(c.filepath, tContent, 0644)
		if err != nil {
			return err
		}

		return c.Reload()
	}

	if err != nil {
		return
	}

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&c); err != nil {
		return
	}

	return c.Parse()
}

func (c *Config) Save() (err error) {
	content, err := yaml.Marshal(&c)
	if err != nil {
		return
	}

	return os.WriteFile(c.filepath, content, 0644)
}
