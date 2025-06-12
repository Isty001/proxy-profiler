package internal

import (
	"fmt"
	"os"

	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

type DestinationItem struct {
	Scheme string
	Host   string
}

type Config struct {
	Proxy struct {
		Port int
		Tls  struct {
			Cert string
			Key  string
		}
		Destination struct {
			InsecureSkipVerify         bool
			DefaultDestination         DestinationItem
			SourceToDestinationHostMap map[string]DestinationItem
		}
	}
	Metrics struct {
		RequestDurationMilliseconds struct {
			Buckets []float64
		}
	}
}

func ReadConfig(file string, defaultFile string) (Config, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigType("yaml")
	v.SetConfigFile(defaultFile)

	err := v.ReadInConfig()

	if err != nil {
		return Config{}, err
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		slog.Noticef("%s file not found, using the default settings")
	} else {
		v.SetConfigFile(file)
		v.MergeInConfig()
	}

	config := Config{}
	err = v.Unmarshal(&config)

	if err != nil {
		return Config{}, err
	}

	err = config.validate()
	slog.Infof("%+v", config)

	if err != nil {
		return Config{}, err
	}

	slog.Infof("%+v", config)

	return config, nil
}

func (c *Config) validate() error {
	var message string

	if 0 == c.Proxy.Port {
		message += "proxy.port must be a valid port; "
	}

	if "" == c.Proxy.Destination.DefaultDestination.Host {
		message += "proxy.destination.defaultDestination.host must be set; "
	}

	if "" == c.Proxy.Destination.DefaultDestination.Scheme {
		c.Proxy.Destination.DefaultDestination.Scheme = "http"
	}

	for i, dest := range c.Proxy.Destination.SourceToDestinationHostMap {
		if "" == dest.Host {
			message += "proxy.destination.sourceToDestinationHostMap." + i + ".host must be set"
		}

		if "" == dest.Scheme  {
			dest.Scheme = "http"

			c.Proxy.Destination.SourceToDestinationHostMap[i] = dest
		}
	}

	if "" != message {
		return fmt.Errorf("Invalid config: " + message)
	}

	return nil
}
