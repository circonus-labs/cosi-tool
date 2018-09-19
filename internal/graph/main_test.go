// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package graph

import (
	"errors"
	"strings"
	"testing"

	"github.com/circonus-labs/circonus-gometrics/api"
)

func genMockClient() *APIMock {
	return &APIMock{
		CreateGraphFunc: func(cfg *api.Graph) (*api.Graph, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			return cfg, nil
		},
		DeleteGraphByCIDFunc: func(cid api.CIDType) (bool, error) {
			if *cid == "/graph/123" {
				return true, nil
			} else if *cid == "error" {
				return false, errors.New("forced mock api call error")
			}
			return false, nil
		},
		FetchGraphFunc: func(cid api.CIDType) (*api.Graph, error) {
			if *cid == "/graph/000" {
				return nil, errors.New("forced mock api call error")
			}
			b := api.Graph{CID: *cid}
			return &b, nil
		},
		SearchGraphsFunc: func(searchCriteria *api.SearchQueryType, filterCriteria *api.SearchFilterType) (*[]api.Graph, error) {
			q := string(*searchCriteria)
			if strings.Contains(q, "apierror") {
				return nil, errors.New(q)
			} else if strings.Contains(q, "none") {
				return &[]api.Graph{}, nil
			} else if strings.Contains(q, "multi") {
				return &[]api.Graph{{CID: "1"}, {CID: "2"}}, nil
			}
			return &[]api.Graph{{CID: "/graph/123"}}, nil
		},
		UpdateGraphFunc: func(cfg *api.Graph) (*api.Graph, error) {
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
