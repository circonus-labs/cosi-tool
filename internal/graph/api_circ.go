// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graph

//go:generate moq -out api_circ_test.go . CircAPI

import circapi "github.com/circonus-labs/go-apiclient"

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateGraph(cfg *circapi.Graph) (*circapi.Graph, error)
	DeleteGraphByCID(cid circapi.CIDType) (bool, error)
	FetchGraph(cid circapi.CIDType) (*circapi.Graph, error)
	SearchGraphs(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.Graph, error)
	UpdateGraph(cfg *circapi.Graph) (*circapi.Graph, error)
}
