// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/worksheet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// worksheetDeleteCmd represents the delete command
var worksheetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a worksheet from Circonus",
	Long: `Delete a worksheet from the Circonus system using a configuration
file or worksheet ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(worksheet.KeyCID)
		in := viper.GetString(worksheet.KeyInFile)

		return worksheet.Delete(client, id, in)
	},
}

func init() {
	worksheetCmd.AddCommand(worksheetDeleteCmd)

	{
		const (
			key         = worksheet.KeyCID
			longOpt     = "id"
			description = "Worksheet ID"
		)

		worksheetDeleteCmd.Flags().String(longOpt, worksheet.CIDDefault, description)
		viper.BindPFlag(key, worksheetDeleteCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = worksheet.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of worksheet to send to Circonus API"
		)

		worksheetDeleteCmd.Flags().StringP(longOpt, shortOpt, worksheet.InFileDefault, description)
		viper.BindPFlag(key, worksheetDeleteCmd.Flags().Lookup(longOpt))
	}
}
