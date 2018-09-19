// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"os"

	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/circonus-labs/cosi-tool/internal/worksheet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// worksheetListCmd represents the list command
var worksheetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List worksheets",
	Long:  `List cosi worksheets on the local system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		uiURL := viper.GetString(config.KeyUIBaseURL)
		quiet := viper.GetBool(worksheet.KeyQuiet)
		long := viper.GetBool(worksheet.KeyLong)

		return worksheet.List(client, os.Stdout, uiURL, defaults.RegPath, quiet, long)
	},
}

func init() {
	worksheetCmd.AddCommand(worksheetListCmd)

	{
		const (
			key         = worksheet.KeyQuiet
			shortOpt    = "q"
			longOpt     = "quiet"
			description = "no header lines"
		)

		worksheetListCmd.Flags().BoolP(longOpt, shortOpt, worksheet.QuietDefault, description)
		viper.BindPFlag(key, worksheetListCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = worksheet.KeyLong
			shortOpt    = "l"
			longOpt     = "long"
			description = "long listings"
		)

		worksheetListCmd.Flags().BoolP(longOpt, shortOpt, worksheet.LongDefault, description)
		viper.BindPFlag(key, worksheetListCmd.Flags().Lookup(longOpt))
	}
}
