// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/ruleset"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rulesetDeleteCmd represents the delete command
var rulesetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a ruleset from Circonus",
	Long: `Delete a ruleset from the Circonus system using a configuration
file or ruleset ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(ruleset.KeyCID)
		in := viper.GetString(ruleset.KeyInFile)

		return ruleset.Delete(client, id, in)
	},
}

func init() {
	rulesetCmd.AddCommand(rulesetDeleteCmd)

	{
		const (
			key         = ruleset.KeyCID
			longOpt     = "id"
			description = "Ruleset ID"
		)

		rulesetDeleteCmd.Flags().String(longOpt, ruleset.DefaultCID, description)
		_ = viper.BindPFlag(key, rulesetDeleteCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = ruleset.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of ruleset to send to Circonus API"
		)

		rulesetDeleteCmd.Flags().StringP(longOpt, shortOpt, ruleset.DefaultInFile, description)
		_ = viper.BindPFlag(key, rulesetDeleteCmd.Flags().Lookup(longOpt))
	}
}
