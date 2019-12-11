// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package checks

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	circapi "github.com/circonus-labs/go-apiclient"
	circapiconf "github.com/circonus-labs/go-apiclient/config"
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
	// set trap specific settings if not set in template
	if cfg.Config == nil {
		cfg.Config = circapi.CheckBundleConfig{}
	}
	if cfg.Type == "httptrap" {
		if val, ok := cfg.Config[circapiconf.AsyncMetrics]; !ok || val == "" {
			cfg.Config[circapiconf.AsyncMetrics] = "true"
		}

		if val, ok := cfg.Config[circapiconf.Secret]; !ok || val == "" {
			s, err := genSecret()
			if err != nil {
				s = "myS3cr3t"
			}
			cfg.Config[circapiconf.Secret] = s
		}
	}
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
	// set display name if configured in custom option
	if c.config.Checks.Group.DisplayName != "" {
		cfg.DisplayName = c.config.Checks.Group.DisplayName
	}
	// default to using metric_filters
	if len(c.config.Checks.Group.MetricFilters) > 0 {
		cfg.MetricFilters = c.config.Checks.Group.MetricFilters
	} else {
		cfg.MetricFilters = [][]string{{"deny", "^$", ""}, {"allow", "^.+$", ""}}
	}
	// add placeholder metric
	// cfg.Metrics = append(cfg.Metrics, circapi.CheckBundleMetric{Name: "cosi_placeholder", Status: "active", Type: "numeric"})

	return c.createCheck(checkID, cfg)
}

func genSecret() (string, error) {
	hash := sha256.New()
	x := make([]byte, 2048)
	if _, err := rand.Read(x); err != nil {
		return "", err
	}
	if _, err := hash.Write(x); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil))[0:16], nil
}
