// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheet

import (
	"errors"
	"strings"

	circapi "github.com/circonus-labs/go-apiclient"
)

func genMockClient() *APIMock {
	return &APIMock{
		CreateWorksheetFunc: func(cfg *circapi.Worksheet) (*circapi.Worksheet, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			return cfg, nil
		},
		DeleteWorksheetByCIDFunc: func(cid circapi.CIDType) (bool, error) {
			if *cid == "/worksheet/123" {
				return true, nil
			} else if *cid == "error" {
				return false, errors.New("forced mock api call error")
			}
			return false, nil
		},
		FetchWorksheetFunc: func(cid circapi.CIDType) (*circapi.Worksheet, error) {
			if *cid == "/worksheet/000" {
				return nil, errors.New("forced mock api call error")
			}
			b := circapi.Worksheet{CID: *cid}
			return &b, nil
		},
		SearchWorksheetsFunc: func(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.Worksheet, error) {
			q := string(*searchCriteria)
			if strings.Contains(q, "apierror") {
				return nil, errors.New(q)
			} else if strings.Contains(q, "none") {
				return &[]circapi.Worksheet{}, nil
			} else if strings.Contains(q, "multi") {
				return &[]circapi.Worksheet{{CID: "1"}, {CID: "2"}}, nil
			}
			return &[]circapi.Worksheet{{CID: "/worksheet/123"}}, nil
		},
		UpdateWorksheetFunc: func(cfg *circapi.Worksheet) (*circapi.Worksheet, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			return cfg, nil
		},
	}
}
