// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graph

//go:generate moq -out api_test.go . API

import "github.com/circonus-labs/circonus-gometrics/api"

// API interface abstraction of circonus api (for mocking)
type API interface {
	CreateGraph(cfg *api.Graph) (*api.Graph, error)
	DeleteGraphByCID(cid api.CIDType) (bool, error)
	FetchGraph(cid api.CIDType) (*api.Graph, error)
	SearchGraphs(searchCriteria *api.SearchQueryType, filterCriteria *api.SearchFilterType) (*[]api.Graph, error)
	UpdateGraph(cfg *api.Graph) (*api.Graph, error)
}
