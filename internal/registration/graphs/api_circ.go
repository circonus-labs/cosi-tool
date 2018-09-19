// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graphs

import "github.com/circonus-labs/circonus-gometrics/api"

//go:generate moq -out api_circ_test.go . CircAPI

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateGraph(cfg *api.Graph) (*api.Graph, error)
	DeleteGraphByCID(cid api.CIDType) (bool, error)
	FetchGraph(cid api.CIDType) (*api.Graph, error)
	SearchGraphs(searchCriteria *api.SearchQueryType, filterCriteria *api.SearchFilterType) (*[]api.Graph, error)
	UpdateGraph(cfg *api.Graph) (*api.Graph, error)
}
