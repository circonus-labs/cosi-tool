// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package checks

import (
	circapi "github.com/circonus-labs/go-apiclient"
)

func (c *Checks) createGroupCheck() (*circapi.CheckBundle, error) {
	cfgType := "check"
	cfgName := "group"
	checkID := cfgType + "-" + cfgName

	// set up the template expansion data
	type templateVars struct {
		HostName string
		GroupID  string
	}
	tvars := templateVars{
		HostName: c.config.Host.Name,
		GroupID:  c.config.Checks.Group.ID,
	}

	cfg, err := c.parseTemplateConfig(cfgType, cfgName, tvars)
	if err != nil {
		return nil, err
	}

	if cfg.Type == "" {
		cfg.Type = "httptrap"
	}

	//
	// add cosi elements and apply any custom options config items
	//

	// set broker
	cfg.Brokers = []string{c.config.Checks.Group.BrokerID}
	// set config.url = agenturl
	cfg.Config = circapi.CheckBundleConfig{}
	// add tags
	if len(c.config.Common.Tags) > 0 {
		cfg.Tags = append(cfg.Tags, c.config.Common.Tags...)
	}
	if len(c.config.Checks.Group.Tags) > 0 {
		cfg.Tags = append(cfg.Tags, c.config.Checks.Group.Tags...)
	}
	cfg.Tags = append(cfg.Tags, "group:"+c.config.Checks.Group.ID)
	// add note
	notes := c.config.Common.Notes
	if cfg.Notes != nil {
		notes += *cfg.Notes
	}
	cfg.Notes = &notes
	// add placeholder metric
	cfg.Metrics = append(cfg.Metrics, circapi.CheckBundleMetric{Name: "cosi_placeholder", Status: "active", Type: "numeric"})
	// set display name if configured in custom option
	if c.config.Checks.Group.DisplayName != "" {
		cfg.DisplayName = c.config.Checks.Group.DisplayName
	}

	return c.createCheck(checkID, cfg)
}
