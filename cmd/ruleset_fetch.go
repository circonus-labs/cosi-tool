// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/circonus-labs/cosi-tool/internal/ruleset"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rulesetFetchCmd represents the fetch command
var rulesetFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch an existing ruleset from API",
	Long:  `Call Circonus API to fetch an existing ruleset and optionally save returned API object.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(ruleset.KeyCID)
		out := viper.GetString(ruleset.KeyOutFile)
		force := viper.GetBool(ruleset.KeyForce)

		rs, err := ruleset.Fetch(client, id)
		if err != nil {
			return err
		}

		return regfiles.Save(out, rs, force)
	},
}

func init() {
	rulesetCmd.AddCommand(rulesetFetchCmd)

	{
		const (
			key         = ruleset.KeyCID
			longOpt     = "id"
			description = "Ruleset ID"
		)

		rulesetFetchCmd.Flags().String(longOpt, ruleset.DefaultCID, description)
		_ = viper.BindPFlag(key, rulesetFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = ruleset.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		rulesetFetchCmd.Flags().StringP(longOpt, shortOpt, ruleset.DefaultOutFile, description)
		_ = viper.BindPFlag(key, rulesetFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = ruleset.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		rulesetFetchCmd.Flags().Bool(longOpt, ruleset.DefaultForce, description)
		_ = viper.BindPFlag(key, rulesetFetchCmd.Flags().Lookup(longOpt))
	}
}
