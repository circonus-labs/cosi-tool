// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/spf13/cobra"
)

// brokerCmd represents the broker command
var brokerCmd = &cobra.Command{
	Use:   "broker",
	Short: "Information about Circonus Brokers",
	Long:  `Obtain information on Circonus Brokers available to the token being used.`,
}

func init() {
	RootCmd.AddCommand(brokerCmd)
}
