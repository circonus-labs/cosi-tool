// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agent

import (
	"errors"
	"testing"

	agentapi "github.com/circonus-labs/circonus-agent/api"
)

func genMockClient() *APIMock {
	return &APIMock{
		InventoryFunc: func() (*agentapi.Inventory, error) {
			return &agentapi.Inventory{
				agentapi.Plugin{ID: "foo`bar", Name: "foo", Instance: "bar"},
				agentapi.Plugin{ID: "foo`baz", Name: "foo", Instance: "baz"},
				agentapi.Plugin{ID: "qux", Name: "qux", Instance: ""},
			}, nil
		},
		MetricsFunc: func(pluginID string) (*agentapi.Metrics, error) {
			switch pluginID {
			case "zero":
				return &agentapi.Metrics{}, nil
			case "error":
				return nil, errors.New("forced mock api error")
			default:
				return &agentapi.Metrics{"foo`bar": agentapi.Metric{Value: 0, Type: "n"}}, nil
			}
		},
	}
}

func TestActivePlugins(t *testing.T) {
	t.Log("Testing ActivePlugins")

	c := &Client{apiClient: genMockClient()}

	p, err := c.ActivePluginList()
	if err != nil {
		t.Fatalf("unexpected error (%s)", err)
	}
	if p == nil {
		t.Fatalf("expected non-nil plugin list")
	}
	if len(*p) != 2 {
		t.Fatalf("expected 2 plugins (%#v)", p)
	}
}

func TestAvailableMetrics(t *testing.T) {
	t.Log("Testing AvailableMetrics")

	tests := []struct {
		name               string
		pluginID           string
		expectedNumMetrics int
		shouldFail         bool
		expectedErrorMsg   string
	}{
		{"invalid (error)", "error", 0, true, "forced mock api error"},
		{"zero metrics", "zero", 0, false, ""},
		{"valid", "", 1, false, ""},
	}

	c := &Client{apiClient: genMockClient()}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			m, err := c.AvailableMetrics(tst.pluginID)
			if tst.shouldFail {
				if err == nil {
					t.Fatal("expected error")
				} else if err.Error() != tst.expectedErrorMsg {
					t.Fatalf("unexpected error (%s)", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error (%s)", err)
				} else if len(*m) != tst.expectedNumMetrics {
					t.Fatalf("expected %d metric(s) %#v", tst.expectedNumMetrics, m)
				}
			}
		})
	}
}
