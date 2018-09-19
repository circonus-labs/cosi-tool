// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Package registration provides methods for interacting with local cosi
// registration files.
package registration

import (
	"path/filepath"
	"time"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	circapi "github.com/circonus-labs/circonus-gometrics/api"
	cosiapi "github.com/circonus-labs/cosi-server/api"
	"github.com/circonus-labs/cosi-tool/internal/config"
	"github.com/circonus-labs/cosi-tool/internal/config/defaults"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/dashboards"
	"github.com/circonus-labs/cosi-tool/internal/registration/graphs"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/registration/rulesets"
	"github.com/circonus-labs/cosi-tool/internal/registration/worksheets"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	jsonCheckType = "json"
	trapCheckType = "httptrap"

	// KeyTemplateList allows user to limit the assets created based on a
	// specific list of templates. If it is empty, the default list will
	// be used. Note that the system check is always implied during
	// registration as no visuals could be created without a check.
	KeyTemplateList = "register.template_list"
)

// Registration defines the registration client
type Registration struct {
	availableMetrics      *agentapi.Metrics
	cliCirc               CircAPI
	cliCosi               CosiAPI
	config                *options.Options
	dashboardList         map[string]*circapi.Dashboard
	graphList             map[string]*circapi.Graph
	regDir                string
	rulesetList           map[string]*circapi.RuleSet
	templateList          map[string]bool // a list of templates used to limit what objects are created
	worksheetList         map[string]*circapi.Worksheet
	templates             *templates.Templates
	maxBrokerResponseTime time.Duration
	logger                zerolog.Logger
}

// Config defines the options available for --regconf config file command line option.
type Config struct {
	Checks     `json:"checks" toml:"checks" yaml:"checks"`
	Dashboards `json:"dashboards" toml:"dashboards" yaml:"dashboards"`
	Graphs     `json:"graphs" toml:"graphs" yaml:"graphs"`
	Host       `json:"host" toml:"host" yaml:"host"`
	Worksheets `json:"worksheets" toml:"worksheets" yaml:"worksheets"`
}

// Host defines the host overrides for registration
type Host struct {
	IP   string `json:"ip" toml:"ip" yaml:"ip"`
	Name string `json:"name" toml:"name" yaml:"name"`
}

// Checks defines the checks supporting overrides
type Checks struct {
	Group  GroupCheck  `json:"group" toml:"group" yaml:"group"`
	System SystemCheck `json:"system" toml:"system" yaml:"system"`
}

// SystemCheck defines the system check overrides for registration
type SystemCheck struct {
	DisplayName string `json:"display_name" toml:"display_name" yaml:"display_name"`
	Target      string `json:"target" toml:"target" yaml:"target"`
}

// GroupCheck defines the group check overrides for registration
type GroupCheck struct {
	Create      bool   `json:"create" toml:"create" yaml:"create"`
	DisplayName string `json:"display_name" toml:"display_name" yaml:"display_name"`
	ID          string `json:"id" toml:"id" yaml:"id"`
}

// Dashboards defines the dashbaords supporting overrides
type Dashboards struct {
	System SystemDashboard `json:"system" toml:"system" yaml:"system"`
}

// SystemDashboard defines the system dashboard overrides for registration
type SystemDashboard struct {
	Create bool   `json:"create" toml:"create" yaml:"create"`
	Title  string `json:"title" toml:"title" yaml:"title"`
}

// Graphs defines the graphs supporting overrides
type Graphs struct {
	Configs map[string]map[string]Graph `json:"configs" toml:"configs" yaml:"configs"`
	Exclude []string                    `json:"exclude" toml:"exclude" yaml:"exclude"`
	Include []string                    `json:"include" toml:"include" yaml:"include"`
}

// Graph defines the generic graph overrides for registration
type Graph struct {
	Title string `json:"title" toml:"title" yaml:"title"`
}

// Worksheets defines the worksheets supporting overrides
type Worksheets struct {
	System SystemWorksheet `json:"system" toml:"system" yaml:"system"`
}

