// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package dashboard

//go:generate moq -out api_circ_test.go . CircAPI

import circapi "github.com/circonus-labs/go-apiclient"

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateDashboard(cfg *circapi.Dashboard) (*circapi.Dashboard, error)
	DeleteDashboardByCID(cid circapi.CIDType) (bool, error)
	FetchDashboard(cid circapi.CIDType) (*circapi.Dashboard, error)
	SearchDashboards(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.Dashboard, error)
	UpdateDashboard(cfg *circapi.Dashboard) (*circapi.Dashboard, error)
}
