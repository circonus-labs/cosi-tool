// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package check

//go:generate moq -out api_test.go . API

import "github.com/circonus-labs/circonus-gometrics/api"

// API interface abstraction of circonus api (for mocking)
type API interface {
	CreateCheckBundle(cfg *api.CheckBundle) (*api.CheckBundle, error)
	DeleteCheckBundleByCID(cid api.CIDType) (bool, error)
	FetchCheckBundle(cid api.CIDType) (*api.CheckBundle, error)
	SearchCheckBundles(searchCriteria *api.SearchQueryType, filterCriteria *map[string][]string) (*[]api.CheckBundle, error)
	UpdateCheckBundle(cfg *api.CheckBundle) (*api.CheckBundle, error)
}