// SystemWorksheet defines the system worksheet overrides for registration
type SystemWorksheet struct {
	Create bool   `json:"create" toml:"create" yaml:"create"`
	Title  string `json:"title" toml:"title" yaml:"title"`
}

// New creates a new registration client
func New(circClient CircAPI) (*Registration, error) {
	if circClient == nil {
		return nil, errors.New("invalid state, nil Circonus API client")
	}

	cosiClient, err := cosiapi.New(&cosiapi.Config{
		OSType:    viper.GetString(config.KeySystemOSType),
		OSDistro:  viper.GetString(config.KeySystemOSDistro),
		OSVersion: viper.GetString(config.KeySystemOSVersion),
		SysArch:   viper.GetString(config.KeySystemArch),
		CosiURL:   viper.GetString(config.KeyCosiURL),
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating cosi API client")
	}

	t, err := templates.New(cosiClient)
	if err != nil {
		return nil, errors.New("creating templates cosi API client")
	}

	r := &Registration{
		cliCirc:               circClient,
		cliCosi:               cosiClient,
		templates:             t,
		dashboardList:         make(map[string]*circapi.Dashboard),
		graphList:             make(map[string]*circapi.Graph),
		regDir:                defaults.RegPath,
		rulesetList:           make(map[string]*circapi.RuleSet),
		templateList:          make(map[string]bool),
		worksheetList:         make(map[string]*circapi.Worksheet),
		maxBrokerResponseTime: time.Millisecond * 500, // configurable option?
		logger:                log.With().Str("cmd", "register").Logger(),
	}

	// configure and finalize registration setup
	if err := r.configure(); err != nil {
		return nil, err
	}

	return r, nil
}

// Register initiates registering the system
func (r *Registration) Register() error {
	var gi *map[string]graphs.GraphInfo
	var ci *checks.CheckInfo

	c, err := checks.New(&checks.Options{
		Client:    r.cliCirc,
		Config:    r.config,
		RegDir:    r.regDir,
		Templates: r.templates,
	})
	if err != nil {
		return err
	}

	{ // create checks
		if err := c.Register(); err != nil {
			return err
		}

		ci, err = c.GetCheckInfo("system")
		if err != nil {
			return errors.Wrap(err, "unable to get system check info")
		}
	}

	{ // create graphs
		g, err := graphs.New(&graphs.Options{
			Client:    r.cliCirc,
			Config:    r.config,
			RegDir:    r.regDir,
			Templates: r.templates,
			CheckInfo: ci,
			Metrics:   r.availableMetrics,
		})
		if err != nil {
			return err
		}

		if err = g.Register(r.templateList); err != nil {
			return err
		}

		// enable any new metrics
		if err := c.UpdateSystemCheck(g.GetMetricList()); err != nil {
			return err
		}

		gi, err = g.GetGraphInfo()
		if err != nil {
			return err
		}
	}

	{ // create worksheet(s)
		w, err := worksheets.New(&worksheets.Options{
			Client:    r.cliCirc,
			Config:    r.config,
			RegDir:    r.regDir,
			Templates: r.templates,
		})
		if err != nil {
			return err
		}
		if err := w.Register(r.templateList); err != nil {
			return err
		}
	}

	{ // create dashboard(s)
		d, err := dashboards.New(&dashboards.Options{
			Client:    r.cliCirc,
			Config:    r.config,
			RegDir:    r.regDir,
			Templates: r.templates,
			CheckInfo: ci,
			GraphInfo: gi,
			Metrics:   r.availableMetrics,
		})
		if err != nil {
			return err
		}
		if err := d.Register(r.templateList); err != nil {
			return err
		}

		// enable any new metrics
		if err := c.UpdateSystemCheck(d.GetMetricList()); err != nil {
			return err
		}
	}

	{ // create ruleset(s)
		rs, err := rulesets.New(&rulesets.Options{
			CheckInfo:  ci,
			Client:     r.cliCirc,
			Config:     r.config,
			RegDir:     r.regDir,
			RulesetDir: filepath.Join(defaults.BasePath, "rulesets"),
		})
		if err != nil {
			return err
		}
		if err := rs.Register(); err != nil {
			return err
		}
	}

	return nil
}
