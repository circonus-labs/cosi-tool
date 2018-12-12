// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	cosiapi "github.com/circonus-labs/cosi-server/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/rs/zerolog"
)

func TestCreateVariableGraphs(t *testing.T) {
	t.Log("Testing createVariableGraphs")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	// do a little housekeeping
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "registration-graph-ignore-variable-") {
			os.Remove(filepath.Join("testdata", file.Name()))
		}
	}

	existValid := cosiapi.TemplateConfig{
		Variable: true,
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Variable: false, MetricRx: "^item1`([^``]+)", Template: `{"metric_name":"{{.MetricName}}"}`},
		},
	}
	existError := cosiapi.TemplateConfig{
		Variable: true,
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Variable: false, MetricRx: "^item2`([^``]+)", Template: `{"metric_name":"{{.MetricName}}"}`},
		},
	}

	okMixed := cosiapi.TemplateConfig{
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Template: `{"metric_name":"static"}`},
			{Variable: true, MetricRx: "^foo`([^``]+)", Template: `{"metric_name":"{{.MetricName}}"}`},
		},
	}

	tests := []struct {
		name        string
		templateID  string
		graphName   string
		cfg         *cosiapi.TemplateConfig
		filters     *globalFilters
		shouldFail  bool
		expectedErr string
	}{
		{"invalid template id (empty)", "", "", nil, nil, true, "invalid template id (empty)"},
		{"invalid graph name (empty)", "graph-test", "", nil, nil, true, "invalid graph name (empty)"},
		{"invalid config (nil)", "graph-test", "foo", nil, nil, true, "invalid graph config (nil)"},
		{"invalid filters (nil)", "graph-test", "foo", &cosiapi.TemplateConfig{}, nil, true, "invalid global filters (nil)"},
		{"reg exists (parse err)", "graph-test", "item", &existError, &globalFilters{}, true, "loading registration-graph-test-item-error: parsing registration (testdata/registration-graph-test-item-error.json): unexpected end of JSON input"},
		{"reg exists", "graph-test", "item", &existValid, &globalFilters{}, false, ""},
		{"variable template", "graph-ignore-variable", "ok_mixed", &okMixed, &globalFilters{}, false, ""},
	}

	g, err := New(&Options{
		CheckInfo: &checks.CheckInfo{CheckID: 1234},
		Client:    genMockCircAPI(),
		Config: &options.Options{
			Host: options.Host{Name: "foo"},
		},
		Metrics: &agentapi.Metrics{
			"foo`bar":     agentapi.Metric{Type: "n", Value: 0},
			"foo`baz":     agentapi.Metric{Type: "n", Value: 0},
			"foo`qux":     agentapi.Metric{Type: "n", Value: 0},
			"baz`qux":     agentapi.Metric{Type: "n", Value: 1},
			"item1`valid": agentapi.Metric{Type: "n", Value: 1},
			"item2`error": agentapi.Metric{Type: "n", Value: 1},
		},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
	})
	if err != nil {
		t.Fatalf("unable to create graphs object (%s)", err)
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			// t.Parallel() - creates run serially, not concurrently
			err := g.createVariableGraphs(tst.templateID, tst.graphName, tst.cfg, tst.filters)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				} else if err.Error() != tst.expectedErr {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				}
			}

		})
	}

}
