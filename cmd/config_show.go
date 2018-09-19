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

	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configShowCmd represents the show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Long:  `Display current configuration and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("config.show.format") == "" {
			log.Fatal().Msg("invalid config format specified, see --format")
		}
		viper.Set(config.KeyConfigFormat, viper.GetString("config.show.format"))
		if err := config.DumpConfig(os.Stdout); err != nil {
			log.Fatal().Err(err).Msg("cosi config show")
		}
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)

	{
		const (
			key         = "config.show.format"
			shortOpt    = "f"
			longOpt     = "format"
			description = "Config format (json|toml|yaml)"
		)

		configShowCmd.Flags().StringP(longOpt, shortOpt, "", description)
		viper.BindPFlag(key, configShowCmd.Flags().Lookup(longOpt))
	}
}
