// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/circonus-labs/cosi-tool/internal/reset"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset system - remove COSI created artifacts",
	Long: `Reset will delete all COSI registration created artifacts.
Checks, graphs, worksheets, rulesets, dashboards and associated registration files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return reset.Reset(client, defaults.RegPath, viper.GetBool(reset.KeyForce))
	},
}

func init() {
	RootCmd.AddCommand(resetCmd)

	{
		const (
			key         = reset.KeyForce
			longOpt     = "force"
			description = "Do not prompt for confirmation"
		)

		resetCmd.Flags().Bool(longOpt, reset.DefaultForce, description)
		viper.BindPFlag(key, resetCmd.Flags().Lookup(longOpt))
	}
}
