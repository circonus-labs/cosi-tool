// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheets

import "github.com/circonus-labs/circonus-gometrics/api"

//go:generate moq -out api_circ_test.go . CircAPI

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateWorksheet(cfg *api.Worksheet) (*api.Worksheet, error)
	DeleteWorksheetByCID(cid api.CIDType) (bool, error)
	FetchWorksheet(cid api.CIDType) (*api.Worksheet, error)
	SearchWorksheets(searchCriteria *api.SearchQueryType, filterCriteria *api.SearchFilterType) (*[]api.Worksheet, error)
	UpdateWorksheet(cfg *api.Worksheet) (*api.Worksheet, error)
}
