// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/spf13/cobra"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Manage COSI registered dashboard(s)",
	Long:  `Intended for managing local COSI dashboards.`,
}

func init() {
	RootCmd.AddCommand(dashboardCmd)
}
