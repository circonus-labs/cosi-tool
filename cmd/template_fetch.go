// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"github.com/circonus-labs/cosi-server/api"
	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// templateFetchCmd represents the fetch command
var templateFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch an existing template from COSI API",
	Long: `Call COSI API to fetch an existing template and optionally save returned object.

Fetch template with an ID of check-system:
    cosi template fetch --id=check-system

Fetch all templates applicable to current system:
    cosi template fetch --all
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetString(templates.KeyID)
		out := viper.GetString(templates.KeyOutFile)
		force := viper.GetBool(templates.KeyForce)

		osType := viper.GetString(config.KeySystemOSType)
		osDist := viper.GetString(config.KeySystemOSDistro)
		osVers := viper.GetString(config.KeySystemOSVersion)
		sysArch := viper.GetString(config.KeySystemArch)
		cosiURL := viper.GetString(config.KeyCosiURL)

		client, err := api.New(&api.Config{
			OSType:    osType,
			OSDistro:  osDist,
			OSVersion: osVers,
			SysArch:   sysArch,
			CosiURL:   cosiURL,
		})
		if err != nil {
			return errors.Wrap(err, "creating cosi-server client")
		}

		tc, err := templates.New(client)
		if err != nil {
			return err
		}

		t, err := tc.Fetch(id)
		if err != nil {
			return err
		}

		return regfiles.Save(out, t, force)
	},
}

func init() {
	templateCmd.AddCommand(templateFetchCmd)

	{
		const (
			key         = templates.KeyID
			longOpt     = "id"
			description = "Template ID (type-name e.g. check-system, graph-cpu)"
		)

		templateFetchCmd.Flags().String(longOpt, templates.IDDefault, description)
		_ = viper.BindPFlag(key, templateFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = templates.KeyOutFile
			shortOpt    = "o"
			longOpt     = "out"
			description = "Save response from COSI API to output file (if empty, default={regdir}/template-<id>.json)"
		)

		templateFetchCmd.Flags().StringP(longOpt, shortOpt, templates.OutFileDefault, description)
		_ = viper.BindPFlag(key, templateFetchCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = templates.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite output file)"
		)

		templateFetchCmd.Flags().Bool(longOpt, templates.ForceDefault, description)
		_ = viper.BindPFlag(key, templateFetchCmd.Flags().Lookup(longOpt))
	}
}
