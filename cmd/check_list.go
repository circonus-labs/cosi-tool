// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"os"

	"github.com/circonus-labs/cosi-tool/internal/check"
	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkListCmd represents the list command
var checkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List checks",
	Long:  `List cosi checks on the local system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		uiURL := viper.GetString(config.KeyUIBaseURL)
		quiet := viper.GetBool(check.KeyQuiet)
		verify := viper.GetBool(check.KeyVerify)
		long := viper.GetBool(check.KeyLong)

		return check.List(client, os.Stdout, uiURL, defaults.RegPath, quiet, verify, long)
	},
}

func init() {
	checkCmd.AddCommand(checkListCmd)

	{
		const (
			key         = check.KeyQuiet
			shortOpt    = "q"
			longOpt     = "quiet"
			description = "no header lines"
		)

		checkListCmd.Flags().BoolP(longOpt, shortOpt, check.DefaultQuiet, description)
		viper.BindPFlag(key, checkListCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyLong
			shortOpt    = "l"
			longOpt     = "long"
			description = "long listings"
		)

		checkListCmd.Flags().BoolP(longOpt, shortOpt, check.DefaultLong, description)
		viper.BindPFlag(key, checkListCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyVerify
			shortOpt    = "v"
			longOpt     = "verify"
			description = "verify local check using Circonus API"
		)

		checkListCmd.Flags().BoolP(longOpt, shortOpt, check.DefaultVerify, description)
		viper.BindPFlag(key, checkListCmd.Flags().Lookup(longOpt))
	}
}
