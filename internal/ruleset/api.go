// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package ruleset

//go:generate moq -out api_test.go . API

import "github.com/circonus-labs/circonus-gometrics/api"

// API interface abstraction of circonus api (for mocking)
type API interface {
	CreateRuleSet(cfg *api.RuleSet) (*api.RuleSet, error)
	DeleteRuleSetByCID(cid api.CIDType) (bool, error)
	FetchRuleSet(cid api.CIDType) (*api.RuleSet, error)
	SearchRuleSets(searchCriteria *api.SearchQueryType, filterCriteria *api.SearchFilterType) (*[]api.RuleSet, error)
	UpdateRuleSet(cfg *api.RuleSet) (*api.RuleSet, error)
}
