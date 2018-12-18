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

func TestCreateStaticGraph(t *testing.T) {
	t.Log("Testing createStaticGraph")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	// do a little housekeeping
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "registration-graph-ignore-static-") {
			os.Remove(filepath.Join("testdata", file.Name()))
		}
	}

	badDPVar := cosiapi.TemplateConfig{
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Template: `{"metric_name":"{{.BadName}}"}`},
		},
	}

	badVDPRx := cosiapi.TemplateConfig{
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Variable: true, MetricRx: "", Template: `{"metric_name":"{{.MetricName}}"}`},
		},
	}
	badVDPVar := cosiapi.TemplateConfig{
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Variable: true, MetricRx: "^foo`([^``]+)", Template: `{"metric_name":"{{.BadName}}"}`},
		},
	}
	badVDPMulti := cosiapi.TemplateConfig{
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Variable: true, MetricRx: "^(foo).*", Template: `{"metric_name":"{{.MetricName}}"}`},
		},
	}

	okStatic := cosiapi.TemplateConfig{
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Template: `{"metric_name":"static"}`},
		},
	}
	okStaticST := cosiapi.TemplateConfig{
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Template: `{"metric_name":"baf_ding"}`},
		},
	}
	okVariable := cosiapi.TemplateConfig{
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Variable: true, MetricRx: "^foo`([^``]+)", Template: `{"check_id":{{.CheckID}}, "metric_name":"{{.MetricName}}"}`},
		},
	}
	okVariableST := cosiapi.TemplateConfig{
		Template: `{"title":"{{.HostName}} graph"}`,
		Datapoints: []cosiapi.TemplateDatapoint{
			{Variable: true, MetricRx: "^baf_(ding|dong)", Template: `{"check_id":{{.CheckID}}, "metric_name":"{{.MetricName}}"}`},
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
		{"reg exists (parse err)", "graph-test", "error", &cosiapi.TemplateConfig{}, &globalFilters{}, true, "loading registration-graph-test-error: parsing registration (testdata/registration-graph-test-error.json): unexpected end of JSON input"},
		{"reg exists", "graph-test", "valid", &cosiapi.TemplateConfig{}, &globalFilters{}, false, ""},
		{"empty template", "graph-test", "bad", &cosiapi.TemplateConfig{}, &globalFilters{}, true, "parsing graph template: invalid template config (empty)"},
		{"static template (bad template var)", "graph-test", "bad_dp_var", &badDPVar, &globalFilters{}, true, `executing template: template: graph-test-bad_dp_var-0:1:18: executing "graph-test-bad_dp_var-0" at <.BadName>: can't evaluate field BadName in type struct { HostName string; CheckID uint; NumCPU int }`},
		{"static template", "graph-ignore-static", "ok_static", &okStatic, &globalFilters{}, false, ""},
		{"static template w/ST", "graph-ignore-static", "ok_static_st", &okStaticST, &globalFilters{}, false, ""},
		{"variable template (bad dp config)", "graph-test", "bad_dp_rx", &badVDPRx, &globalFilters{}, true, `invalid variable datapoint graph-test-bad_dp_rx-bad_dp_rx:0 regex (empty)`},
		{"variable template (bad dp var)", "graph-test", "bad_dp_rx", &badVDPVar, &globalFilters{}, true, `executing template: template: graph-test-bad_dp_rx-0:1:18: executing "graph-test-bad_dp_rx-0" at <.BadName>: can't evaluate field BadName in type struct { HostName string; CheckID uint; NumCPU int; Item string; ItemIndex int; MetricName string }`},
		{"variable template (multimetric)", "graph-test", "bad_multi_metric", &badVDPMulti, &globalFilters{}, true, `invalid variable datapoint graph-test-bad_multi_metric-bad_multi_metric:0 regex (matched>1 metrics)`},
		{"variable template", "graph-ignore-static", "ok_variable", &okVariable, &globalFilters{}, false, ""},
		{"variable template w/ST", "graph-ignore-static", "ok_variable_st", &okVariableST, &globalFilters{}, false, ""},
	}

	g, err := New(&Options{
		CheckInfo: &checks.CheckInfo{CheckID: 1234},
		Client:    genMockCircAPI(),
		Config: &options.Options{
			Host: options.Host{Name: "foo"},
		},
		Metrics: &agentapi.Metrics{
			"foo`bar": agentapi.Metric{Type: "n", Value: 0},
			"foo`baz": agentapi.Metric{Type: "n", Value: 0},
			"foo`qux": agentapi.Metric{Type: "n", Value: 0},
			"baz`qux": agentapi.Metric{Type: "n", Value: 1},
			"baf_ding|ST[b\"YXJjaA==\":b\"eDg2XzY0\",b\"ZGlzdHJv\":b\"dWJ1bnR1LTE4LjA0\",b\"b3M=\":b\"bGludXg=\"]": agentapi.Metric{Type: "i", Value: 1},
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
			err := g.createStaticGraph(tst.templateID, tst.graphName, tst.cfg, tst.filters)
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
