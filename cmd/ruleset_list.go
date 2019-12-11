// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"os"

	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/circonus-labs/cosi-tool/internal/ruleset"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rulesetListCmd represents the list command
var rulesetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List rulesets",
	Long:  `List cosi rulesets on the local system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		uiURL := viper.GetString(config.KeyUIBaseURL)
		quiet := viper.GetBool(ruleset.KeyQuiet)
		long := viper.GetBool(ruleset.KeyLong)

		return ruleset.List(client, os.Stdout, uiURL, defaults.RegPath, quiet, long)
	},
}

func init() {
	rulesetCmd.AddCommand(rulesetListCmd)

	{
		const (
			key         = ruleset.KeyQuiet
			shortOpt    = "q"
			longOpt     = "quiet"
			description = "no header lines"
		)

		rulesetListCmd.Flags().BoolP(longOpt, shortOpt, ruleset.DefaultQuiet, description)
		_ = viper.BindPFlag(key, rulesetListCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = ruleset.KeyLong
			shortOpt    = "l"
			longOpt     = "long"
			description = "long listings"
		)

		rulesetListCmd.Flags().BoolP(longOpt, shortOpt, ruleset.DefaultLong, description)
		_ = viper.BindPFlag(key, rulesetListCmd.Flags().Lookup(longOpt))
	}
}
