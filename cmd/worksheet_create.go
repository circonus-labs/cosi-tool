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

// worksheetCreateCmd represents the create command
var worksheetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a worksheet from a configuration file",
	Long:  `Use Circonus API to create a worksheet from a valid configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		in := viper.GetString(worksheet.KeyInFile)
		out := viper.GetString(worksheet.KeyOutFile)
		force := viper.GetBool(worksheet.KeyForce)

		return worksheet.CreateFromFile(client, in, out, force)
	},
}

func init() {
	worksheetCmd.AddCommand(worksheetCreateCmd)

	{
		const (
			key         = worksheet.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of worksheet to send to Circonus API"
		)

		worksheetCreateCmd.Flags().StringP(longOpt, shortOpt, worksheet.InFileDefault, description)
		_ = viper.BindPFlag(key, worksheetCreateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = worksheet.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		worksheetCreateCmd.Flags().StringP(longOpt, shortOpt, worksheet.OutFileDefault, description)
		_ = viper.BindPFlag(key, worksheetCreateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = worksheet.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		worksheetCreateCmd.Flags().Bool(longOpt, worksheet.ForceDefault, description)
		_ = viper.BindPFlag(key, worksheetCreateCmd.Flags().Lookup(longOpt))
	}
}
