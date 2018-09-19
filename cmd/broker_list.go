// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"os"

	"github.com/circonus-labs/cosi-tool/internal/broker"
	"github.com/spf13/cobra"
)

// brokerListCmd represents the list command
var brokerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available brokers",
	Long:  `Use Circonus API to fetch list of available brokers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return broker.DisplayList(client, os.Stdout)
	},
}

func init() {
	brokerCmd.AddCommand(brokerListCmd)
}
