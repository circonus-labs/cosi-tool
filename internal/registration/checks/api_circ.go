// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package checks

import "github.com/circonus-labs/go-apiclient"

//go:generate moq -out api_circ_test.go . CircAPI

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateCheckBundle(cfg *apiclient.CheckBundle) (*apiclient.CheckBundle, error)
	DeleteCheckBundleByCID(cid apiclient.CIDType) (bool, error)
	FetchCheckBundle(cid apiclient.CIDType) (*apiclient.CheckBundle, error)
	SearchCheckBundles(searchCriteria *apiclient.SearchQueryType, filterCriteria *apiclient.SearchFilterType) (*[]apiclient.CheckBundle, error)
	UpdateCheckBundle(cfg *apiclient.CheckBundle) (*apiclient.CheckBundle, error)
}
