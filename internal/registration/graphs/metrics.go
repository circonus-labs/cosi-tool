// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"regexp"
	"strings"

	cosiapi "github.com/circonus-labs/cosi-server/api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type dpMetric struct {
	index  uint   // position in the graph
	metric string // metric name to be used
}

func (g *Graphs) getMatchingMetrics(datapoints []cosiapi.TemplateDatapoint, gf *globalFilters) (map[string][]*dpMetric, error) {
	if len(datapoints) == 0 {
		return nil, errors.New("invalid graph (zero datapoints)")
	}
	if len(*g.metrics) == 0 {
		return nil, errors.New("invalid instance/agent (zero available metrics)")
	}

	items := make(map[string][]*dpMetric)
	for idx, datapoint := range datapoints {
		if datapoint.MetricRx == "" {
			continue
		}

		graphInclude := compileFilters(datapoint.Filter.Include)
		if len(graphInclude) == 0 {
			graphInclude = gf.include
		}

		graphExclude := compileFilters(datapoint.Filter.Exclude)
		if len(graphExclude) == 0 {
			graphExclude = gf.exclude
		}

		metricRx, err := regexp.Compile(datapoint.MetricRx)
		if err != nil {
			return nil, errors.Wrap(err, "invalid metric_regex")
		}

		if metricRx.NumSubexp() != 1 {
			return nil, errors.Errorf("invalid regex, need 1 subexpression (%s)", datapoint.MetricRx)
		}

		for fullMetricName := range *g.metrics {
			// determine the 'short' metric name, full metric name sans any stream tags
			shortMetricName := fullMetricName
			if idx := strings.Index(fullMetricName, "|ST["); idx > -1 {
				shortMetricName = fullMetricName[0:idx]
			}
			if shortMetricName == "" {
				continue
			}

			m := metricRx.FindStringSubmatch(shortMetricName)
			if m == nil {
				continue
			}
			if len(m) != 2 {
				log.Warn().Strs("match", m).Msg("invalid match result")
				continue
			}
			item := m[1]
			if item == "" {
				continue
			}

			keepMetric := true
			if len(graphInclude) > 0 {
				keepMetric = false
				for _, rx := range graphInclude {
					if rx.MatchString(item) {
						keepMetric = true
						break
					}
				}
			}
			if keepMetric && len(graphExclude) > 0 {
				for _, rx := range graphExclude {
					if rx.MatchString(item) {
						keepMetric = false
						break
					}
				}
			}
			if keepMetric {
				items[item] = append(items[item], &dpMetric{uint(idx), fullMetricName})
			}
		}
	}

	return items, nil
}
