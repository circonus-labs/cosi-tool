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
	circapi "github.com/circonus-labs/go-apiclient"
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

	// KeyShowConfig flags the registration options config should be dumped and
	// `cosi register` should then exit.
	KeyShowConfig = "register.show_config"
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

    r.logger.Info().Msg("registration complete")

	return nil
}
