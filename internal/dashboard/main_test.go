// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboard

import (
	"errors"
	"strings"
	"testing"

	"github.com/circonus-labs/circonus-gometrics/api"
)

func genMockClient() *APIMock {
	return &APIMock{
		CreateDashboardFunc: func(cfg *api.Dashboard) (*api.Dashboard, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			return cfg, nil
		},
		DeleteDashboardByCIDFunc: func(cid api.CIDType) (bool, error) {
			if *cid == "/dashboard/123" {
				return true, nil
			} else if *cid == "error" {
				return false, errors.New("forced mock api call error")
			}
			return false, nil
		},
		FetchDashboardFunc: func(cid api.CIDType) (*api.Dashboard, error) {
			if *cid == "/dashboard/000" {
				return nil, errors.New("forced mock api call error")
			}
			b := api.Dashboard{CID: *cid}
			return &b, nil
		},
		SearchDashboardsFunc: func(searchCriteria *api.SearchQueryType, filterCriteria *api.SearchFilterType) (*[]api.Dashboard, error) {
			q := string(*searchCriteria)
			if strings.Contains(q, "apierror") {
				return nil, errors.New(q)
			} else if strings.Contains(q, "none") {
				return &[]api.Dashboard{}, nil
			} else if strings.Contains(q, "multi") {
				return &[]api.Dashboard{{CID: "1"}, {CID: "2"}}, nil
			}
			return &[]api.Dashboard{{CID: "/dashboard/123"}}, nil
		},
		UpdateDashboardFunc: func(cfg *api.Dashboard) (*api.Dashboard, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			return cfg, nil
		},
	}
}

func Test(t *testing.T) {
	t.Log("placeholder")
}
