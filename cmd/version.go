// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"fmt"

	"github.com/circonus-labs/cosi-tool/internal/release"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version and exit",
	Long:  "Display version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("short") {
			fmt.Printf("%s\n", release.VERSION)
			return
		}
		fmt.Printf("%s %s - commit: %s, date: %s, tag: %s\n", release.NAME, release.VERSION, release.COMMIT, release.DATE, release.TAG)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP("short", "s", false, "Short, version number only")
	_ = viper.BindPFlag("short", versionCmd.Flags().Lookup("short"))
}
