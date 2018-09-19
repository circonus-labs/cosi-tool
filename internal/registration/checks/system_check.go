// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package checks

import (
	circapi "github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/spf13/viper"
)

func (c *Checks) createSystemCheck() (*circapi.CheckBundle, error) {
	cfgType := "check"
	cfgName := "system"
	checkID := cfgType + "-" + cfgName

	// set up the template expansion data
	type templateVars struct {
		HostName   string
		HostTarget string
	}
	tvars := templateVars{
		HostName:   c.config.Host.Name,
		HostTarget: c.config.Checks.System.Target,
	}

	cfg, err := c.parseTemplateConfig(cfgType, cfgName, tvars)
	if err != nil {
		return nil, err
	}

	if cfg.Type == "" {
		cfg.Type = "json:nad"
	}

	//
	// add cosi elements and apply any custom options config items
	//

	// set broker
	cfg.Brokers = []string{c.config.Checks.System.BrokerID}
	// set config.url = agenturl
	cfg.Config = circapi.CheckBundleConfig{"url": viper.GetString(config.KeyAgentURL)}
	// add tags
	if len(c.config.Common.Tags) > 0 {
		cfg.Tags = append(cfg.Tags, c.config.Common.Tags...)
	}
	if len(c.config.Checks.System.Tags) > 0 {
		cfg.Tags = append(cfg.Tags, c.config.Checks.System.Tags...)
	}
	// add note
	notes := c.config.Common.Notes
	if cfg.Notes != nil {
		notes += *cfg.Notes
	}
	cfg.Notes = &notes
	// add placeholder metric
	cfg.Metrics = append(cfg.Metrics, circapi.CheckBundleMetric{Name: "cosi_placeholder", Status: "active", Type: "numeric"})
	// set display name if configured in custom option
	if c.config.Checks.System.DisplayName != "" {
		cfg.DisplayName = c.config.Checks.System.DisplayName
	}

	return c.createCheck(checkID, cfg)
}
