// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import (
	"errors"
	"testing"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	circapi "github.com/circonus-labs/circonus-gometrics/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	"github.com/rs/zerolog"
)

func genMockCircAPI() CircAPI {
	return &CircAPIMock{
		CreateGraphFunc: func(cfg *circapi.Graph) (*circapi.Graph, error) {
			if cfg.CID == "error" {
				return nil, errors.New("forced mock api error")
			}
			return cfg, nil
		},
	}
}

func TestRegister(t *testing.T) {
	t.Log("Testing Register")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	// Metrics: &agentapi.Metrics{
	// 	"foo`aaa`m1":     agentapi.Metric{Type: "n", Value: 0},
	// 	"foo`aaa`m2":     agentapi.Metric{Type: "n", Value: 0},
	// 	"foo`bbb`m1":     agentapi.Metric{Type: "n", Value: 0},
	// 	"foo`bbb`m2":     agentapi.Metric{Type: "n", Value: 0},
	// 	"foo`exclude":    agentapi.Metric{Type: "n", Value: 0},
	// 	"bar`include`m1": agentapi.Metric{Type: "n", Value: 0},
	// 	"bar`ccc`m1":     agentapi.Metric{Type: "n", Value: 0},
	// 	"bar`ccc`m2":     agentapi.Metric{Type: "n", Value: 0},
	// 	"bar`ddd`m1":     agentapi.Metric{Type: "n", Value: 0},
	// 	"bar`ddd`m2":     agentapi.Metric{Type: "n", Value: 0},
	// },

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

	empty := map[string]bool{}
	nographs := map[string]bool{"worksheet-system": true, "dashboard-system": true}
	list := map[string]bool{"graph-ignore": true, "graph-disabled": false, "dashboard-system": true}

	{
		t.Log("empty")
		err := g.Register(empty)
		if err == nil {
			t.Fatal("epxected error")
		}
		if err.Error() != "invalid list (empty)" {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
	{
		t.Log("nographs")
		err := g.Register(nographs)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
	{
		t.Log("list")
		err := g.Register(list)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
}
