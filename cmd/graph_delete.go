// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/graph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// graphDeleteCmd represents the delete command
var graphDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a graph from Circonus",
	Long: `Delete a graph from the Circonus system using a configuration
file or graph ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(graph.KeyCID)
		in := viper.GetString(graph.KeyInFile)

		return graph.Delete(client, id, in)
	},
}

func init() {
	graphCmd.AddCommand(graphDeleteCmd)

	{
		const (
			key         = graph.KeyCID
			longOpt     = "id"
			description = "Graph ID"
		)

		graphDeleteCmd.Flags().String(longOpt, graph.CIDDefault, description)
		viper.BindPFlag(key, graphDeleteCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = graph.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of graph to send to Circonus API"
		)

		graphDeleteCmd.Flags().StringP(longOpt, shortOpt, graph.InFileDefault, description)
		viper.BindPFlag(key, graphDeleteCmd.Flags().Lookup(longOpt))
	}
}
