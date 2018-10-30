// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package worksheets

import circapi "github.com/circonus-labs/go-apiclient"

//go:generate moq -out api_circ_test.go . CircAPI

// CircAPI interface abstraction of circonus api (for mocking)
type CircAPI interface {
	CreateWorksheet(cfg *circapi.Worksheet) (*circapi.Worksheet, error)
	DeleteWorksheetByCID(cid circapi.CIDType) (bool, error)
	FetchWorksheet(cid circapi.CIDType) (*circapi.Worksheet, error)
	SearchWorksheets(searchCriteria *circapi.SearchQueryType, filterCriteria *circapi.SearchFilterType) (*[]circapi.Worksheet, error)
	UpdateWorksheet(cfg *circapi.Worksheet) (*circapi.Worksheet, error)
}
