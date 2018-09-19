// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// worksheetCmd represents the worksheet command
var worksheetCmd = &cobra.Command{
	Use:   "worksheet",
	Short: "Manage COSI registered worksheet(s)",
	Long:  `Intended for managing local COSI worksheets.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("worksheet called")
	},
}

func init() {
	RootCmd.AddCommand(worksheetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// worksheetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// worksheetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
