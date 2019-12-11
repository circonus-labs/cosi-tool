// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"os"

	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/circonus-labs/cosi-tool/internal/dashboard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dashboardListCmd represents the list command
var dashboardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List dashboards",
	Long:  `List cosi dashboards on the local system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		uiURL := viper.GetString(config.KeyUIBaseURL)
		quiet := viper.GetBool(dashboard.KeyQuiet)
		verify := viper.GetBool(dashboard.KeyVerify)
		long := viper.GetBool(dashboard.KeyLong)

		return dashboard.List(client, os.Stdout, uiURL, defaults.RegPath, quiet, verify, long)
	},
}

func init() {
	dashboardCmd.AddCommand(dashboardListCmd)

	{
		const (
			key         = dashboard.KeyQuiet
			shortOpt    = "q"
			longOpt     = "quiet"
			description = "no header lines"
		)

		dashboardListCmd.Flags().BoolP(longOpt, shortOpt, dashboard.QuietDefault, description)
		_ = viper.BindPFlag(key, dashboardListCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = dashboard.KeyLong
			shortOpt    = "l"
			longOpt     = "long"
			description = "long listings"
		)

		dashboardListCmd.Flags().BoolP(longOpt, shortOpt, dashboard.LongDefault, description)
		_ = viper.BindPFlag(key, dashboardListCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = dashboard.KeyVerify
			shortOpt    = "v"
			longOpt     = "verify"
			description = "verify local dashboard using Circonus API"
		)

		dashboardListCmd.Flags().BoolP(longOpt, shortOpt, dashboard.VerifyDefault, description)
		_ = viper.BindPFlag(key, dashboardListCmd.Flags().Lookup(longOpt))
	}
}
