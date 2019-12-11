// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Package graphs handles the building of graphs
package graphs

import (
	"strings"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Graphs defines the graphs registration instance
type Graphs struct {
	graphList        map[string]circapi.Graph
	client           CircAPI
	config           *options.Options
	regDir           string
	templates        *templates.Templates
	checkInfo        *checks.CheckInfo
	metrics          *agentapi.Metrics
	shortMetricNames map[string]string // metric names w/o stream tags - value is full metric name (can be used as key into Graphs.metrics)
	regFiles         *[]string
	logger           zerolog.Logger
}

// Options defines the settings required to create a new graphs instance
type Options struct {
	CheckInfo *checks.CheckInfo
	Client    CircAPI
	Config    *options.Options
	Metrics   *agentapi.Metrics
	RegDir    string
	Templates *templates.Templates
}

// GraphInfo holds details needed for dashboards
type GraphInfo struct {
	CID  string
	UUID string
}

// New creates a new Graphs instance
func New(o *Options) (*Graphs, error) {
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
	if o.Metrics == nil {
		return nil, errors.New("invalid metrics (nil)")
	}
	if len(*o.Metrics) == 0 {
		return nil, errors.New("invalid metrics (zero)")
	}

	regs, err := regfiles.Find(o.RegDir, "graph")
	if err != nil {
		return nil, errors.Wrap(err, "finding graph registrations")
	}

	g := Graphs{
		graphList:        make(map[string]circapi.Graph),
		client:           o.Client,
		config:           o.Config,
		regDir:           o.RegDir,
		templates:        o.Templates,
		checkInfo:        o.CheckInfo,
		metrics:          o.Metrics,
		shortMetricNames: make(map[string]string),
		regFiles:         regs,
		logger:           log.With().Str("cmd", "register.graphs").Logger(),
	}

	for fullMetricName := range *g.metrics {
		shortMetricName := fullMetricName
		if idx := strings.Index(fullMetricName, "|ST["); idx > -1 {
			shortMetricName = fullMetricName[0:idx]
		}
		if shortMetricName == "" {
			continue
		}
		g.shortMetricNames[shortMetricName] = fullMetricName
	}

	return &g, nil
}

// Register creates graphs using the Circonus API
func (g *Graphs) Register(list map[string]bool) error {
	if len(list) == 0 {
		return errors.New("invalid list (empty)")
	}

	// list of _all_ things to create during registration
	// filter out non-graph items
	graphList := make(map[string]bool)
	for k, v := range list {
		if strings.HasPrefix(k, "graph-") {
			graphList[k] = v
		}
	}
	if len(graphList) == 0 {
		// this technically isn't an error. while it will render any
		// dashboard and worksheet useless it is completely within
		// the user's power to provide a list of templates with no graphs
		// while still having the dashboard and worksheet built...
		g.logger.Warn().Msg("0 graphs found in list, not building ANY graphs")
		return nil
	}

	for id, create := range graphList {
		if !create {
			g.logger.Warn().Str("id", id).Msg("Skipping, graph disabled")
			continue
		}
		err := g.create(id)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetGraphInfo returns a list of graphs and details on each graph to be used by dashbaords
func (g *Graphs) GetGraphInfo() (*map[string]GraphInfo, error) {
	if len(g.graphList) == 0 {
		return nil, errors.Errorf("invalid state (0 graphs)")
	}
	gi := make(map[string]GraphInfo)
	for gid, v := range g.graphList {
		g.logger.Debug().Str("id", gid).Str("cid", v.CID).Msg("adding graph to graph info")
		gi[gid] = GraphInfo{
			CID:  v.CID,
			UUID: strings.Replace(v.CID, "/graph/", "", 1),
		}
	}
	return &gi, nil
}

// GetMetricList returns a list of metrics used in graphs
func (g *Graphs) GetMetricList() *map[string]string {
	metrics := make(map[string]string)
	for _, graph := range g.graphList {
		for _, dp := range graph.Datapoints {
			metrics[dp.MetricName] = dp.MetricType
		}
	}
	return &metrics
}
