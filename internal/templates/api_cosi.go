// Copyright © 2018 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package templates

import cosiapi "github.com/circonus-labs/cosi-server/api"

//go:generate moq -out api_cosi_test.go . API

// CosiAPI interface abstraction of cosi server api (for mocking)
type CosiAPI interface {
	FetchTemplate(id string) (*cosiapi.Template, error)
}
