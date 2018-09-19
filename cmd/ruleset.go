// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/spf13/cobra"
)

// rulesetCmd represents the ruleset command
var rulesetCmd = &cobra.Command{
	Use:   "ruleset",
	Short: "Manage rulesets for the system check",
	Long:  `Intended for managing rulesets for the system check.`,
}

func init() {
	RootCmd.AddCommand(rulesetCmd)
}
