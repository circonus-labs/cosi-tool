// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboards

import (
	"path"
	"strings"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/graphs"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Dashboards defines the registration instance
type Dashboards struct {
	dashList  map[string]*circapi.Dashboard
	client    CircAPI
	config    *options.Options
	regDir    string
	templates *templates.Templates
	checkInfo *checks.CheckInfo
	graphInfo map[string]graphs.GraphInfo
	metrics   *agentapi.Metrics
	regFiles  *[]string
	logger    zerolog.Logger
}

// Options defines the settings required to create a new instance
type Options struct {
	Client    CircAPI
	Config    *options.Options
	RegDir    string
	Templates *templates.Templates
	CheckInfo *checks.CheckInfo
	GraphInfo *map[string]graphs.GraphInfo
	Metrics   *agentapi.Metrics
}

// New creates a new Dashboards instance
func New(o *Options) (*Dashboards, error) {
	if o == nil {
		return nil, errors.New("invalid options (nil)")
	}
	if o.Client == nil {
		return nil, errors.New("invalid client (nil)")
	}
	if o.Config == nil {
		return nil, errors.New("invalid options config (nil)")
	}
	if o.RegDir == "" {
		return nil, errors.New("invalid reg dir (empty)")
	}
	if o.Templates == nil {
		return nil, errors.New("invalid templates (nil)")
	}
	if o.CheckInfo == nil {
		return nil, errors.New("invalid check info (nil)")
	}
	if o.GraphInfo == nil {
		return nil, errors.New("invalid graph info (nil)")
	}
	if o.Metrics == nil {
		return nil, errors.New("invalid metrics (nil)")
	}
	if len(*o.Metrics) == 0 {
		return nil, errors.New("invalid metrics (zero)")
	}

	regs, err := regfiles.Find(o.RegDir, "dashboard")
	if err != nil {
		return nil, errors.Wrap(err, "finding dashboard registrations")
	}

	d := Dashboards{
		dashList:  make(map[string]*circapi.Dashboard),
		client:    o.Client,
		config:    o.Config,
		regDir:    o.RegDir,
		templates: o.Templates,
		checkInfo: o.CheckInfo,
		graphInfo: *o.GraphInfo,
		metrics:   o.Metrics,
		regFiles:  regs,
		logger:    log.With().Str("cmd", "register.dashboards").Logger(),
	}

	return &d, nil
}

// Register creates dashboards using the Circonus API
func (d *Dashboards) Register(list map[string]bool) error {
	if len(list) == 0 {
		return errors.New("invalid list (empty)")
	}

	// list of _all_ things to create during registration
	// filter out non-dashboard items
	dashList := make(map[string]bool)
	for k, v := range list {
		if strings.HasPrefix(k, "dashboard-") {
			dashList[k] = v
		}
	}
	if len(dashList) == 0 {
		d.logger.Warn().Msg("0 dashboards found in list, not building ANY dashboards")
		return nil
	}

	for id, create := range dashList {
		if !create {
			continue
		}

		if loaded, err := d.checkForRegistration(id); err != nil {
			return err
		} else if loaded {
			continue
		}

		if err := d.create(id); err != nil {
			return err
		}
	}

	return nil
}

// GetMetricList returns a list of metrics used in dashboards
func (d *Dashboards) GetMetricList() *map[string]string {
	metrics := make(map[string]string)
	for _, dash := range d.dashList {
		for _, widget := range dash.Widgets {
			if widget.Type == "gauge" { // TODO: is 'gauge' the only type of widget that directly uses a metric?
				metrics[widget.Settings.MetricName] = "numeric" // TODO: need to lookup the type somewhere
			}
		}
	}
	return &metrics
}

// checkForRegistration looks through existing registration files and
// if the graphID is found, it is loaded. Returns a boolean indicating
// if the registration was found+loaded successfully or an error
func (d *Dashboards) checkForRegistration(id string) (bool, error) {
	if id == "" {
		return false, errors.New("invalid dashboard id (empty)")
	}

	regFileSig := "registration-" + id
	for _, rf := range *d.regFiles {
		if !strings.Contains(rf, regFileSig) {
			continue
		}
		var dash circapi.Dashboard
		found, err := regfiles.Load(path.Join(d.regDir, rf), &dash)
		if err != nil {
			return false, errors.Wrapf(err, "loading %s", regFileSig)
		}
		if found {
			d.dashList[regFileSig] = &dash
			return found, nil
		}
		break // we already found it but there was an issue
	}

	return false, nil
}
