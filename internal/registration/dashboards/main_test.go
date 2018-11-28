// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboards

import (
	"testing"

	agentapi "github.com/circonus-labs/circonus-agent/api"
	"github.com/circonus-labs/cosi-tool/internal/registration/checks"
	"github.com/circonus-labs/cosi-tool/internal/registration/graphs"
	"github.com/circonus-labs/cosi-tool/internal/registration/options"
	"github.com/circonus-labs/cosi-tool/internal/templates"
	circapi "github.com/circonus-labs/go-apiclient"
	"github.com/rs/zerolog"
)

func genMockCircAPI() CircAPI {
	return &CircAPIMock{
		CreateDashboardFunc: func(cfg *circapi.Dashboard) (*circapi.Dashboard, error) {
			return cfg, nil
		},
		DeleteDashboardByCIDFunc: func(cid circapi.CIDType) (bool, error) {
			panic("TODO: mock out the DeleteDashboardByCID method")
		},
		FetchDashboardFunc: func(cid circapi.CIDType) (*circapi.Dashboard, error) {
			panic("TODO: mock out the FetchDashboard method")
		},
		SearchDashboardsFunc: func(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.Dashboard, error) {
			panic("TODO: mock out the SearchDashboards method")
		},
		UpdateDashboardFunc: func(cfg *circapi.Dashboard) (*circapi.Dashboard, error) {
			panic("TODO: mock out the UpdateDashboard method")
		},
	}
}

func TestCheckForRegistration(t *testing.T) {
	t.Log("Testing checkForRegistration")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	tests := []struct {
		name        string
		id          string
		shouldFind  bool
		shouldFail  bool
		expectedErr string
	}{
		{"invalid id (empty)", "", false, true, "invalid dashboard id (empty)"},
		{"missing", "dashboard-missing", false, false, ""},
		{"error", "dashboard-error", false, true, "loading registration-dashboard-error: parsing registration (testdata/registration-dashboard-error.json): unexpected end of JSON input"},
		{"valid", "dashboard-valid", true, false, ""},
	}

	d, err := New(&Options{
		Client:    genMockCircAPI(),
		Config:    &options.Options{},
		RegDir:    "testdata",
		Templates: &templates.Templates{},
		CheckInfo: &checks.CheckInfo{CheckID: 1234},
		GraphInfo: &map[string]graphs.GraphInfo{},
		Metrics:   &agentapi.Metrics{"test": {}},
	})
	if err != nil {
		t.Fatalf("unable to create dashboards object (%s)", err)
	}
	d.regFiles = &[]string{
		"registration-dashboard-missing.json",
		"registration-dashboard-error.json",
		"registration-dashboard-valid.json",
	}
	d.dashList = make(map[string]*circapi.Dashboard)

	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			found, err := d.checkForRegistration(tst.id)
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
				if !found {
					t.Fatal("expected true")
				}
			} else {
				if found {
					t.Fatal("expected false")
				}
			}
		})
	}
}
