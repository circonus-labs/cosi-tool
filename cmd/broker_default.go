// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"os"

	"github.com/circonus-labs/cosi-tool/internal/broker"
	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// brokerDefaultCmd represents the default command
var brokerDefaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Show default broker",
	Long:  `Use COSI API to show the default broker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cosiURL := viper.GetString(config.KeyCosiURL)
		optsFile := viper.GetString(config.KeyHostOptionsFile)
		cosiBID := viper.GetInt(config.KeyHostBrokerID)

		return broker.Default(cosiURL, os.Stdout, uint(cosiBID), optsFile)
	},
}

func init() {
	brokerCmd.AddCommand(brokerDefaultCmd)
}
