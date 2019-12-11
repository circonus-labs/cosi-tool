// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"os"

	"github.com/circonus-labs/cosi-tool/internal/broker"
	"github.com/circonus-labs/cosi-tool/internal/check"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// brokerShowCmd represents the show command
var brokerShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show information for a specific broker",
	Long:  `Use Circonus API to fetch information for a specific broker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(broker.KeyCID)

		return broker.Show(client, os.Stdout, id)
	},
}

func init() {
	brokerCmd.AddCommand(brokerShowCmd)

	{
		const (
			key         = broker.KeyCID
			longOpt     = "id"
			description = "Broker ID"
		)

		brokerShowCmd.Flags().String(longOpt, check.DefaultCID, description)
		_ = viper.BindPFlag(key, brokerShowCmd.Flags().Lookup(longOpt))
	}
}
