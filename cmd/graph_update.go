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

// graphUpdateCmd represents the update command
var graphUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a graph using configuration file",
	Long:  `Use Circonus API to update a graph from a valid graph configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		in := viper.GetString(graph.KeyInFile)
		out := viper.GetString(graph.KeyOutFile)
		force := viper.GetBool(graph.KeyForce)

		return graph.Update(client, in, out, force)
	},
}

func init() {
	graphCmd.AddCommand(graphUpdateCmd)

	{
		const (
			key         = graph.KeyInFile
			shortOpt    = "i"
			longOpt     = "in"
			description = "Configuration of graph to send to Circonus API"
		)

		graphUpdateCmd.Flags().StringP(longOpt, shortOpt, graph.InFileDefault, description)
		_ = viper.BindPFlag(key, graphUpdateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = graph.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		graphUpdateCmd.Flags().StringP(longOpt, shortOpt, graph.OutFileDefault, description)
		_ = viper.BindPFlag(key, graphUpdateCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = graph.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		graphUpdateCmd.Flags().Bool(longOpt, graph.ForceDefault, description)
		_ = viper.BindPFlag(key, graphUpdateCmd.Flags().Lookup(longOpt))
	}
}
