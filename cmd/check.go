// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Manage COSI registered check(s)",
	Long:  `Intended for managing local COSI checks.`,
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
