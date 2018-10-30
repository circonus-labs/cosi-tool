// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package reset

//go:generate moq -out api_circ_test.go . CircAPI

import circapi "github.com/circonus-labs/go-apiclient"

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	DeleteCheckBundleByCID(cid circapi.CIDType) (bool, error)
	DeleteDashboardByCID(cid circapi.CIDType) (bool, error)
	DeleteGraphByCID(cid circapi.CIDType) (bool, error)
	DeleteRuleSetByCID(cid circapi.CIDType) (bool, error)
	DeleteWorksheetByCID(cid circapi.CIDType) (bool, error)
}
