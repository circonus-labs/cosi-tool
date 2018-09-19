// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboard

//go:generate moq -out api_test.go . API

import "github.com/circonus-labs/circonus-gometrics/api"

// API interface abstraction of circonus api (for mocking)
type API interface {
	CreateDashboard(cfg *api.Dashboard) (*api.Dashboard, error)
	DeleteDashboardByCID(cid api.CIDType) (bool, error)
	FetchDashboard(cid api.CIDType) (*api.Dashboard, error)
	SearchDashboards(searchCriteria *api.SearchQueryType, filterCriteria *api.SearchFilterType) (*[]api.Dashboard, error)
	UpdateDashboard(cfg *api.Dashboard) (*api.Dashboard, error)
}
