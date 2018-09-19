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

// rulesetCreateCmd represents the create command
var rulesetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a ruleset from a configuration file",
	Long:  `Use Circonus API to create a graph from a valid configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		in := viper.GetString(ruleset.KeyInFile)
		out := viper.GetString(ruleset.KeyOutFile)
		force := viper.GetBool(ruleset.KeyForce)

		return ruleset.CreateFromFile(client, in, out, force)

	},
}

func init() {
	rulesetCmd.AddCommand(rulesetCreateCmd)

	{
		const (
			key         = ruleset.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of ruleset to send to Circonus API"
		)

		rulesetCreateCmd.Flags().StringP(longOpt, shortOpt, ruleset.DefaultInFile, description)
		viper.BindPFlag(key, rulesetCreateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = ruleset.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		rulesetCreateCmd.Flags().StringP(longOpt, shortOpt, ruleset.DefaultOutFile, description)
		viper.BindPFlag(key, rulesetCreateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = ruleset.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		rulesetCreateCmd.Flags().Bool(longOpt, ruleset.DefaultForce, description)
		viper.BindPFlag(key, rulesetCreateCmd.Flags().Lookup(longOpt))
	}
}
