// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"os"

	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/circonus-labs/cosi-tool/internal/graph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// graphListCmd represents the list command
var graphListCmd = &cobra.Command{
	Use:   "list",
	Short: "List graphs",
	Long:  `List cosi graphs on the local system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		uiURL := viper.GetString(config.KeyUIBaseURL)
		quiet := viper.GetBool(graph.KeyQuiet)
		long := viper.GetBool(graph.KeyLong)

		return graph.List(client, os.Stdout, uiURL, defaults.RegPath, quiet, long)
	},
}

func init() {
	graphCmd.AddCommand(graphListCmd)

	{
		const (
			key         = graph.KeyQuiet
			shortOpt    = "q"
			longOpt     = "quiet"
			description = "no header lines"
		)

		graphListCmd.Flags().BoolP(longOpt, shortOpt, graph.QuietDefault, description)
		viper.BindPFlag(key, graphListCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = graph.KeyLong
			shortOpt    = "l"
			longOpt     = "long"
			description = "long listings"
		)

		graphListCmd.Flags().BoolP(longOpt, shortOpt, graph.LongDefault, description)
		viper.BindPFlag(key, graphListCmd.Flags().Lookup(longOpt))
	}
}
