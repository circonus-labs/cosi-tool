// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package check

import (
	"errors"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api"
)

func genMockClient() *APIMock {
	return &APIMock{
		CreateCheckBundleFunc: func(cfg *api.CheckBundle) (*api.CheckBundle, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			ok := "OK"
			cfg.Notes = &ok
			return cfg, nil
		},
		DeleteCheckBundleByCIDFunc: func(cid api.CIDType) (bool, error) {
			if *cid == "/check_bundle/123" {
				return true, nil
			} else if *cid == "error" {
				return false, errors.New("forced mock api call error")
			}
			return false, nil
		},
		FetchCheckBundleFunc: func(cid api.CIDType) (*api.CheckBundle, error) {
			if *cid == "/check_bundle/000" {
				return nil, errors.New("forced mock api call error")
			}
			b := api.CheckBundle{CID: *cid}
			return &b, nil
		},
		SearchCheckBundlesFunc: func(searchCriteria *api.SearchQueryType, filterCriteria *map[string][]string) (*[]api.CheckBundle, error) {
			q := string(*searchCriteria)
			if strings.Contains(q, "apierror") {
				return nil, errors.New(q)
			} else if strings.Contains(q, "none") {
				return &[]api.CheckBundle{}, nil
			} else if strings.Contains(q, "multi") {
				return &[]api.CheckBundle{{CID: "1"}, {CID: "2"}}, nil
			}
			return &[]api.CheckBundle{{CID: "/check_bundle/123"}}, nil
		},
		UpdateCheckBundleFunc: func(cfg *api.CheckBundle) (*api.CheckBundle, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			ok := "OK"
			cfg.Notes = &ok
			return cfg, nil
		},
	}
}
