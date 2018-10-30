// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package ruleset

//go:generate moq -out api_test.go . CircAPI

import circapi "github.com/circonus-labs/go-apiclient"

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateRuleSet(cfg *circapi.RuleSet) (*circapi.RuleSet, error)
	DeleteRuleSetByCID(cid circapi.CIDType) (bool, error)
	FetchRuleSet(cid circapi.CIDType) (*circapi.RuleSet, error)
	SearchRuleSets(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.RuleSet, error)
	UpdateRuleSet(cfg *circapi.RuleSet) (*circapi.RuleSet, error)
}
