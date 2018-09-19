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

// dashboardCreateCmd represents the create command
var dashboardCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a dashboard from a configuration file",
	Long:  `Use Circonus API to create a dashboard from a valid configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		in := viper.GetString(dashboard.KeyInFile)
		out := viper.GetString(dashboard.KeyOutFile)
		force := viper.GetBool(dashboard.KeyForce)

		return dashboard.CreateFromFile(client, in, out, force)
	},
}

func init() {
	dashboardCmd.AddCommand(dashboardCreateCmd)

	{
		const (
			key         = dashboard.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of dashboard to send to Circonus API"
		)

		dashboardCreateCmd.Flags().StringP(longOpt, shortOpt, dashboard.InFileDefault, description)
		viper.BindPFlag(key, dashboardCreateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = dashboard.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		dashboardCreateCmd.Flags().StringP(longOpt, shortOpt, dashboard.OutFileDefault, description)
		viper.BindPFlag(key, dashboardCreateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = dashboard.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		dashboardCreateCmd.Flags().Bool(longOpt, dashboard.ForceDefault, description)
		viper.BindPFlag(key, dashboardCreateCmd.Flags().Lookup(longOpt))
	}
}
