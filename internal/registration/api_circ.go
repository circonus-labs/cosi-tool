// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package registration

import "github.com/circonus-labs/circonus-gometrics/api"

//go:generate moq -out api_circ_test.go . CircAPI

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateCheckBundle(cfg *api.CheckBundle) (*api.CheckBundle, error)
	CreateDashboard(cfg *api.Dashboard) (*api.Dashboard, error)
	CreateGraph(cfg *api.Graph) (*api.Graph, error)
	CreateRuleSet(cfg *api.RuleSet) (*api.RuleSet, error)
	CreateWorksheet(cfg *api.Worksheet) (*api.Worksheet, error)
	DeleteCheckBundleByCID(cid api.CIDType) (bool, error)
	DeleteDashboardByCID(cid api.CIDType) (bool, error)
	DeleteGraphByCID(cid api.CIDType) (bool, error)
	DeleteWorksheetByCID(cid api.CIDType) (bool, error)
	FetchCheckBundle(cid api.CIDType) (*api.CheckBundle, error)
	FetchBroker(cid api.CIDType) (*api.Broker, error)
	FetchBrokers() (*[]api.Broker, error)
	FetchDashboard(cid api.CIDType) (*api.Dashboard, error)
	FetchGraph(cid api.CIDType) (*api.Graph, error)
	FetchWorksheet(cid api.CIDType) (*api.Worksheet, error)
	SearchCheckBundles(searchCriteria *api.SearchQueryType, filterCriteria *map[string][]string) (*[]api.CheckBundle, error)
	SearchDashboards(searchCriteria *api.SearchQueryType, filterCriteria *api.SearchFilterType) (*[]api.Dashboard, error)
	SearchGraphs(searchCriteria *api.SearchQueryType, filterCriteria *api.SearchFilterType) (*[]api.Graph, error)
	SearchWorksheets(searchCriteria *api.SearchQueryType, filterCriteria *api.SearchFilterType) (*[]api.Worksheet, error)
	UpdateCheckBundle(cfg *api.CheckBundle) (*api.CheckBundle, error)
	UpdateDashboard(cfg *api.Dashboard) (*api.Dashboard, error)
	UpdateGraph(cfg *api.Graph) (*api.Graph, error)
	UpdateWorksheet(cfg *api.Worksheet) (*api.Worksheet, error)
}
