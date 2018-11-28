// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package check

//go:generate moq -out api_circ_test.go . CircAPI

import circapi "github.com/circonus-labs/go-apiclient"

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateCheckBundle(cfg *circapi.CheckBundle) (*circapi.CheckBundle, error)
	DeleteCheckBundleByCID(cid circapi.CIDType) (bool, error)
	FetchCheckBundle(cid circapi.CIDType) (*circapi.CheckBundle, error)
	SearchCheckBundles(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.CheckBundle, error)
	UpdateCheckBundle(cfg *circapi.CheckBundle) (*circapi.CheckBundle, error)
}
