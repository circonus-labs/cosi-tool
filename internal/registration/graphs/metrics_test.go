// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"regexp"
	"testing"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	cosiapi "github.com/circonus-labs/cosi-server/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/rs/zerolog"
)

func TestGetMatchingMetrics(t *testing.T) {
	t.Log("Testing getMatchingMetrics")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	zeroMetrics := &agentapi.Metrics{}
	zeroDatapoints := []cosiapi.TemplateDatapoint{}
	emptyGloabFilters := &globalFilters{}

	metrics := &agentapi.Metrics{
		"foo`bar": agentapi.Metric{Type: "n", Value: 0},
		"foo`baz": agentapi.Metric{Type: "n", Value: 0},
		"foo`qux": agentapi.Metric{Type: "n", Value: 0},
		"baz`qux": agentapi.Metric{Type: "n", Value: 1},
	}
	noMatchDatapoint := []cosiapi.TemplateDatapoint{
		{MetricRx: `^not_found(.+)$`}, // ensure it does not match any metric
		{}, // e.g. a static datapoint, should be skipped and not break anything
	}

	g, err := New(&Options{
		CheckInfo: &checks.CheckInfo{CheckID: 1234},
		Client:    genMockCircAPI(),
		Config: &options.Options{
			Host: options.Host{Name: "foo"},
		},
		Metrics:   &agentapi.Metrics{"test": {}},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
	})
	if err != nil {
		t.Fatalf("unable to create graphs object (%s)", err)
	}

	// no datapoints
	{
		g.metrics = zeroMetrics

		_, err := g.getMatchingMetrics(zeroDatapoints, emptyGloabFilters)
		if err == nil {
			t.Fatal("expected error")
		} else if err.Error() != "invalid graph (zero datapoints)" {
			t.Fatalf("unexpected error (%s)", err)
		}

	}
	// no available metrics
	{
		g.metrics = zeroMetrics

		_, err := g.getMatchingMetrics([]cosiapi.TemplateDatapoint{{}}, emptyGloabFilters)
		if err == nil {
			t.Fatal("expected error")
		} else if err.Error() != "invalid instance/agent (zero available metrics)" {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
	// bad dp regex (syntax)
	{
		g.metrics = metrics

		_, err := g.getMatchingMetrics([]cosiapi.TemplateDatapoint{{MetricRx: `^(bar]$`}}, emptyGloabFilters)
		if err == nil {
			t.Fatal("expected error")
		} else if err.Error() != "invalid metric_regex: error parsing regexp: missing closing ): `^(bar]$`" {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
	// bad dp regex (no subexpression)
	{
		g.metrics = metrics

		_, err := g.getMatchingMetrics([]cosiapi.TemplateDatapoint{{MetricRx: `^bar$`}}, emptyGloabFilters)
		if err == nil {
			t.Fatal("expected error")
		} else if err.Error() != "invalid regex, need 1 subexpression (^bar$)" {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
	// no matching metrics
	{
		g.metrics = metrics

		m, err := g.getMatchingMetrics(noMatchDatapoint, emptyGloabFilters)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		} else if len(m) > 0 {
			t.Fatalf("unexpected zero matching metrics (%#v)", m)
		}
	}
	// global include
	{
		dp := []cosiapi.TemplateDatapoint{
			{MetricRx: "^foo`([^`]+)"}, // will match 3
		}
		gf := &globalFilters{include: []*regexp.Regexp{regexp.MustCompile("bar")}} // will include 1 of 3

		g.metrics = metrics

		m, err := g.getMatchingMetrics(dp, gf)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		} else if len(m) != 1 {
			t.Fatal("unexpected 1 matching metric")
		}
	}
	// global exclude
	{
		dp := []cosiapi.TemplateDatapoint{
			{MetricRx: "^foo`([^`]+)"}, // will match 3
		}
		gf := &globalFilters{exclude: []*regexp.Regexp{regexp.MustCompile("bar"), regexp.MustCompile("baz")}} // will exclude 2 of 3

		g.metrics = metrics

		m, err := g.getMatchingMetrics(dp, gf)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		} else if len(m) != 1 {
			t.Fatal("unexpected 1 matching metric")
		}
	}
	// dp include
	{
		dp := []cosiapi.TemplateDatapoint{
			{
				MetricRx: "^foo`([^`]+)", // will match 3
				Filter: cosiapi.TemplateFilter{
					Include: []string{"bar"},
				},
			},
		}

		g.metrics = metrics

		m, err := g.getMatchingMetrics(dp, emptyGloabFilters)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		} else if len(m) != 1 {
			t.Fatal("unexpected 1 matching metric")
		}
	}
	// dp exclude
	{
		dp := []cosiapi.TemplateDatapoint{
			{
				MetricRx: "^foo`([^`]+)", // will match 3
				Filter: cosiapi.TemplateFilter{
					Exclude: []string{"bar", "baz"},
				},
			},
		}

		g.metrics = metrics

		m, err := g.getMatchingMetrics(dp, emptyGloabFilters)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		} else if len(m) != 1 {
			t.Fatal("unexpected 1 matching metric")
		}
	}
	// multi-match dp regex
	{
		dp := []cosiapi.TemplateDatapoint{
			{MetricRx: "^foo`([^`]+)"}, // will match 3
		}

		g.metrics = metrics

		m, err := g.getMatchingMetrics(dp, emptyGloabFilters)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		} else if len(m) != 3 {
			t.Fatal("unexpected 3 matching metrics")
		}
	}
}
