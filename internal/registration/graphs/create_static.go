// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"fmt"
	"runtime"

	cosiapi "github.com/circonus-labs/cosi-server/api"
	"github.com/pkg/errors"
)

// createStaticGraph builds static graphs with static or variable datapoints
// e.g. see cpu template for static datapoints and diskstats template for variable datapoints
func (g *Graphs) createStaticGraph(templateID, graphName string, cfg *cosiapi.TemplateConfig, gf *globalFilters) error {
	if templateID == "" {
		return errors.New("invalid template id (empty)")
	}
	if graphName == "" {
		return errors.New("invalid graph name (empty)")
	}
	if cfg == nil {
		return errors.New("invalid graph config (nil)")
	}
	if gf == nil {
		return errors.New("invalid global filters (nil)")
	}

	graphID := templateID + "-" + graphName

	g.logger.Info().Str("id", graphID).Msg("building graph")

	// 1. check for registration file
	loaded, err := g.checkForRegistration(graphID)
	if err != nil {
		return err
	}
	if loaded {
		g.logger.Info().Str("id", graphID).Msg("registration found and loaded")
		return nil
	}

	gtvars := struct {
		HostName string
		CheckID  uint
		NumCPU   int
	}{
		g.config.Host.Name,
		g.checkInfo.CheckID,
		runtime.NumCPU(),
	}
	// 2. build base graph config
	graph, err := parseGraphTemplate(graphID, cfg.Template, gtvars)
	if err != nil {
		return errors.Wrap(err, "parsing graph template")
	}

	// 2b. add datapoints to base graph config
	for dpIdx, dpConfig := range cfg.Datapoints {
		// static datapoint
		if !dpConfig.Variable {
			dp, err := parseDatapointTemplate(fmt.Sprintf("%s-%d", graphID, dpIdx), dpConfig.Template, gtvars)
			if err != nil {
				return err
			}
			// compensate for agent returning metrics with dynamic stream tags and template has static metric name w/o stream tags
			// use the short metric name (w/o stream tags) to lookup up the full metric name with dynamic stream tags
			// if the template metric name is not found in the short metric name list, the original metric name in the template will be used
			if fullMetricName, ok := g.shortMetricNames[dp.MetricName]; ok {
				dp.MetricName = fullMetricName
			}
			graph.Datapoints = append(graph.Datapoints, *dp)
			continue
		}

		// variable datapoint - one graph, multiple datapoints based on a variable component (the "item")
		if dpConfig.MetricRx == "" {
			return errors.Errorf("invalid variable datapoint %s-%s:%d regex (empty)", graphID, graphName, dpIdx)
		}
		items, err := g.getMatchingMetrics([]cosiapi.TemplateDatapoint{dpConfig}, gf)
		if err != nil {
			return err
		}
		for item, metrics := range items {
			if len(metrics) > 1 {
				return errors.Errorf("invalid variable datapoint %s-%s:%d regex (matched>1 metrics)", graphID, graphName, dpIdx)
			}
			dtvars := struct {
				HostName   string
				CheckID    uint
				NumCPU     int
				Item       string
				ItemIndex  int
				MetricName string
			}{
				g.config.Host.Name,
				g.checkInfo.CheckID,
				runtime.NumCPU(),
				item,
				dpIdx,
				(*metrics[0]).metric,
			}
			dp, err := parseDatapointTemplate(fmt.Sprintf("%s-%d", graphID, dpIdx), dpConfig.Template, dtvars)
			if err != nil {
				return err
			}
			graph.Datapoints = append(graph.Datapoints, *dp)
		}
	}

	notes := g.config.Common.Notes
	if graph.Notes != nil {
		notes += *graph.Notes
	}
	graph.Notes = &notes
	if len(g.config.Common.Tags) > 0 {
		graph.Tags = append(graph.Tags, g.config.Common.Tags...)
	}

	// 3. create graph
	return g.createGraph(graphID, graph)
}
