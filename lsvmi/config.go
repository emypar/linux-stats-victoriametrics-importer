// Configuration for the stats collector

package lsvmi

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

// The configuration is stored in one object to make it easy to load it from a
// file. Most of the configuration parameters are based on the file settings and
// a few can be overridden by command line arguments.
//
// The decreasing order of precedence for parameter values:
//   - command line arg (if applicable)
//   - config file
//   - built-in default
//
// Notes:
//  1. Each component will have its specific configuration, which may be
//     defined elsewhere, for instance in the files(s) providing the implementation.
// 2. The object will be loaded from a YAML file, therefore all configuration
//    parameters should be public and they should have tag annotations.

type LsvmiConfig struct {
	HttpEndpointPoolConfig *HttpEndpointPoolConfig `yaml:"http_endpoint_pool_config"`
	LoggerConfig           *LoggerConfig           `yaml:"log_config"`
}

var lsvmiConfigFile = flag.String(
	"config",
	"",
	`Config file to load`,
)

var LsvmiCfg *LsvmiConfig

func NewLsvmiConfig() *LsvmiConfig {
	return &LsvmiConfig{
		HttpEndpointPoolConfig: DefaultHttpEndpointPoolConfig(),
	}
}

func LoadLsvmiConfig(cfgFile string) (*LsvmiConfig, error) {
	f, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	cfg := NewLsvmiConfig()
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("file: %q: %v", cfgFile, err)
	}
	return cfg, nil
}

func LoadLsvmiConfigFromArgs() error {
	var err error
	if *lsvmiConfigFile != "" {
		var cfg *LsvmiConfig
		cfg, err = LoadLsvmiConfig(*lsvmiConfigFile)
		if err == nil {
			LsvmiCfg = cfg
		}
	} else {
		err = fmt.Errorf("config file arg not provided")
	}
	return err
}