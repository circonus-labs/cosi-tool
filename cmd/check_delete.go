// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/check"
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkDeleteCmd represents the delete command
var checkDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a check from Circonus",
	Long: `Delete a check from the Circonus system using a configuration file,
check bundle ID, or cosi check type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(check.KeyCID)
		checkType := viper.GetString(check.KeyType)
		in := viper.GetString(check.KeyInFile)

		return check.Delete(client, defaults.RegPath, id, checkType, in)
	},
}

func init() {
	checkCmd.AddCommand(checkDeleteCmd)

	{
		const (
			key         = check.KeyCID
			longOpt     = "id"
			description = "Check Bundle ID"
		)

		checkDeleteCmd.Flags().String(longOpt, check.DefaultCID, description)
		_ = viper.BindPFlag(key, checkDeleteCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyType
			shortOpt    = "k"
			longOpt     = "type"
			description = "Check type (e.g. system)"
		)

		checkDeleteCmd.Flags().StringP(longOpt, shortOpt, check.DefaultType, description)
		_ = viper.BindPFlag(key, checkDeleteCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of check to send to Circonus API"
		)

		checkDeleteCmd.Flags().StringP(longOpt, shortOpt, check.DefaultInFile, description)
		_ = viper.BindPFlag(key, checkDeleteCmd.Flags().Lookup(longOpt))
	}
}
