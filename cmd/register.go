// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"os"

	"github.com/circonus-labs/cosi-tool/internal/registration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "COSI registration of this system",
	Long: `Register this system using COSI method.

Create a system check and optional group check.
Create graphs for known plugins.
Create system worksheet.
Create system dashboard.
Create rulesets for system check.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := registration.New(client)
		if err != nil {
			return err
		}
		if err := r.Register(); err != nil {
			logger := log.With().Str("cmd", "register").Logger()
			logger.Fatal().Err(err).Msg("unable to complete registration")
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(registerCmd)

	{
		const (
			key         = registration.KeyTemplateList
			longOpt     = "templates"
			description = "Template ID list (type-name[,type-name,...] e.g. check-system,graph-cpu)"
		)
		registerCmd.Flags().StringSlice(longOpt, []string{}, description)
		_ = viper.BindPFlag(key, registerCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = registration.KeyShowConfig
			longOpt     = "show-config"
			defaultFmt  = ""
			description = "Show registration options configuration using format yaml|json|toml"
		)
		registerCmd.Flags().String(longOpt, defaultFmt, description)
		_ = viper.BindPFlag(key, registerCmd.Flags().Lookup(longOpt))
	}
}
