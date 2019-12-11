// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"fmt"
	"path"
	"strings"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	"github.com/circonus-labs/cosi-server/api"
	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// templateFetchAllCmd represents the fetch command
var templateFetchAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Fetch all known templates from COSI API",
	Long: `Call COSI API to fetch all known templates and save returned objects.

If no list is supplied the default templates of check-system, dashboard-system,
worksheet-system, and any graph templates needed by analyzing the metrics
available from the agent.

Exactly what templates are available is dependent on the OS Type, distro,
version, system architecture, and whether there is a template available for the
specific metric plugin.

Templates are saved in the registration directory with a prefix of 'template-'
and an extension of '.json'.

e.g. /opt/circonus/cosi/registration/template-check-system.json

Fetch default templates:
    cosi template fetch all

Fetch templates listed:
    cosi template fetch all --list=check-system,graph-cpu
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		list := viper.GetStringSlice(templates.KeyIDList)
		show := viper.GetBool(templates.KeyShow)
		force := viper.GetBool(templates.KeyForce)

		if len(list) == 0 {
			ac, err := agentapi.New(viper.GetString(config.KeyAgentURL))
			if err != nil {
				return errors.New("creating agent API client")
			}

			metrics, err := ac.Metrics("")
			if err != nil {
				return errors.Wrap(err, "fetching available metrics from agent")
			}

			l, err := templates.DefaultTemplateList(metrics)
			if err != nil {
				return err
			}
			list = *l
		}

		if show {
			for i, t := range list {
				fmt.Printf("%02d - %s\n", i, t)
			}
			return nil
		}

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

		tt, err := tc.FetchAll(list)
		if err != nil {
			return err
		}
		for _, t := range *tt {
			if t.Err != nil {
				log.Warn().Str("template", list[t.IDX]).Err(t.Err).Msg("fetching template")
				continue
			}
			tfile := path.Join(defaults.RegPath, strings.Join([]string{"template", t.Template.Type, t.Template.Name}, "-")+".json")
			if err := regfiles.Save(tfile, t.Template, force); err != nil {
				log.Warn().Str("template", list[t.IDX]).Err(t.Err).Msg("saving template")
				continue
			}
		}

		return nil
	},
}

func init() {
	templateFetchCmd.AddCommand(templateFetchAllCmd)

	{
		const (
			key         = templates.KeyIDList
			longOpt     = "list"
			description = "Template ID list (type-name[,type-name,...] e.g. check-system,graph-cpu)"
		)

		templateFetchAllCmd.Flags().StringSlice(longOpt, []string{}, description)
		_ = viper.BindPFlag(key, templateFetchAllCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = templates.KeyForce
			longOpt     = "force"
			description = "Force save (overwrite files)"
		)

		templateFetchAllCmd.Flags().Bool(longOpt, templates.ForceDefault, description)
		_ = viper.BindPFlag(key, templateFetchAllCmd.Flags().Lookup(longOpt))
	}

	{
		const (
			key         = templates.KeyShow
			longOpt     = "show"
			description = "Show list of templates fetch would attempt to retrieve"
		)

		templateFetchAllCmd.Flags().Bool(longOpt, templates.ShowDefault, description)
		_ = viper.BindPFlag(key, templateFetchAllCmd.Flags().Lookup(longOpt))
	}
}
