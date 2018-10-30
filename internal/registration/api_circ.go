// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package registration

import circapi "github.com/circonus-labs/go-apiclient"

//go:generate moq -out api_circ_test.go . CircAPI

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateCheckBundle(cfg *circapi.CheckBundle) (*circapi.CheckBundle, error)
	CreateDashboard(cfg *circapi.Dashboard) (*circapi.Dashboard, error)
	CreateGraph(cfg *circapi.Graph) (*circapi.Graph, error)
	CreateRuleSet(cfg *circapi.RuleSet) (*circapi.RuleSet, error)
	CreateWorksheet(cfg *circapi.Worksheet) (*circapi.Worksheet, error)
	DeleteCheckBundleByCID(cid circapi.CIDType) (bool, error)
	DeleteDashboardByCID(cid circapi.CIDType) (bool, error)
	DeleteGraphByCID(cid circapi.CIDType) (bool, error)
	DeleteWorksheetByCID(cid circapi.CIDType) (bool, error)
	FetchCheckBundle(cid circapi.CIDType) (*circapi.CheckBundle, error)
	FetchBroker(cid circapi.CIDType) (*circapi.Broker, error)
	FetchBrokers() (*[]circapi.Broker, error)
	FetchDashboard(cid circapi.CIDType) (*circapi.Dashboard, error)
	FetchGraph(cid circapi.CIDType) (*circapi.Graph, error)
	FetchWorksheet(cid circapi.CIDType) (*circapi.Worksheet, error)
	SearchCheckBundles(searchCriteria *circapi.SearchQueryType, filterCriteria *map[string][]string) (*[]circapi.CheckBundle, error)
	SearchDashboards(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.Dashboard, error)
	SearchGraphs(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.Graph, error)
	SearchWorksheets(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.Worksheet, error)
	UpdateCheckBundle(cfg *circapi.CheckBundle) (*circapi.CheckBundle, error)
	UpdateDashboard(cfg *circapi.Dashboard) (*circapi.Dashboard, error)
	UpdateGraph(cfg *circapi.Graph) (*circapi.Graph, error)
	UpdateWorksheet(cfg *circapi.Worksheet) (*circapi.Worksheet, error)
}
