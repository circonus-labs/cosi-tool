// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Manage COSI registered graph(s)",
	Long:  `Intended for managing local COSI graphs.`,
}

func init() {
	RootCmd.AddCommand(graphCmd)
}
