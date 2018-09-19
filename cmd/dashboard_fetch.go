// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/dashboard"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dashboardFetchCmd represents the fetch command
var dashboardFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch an existing dashboard from API",
	Long: `Call Circonus API to fetch an existing dashboard and optionally save returned API object.

Fetch dashboard with an ID of 123:
    cosi dashboard fetch --id=123

Fetch dashboard with a title of 'foo bar baz':
    cosi dashboard fetch --title="foo bar baz"
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(dashboard.KeyCID)
		title := viper.GetString(dashboard.KeyTitle)
		out := viper.GetString(dashboard.KeyOutFile)
		force := viper.GetBool(dashboard.KeyForce)

		d, err := dashboard.Fetch(client, id, title)
		if err != nil {
			return err
		}
		return regfiles.Save(out, d, force)
	},
}

func init() {
	dashboardCmd.AddCommand(dashboardFetchCmd)

	{
		const (
			key         = dashboard.KeyCID
			longOpt     = "id"
			description = "Dashboard ID"
		)

		dashboardFetchCmd.Flags().String(longOpt, dashboard.CIDDefault, description)
		viper.BindPFlag(key, dashboardFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = dashboard.KeyTitle
			shortOpt    = "t"
			longOpt     = "title"
			description = "Dashboard title"
		)

		dashboardFetchCmd.Flags().StringP(longOpt, shortOpt, dashboard.TitleDefault, description)
		viper.BindPFlag(key, dashboardFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = dashboard.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		dashboardFetchCmd.Flags().StringP(longOpt, shortOpt, dashboard.OutFileDefault, description)
		viper.BindPFlag(key, dashboardFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = dashboard.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		dashboardFetchCmd.Flags().Bool(longOpt, dashboard.ForceDefault, description)
		viper.BindPFlag(key, dashboardFetchCmd.Flags().Lookup(longOpt))
	}
}
