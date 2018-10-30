// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package ruleset

import (
	"errors"
	"strings"

	circapi "github.com/circonus-labs/go-apiclient"
)

func genMockClient() *APIMock {
	return &APIMock{
		CreateRuleSetFunc: func(cfg *circapi.RuleSet) (*circapi.RuleSet, error) {
			if strings.Contains(cfg.CID, "error") {
				return nil, errors.New("forced mock api call error")
			}
			return cfg, nil
		},
		DeleteRuleSetByCIDFunc: func(cid circapi.CIDType) (bool, error) {
			if *cid == "/rule_set/123_test_metric" {
				return true, nil
			} else if *cid == "error" {
				return false, errors.New("forced mock api call error")
			}
			return false, nil
		},
		FetchRuleSetFunc: func(cid circapi.CIDType) (*circapi.RuleSet, error) {
			if *cid == "/rule_set/000_error" {
				return nil, errors.New("forced mock api call error")
			}
			b := circapi.RuleSet{CID: *cid}
			return &b, nil
		},
		SearchRuleSetsFunc: func(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.RuleSet, error) {
			panic("TODO: mock out the SearchRuleSets method")
		},
		UpdateRuleSetFunc: func(cfg *circapi.RuleSet) (*circapi.RuleSet, error) {
			panic("TODO: mock out the UpdateRuleSet method")
		},
	}
}
