// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package registration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/circonus-labs/go-apiclient"
	"github.com/pkg/errors"
)

// helpers

// genMockAgent creates a mock agent server and returns the server and a closer
// usage:
//   ta, taClose := testAgent(t)
//   defer taClose()
//   ... perform tests using ta server
func genMockAgent(t *testing.T) (*httptest.Server, func()) {
	ta := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rp := r.URL.Path
		fmt.Fprintln(w, "not implemented yet", rp)
	}))
	return ta, func() { ta.Close() }
}

func genMockCircAPI() *CircAPIMock {
	return &CircAPIMock{
		CreateCheckBundleFunc: func(cfg *apiclient.CheckBundle) (*apiclient.CheckBundle, error) {
			panic("TODO: mock out the CreateCheckBundle method")
		},
		CreateDashboardFunc: func(cfg *apiclient.Dashboard) (*apiclient.Dashboard, error) {
			panic("TODO: mock out the CreateDashboard method")
		},
		CreateGraphFunc: func(cfg *apiclient.Graph) (*apiclient.Graph, error) {
			panic("TODO: mock out the CreateGraph method")
		},
		CreateWorksheetFunc: func(cfg *apiclient.Worksheet) (*apiclient.Worksheet, error) {
			panic("TODO: mock out the CreateWorksheet method")
		},
		DeleteCheckBundleByCIDFunc: func(cid apiclient.CIDType) (bool, error) {
			panic("TODO: mock out the DeleteCheckBundleByCID method")
		},
		FetchBrokerFunc: func(cid apiclient.CIDType) (*apiclient.Broker, error) {
			panic("TODO: mock out the FetchBroker method")
		},
		FetchBrokersFunc: func() (*[]apiclient.Broker, error) {
			return &[]apiclient.Broker{}, nil
		},
		SearchCheckBundlesFunc: func(searchCriteria *apiclient.SearchQueryType, filterCriteria *map[string][]string) (*[]apiclient.CheckBundle, error) {
			panic("TODO: mock out the SearchCheckBundles method")
		},
		UpdateCheckBundleFunc: func(cfg *apiclient.CheckBundle) (*apiclient.CheckBundle, error) {
			panic("TODO: mock out the UpdateCheckBundle method")
		},
	}
}

func genMockCosiAPI() *CosiAPIMock {
	return &CosiAPIMock{
		FetchBrokerFunc: func(checkType string) (string, error) {
			if checkType == jsonCheckType {
				return "/broker/1", nil
			} else if checkType == "httptrap" {
				return "/broker/2", nil
			} else {
				return "", errors.Errorf("unknown check type (%s)", checkType)
			}
		},
	}
}
