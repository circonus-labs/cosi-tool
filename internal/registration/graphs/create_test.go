// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"testing"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/rs/zerolog"
)

func TestCreate(t *testing.T) {
	t.Log("Testing Create")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	g, err := New(&Options{
		CheckInfo: &checks.CheckInfo{CheckID: 1234},
		Client:    genMockCircAPI(),
		Config: &options.Options{
			Host: options.Host{Name: "foo"},
		},
		Metrics: &agentapi.Metrics{
			"foo`aaa`m1":     agentapi.Metric{Type: "n", Value: 0},
			"foo`aaa`m2":     agentapi.Metric{Type: "n", Value: 0},
			"foo`bbb`m1":     agentapi.Metric{Type: "n", Value: 0},
			"foo`bbb`m2":     agentapi.Metric{Type: "n", Value: 0},
			"foo`exclude":    agentapi.Metric{Type: "n", Value: 0},
			"bar`include`m1": agentapi.Metric{Type: "n", Value: 0},
			"bar`ccc`m1":     agentapi.Metric{Type: "n", Value: 0},
			"bar`ccc`m2":     agentapi.Metric{Type: "n", Value: 0},
			"bar`ddd`m1":     agentapi.Metric{Type: "n", Value: 0},
			"bar`ddd`m2":     agentapi.Metric{Type: "n", Value: 0},
		},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
	})
	if err != nil {
		t.Fatalf("unable to create graphs object (%s)", err)
	}

	{
		t.Log("empty id")
		err := g.create("")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "invalid id (empty)" {
			t.Fatalf("unexpected error (%s)", err)
		}
	}

	{
		t.Log("graph-ignore")
		err := g.create("graph-ignore")
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
}

func TestCheckForRegistration(t *testing.T) {
	t.Log("Testing checkForRegistration")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	g, err := New(&Options{
		CheckInfo: &checks.CheckInfo{CheckID: 1234},
		Client:    genMockCircAPI(),
		Config: &options.Options{
			Host: options.Host{Name: "foo"},
		},
		Metrics: &agentapi.Metrics{
			"foo`aaa`m1":     agentapi.Metric{Type: "n", Value: 0},
			"foo`aaa`m2":     agentapi.Metric{Type: "n", Value: 0},
			"foo`bbb`m1":     agentapi.Metric{Type: "n", Value: 0},
			"foo`bbb`m2":     agentapi.Metric{Type: "n", Value: 0},
			"foo`exclude":    agentapi.Metric{Type: "n", Value: 0},
			"bar`include`m1": agentapi.Metric{Type: "n", Value: 0},
			"bar`ccc`m1":     agentapi.Metric{Type: "n", Value: 0},
			"bar`ccc`m2":     agentapi.Metric{Type: "n", Value: 0},
			"bar`ddd`m1":     agentapi.Metric{Type: "n", Value: 0},
			"bar`ddd`m2":     agentapi.Metric{Type: "n", Value: 0},
		},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
	})
	if err != nil {
		t.Fatalf("unable to create graphs object (%s)", err)
	}

	tests := []struct {
		name        string
		graphID     string
		shouldFind  bool
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", false, true, "invalid graph id (empty)"},
		{"missing", "graph-missing", false, false, ""},
		{"error (parsing)", "graph-error", false, true, "loading registration-graph-error: parsing registration (testdata/registration-graph-error.json): unexpected end of JSON input"},
		{"valid", "graph-valid", true, false, ""},
	}

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			// t.Parallel() - creates run serially, not concurrently
			loaded, err := g.checkForRegistration(tst.graphID)
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
			if tst.shouldFind {
				if !loaded {
					t.Fatal("expected found+loaded")
				}
			} else {
				if loaded {
					t.Fatal("expected not found/loaded")
				}
			}
		})
	}
}
