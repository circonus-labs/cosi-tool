// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-tool/internal/check"
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkFetchCmd represents the fetch command
var checkFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch an existing check bundle from API",
	Long: `Call Circonus API to fetch an existing check bundle and optionally save returned API object.

Fetch check bundle with an ID of 123:
    cosi check fetch --id=123

Fetch check bundle with a display name of 'foo bar baz':
    cosi check fetch --name="foo bar baz"

Fetch check bundle with target host IP 10.1.123.36:
    cosi check fetch --target=10.1.123.36

Refresh cosi system check after changes made in UI (as run from /opt/circonus/cosi):
    bin/cosi check fetch --type=system --out=registration/registration-check-system.json --force
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(check.KeyCID)
		checkType := viper.GetString(check.KeyType)
		name := viper.GetString(check.KeyName)
		target := viper.GetString(check.KeyTarget)
		out := viper.GetString(check.KeyOutFile)
		force := viper.GetBool(check.KeyForce)

		bundle, err := check.Fetch(client, defaults.RegPath, id, checkType, name, target)
		if err != nil {
			return err
		}

		return regfiles.Save(out, bundle, force)
	},
}

func init() {
	checkCmd.AddCommand(checkFetchCmd)

	{
		const (
			key         = check.KeyCID
			longOpt     = "id"
			description = "Check Bundle ID"
		)

		checkFetchCmd.Flags().String(longOpt, check.DefaultCID, description)
		viper.BindPFlag(key, checkFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyType
			shortOpt    = "k"
			longOpt     = "type"
			description = "Check type (e.g. system)"
		)

		checkFetchCmd.Flags().StringP(longOpt, shortOpt, check.DefaultType, description)
		viper.BindPFlag(key, checkFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyName
			shortOpt    = "n"
			longOpt     = "display-name"
			description = "Check display name"
		)

		checkFetchCmd.Flags().StringP(longOpt, shortOpt, check.DefaultName, description)
		viper.BindPFlag(key, checkFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyTarget
			shortOpt    = "t"
			longOpt     = "target"
			description = "Check target host (IP or FQDN)"
		)

		checkFetchCmd.Flags().StringP(longOpt, shortOpt, check.DefaultTarget, description)
		viper.BindPFlag(key, checkFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from Circonus API to output file"
		)

		checkFetchCmd.Flags().StringP(longOpt, shortOpt, check.DefaultOutFile, description)
		viper.BindPFlag(key, checkFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = check.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		checkFetchCmd.Flags().Bool(longOpt, check.DefaultForce, description)
		viper.BindPFlag(key, checkFetchCmd.Flags().Lookup(longOpt))
	}
}
