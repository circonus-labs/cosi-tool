// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/dashboard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dashboardDeleteCmd represents the delete command
var dashboardDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a dashboard from Circonus",
	Long: `Delete a dashboard from the Circonus system using a configuration
file or dashboard ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(dashboard.KeyCID)
		in := viper.GetString(dashboard.KeyInFile)

		return dashboard.Delete(client, id, in)
	},
}

func init() {
	dashboardCmd.AddCommand(dashboardDeleteCmd)

	{
		const (
			key         = dashboard.KeyCID
			longOpt     = "id"
			description = "Dashboard ID"
		)

		dashboardDeleteCmd.Flags().String(longOpt, dashboard.CIDDefault, description)
		_ = viper.BindPFlag(key, dashboardDeleteCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = dashboard.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of dashboard to send to Circonus API"
		)

		dashboardDeleteCmd.Flags().StringP(longOpt, shortOpt, dashboard.InFileDefault, description)
		_ = viper.BindPFlag(key, dashboardDeleteCmd.Flags().Lookup(longOpt))
	}
}
