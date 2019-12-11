// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboard

import (
	"errors"
	"strings"
	"testing"

	circapi "github.com/circonus-labs/go-apiclient"
)

func genMockClient() *APIMock {
	return &APIMock{
		CreateDashboardFunc: func(cfg *circapi.Dashboard) (*circapi.Dashboard, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			return cfg, nil
		},
		DeleteDashboardByCIDFunc: func(cid circapi.CIDType) (bool, error) {
			if *cid == "/dashboard/123" {
				return true, nil
			} else if *cid == "error" {
				return false, errors.New("forced mock api call error")
			}
			return false, nil
		},
		FetchDashboardFunc: func(cid circapi.CIDType) (*circapi.Dashboard, error) {
			if *cid == "/dashboard/000" {
				return nil, errors.New("forced mock api call error")
			}
			b := circapi.Dashboard{CID: *cid}
			return &b, nil
		},
		SearchDashboardsFunc: func(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.Dashboard, error) {
			q := string(*searchCriteria)
			switch {
			case strings.Contains(q, "apierror"):
				return nil, errors.New(q)
			case strings.Contains(q, "none"):
				return &[]circapi.Dashboard{}, nil
			case strings.Contains(q, "multi"):
				return &[]circapi.Dashboard{{CID: "1"}, {CID: "2"}}, nil
			default:
				return &[]circapi.Dashboard{{CID: "/dashboard/123"}}, nil
			}
		},
		UpdateDashboardFunc: func(cfg *circapi.Dashboard) (*circapi.Dashboard, error) {
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
