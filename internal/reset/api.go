// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package reset

//go:generate moq -out api_test.go . API

import "github.com/circonus-labs/circonus-gometrics/api"

// API interface abstraction of circonus api (for mocking)
type API interface {
	DeleteCheckBundleByCID(cid api.CIDType) (bool, error)
	DeleteDashboardByCID(cid api.CIDType) (bool, error)
	DeleteGraphByCID(cid api.CIDType) (bool, error)
	DeleteRuleSetByCID(cid api.CIDType) (bool, error)
	DeleteWorksheetByCID(cid api.CIDType) (bool, error)
}
