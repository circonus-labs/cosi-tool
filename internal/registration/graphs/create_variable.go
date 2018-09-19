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

// createVariableGraphs build variable graphs, one graph per "Item" found
// matching a given metric_regex - datapoint may not be variable and must have a metric_regex
// e.g. see df and disk graph templates
func (g *Graphs) createVariableGraphs(templateID, graphName string, cfg *cosiapi.TemplateConfig, gf *globalFilters) error {
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

	// 1. gather full list of items before making n copies of template (one for each item)
	items, err := g.getMatchingMetrics(cfg.Datapoints, gf)
	if err != nil {
		return errors.Wrapf(err, "gathering datapoint items for %s configs.%s", templateID, graphName)
	}

	if len(items) == 0 {
		g.logger.Warn().Str("template_id", templateID).Msg("0 metrics match regexes for datapoints, skipping")
		return nil
	}

	// one graph per "item"
	for item, metrics := range items {
		graphID := templateID + "-" + graphName + "-" + item
		g.logger.Info().Str("id", graphID).Msg("building graph")
		// 2. check for registration file
		loaded, err := g.checkForRegistration(graphID)
		if err != nil {
			return err
		}
		if loaded {
			g.logger.Info().Str("id", graphID).Msg("registration found and loaded")
			continue
		}
		gtvars := struct {
			HostName string
			CheckID  uint
			NumCPU   int
			Item     string
		}{
			g.config.Host.Name,
			g.checkInfo.CheckID,
			runtime.NumCPU(),
			item,
		}
		// 3. build base graph config (based on the "item")
		graph, err := parseGraphTemplate(graphID, cfg.Template, gtvars)
		if err != nil {
			return err
		}

		// 3b. add datapoints to base graph
		for dpIdx, dpConfig := range cfg.Datapoints {
			// static datapoint
			if dpConfig.MetricRx == "" {
				dp, err := parseDatapointTemplate(fmt.Sprintf("%s-%d", graphID, dpIdx), dpConfig.Template, gtvars)
				if err != nil {
					return err
				}
				graph.Datapoints = append(graph.Datapoints, *dp)
				continue
			}

			// datapoint based on graph "item"
			var metricName string
			for _, metric := range metrics {
				if (*metric).index == uint(dpIdx) {
					metricName = (*metric).metric
					break
				}
			}
			if metricName == "" {
				return errors.Errorf("unable to find correct metric idx:%d metrics:%#v", dpIdx, metrics)
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
				metricName,
			}
			dp, err := parseDatapointTemplate(fmt.Sprintf("%s-%d", graphID, dpIdx), dpConfig.Template, dtvars)
			if err != nil {
				return err
			}
			graph.Datapoints = append(graph.Datapoints, *dp)
		}

		notes := g.config.Common.Notes
		if graph.Notes != nil {
			notes += *graph.Notes
		}
		graph.Notes = &notes
		if len(g.config.Common.Tags) > 0 {
			graph.Tags = append(graph.Tags, g.config.Common.Tags...)
		}

		// 4. create graph
		if err := g.createGraph(graphID, graph); err != nil {
			return err
		}
	}

	return nil
}
