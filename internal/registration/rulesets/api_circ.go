// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package rulesets

import circapi "github.com/circonus-labs/go-apiclient"

//go:generate moq -out api_circ_test.go . CircAPI

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateRuleSet(cfg *circapi.RuleSet) (*circapi.RuleSet, error)
}
