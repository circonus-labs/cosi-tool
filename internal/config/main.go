// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

import (
	"encoding/json"
	"expvar"
	"fmt"
	"io"

	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// Broker defines the broker and type to use when creating a check
type Broker struct {
	ID   string `json:"id" toml:"id" yaml:"id"`
	Type string `json:"type" toml:"type" yaml:"type"`
}

// Log defines the running config.log structure
type Log struct {
	Level  string `json:"level" toml:"level" yaml:"level"`
	Pretty bool   `json:"pretty" toml:"pretty" yaml:"pretty"`
}

// API defines the api settings
type API struct {
	App    string `json:"app" toml:"app" yaml:"app"`
	CAFile string `mapstructure:"ca_file" json:"ca_file" toml:"ca_file" yaml:"ca_file"`
	Key    string `json:"key" toml:"key" yaml:"key"`
	URL    string `json:"url" toml:"url" yaml:"url"`
}

// Agent defines the agent settings
type Agent struct {
	Mode string `json:"mode" toml:"mode" yaml:"mode"`
	URL  string `json:"url" toml:"url" yaml:"url"`
}

// System defines the various system settings
type System struct {
	Arch      string `json:"arch" toml:"arch" yaml:"arch"`
	DMI       string `json:"dmi" toml:"dmi" yaml:"dmi"`
	OSDistro  string `mapstructure:"os_dist" json:"os_dist" toml:"os_dist" yaml:"os_dist"`
	OSType    string `mapstructure:"os_type" json:"os_type" toml:"os_type" yaml:"os_type"`
	OSVersion string `mapstructure:"os_vers" json:"os_vers" toml:"os_vers" yaml:"os_vers"`
}

// Checks defines the host specific settings
type Checks struct {
	Broker  Broker `json:"broker" toml:"broker" yaml:"broker"`
	GroupID string `mapstructure:"group_id" json:"group_id" toml:"group_id" yaml:"group_id"`
	Target  string `json:"target" toml:"target" yaml:"target"`
}

// Config defines the cosi configuration (created by cosi-install or manually)
type Config struct {
	Agent     Agent  `json:"agent" toml:"agent" yaml:"agent"`
	API       API    `json:"api" toml:"api" yaml:"api"`
	BaseUIURL string `mapstructure:"base_ui_url" json:"base_ui_url" toml:"base_ui_url" yaml:"base_ui_url"`
	CosiURL   string `mapstructure:"cosi_url" json:"cosi_url" toml:"cosi_url" yaml:"cosi_url"`
	Debug     bool   `json:"debug" toml:"debug" yaml:"debug"`
	Checks    Checks `json:"checks" toml:"checks" yaml:"checks"`
	Log       Log    `json:"log" toml:"log" yaml:"log"`
	RegConf   string `mapstructure:"reg_conf" json:"reg_conf" toml:"reg_conf" yaml:"reg_conf"`
	System    System `json:"system" toml:"system" yaml:"system"`
}

// // OldConfig defines the old cosi configuration file format, for conversion
// type OldConfig struct {
// 	// backwards compatible options
// 	APIKey      string `mapstructure:"api_key" json:"api_key" yaml:"api_key" toml:"api_key"`
// 	APIApp      string `mapstructure:"api_app" json:"api_app" yaml:"api_app" toml:"api_app"`
// 	APIURL      string `mapstructure:"api_url" json:"api_url" yaml:"api_url" toml:"api_url"`
// 	CAFile      string `mapstructure:"api_ca_file" json:"api_ca_file" yaml:"api_ca_file" toml:"api_ca_file"`
// 	CosiURL     string `mapstructure:"cosi_url" json:"cosi_url" yaml:"cosi_url" toml:"cosi_url"`
// 	AgentMode   string `mapstructure:"agent_mode" json:"agent_mode" yaml:"agent_mode" toml:"agent_mode"`
// 	AgentURL    string `mapstructure:"agent_url" json:"agent_url" yaml:"agent_url" toml:"agent_url"`
// 	OptsFile    string `mapstructure:"custom_options_file" json:"custom_options_file" yaml:"custom_options_file" toml:"custom_options_file"`
// 	CheckTarget string `mapstructure:"cosi_host_target" json:"cosi_host_target" yaml:"cosi_host_target" toml:"cosi_host_target"`
// 	BrokerID    string `mapstructure:"cosi_broker_id" json:"cosi_broker_id" yaml:"cosi_broker_id" toml:"cosi_broker_id"`
// 	BrokerType  string `mapstructure:"cosi_broker_type" json:"cosi_broker_type" yaml:"cosi_broker_type" toml:"cosi_broker_type"`
// 	OSType      string `mapstructure:"cosi_os_type" json:"cosi_os_type" yaml:"cosi_os_type" toml:"cosi_os_type"`
// 	OSDistro    string `mapstructure:"cosi_os_dist" json:"cosi_os_dist" yaml:"cosi_os_dist" toml:"cosi_os_dist"`
// 	OSVers      string `mapstructure:"cosi_os_vers" json:"cosi_os_vers" yaml:"cosi_os_vers" toml:"cosi_os_vers"`
// 	SysArch     string `mapstructure:"cosi_os_arch" json:"cosi_os_arch" yaml:"cosi_os_arch" toml:"cosi_os_arch"`
// 	DMI         string `mapstructure:"cosi_os_dmi" json:"cosi_os_dmi" yaml:"cosi_os_dmi" toml:"cosi_os_dmi"`
// 	GroupID     string `mapstructure:"cosi_group_id" json:"cosi_group_id" yaml:"cosi_group_id" toml:"cosi_group_id"`
// 	UIBaseURL   string `mapstructure:"_ui_base_url" json:"_ui_base_url" yaml:"_ui_base_url" toml:"_ui_base_url"`
// }

const (
	// KeyAPITokenKey circonus api token key
	KeyAPITokenKey = "api.key"

	// KeyAPITokenApp circonus api token key application name
	KeyAPITokenApp = "api.app"

	// KeyAPIURL custom circonus api url (e.g. inside)
	KeyAPIURL = "api.url"

	// KeyAPICAFile custom ca for circonus api (e.g. inside)
	KeyAPICAFile = "api.ca_file"

	// KeyCosiURL for cosi server
	KeyCosiURL = "cosi_url"

	// KeyUIBaseURL base ui url for account
	KeyUIBaseURL = "base_ui_url"

	// KeyAgentMode defines the mode of the agent
	KeyAgentMode = "agent.mode"

	// KeyAgentURL defines the url of the local agent
	KeyAgentURL = "agent.url"

	// // KeyHostOptionsFile defines a json file with custom options
	// KeyHostOptionsFile = "host.options_file"

	// KeyHostTarget defines the target host for the check
	KeyHostTarget = "checks.target"

	// KeyHostBrokerID defines the broker to use when creating a check
	KeyHostBrokerID = "checks.broker.id"

	// KeyHostBrokerType defines the 'type' of broker to use (any|enterprise)
	KeyHostBrokerType = "checks.broker.type"

	// KeyHostGroupID defines the group ID (if this system will be used
	// in a group for statsd metrics, this ID will be used to find/create
	// the group check)
	KeyHostGroupID = "checks.group_id"

	// KeySystemOSType defines the type of OS
	KeySystemOSType = "system.os_type"

	// KeySystemOSDistro defines the OS distribution
	KeySystemOSDistro = "system.os_dist"

	// KeySystemOSVersion defines the OS distribution version
	KeySystemOSVersion = "system.os_vers"

	// KeySystemArch defines the system architecture
	KeySystemArch = "system.arch"

	// KeySystemDMI defines the dmi (aws thing)
	KeySystemDMI = "system.dmi"

	//
	// Registration
	//

	// KeyRegConf defines the registration options configuration file
	KeyRegConf = "register.config"

	//
	// generic flags
	//

	// KeyDebug enables debug messages
	KeyDebug = "debug"

	// KeyLogLevel logging level (panic, fatal, error, warn, info, debug, disabled)
	KeyLogLevel = "log.level"

	// KeyLogPretty output formatted log lines (for running in foreground)
	KeyLogPretty = "log.pretty"

	// KeyConfigFormat format for 'config show'
	KeyConfigFormat = "cfg_format"
	// KeyConfigForce forces overwite of existing config on `config init`
	KeyConfigForce = "cfg_force"

	// private, not in main config
	// KeyCosiID is a UUID generated for this system - use in registration
	// to uniquely identify the assets (graphs, etc.) created by this registration
	KeyCosiID = "cosi_id"
)

// getConfig dumps the current configuration and returns it
func getConfig() (*Config, error) {
	var cfg *Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "parsing config")
	}

	return cfg, nil
}

// DumpConfig prints the running configuration
func DumpConfig(w io.Writer) error {
	var cfg *Config
	var err error
	var data []byte

	cfg, err = getConfig()
	if err != nil {
		return err
	}

	format := viper.GetString(KeyConfigFormat)

	log.Debug().Str("format", format).Msg("config-show")

	switch format {
	case "json":
		data, err = json.MarshalIndent(cfg, " ", "  ")
		if err != nil {
			return errors.Wrap(err, "formatting config (json)")
		}
	case "yaml":
		data, err = yaml.Marshal(cfg)
		if err != nil {
			return errors.Wrap(err, "formatting config (yaml)")
		}
	case "toml":
		data, err = toml.Marshal(*cfg)
		if err != nil {
			return errors.Wrap(err, "formatting config (toml)")
		}
	default:
		return errors.Errorf("unknown config format '%s'", format)
	}

	_, err = fmt.Fprintf(w, "\n%s\n", data)
	return err
}

// StatConfig adds the running config to the app stats
func StatConfig() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	expvar.Publish("config", expvar.Func(func() interface{} {
		return &cfg
	}))

	return nil
}
