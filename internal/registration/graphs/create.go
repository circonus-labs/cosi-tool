// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"path"
	"strings"

	"github.com/circonus-labs/cosi-tool/internal/graph"
	"github.com/circonus-labs/cosi-tool/internal/registration/regfiles"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

/* There are primarily three variations of graph templates:

   1. Static - single graph, template contains static metric names
   2. Variable metrics - single graph, 1-n variable datapoints
   3. Variable graphs - 1-n variable graphs, static datapoints
*/

// creates graph(s) for passed id using the Circonus API
func (g *Graphs) create(id string) error {
	if id == "" {
		return errors.Errorf("invalid id (empty)")
	}

	t, found, err := g.templates.Load(g.regDir, id)
	if err != nil {
		if !found {
			g.logger.Warn().Str("id", id).Msg("no template found")
			return nil
		}
		return errors.Wrap(err, "loading template")
	}

	if len(t.Configs) == 0 {
		return errors.Errorf("%s invalid template (no configs)", id)
	}

	gf := &globalFilters{
		include: compileFilters(t.Filter.Include),
		exclude: compileFilters(t.Filter.Exclude),
	}

	for graphName, cfg := range t.Configs {
		creator := g.createStaticGraph
		if cfg.Variable {
			creator = g.createVariableGraphs
		}
		if err := creator(id, graphName, &cfg, gf); err != nil {
			return err
		}
	}

	return nil
}

func (g *Graphs) createGraph(graphID string, cfg *circapi.Graph) error {
	// map short metric names to full agent metric names with dynamic stream tags
	for dpIdx, dp := range cfg.Datapoints {
		if fullMetricName, ok := g.shortMetricNames[dp.MetricName]; ok {
			// use the short metric name for display (otherwise the graph legend displays metric names with base64 encoded stream tags)
			if dp.Name == "" && dp.MetricName != fullMetricName {
				cfg.Datapoints[dpIdx].Name = dp.MetricName
			}
			g.logger.Debug().Str("graph_id", graphID).Str("graph_metric_name", dp.MetricName).Str("full_name", fullMetricName).Msg("metric mapped")
			cfg.Datapoints[dpIdx].MetricName = fullMetricName
		}
	}

	if e := log.Debug(); e.Enabled() {
		cfgFile := path.Join(g.regDir, "config-"+graphID+".json")
		g.logger.Debug().Str("cfg_file", cfgFile).Msg("saving registration config")
		if err := regfiles.Save(cfgFile, cfg, true); err != nil {
			return errors.Wrapf(err, "saving config (%s)", cfgFile)
		}
	}

	graph, err := graph.Create(g.client, cfg)
	if err != nil {
		return err
	}
	g.graphList[graphID] = *graph
	return regfiles.Save(path.Join(g.regDir, "registration-"+graphID+".json"), graph, true)
}

// checkForRegistration looks through existing registration files and
// if the graphID is found, it is loaded. Returns a boolean indicating
// if the registration was found+loaded successfully or an error
func (g *Graphs) checkForRegistration(graphID string) (bool, error) {
	// logger := log.With().Str("cmd", "register.graphs").Logger()
	if graphID == "" {
		return false, errors.New("invalid graph id (empty)")
	}

	regFileSig := "registration-" + graphID
	for _, rf := range *g.regFiles {
		if !strings.Contains(rf, regFileSig) {
			continue
		}
		var graph circapi.Graph
		found, err := regfiles.Load(path.Join(g.regDir, rf), &graph)
		if err != nil {
			return false, errors.Wrapf(err, "loading %s", regFileSig)
		}
		if found {
			g.graphList[graphID] = graph
			return found, nil
		}
		break // we already found it but there was an issue
	}

	return false, nil
}
