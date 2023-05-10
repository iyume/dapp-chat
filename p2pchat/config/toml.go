package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
	"github.com/naoina/toml"
)

type gethConfig struct {
	Eth  ethconfig.Config
	Node node.Config
	// Ethstats and Metrics will never be enabled
}

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		if field == "Ethstats" || field == "Metrics" {
			return fmt.Errorf("field '%s' is not supported", field)
		}
		return fmt.Errorf("field '%s' is not defined in %s", field, rt.String())
	},
}

// LoadTomlConfig is unstable,
// See https://github.com/celo-org/celo-blockchain/issues/1066
func LoadTomlConfig(file string) (*gethConfig, error) {
	cfg := &gethConfig{ethconfig.Defaults, node.DefaultConfig}
	f, err := os.Open(file)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	err = tomlSettings.NewDecoder(f).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return cfg, err
}
