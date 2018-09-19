// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pluginCmd represents the plugin command
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage specific NAD plugins",
	Long: `Intended for managing a subset of NAD plugins
which require additional steps and/or local
configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("'plugin' has not been implemented at this time.", "It is on the schedule for implemenation in the near future.")
	},
}

func init() {
	RootCmd.AddCommand(pluginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pluginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pluginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
