// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage COSI templates",
	Long:  `Intended for managing local templates for COSI.`,
}

func init() {
	RootCmd.AddCommand(templateCmd)
}
