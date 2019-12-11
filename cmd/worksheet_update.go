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

// worksheetUpdateCmd represents the update command
var worksheetUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a worksheet using configuration file",
	Long:  `Use Circonus API to update a worksheet from a valid worksheet configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		in := viper.GetString(worksheet.KeyInFile)
		out := viper.GetString(worksheet.KeyOutFile)
		force := viper.GetBool(worksheet.KeyForce)

		return worksheet.Update(client, in, out, force)
	},
}

func init() {
	worksheetCmd.AddCommand(worksheetUpdateCmd)

	{
		const (
			key         = worksheet.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of worksheet to send to Circonus API"
		)

		worksheetUpdateCmd.Flags().StringP(longOpt, shortOpt, worksheet.InFileDefault, description)
		_ = viper.BindPFlag(key, worksheetUpdateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = worksheet.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		worksheetUpdateCmd.Flags().StringP(longOpt, shortOpt, worksheet.OutFileDefault, description)
		_ = viper.BindPFlag(key, worksheetUpdateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = worksheet.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		worksheetUpdateCmd.Flags().Bool(longOpt, worksheet.ForceDefault, description)
		_ = viper.BindPFlag(key, worksheetUpdateCmd.Flags().Lookup(longOpt))
	}
}
