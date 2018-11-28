// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package check

import (
	"errors"
	"strings"

	circapi "github.com/circonus-labs/go-apiclient"
)

func genMockClient() *APIMock {
	return &APIMock{
		CreateCheckBundleFunc: func(cfg *circapi.CheckBundle) (*circapi.CheckBundle, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			ok := "OK"
			cfg.Notes = &ok
			return cfg, nil
		},
		DeleteCheckBundleByCIDFunc: func(cid circapi.CIDType) (bool, error) {
			if *cid == "/check_bundle/123" {
				return true, nil
			} else if *cid == "error" {
				return false, errors.New("forced mock api call error")
			}
			return false, nil
		},
		FetchCheckBundleFunc: func(cid circapi.CIDType) (*circapi.CheckBundle, error) {
			if *cid == "/check_bundle/000" {
				return nil, errors.New("forced mock api call error")
			}
			b := circapi.CheckBundle{CID: *cid}
			return &b, nil
		},
		SearchCheckBundlesFunc: func(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.CheckBundle, error) {
			q := string(*searchCriteria)
			if strings.Contains(q, "apierror") {
				return nil, errors.New(q)
			} else if strings.Contains(q, "none") {
				return &[]circapi.CheckBundle{}, nil
			} else if strings.Contains(q, "multi") {
				return &[]circapi.CheckBundle{{CID: "1"}, {CID: "2"}}, nil
			}
			return &[]circapi.CheckBundle{{CID: "/check_bundle/123"}}, nil
		},
		UpdateCheckBundleFunc: func(cfg *circapi.CheckBundle) (*circapi.CheckBundle, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			ok := "OK"
			cfg.Notes = &ok
			return cfg, nil
		},
	}
}
