// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/graph"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// graphFetchCmd represents the fetch command
var graphFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch an existing graph from API",
	Long: `Call Circonus API to fetch an existing graph and optionally save returned API object.

Fetch graph with an ID of 123:
    cosi graph fetch --id=123

Fetch graph with a title of 'foo bar baz':
    cosi graph fetch --title="foo bar baz"
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(graph.KeyCID)
		title := viper.GetString(graph.KeyTitle)
		out := viper.GetString(graph.KeyOutFile)
		force := viper.GetBool(graph.KeyForce)

		g, err := graph.Fetch(client, id, title)
		if err != nil {
			return err
		}

		return regfiles.Save(out, g, force)
	},
}

func init() {
	graphCmd.AddCommand(graphFetchCmd)

	{
		const (
			key         = graph.KeyCID
			longOpt     = "id"
			description = "Graph ID"
		)

		graphFetchCmd.Flags().String(longOpt, graph.CIDDefault, description)
		_ = viper.BindPFlag(key, graphFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = graph.KeyTitle
			shortOpt    = "t"
			longOpt     = "title"
			description = "Graph title"
		)

		graphFetchCmd.Flags().StringP(longOpt, shortOpt, graph.TitleDefault, description)
		_ = viper.BindPFlag(key, graphFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = graph.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		graphFetchCmd.Flags().StringP(longOpt, shortOpt, graph.OutFileDefault, description)
		_ = viper.BindPFlag(key, graphFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = graph.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		graphFetchCmd.Flags().Bool(longOpt, graph.ForceDefault, description)
		_ = viper.BindPFlag(key, graphFetchCmd.Flags().Lookup(longOpt))
	}
}
