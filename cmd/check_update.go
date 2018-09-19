// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/check"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkUpdateCmd represents the update command
var checkUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a check using configuration file",
	Long:  `Use Circonus API to update a check from a valid check configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		in := viper.GetString(check.KeyInFile)
		out := viper.GetString(check.KeyOutFile)
		force := viper.GetBool(check.KeyForce)

		return check.UpdateFromFile(client, in, out, force)
	},
}

func init() {
	checkCmd.AddCommand(checkUpdateCmd)

	{
		const (
			key         = check.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of check to send to Circonus API"
		)

		checkUpdateCmd.Flags().StringP(longOpt, shortOpt, check.DefaultInFile, description)
		viper.BindPFlag(key, checkUpdateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		checkUpdateCmd.Flags().StringP(longOpt, shortOpt, check.DefaultOutFile, description)
		viper.BindPFlag(key, checkUpdateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		checkUpdateCmd.Flags().Bool(longOpt, check.DefaultForce, description)
		viper.BindPFlag(key, checkUpdateCmd.Flags().Lookup(longOpt))
	}
}
