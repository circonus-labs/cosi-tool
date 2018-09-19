// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configInitCmd represents the init command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an initial default configuration file",
	Long: `Create initial configuration file using the defaults and any
settings from command line arguments/flags and environment settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if cfgFile == "" {
			return errors.Errorf("config file not set, see --config")
		}

		key := "config.init.format"
		if viper.GetString(key) == "" {
			format := filepath.Ext(cfgFile)
			if format == "" {
				return errors.Errorf("invalid file name, no extension %s", cfgFile)
			}
			viper.Set(config.KeyConfigFormat, format[1:])
		}

		if ok, err := regexp.MatchString(`^(yaml|json|toml)$`, viper.GetString(config.KeyConfigFormat)); err != nil {
			return errors.Wrap(err, "config file format regexp")
		} else if !ok {
			return errors.Errorf("invalid/unknown format (%s) for configuration, see --format", viper.GetString(config.KeyConfigFormat))
		}

		if s, err := os.Stat(cfgFile); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			if !s.Mode().IsRegular() {
				return errors.Errorf("config (%s) not a regular file", cfgFile)
			}
			if !viper.GetBool(config.KeyConfigForce) {
				return errors.Errorf("config (%s) exists, use --force to overwrite", cfgFile)
			}
		}

		f, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		return config.DumpConfig(f)
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)

	{
		const (
			key         = "config.init.format"
			shortOpt    = "f"
			longOpt     = "format"
			description = "Config format (json|toml|yaml)"
		)

		configInitCmd.Flags().StringP(longOpt, shortOpt, "", description)
		viper.BindPFlag(key, configInitCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = config.KeyConfigForce
			longOpt     = "force"
			description = "Force overwrite of existing configuration"
		)

		configInitCmd.Flags().Bool(longOpt, false, description)
		viper.BindPFlag(key, configInitCmd.Flags().Lookup(longOpt))
	}
}
