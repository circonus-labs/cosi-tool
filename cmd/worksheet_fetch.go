// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/circonus-labs/cosi-tool/internal/worksheet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// worksheetFetchCmd represents the fetch command
var worksheetFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch an existing worksheet from API",
	Long: `Call Circonus API to fetch an existing worksheet and optionally save returned API object.

Fetch worksheet with an ID of 123:
    cosi worksheet fetch --id=123

Fetch worksheet with a title of 'foo bar baz':
    cosi worksheet fetch --title="foo bar baz"
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(worksheet.KeyCID)
		title := viper.GetString(worksheet.KeyTitle)
		out := viper.GetString(worksheet.KeyOutFile)
		force := viper.GetBool(worksheet.KeyForce)

		w, err := worksheet.Fetch(client, id, title)
		if err != nil {
			return err
		}
		return regfiles.Save(out, w, force)
	},
}

func init() {
	worksheetCmd.AddCommand(worksheetFetchCmd)

	{
		const (
			key         = worksheet.KeyCID
			longOpt     = "id"
			description = "Worksheet ID"
		)

		worksheetFetchCmd.Flags().String(longOpt, worksheet.CIDDefault, description)
		viper.BindPFlag(key, worksheetFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = worksheet.KeyTitle
			shortOpt    = "t"
			longOpt     = "title"
			description = "Worksheet title"
		)

		worksheetFetchCmd.Flags().StringP(longOpt, shortOpt, worksheet.TitleDefault, description)
		viper.BindPFlag(key, worksheetFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = worksheet.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		worksheetFetchCmd.Flags().StringP(longOpt, shortOpt, worksheet.OutFileDefault, description)
		viper.BindPFlag(key, worksheetFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = worksheet.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		worksheetFetchCmd.Flags().Bool(longOpt, worksheet.ForceDefault, description)
		viper.BindPFlag(key, worksheetFetchCmd.Flags().Lookup(longOpt))
	}
}
