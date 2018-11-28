// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package reset

import apiclient "github.com/circonus-labs/go-apiclient"

func genMockClient() *CircAPIMock {
	return &CircAPIMock{
		DeleteCheckBundleByCIDFunc: func(cid apiclient.CIDType) (bool, error) {
			panic("TODO: mock out the DeleteCheckBundleByCID method")
		},
		DeleteDashboardByCIDFunc: func(cid apiclient.CIDType) (bool, error) {
			panic("TODO: mock out the DeleteDashboardByCID method")
		},
		DeleteGraphByCIDFunc: func(cid apiclient.CIDType) (bool, error) {
			panic("TODO: mock out the DeleteGraphByCID method")
		},
		DeleteRuleSetByCIDFunc: func(cid apiclient.CIDType) (bool, error) {
			panic("TODO: mock out the DeleteRuleSetByCID method")
		},
		DeleteWorksheetByCIDFunc: func(cid apiclient.CIDType) (bool, error) {
			panic("TODO: mock out the DeleteWorksheetByCID method")
		},
	}
}
